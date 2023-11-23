package routes

import (
	"Nestar/controllers/weakpass"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", weakpass.MainPageController)
	r.GET("/index.html", weakpass.IndexController)

	r.GET("/manage.html", weakpass.ManageController)
	r.GET("/manage.html/:systemName", weakpass.ManageSearchController)
	r.GET("/universal.html", weakpass.UniversalController)
	r.GET("/classification.html", weakpass.ClassificationController)
	r.GET("/classification.html/:systemName", weakpass.ClassificationSearchController)
	r.GET("/system.html/:systemName", weakpass.SystemController)
	r.GET("/allin.html", weakpass.AllinController)
	r.GET("/tools.html", weakpass.ToolsController)
	r.GET("/copyAndDownData.html/:systemName/:where", weakpass.CopyAndDownDataController)
	//r.GET("/FillAndCopyData.html/:systemName/:where", weakpass.FillAndCopyDataController)
	r.GET("/AddData.html/:systemName/:where/:clazz", weakpass.AddDataController)

	r.POST("/api/copyData", weakpass.CopyDataController)
	r.POST("/api/addData", weakpass.AddDataApiController)
	r.POST("/api/sortData", weakpass.SortDataApiController)
	r.POST("/api/addSystem", weakpass.AddSystemApiController)
}
