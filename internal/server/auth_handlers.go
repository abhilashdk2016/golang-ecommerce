package server

import (
	"github.com/abhilashdk2016/golang-ecommerce/internal/dto"
	"github.com/abhilashdk2016/golang-ecommerce/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authService.Register(&req)
	if err != nil {
		utils.BadRequestResponse(c, "registration failed", err)
		return
	}

	utils.CreatedResponse(c, "user registered successfully", response)
}

func (s *Server) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authService.Login(&req)
	if err != nil {
		utils.BadRequestResponse(c, "login failed", err)
		return
	}

	utils.SuccessResponse(c, "logged in successfully", response)
}

func (s *Server) refreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authService.RefreshToken(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "token refresh failed")
		return
	}

	utils.SuccessResponse(c, "token refreshed successfully", response)
}

func (s *Server) logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	err := s.authService.Logout(req.RefreshToken)
	if err != nil {
		utils.InternalServerErrorResponse(c, "logout failed", err)
		return
	}

	utils.SuccessResponse(c, "logged out successfully", nil)
}
