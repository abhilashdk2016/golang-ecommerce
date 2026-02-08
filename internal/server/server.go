package server

import (
	"net/http"

	"github.com/abhilashdk2016/golang-ecommerce/internal/config"
	"github.com/abhilashdk2016/golang-ecommerce/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	config         *config.Config
	db             *gorm.DB
	logger         *zerolog.Logger
	authService    *services.AuthService
	productService *services.ProductService
	userService    *services.UserService
	uploadService  *services.UploadService
}

func New(
	cfg *config.Config,
	db *gorm.DB,
	logger *zerolog.Logger,
	authService *services.AuthService,
	productService *services.ProductService,
	userService *services.UserService,
	uploadService *services.UploadService,
) *Server {
	return &Server{
		config:         cfg,
		db:             db,
		logger:         logger,
		authService:    authService,
		productService: productService,
		userService:    userService,
		uploadService:  uploadService,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()

	// Add middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	// Add routes
	router.GET("/health", s.healthCheck)
	router.Static("/uploads", "./uploads")
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			authRoutes := auth
			authRoutes.POST("/register", s.register)
			authRoutes.POST("/login", s.login)
			authRoutes.POST("/refresh", s.refreshToken)
			authRoutes.POST("/logout", s.logout)
		}
		protected := api.Group("/")
		protected.Use(s.authMiddleware())
		{
			users := protected.Group("/users")
			{
				userRoutes := users
				userRoutes.GET("/profile", s.getProfile)
				userRoutes.PUT("/profile", s.updateProfile)
			}
		}
		categories := protected.Group("/categories")
		{
			categoryRoute := categories
			categoryRoute.POST("/", s.adminMiddleware(), s.createCategory)
			categoryRoute.PUT("/:id", s.adminMiddleware(), s.updateCategory)
			categoryRoute.DELETE("/:id", s.adminMiddleware(), s.deleteCategory)
		}

		products := protected.Group("/products")
		{
			productRoutes := products
			productRoutes.POST("/", s.adminMiddleware(), s.createProduct)
			productRoutes.PUT("/:id", s.adminMiddleware(), s.updateProduct)
			productRoutes.DELETE("/:id", s.adminMiddleware(), s.deleteProduct)
			productRoutes.POST("/:id/images", s.adminMiddleware(), s.uploadProductImage)
		}
		api.GET("/categories", s.getCategories)
		api.GET("/products", s.getProducts)
		api.GET("/products/:id", s.getProduct)
	}
	return router
}

func (s *Server) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
