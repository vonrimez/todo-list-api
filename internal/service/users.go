package service

import (
	"github.com/vonrimez/TaskAPI/internal/common/auth"
	"github.com/vonrimez/TaskAPI/internal/common/domain"
	"github.com/vonrimez/TaskAPI/internal/common/models"
)

type UserRepository interface {
	CreateUser(models.UserRegisterInput) (*models.UserOutput, error)
	LoginUser(models.UserLoginInput) (*models.UserOutput, error)
}

type UserService struct {
	repo      UserRepository
	JWTSecret string
}

func NewUserService(repo UserRepository, jwtSecret string) *UserService {
	return &UserService{repo: repo, JWTSecret: jwtSecret}
}

func (s *UserService) Register(inputUser models.UserRegisterInput) (*models.User, error) {
	encryptedPass, err := auth.HashPassword(*inputUser.Pass)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}
	*inputUser.Pass = encryptedPass

	outputUser, appErr := s.repo.CreateUser(inputUser)
	if appErr != nil {
		return nil, appErr
	}

	jwt, err := auth.GetJWT(outputUser.ID, s.JWTSecret)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}

	return &models.User{
		ID:   outputUser.ID,
		Name: outputUser.Name,
		JWT:  jwt,
	}, nil
}

func (s *UserService) Login(inputUser models.UserLoginInput) (*models.User, error) {
	outputUser, appErr := s.repo.LoginUser(inputUser)
	if appErr != nil {
		return nil, appErr
	}
	if !auth.IsCorrectPassword(*inputUser.Pass, outputUser.Pass) {
		return nil, domain.NewBadRequestError("invalid email or password")
	}

	jwt, err := auth.GetJWT(outputUser.ID, s.JWTSecret)
	if err != nil {
		return nil, domain.NewInternalError(err)
	}

	return &models.User{
		ID:   outputUser.ID,
		Name: outputUser.Name,
		JWT:  jwt,
	}, nil
}
