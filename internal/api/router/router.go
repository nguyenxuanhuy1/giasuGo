package router

import (
	"traingolang/internal/api/handler"
	"traingolang/internal/auth"
	"traingolang/internal/config"
	"traingolang/internal/repository"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	postRepo := repository.NewPostRepo(config.DB)
	imageRepo := repository.NewImageRepository(config.DB)
	r.POST("/analyze", auth.LimitUploadSize(1<<20), handler.AnalyzeImage)
	// r.POST("/api/cv/ocr", auth.LimitUploadSize(1<<20), handler.AnalyzeCV)
	r.POST("/api/user/register", handler.Register)
	r.POST("/api/user/login", handler.Login)
	r.POST("/api/search/post", handler.SearchPostsHandler(postRepo))
	r.GET("/api/posts/options", handler.GetPostOptionsHandler(postRepo))

	api := r.Group("/api")
	{
		api.POST("/match/create", handler.CreateMatch)
		api.POST("/match/join", handler.JoinMatch)
		api.GET("/profile", handler.Profile)
		api.POST("/upload", handler.UploadHandler)

		api.POST(
			"/create/post",
			auth.AdminOnly(),
			auth.LimitUploadSize(1<<20),
			handler.CreatePost(postRepo, imageRepo),
		)

		api.POST(
			"/update/post/:id",
			auth.AdminOnly(),
			auth.LimitUploadSize(1<<20),
			handler.UpdatePost(postRepo, imageRepo),
		)

		api.POST(
			"/delete/post/:id",
			auth.AdminOnly(),
			handler.DeletePost(postRepo, imageRepo),
		)
	}

	return r
}
