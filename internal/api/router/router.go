package router

import (
	"traingolang/internal/api/handler"
	"traingolang/internal/auth"
	"traingolang/internal/config"
	"traingolang/internal/repository"
	"traingolang/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Cấu hình CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://congdongonthi.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	postRepo := repository.NewPostRepo(config.DB)
	imageRepo := repository.NewImageRepository(config.DB)
	examService := service.NewExamService(config.DB)

	r.GET("/oauth2/authorization/google", handler.GoogleLogin)
	r.GET("/oauth2/callback/google", handler.GoogleCallback)
	// PUBLIC ROUTES (KHÔNG CẦN TOKEN)
	public := r.Group("/api")
	{
		// Google OAuth routes
		// public.GET("/oauth2/authorization/google", handler.GoogleLogin)
		// public.GET("/oauth2/callback/google", handler.GoogleCallback)
		public.POST("/auth/refresh", handler.RefreshToken)
		public.GET("/posts/options", handler.GetPostOptionsHandler(postRepo))
		public.POST("/exams/public", handler.GetPublicExamsHandler(examService))
	}

	// AUTH ROUTES (BẮT BUỘC TOKEN)
	authGroup := r.Group("/api", auth.Middleware())
	{
		authGroup.POST(
			"/analyze",
			auth.LimitUploadSize(1<<20),
			handler.AnalyzeImage,
		)
		authGroup.POST(
			"/analyze/question",
			handler.AnalyzeQuestion,
		)
		authGroup.POST(
			"/exams/submit",
			handler.SubmitExamHandler(examService),
		)

		authGroup.GET("/user/info", handler.Profile)
		authGroup.POST("/upload", handler.UploadHandler)
		authGroup.GET("/exam-sets/redo/:id", handler.RedoExamHandler(examService))
		authGroup.GET("/exams/history", handler.GetMyExamSetsHandler(examService))

	}

	// ADMIN ROUTES (TOKEN + ADMIN)
	admin := r.Group(
		"/api/admin",
		auth.Middleware(),
		auth.AdminOnly(),
	)
	{
		admin.POST(
			"/create/post",
			auth.LimitUploadSize(1<<20),
			handler.CreatePost(postRepo, imageRepo),
		)

		admin.POST(
			"/update/post/:id",
			auth.LimitUploadSize(1<<20),
			handler.UpdatePost(postRepo, imageRepo),
		)

		admin.POST(
			"/delete/post/:id",
			handler.DeletePost(postRepo, imageRepo),
		)
		admin.POST(
			"/exams/update/:id",
			handler.UpdateExamSetHandler(examService),
		)
	}

	return r
}
