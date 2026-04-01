package repositories

import "gorm.io/gorm"

type Repositories struct {
	TxManager    *TransactionManager
	User         *UserRepository
	PasswordReset *PasswordResetRepository
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
	txManager := NewTransactionManager(db)
	userRepo := NewUserRepository(db)
	passwordResetRepo := NewPasswordResetRepository(db)

	return &Repositories{
		TxManager:     txManager,
		User:          userRepo,
		PasswordReset: passwordResetRepo,
	}, nil
}
