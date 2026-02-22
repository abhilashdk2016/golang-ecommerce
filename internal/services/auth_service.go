package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/abhilashdk2016/golang-ecommerce/internal/config"
	"github.com/abhilashdk2016/golang-ecommerce/internal/dto"
	"github.com/abhilashdk2016/golang-ecommerce/internal/events"
	"github.com/abhilashdk2016/golang-ecommerce/internal/models"
	"github.com/abhilashdk2016/golang-ecommerce/internal/repository"
	"github.com/abhilashdk2016/golang-ecommerce/internal/utils"
)

var _ AuthServiceInterface = (*AuthService)(nil)

type AuthService struct {
	userRepo       repository.UserRepositoryInterface
	cartRepo       repository.CartRepositoryInterface
	config         *config.Config
	eventPublisher events.Publisher
}

func NewAuthService(
	cfg *config.Config,
	eventPublisher events.Publisher,
	userRepo repository.UserRepositoryInterface,
	cartRepo repository.CartRepositoryInterface) *AuthService {
	return &AuthService{
		config:         cfg,
		eventPublisher: eventPublisher,
		userRepo:       userRepo,
		cartRepo:       cartRepo,
	}
}

func (a *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	if _, err := a.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("you cannot register with this email")
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      models.UserRoleCustomer,
	}

	if err := a.userRepo.Create(&user); err != nil {
		return nil, err
	}

	cart := models.Cart{UserID: user.ID}
	if err := a.cartRepo.Create(&cart).Error; err != nil {
		fmt.Println("Unable to create cart...")
	}

	return a.generateAuthResponse(&user)
}

func (a *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := a.userRepo.GetByEmailAndActive(req.Email, true)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return a.generateAuthResponse(user)
}

func (a *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken, a.config.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	refreshToken, err := a.userRepo.GetValidRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	user, err := a.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := a.userRepo.DeleteRefreshTokenByID(refreshToken.ID); err != nil {
		log.Println(err)
		_ = err
	}

	return a.generateAuthResponse(user)
}

func (a *AuthService) Logout(refreshToken string) error {
	return a.userRepo.DeleteRefreshToken(refreshToken)
}

func (a *AuthService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(
		&a.config.JWT,
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(a.config.JWT.RefreshTokenExpiresIn),
	}

	if err := a.userRepo.CreateRefreshToken(&refreshTokenModel); err != nil {
		log.Println(err)
		_ = err
	}

	err = a.eventPublisher.Publish("USER_LOGIN", user, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("unable to publish user login event: %w", err)
	}
	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
