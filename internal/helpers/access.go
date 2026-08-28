package helpers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/reshap0318/go-boilerplate/internal/database"
	"github.com/reshap0318/go-boilerplate/internal/models"
	"gorm.io/gorm"
)

type userAccessData struct {
	permissions map[string]bool
	roles       map[string]bool
}

// accessCacheEntry is the JSON shape stored under the shared "session:%d" key.
// It matches dtos.UserDTO's roles/permissions fields, so it reads the same
// entry AuthLogin writes and writes back a shape AuthValidateToken can still read.
type accessCacheEntry struct {
	Roles []struct {
		Name string `json:"name"`
	} `json:"roles"`
	Permissions []struct {
		Name string `json:"name"`
	} `json:"permissions"`
}

// Access handles permission and role checking.
// Redis is used exclusively when available; the local in-memory cache is the
// fallback only while Redis is down. Either way, a miss falls back to the DB.
type Access struct {
	redis *database.RedisCache
	db    *gorm.DB
	mu    sync.RWMutex
	cache map[uint]*userAccessData
}

// NewAccess creates a new Access instance.
func NewAccess(redis *database.RedisCache, db *gorm.DB) *Access {
	return &Access{
		redis: redis,
		db:    db,
		cache: make(map[uint]*userAccessData),
	}
}

// getUserAccess retrieves user permissions and roles.
// Redis alive: Redis → DB. Redis down: local cache → DB.
func (a *Access) getUserAccess(userID uint) (*userAccessData, bool) {
	if a.redis != nil && a.redis.IsCacheAvailable() {
		return a.getUserAccessFromRedis(userID)
	}
	return a.getUserAccessFromLocal(userID)
}

func (a *Access) getUserAccessFromRedis(userID uint) (*userAccessData, bool) {
	key := fmt.Sprintf("session:%d", userID)

	var cached accessCacheEntry
	if err := a.redis.GetJSON(key, &cached); err == nil {
		return accessDataFromCacheEntry(&cached), true
	}

	user, err := a.findUserWithRolesPermissions(userID)
	if err != nil {
		return nil, false
	}
	data := a.buildAccessDataFromUser(user)

	ttl := time.Duration(GetEnvInt("JWT_EXPIRATION", 24)) * time.Hour
	_ = a.redis.SetJSON(key, accessCacheEntryFromData(data), ttl)

	return data, true
}

func (a *Access) getUserAccessFromLocal(userID uint) (*userAccessData, bool) {
	a.mu.RLock()
	data, ok := a.cache[userID]
	a.mu.RUnlock()
	if ok {
		return data, true
	}

	user, err := a.findUserWithRolesPermissions(userID)
	if err != nil {
		return nil, false
	}
	data = a.buildAccessDataFromUser(user)

	a.mu.Lock()
	a.cache[userID] = data
	a.mu.Unlock()

	return data, true
}

func accessDataFromCacheEntry(cached *accessCacheEntry) *userAccessData {
	data := &userAccessData{
		permissions: make(map[string]bool),
		roles:       make(map[string]bool),
	}
	for _, r := range cached.Roles {
		data.roles[r.Name] = true
	}
	for _, p := range cached.Permissions {
		data.permissions[p.Name] = true
	}
	return data
}

func accessCacheEntryFromData(data *userAccessData) *accessCacheEntry {
	entry := &accessCacheEntry{}
	for role := range data.roles {
		entry.Roles = append(entry.Roles, struct {
			Name string `json:"name"`
		}{Name: role})
	}
	for perm := range data.permissions {
		entry.Permissions = append(entry.Permissions, struct {
			Name string `json:"name"`
		}{Name: perm})
	}
	return entry
}

// findUserWithRolesPermissions fetches user with roles and permissions from DB.
func (a *Access) findUserWithRolesPermissions(userID uint) (*models.User, error) {
	var user models.User
	err := a.db.Preload("Roles.Permissions").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *Access) buildAccessDataFromUser(user *models.User) *userAccessData {
	data := &userAccessData{
		permissions: make(map[string]bool),
		roles:       make(map[string]bool),
	}
	for _, r := range user.Roles {
		data.roles[r.Name] = true
		for _, p := range r.Permissions {
			data.permissions[p.Name] = true
		}
	}
	return data
}

// HasPermission checks if the caller has ANY of the specified permissions.
func (a *Access) HasPermission(ctx context.Context, permissions ...string) bool {
	userID := GetCallerID(ctx)
	if userID == 0 {
		return false
	}

	data, ok := a.getUserAccess(userID)
	if !ok {
		return false
	}

	for _, perm := range permissions {
		if data.permissions[perm] {
			return true
		}
	}

	return false
}

// HasRole checks if the caller has the specified role.
func (a *Access) HasRole(ctx context.Context, role string) bool {
	userID := GetCallerID(ctx)
	if userID == 0 {
		return false
	}

	data, ok := a.getUserAccess(userID)
	if !ok {
		return false
	}

	return data.roles[role]
}

// Invalidate clears the cached access data for a user.
func (a *Access) Invalidate(userID uint) {
	a.mu.Lock()
	delete(a.cache, userID)
	a.mu.Unlock()

	if a.redis != nil && a.redis.IsCacheAvailable() {
		key := fmt.Sprintf("session:%d", userID)
		a.redis.Delete(key)
	}
}
