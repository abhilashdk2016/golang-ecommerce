package server

import (
	"github.com/abhilashdk2016/golang-ecommerce/internal/dto"
	"github.com/abhilashdk2016/golang-ecommerce/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) getProfile(c *gin.Context) {

	userID := c.GetUint("user_id")
	profile, err := s.userService.GetProfile(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, "Profile retrieved successfully", profile)
}

func (s *Server) updateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	profile, err := s.userService.UpdateProfile(userID, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update profile", err)
		return
	}
	utils.SuccessResponse(c, "Profile updated successfully", profile)
}
