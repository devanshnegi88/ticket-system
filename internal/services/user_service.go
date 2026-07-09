package services

import (
	"errors"

	"ticket-system/internal/auth"
	"ticket-system/internal/models"
	"ticket-system/internal/repositories"
)

var (
	// ErrEmailExists is returned when registering with an email that is
	// already taken.
	ErrEmailExists = errors.New("email already registered")
	// ErrInvalidCredentials is returned for any login failure. It is kept
	// deliberately generic (not "user not found" vs "wrong password") to
	// avoid leaking which emails are registered.
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// UserService contains the business logic for user registration and
// authentication.
type UserService struct {
	repo       *repositories.UserRepository
	jwtManager *auth.JWTManager
}

// NewUserService builds a UserService.
func NewUserService(repo *repositories.UserRepository, jwtManager *auth.JWTManager) *UserService {
	return &UserService{repo: repo, jwtManager: jwtManager}
}

// Register creates a new user with a bcrypt-hashed password. Returns
// ErrEmailExists if the email is already taken.
func (s *UserService) Register(email, password string) (*models.User, error) {
	_, err := s.repo.FindByEmail(email)
	if err == nil {
		return nil, ErrEmailExists
	}
	if !errors.Is(err, repositories.ErrNotFound) {
		return nil, err
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{Email: email, PasswordHash: hash}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login verifies credentials and returns a signed JWT on success, or
// ErrInvalidCredentials otherwise.
func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if !auth.CheckPassword(password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}

	return s.jwtManager.Generate(user.ID, user.Email)
}
