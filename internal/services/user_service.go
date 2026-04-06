package services

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/reshap0318/go-boilerplate/internal/dtos"
	"github.com/reshap0318/go-boilerplate/internal/helpers"
	"github.com/reshap0318/go-boilerplate/internal/models"
	"github.com/reshap0318/go-boilerplate/internal/repositories"
)

// UserCreate creates a new user with optional roles.
func (s *Services) UserCreate(ctx context.Context, req dtos.UserRequest) (*dtos.UserDTO, error) {
	s.Logger.LogStart("UserCreate", "Creating user: %s", req.Email)

	exists, err := s.repo.User.Exists(nil, map[string]interface{}{"email": req.Email})
	if err != nil {
		s.Logger.LogEndWithError("UserCreate", "Failed to check email: %v", err)
		return nil, err
	}
	if exists {
		s.Logger.LogEndWithError("UserCreate", "Email already exists: %s", req.Email)
		return nil, helpers.ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.LogEndWithError("UserCreate", "Failed to hash password: %v", err)
		return nil, err
	}

	user := &models.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hashedPassword),
	}

	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		result, err := s.repo.User.Create(tx, user)
		if err != nil {
			return nil, err
		}

		// Assign roles
		var roles []models.Role
		for _, roleID := range req.Roles {
			roles = append(roles, models.Role{ID: roleID})
		}
		if err := tx.Model(&result).Association("Roles").Append(roles); err != nil {
			s.Logger.LogStep("UserCreate", "Failed to assign roles: %v", err)
		}

		// Reload user with roles
		reloaded, err := s.repo.User.FindByID(tx, result.ID, "Roles")
		if err != nil {
			return nil, err
		}

		return reloaded, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("UserCreate", "Failed to create user: %v", err)
		return nil, err
	}

	result := res.(*models.User)
	dto := dtos.ToUserDTO(result)
	s.Logger.LogEnd("UserCreate", "User created: %s (ID: %d)", dto.Email, dto.ID)
	return &dto, nil
}

// UserGetAll returns paginated users with roles.
func (s *Services) UserGetAll(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[models.User], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	if opts.SortBy == "" {
		opts.SortBy = "id"
	}
	if opts.Order == "" {
		opts.Order = "ASC"
	}
	opts.Preloads = []string{"Roles"}

	return s.repo.User.FindAllWithOpts(nil, opts)
}

// UserGetByID returns a user by ID with roles.
func (s *Services) UserGetByID(ctx context.Context, id uint) (*dtos.UserDTO, error) {
	user, err := s.repo.User.FindByID(nil, id, "Roles")
	if err != nil {
		return nil, helpers.ErrNotFound
	}

	dto := dtos.ToUserDTO(user)
	return &dto, nil
}

// UserUpdate updates an existing user with optional roles.
func (s *Services) UserUpdate(ctx context.Context, id uint, req dtos.UserRequest) (*dtos.UserDTO, error) {
	s.Logger.LogStart("UserUpdate", "Updating user ID: %d", id)

	existing, err := s.repo.User.FindByID(nil, id)
	if err != nil {
		s.Logger.LogEndWithError("UserUpdate", "User not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	if existing.Email != req.Email {
		exists, err := s.repo.User.Exists(nil, map[string]interface{}{"email": req.Email})
		if err != nil {
			s.Logger.LogEndWithError("UserUpdate", "Failed to check email: %v", err)
			return nil, err
		}
		if exists {
			s.Logger.LogEndWithError("UserUpdate", "Email already exists: %s", req.Email)
			return nil, helpers.ErrUserExists
		}
	}

	updates := map[string]interface{}{
		"name":  req.Name,
		"email": req.Email,
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.Logger.LogEndWithError("UserUpdate", "Failed to hash password: %v", err)
			return nil, err
		}
		updates["password"] = string(hashedPassword)
	}

	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		result, err := s.repo.User.UpdateMap(tx, &models.User{ID: id}, updates)
		if err != nil {
			return nil, err
		}

		// Replace roles - clear then assign
		if err := tx.Model(&result).Association("Roles").Clear(); err != nil {
			return nil, err
		}

		var roles []models.Role
		for _, roleID := range req.Roles {
			roles = append(roles, models.Role{ID: roleID})
		}
		if err := tx.Model(&result).Association("Roles").Append(roles); err != nil {
			s.Logger.LogStep("UserUpdate", "Failed to assign roles: %v", err)
		}

		// Reload user with roles
		reloaded, err := s.repo.User.FindByID(tx, result.ID, "Roles")
		if err != nil {
			return nil, err
		}

		return reloaded, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("UserUpdate", "Failed to update user: %v", err)
		return nil, err
	}

	result := res.(*models.User)
	dto := dtos.ToUserDTO(result)
	s.Logger.LogEnd("UserUpdate", "User updated: %s (ID: %d)", dto.Email, dto.ID)
	return &dto, nil
}

// UserDelete soft deletes a user and its role associations.
func (s *Services) UserDelete(ctx context.Context, id uint) error {
	s.Logger.LogStart("UserDelete", "Deleting user ID: %d", id)

	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		user := models.User{ID: id}
		if err := tx.Model(&user).Association("Roles").Clear(); err != nil {
			return err
		}
		_, err := s.repo.User.Delete(tx, id)
		return err
	}); err != nil {
		s.Logger.LogEndWithError("UserDelete", "Failed to delete user: %v", err)
		return err
	}

	s.Logger.LogEnd("UserDelete", "User deleted: ID: %d", id)
	return nil
}
