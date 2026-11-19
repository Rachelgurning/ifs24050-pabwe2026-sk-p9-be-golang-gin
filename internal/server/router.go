package server

import (
	"net/http"
	"time"

	"app/docs"
	"app/internal/config"
	"app/internal/features/ai"
	"app/internal/features/auth"
	"app/internal/features/post"
	"app/internal/middleware"
	"app/internal/response"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func NewRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS())

	authH := auth.NewHandler(auth.NewService(auth.NewRepository(db), cfg.AppKey))
	postH := post.NewHandler(post.NewService(post.NewRepository(db)))
	aiH := ai.NewHandler(ai.NewService(
		ai.NewSearchClient(cfg.SearchEndpoint, cfg.OpenAIKey),
		ai.NewOpenAIClient(cfg.OpenAIEndpoint, cfg.OpenAIKey, cfg.OpenAIModel),
	))
	requireAuth := middleware.Auth(cfg.AppKey)
	limiter := middleware.RateLimiter(300, time.Minute)

	r.GET("/", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "Delcom Blog API berjalan", gin.H{"name": "Delcom Blog API", "version": "1.0.0", "docs": "/swagger/index.html"})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/api.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(docs.SwaggerInfo.ReadDoc()))
	})

	api := r.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "API sehat", gin.H{"time": time.Now()})
	})

	public := api.Group("/public")
	public.GET("/categories", postH.Categories)
	public.GET("/posts", postH.PublicList)
	public.GET("/posts/:slug", postH.PublicDetail)

	authG := api.Group("/auth")
	authG.POST("/register", limiter, authH.Register)
	authG.POST("/login", limiter, authH.Login)
	authG.POST("/logout", requireAuth, authH.Logout)
	authG.GET("/profile", requireAuth, authH.Profile)
	authG.PUT("/profile", requireAuth, authH.UpdateProfile)

	posts := api.Group("/posts", requireAuth)
	posts.GET("", postH.MyPosts)
	posts.POST("", postH.Create)
	posts.GET("/stats", postH.Stats)
	posts.POST("/ai-generate", aiH.Generate)
	posts.GET("/:id", postH.Detail)
	posts.PUT("/:id", postH.Update)
	posts.PATCH("/:id/toggle-publish", postH.TogglePublish)
	posts.DELETE("/:id", postH.Delete)

	admin := api.Group("/admin", requireAuth, middleware.RequireRole("admin"))
	admin.GET("/users", postH.AdminUsers)
	admin.GET("/posts", postH.AdminPosts)
	admin.GET("/users/:id/posts", postH.AdminUserPosts)

	r.NoRoute(func(c *gin.Context) {
		response.Fail(c, http.StatusNotFound, "Endpoint not found")
	})
	return r
}
