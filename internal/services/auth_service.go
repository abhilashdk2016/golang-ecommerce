package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/abhilashdk2016/golang-ecommerce/internal/config"
	"github.com/abhilashdk2016/golang-ecommerce/internal/dto"
	"github.com/abhilashdk2016/golang-ecommerce/internal/events"
	"github.com/abhilashdk2016/golang-ecommerce/internal/models"
	"github.com/abhilashdk2016/golang-ecommerce/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db             *gorm.DB
	config         *config.Config
	eventPublisher events.Publisher
}

func NewAuthService(db *gorm.DB, cfg *config.Config, eventPublisher events.Publisher) *AuthService {
	return &AuthService{
		db:             db,
		config:         cfg,
		eventPublisher: eventPublisher,
	}
}

func (a *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	var existingUser models.User
	if err := a.db.Where("email = ?", req.Email).First(&existingUser); err == nil {
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

	if err := a.db.Create(&user).Error; err != nil {
		return nil, err
	}

	cart := models.Cart{UserID: user.ID}
	if err := a.db.Create(&cart).Error; err != nil {
		fmt.Println("Unable to create cart...")
	}

	return a.generateAuthResponse(&user)
}

func (a *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	var user models.User
	if err := a.db.Where("email = ? AND is_active = ?", req.Email, true).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return a.generateAuthResponse(&user)
}

func (a *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken, a.config.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var refreshToken models.RefreshToken
	if err := a.db.Where("token = ? AND expires_at > ?", req.RefreshToken, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	var user models.User
	if err := a.db.First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	a.db.Delete(&refreshToken)

	return a.generateAuthResponse(&user)
}

func (a *AuthService) Logout(refreshToken string) error {
	return a.db.Where("token = ?", refreshToken).Delete(&models.RefreshToken{}).Error
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

	a.db.Create(&refreshTokenModel)

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
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
