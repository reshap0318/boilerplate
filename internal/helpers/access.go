package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/reshap0318/go-boilerplate/internal/database"
	"github.com/reshap0318/go-boilerplate/internal/models"
	"gorm.io/gorm"
)

type userAccessData struct {
	permissions map[string]bool
	roles       map[string]bool
}

// Access handles permission and role checking with 3-tier caching.
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

// getUserAccess retrieves user permissions and roles using 3-tier cache.
// L1: Local in-memory cache → L2: Redis → L3: Database
func (a *Access) getUserAccess(userID uint) (*userAccessData, bool) {
	// L1: Check local cache
	a.mu.RLock()
	data, ok := a.cache[userID]
	a.mu.RUnlock()
	if ok {
		return data, true
	}

	// L2: Check Redis
	if a.redis != nil && a.redis.IsCacheAvailable() {
		var user models.User
		key := fmt.Sprintf("session:%d", userID)
		if err := a.getRedisJSON(key, &user); err == nil && user.ID != 0 {
			data = a.buildAccessDataFromUser(&user)

			a.mu.Lock()
			a.cache[userID] = data
			a.mu.Unlock()

			return data, true
		}
	}

	// L3: Fallback to DB
	user, err := a.findUserWithRolesPermissions(userID)
	if err != nil {
		return nil, false
	}

	data = a.buildAccessDataFromUser(user)

	// Store to Redis (if active)
	if a.redis != nil && a.redis.IsCacheAvailable() {
		key := fmt.Sprintf("session:%d", userID)
		if err := a.setRedisJSON(key, user); err != nil {
			// Redis error, ignore
		}
	}

	a.mu.Lock()
	a.cache[userID] = data
	a.mu.Unlock()

	return data, true
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

func (a *Access) getRedisJSON(key string, dest interface{}) error {
	val, err := a.redis.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (a *Access) setRedisJSON(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return a.redis.Set(key, data, 0)
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
