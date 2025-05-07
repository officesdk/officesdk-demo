package http

import (
	"github.com/gotomicro/ego/server/egin"

	"github.com/officesdk/go-sdk/officesdk"

	"turbo-demo/pkg/invoker"
	"turbo-demo/pkg/server/http/api"
	"turbo-demo/pkg/server/http/callback"
	"turbo-demo/pkg/server/http/middlewares"
	"turbo-demo/ui"
)

func ServeHTTP() *egin.Component {
	r := invoker.Gin
	r.Use(middlewares.CORS())
	apiRouters := r.Group("/showcase")
	{
		// 文件操作
		apiRouters.GET("/files", api.GetFiles)
		apiRouters.GET("/files/:guid", api.GetFile)
		apiRouters.DELETE("/file/:guid", api.DeleteFile)
		apiRouters.POST("/file", api.UploadFile)
		apiRouters.GET("/:guid/page", api.GetPageParams)
	}

	// 为 officesdk 添加鉴权中间件
	authMiddleware := middlewares.Auth()
	r.Use(authMiddleware)
	officesdk.NewServer(officesdk.Config{
		FileProvider: &callback.FileProvider{},
		AIProvider:   &callback.AIProvider{},
		Prefix:       "",
	}, r.Engine)

	r.Use(middlewares.Serve("/", middlewares.EmbedFolder(ui.WebUI, "dist"), false))
	r.Use(middlewares.Serve("/", middlewares.FallbackFileSystem(middlewares.EmbedFolder(ui.WebUI, "dist")), true))
	return r
}
