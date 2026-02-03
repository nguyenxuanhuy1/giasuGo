package router

import (
	"traingolang/internal/api/handler"
	"traingolang/internal/auth"
	"traingolang/internal/config"
	"traingolang/internal/repository"
	"traingolang/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	postRepo := repository.NewPostRepo(config.DB)
	imageRepo := repository.NewImageRepository(config.DB)
	examService := service.NewExamService(config.DB)

	// PUBLIC ROUTES (KHÔNG CẦN TOKEN)
	public := r.Group("/api")
	{
		// Google OAuth routes
		public.GET("/auth/google", handler.GoogleLogin)
		public.GET("/auth/google/callback", handler.GoogleCallback)
		public.POST("/auth/refresh", handler.RefreshToken)
		public.GET("/posts/options", handler.GetPostOptionsHandler(postRepo))

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
		authGroup.POST("/match/create", handler.CreateMatch)
		authGroup.POST("/match/join", handler.JoinMatch)
		authGroup.POST("/upload", handler.UploadHandler)
		authGroup.GET("/exam-sets/redo/:id", handler.RedoExamHandler(examService))
		authGroup.GET("/exams/history", handler.GetMyExamSetsHandler(examService))
		authGroup.GET(
			"/exam-sets/questions/:id",
			handler.GetExamQuestionsHandler(examService),
		)
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
	}

	return r
}
