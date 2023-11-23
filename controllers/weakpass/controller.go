package weakpass

import (
	"Nestar/models"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

// MainPageController 主页的 iframe 框架的控制器
func MainPageController(c *gin.Context) {
	c.HTML(200, "gui/iframe", gin.H{})
}

// IndexController “介绍”的控制器
func IndexController(c *gin.Context) {
	c.HTML(200, "gui/index", gin.H{})
}

// ManageController “系统管理”的控制器
func ManageController(c *gin.Context) {
	count := models.GetSystemCount()
	c.HTML(200, "gui/manage", gin.H{
		"Count": count,
	})
}

func ManageSearchController(c *gin.Context) {
	systemSearch := c.Param("systemName")
	count := models.GetSystemCount()
	search := models.GetSystemSearch(systemSearch)

	if len(search) >= 10 {
		search = search[:10]
	}
	c.HTML(200, "gui/manage", gin.H{
		"Count":  count,
		"Search": search,
	})
}

// UniversalController “通用弱口令”的控制器
func UniversalController(c *gin.Context) {
	index := "universal"
	dataAll := models.GetDataAll(index)
	c.HTML(200, "gui/universal", gin.H{
		"System":  "universal",
		"DataAll": dataAll,
	})
}

// ClassificationController “分类速查库”的控制器
func ClassificationController(c *gin.Context) {
	systems := models.GetSystemUniversalList()
	if len(systems) >= 10 {
		systems = systems[:10]
	}

	c.HTML(200, "gui/classification", gin.H{
		"Systems": systems,
	})
}

func ClassificationSearchController(c *gin.Context) {
	systemSearch := c.Param("systemName")
	systems := models.GetSystemUniversalList()
	search := models.GetSystemSearch(systemSearch)

	if len(systems) >= 10 {
		systems = systems[:10]
	}
	if len(search) >= 10 {
		search = search[:10]
	}
	c.HTML(200, "gui/classification", gin.H{
		"Systems": systems,
		"Search":  search,
	})
}

func SystemController(c *gin.Context) {
	systemName := c.Param("systemName")
	dataAll := models.GetDataAll(systemName)
	c.HTML(200, "gui/system", gin.H{
		"System":  systemName,
		"DataAll": dataAll,
	})
}

// AllinController “all in 汇总”的控制器
func AllinController(c *gin.Context) {
	dataAll := models.GetAllin()
	c.HTML(200, "gui/allin", gin.H{
		"Part": dataAll,
	})
}

// ToolsController “其他工具” 的控制器
func ToolsController(c *gin.Context) {
	c.HTML(200, "gui/tools", gin.H{})
}

// CopyAndDownDataController 复制和下载数据弹窗的控制器
func CopyAndDownDataController(c *gin.Context) {
	systemName := c.Param("systemName")
	where := c.Param("where")
	count := models.GetAllinCount(where)

	c.HTML(200, "gui/copyAndDownData", gin.H{
		"SystemName": systemName,
		"Where":      where,
		"Count":      count,
	})
}

//func FillAndCopyDataController(c *gin.Context) {
//	systemName := c.Param("systemName")
//	where := c.Param("where")
//	count := 10000 // = system count + all count // 从文件中获取数据--> 待编写。
//
//	c.HTML(200, "gui/fillAndCopyData", gin.H{
//		"SystemName": systemName,
//		"Where":      where,
//		"Count":      count,
//	})
//}

// AddDataController 添加数据弹框的控制器
func AddDataController(c *gin.Context) {
	systemName := c.Param("systemName")
	where := c.Param("where")
	clazz := c.Param("clazz")

	c.HTML(200, "gui/AddData", gin.H{
		"System": systemName,
		"Where":  where,
		"Clazz":  clazz,
	})
}

// CopyDataController 复制数据弹框的控制器
func CopyDataController(c *gin.Context) {
	var query models.CopyQuery
	err := c.BindJSON(&query)
	if err != nil {
		log.Printf("CopyDataController 参数绑定错误 %v", err)
	}
	log.Printf("CopyDataController 参数 %v", query)

	data := ""
	// 通过 query 参数 获取 数据
	if query.SystemName == "allin" {
		// 获取 allin 的数据
		data += models.GetAllinData(query)
	}
	data = strings.TrimRight(data, "\n")
	c.JSON(200, gin.H{
		"data": data,
	})
}

// AddDataApiController 添加数据的 API
func AddDataApiController(c *gin.Context) {
	var addDataQuery models.AddDataQuery
	_ = c.BindJSON(&addDataQuery)

	// 一般情况下，业务逻辑都是放到 service 中去写的，我这里就偷懒了，直接在这里写吧
	// 分析查询参数，然后将数据写到对应的文件中
	// systemName 、where、clazz、data
	models.AddData(addDataQuery)

	c.JSON(200, gin.H{
		"data": "添加成功",
	})
}

// SortDataApiController 整理数据的 API
func SortDataApiController(c *gin.Context) {
	var sortQuery models.SortQuery
	_ = c.BindJSON(&sortQuery)

	if sortQuery.SystemName == "allin" {
		models.AddAllinData(sortQuery)
	} else {
		models.SortData(sortQuery)
	}

	c.JSON(200, gin.H{
		"data": "整理成功",
	})
}

func AddSystemApiController(c *gin.Context) {
	var addSystemQuery models.AddSystemQuery
	_ = c.BindJSON(&addSystemQuery)

	info := models.AddSystem(addSystemQuery)

	c.JSON(200, gin.H{
		"data": info,
	})
}
