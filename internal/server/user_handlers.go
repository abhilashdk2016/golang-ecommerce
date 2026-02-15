package server

import (
	"github.com/abhilashdk2016/golang-ecommerce/internal/dto"
	"github.com/abhilashdk2016/golang-ecommerce/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Get user profile
// @Description Get current authenticated user's profile information
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.UserResponse} "Profile retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "User not found"
// @Router /users/profile [get]
func (s *Server) getProfile(c *gin.Context) {

	userID := c.GetUint("user_id")
	profile, err := s.userService.GetProfile(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, "Profile retrieved successfully", profile)
}

// @Summary Update user profile
// @Description Update current authenticated user's profile information
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} utils.Response{data=dto.UserResponse} "Profile updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /users/profile [put]
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
