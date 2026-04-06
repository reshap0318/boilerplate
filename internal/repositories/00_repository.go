package repositories

import "gorm.io/gorm"

type Repositories struct {
	TxManager     *TransactionManager
	User          *UserRepository
	PasswordReset *PasswordResetRepository
	Permission    *PermissionRepository
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
	txManager := NewTransactionManager(db)
	userRepo := NewUserRepository(db)
	passwordResetRepo := NewPasswordResetRepository(db)
	permissionRepo := NewPermissionRepository(db)

	return &Repositories{
		TxManager:     txManager,
		User:          userRepo,
		PasswordReset: passwordResetRepo,
		Permission:    permissionRepo,
	}, nil
}
