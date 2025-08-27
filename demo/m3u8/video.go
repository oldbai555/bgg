package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// VideoSt 视频数据结构
type VideoSt struct {
	Id        uint32 `json:"id" gorm:"primaryKey;autoIncrement"`
	Uuid      string `json:"uuid" gorm:"type:varchar(36);uniqueIndex;not null"`
	PlayerUrl string `json:"player_url" gorm:"type:varchar(500)"`
	Name      string `json:"name" gorm:"type:varchar(200)"`
	GodNum    string `json:"god_num" gorm:"type:varchar(127)"`
}

// 全局数据库实例
var db *gorm.DB

// Response 响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 初始化数据库连接
func initDB() {
	// 数据库连接配置，请根据实际情况修改
	dsn := ""

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	// 自动迁移创建表
	err = db.AutoMigrate(&VideoSt{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
}

// 添加视频接口
func addVideo(c *gin.Context) {
	var video VideoSt

	// 解析 JSON 请求体
	if err := c.ShouldBindJSON(&video); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数格式错误: " + err.Error(),
		})
		return
	}

	// 检查 uuid 是否为空
	if video.Uuid == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "新增失败: uuid 不能为空",
		})
		return
	}

	if !strings.Contains(video.PlayerUrl, video.Uuid) {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "链接错误",
		})
		return
	}

	// 检查数据库中是否已存在相同 uuid 的数据
	var existingVideo VideoSt
	result := db.Where("uuid = ?", video.Uuid).First(&existingVideo)
	if result.Error == nil {
		// 找到了记录，说明 uuid 已存在
		c.JSON(http.StatusConflict, Response{
			Code:    409,
			Message: "新增失败: 该 uuid 已存在",
		})
		return
	} else if result.Error != gorm.ErrRecordNotFound {
		// 数据库查询出错
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "数据库查询失败: " + result.Error.Error(),
		})
		return
	}

	// 插入新数据
	if err := db.Create(&video).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "数据库插入失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "新增成功",
		Data:    video,
	})
}

// 获取所有视频列表接口
func getVideoList(c *gin.Context) {
	var videos []VideoSt

	// 查询所有视频，按 id 正序排序
	if err := db.Order("id ASC").Find(&videos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "数据库查询失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.HTML(http.StatusOK, "list.html", gin.H{
		"videoList": videos,
	})
}
