package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// User 用户模型
type User struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
	CreateAt time.Time `json:"create_at"`
}

// 模拟数据库
var users = []User{
	{ID: 1, Username: "admin", Email: "admin@example.com", Password: "123456", Role: "admin", CreateAt: time.Now()},
	{ID: 2, Username: "user1", Email: "user1@example.com", Password: "123456", Role: "user", CreateAt: time.Now()},
	{ID: 3, Username: "user2", Email: "user2@example.com", Password: "123456", Role: "user", CreateAt: time.Now()},
}

var nextID = 4

func main() {
	r := gin.Default()

	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// 静态文件
	r.Static("/static", "./static")

	// 中间件 - 简单的认证检查
	authMiddleware := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			// 检查session或token，这里简化处理
			if c.Request.URL.Path != "/login" && c.Request.URL.Path != "/api/login" {
				username, exists := c.Get("username")
				if !exists {
					// 检查cookie
					cookie, err := c.Cookie("username")
					if err != nil {
						c.Redirect(http.StatusFound, "/login")
						c.Abort()
						return
					}
					c.Set("username", cookie)
				} else if username == "" {
					c.Redirect(http.StatusFound, "/login")
					c.Abort()
					return
				}
			}
			c.Next()
		}
	}

	// 路由组
	auth := r.Group("/")
	auth.Use(authMiddleware())

	// 登录页面
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "登录",
		})
	})

	// 登录处理
	r.POST("/api/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// 验证用户
		for _, user := range users {
			if user.Username == username && user.Password == password {
				// 设置cookie
				c.SetCookie("username", username, 3600, "/", "", false, true)
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": "登录成功",
				})
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户名或密码错误",
		})
	})

	// 退出登录
	auth.POST("/api/logout", func(c *gin.Context) {
		c.SetCookie("username", "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "退出成功",
		})
	})

	// 首页
	auth.GET("/", func(c *gin.Context) {
		username, _ := c.Get("username")
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title":     "控制台",
			"pageTitle": "控制台",
			"page":      "dashboard",
			"username":  username,
			"userCount": len(users),
		})
	})

	// 用户列表页面
	auth.GET("/users", func(c *gin.Context) {
		username, _ := c.Get("username")
		c.HTML(http.StatusOK, "users.html", gin.H{
			"title":     "用户管理",
			"pageTitle": "用户管理",
			"page":      "users",
			"username":  username,
			"users":     users,
		})
	})

	// API路由
	api := auth.Group("/api")
	{
		// 获取用户列表
		api.GET("/users", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    users,
			})
		})

		// 获取单个用户
		api.GET("/users/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "无效的用户ID",
				})
				return
			}

			for _, user := range users {
				if user.ID == id {
					c.JSON(http.StatusOK, gin.H{
						"success": true,
						"data":    user,
					})
					return
				}
			}

			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "用户不存在",
			})
		})

		// 创建用户
		api.POST("/users", func(c *gin.Context) {
			var newUser User
			if err := c.ShouldBindJSON(&newUser); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "请求参数错误",
				})
				return
			}

			newUser.ID = nextID
			nextID++
			newUser.CreateAt = time.Now()
			users = append(users, newUser)

			c.JSON(http.StatusCreated, gin.H{
				"success": true,
				"message": "用户创建成功",
				"data":    newUser,
			})
		})

		// 更新用户
		api.PUT("/users/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "无效的用户ID",
				})
				return
			}

			var updateUser User
			if err := c.ShouldBindJSON(&updateUser); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "请求参数错误",
				})
				return
			}

			for i, user := range users {
				if user.ID == id {
					updateUser.ID = id
					updateUser.CreateAt = user.CreateAt
					users[i] = updateUser
					c.JSON(http.StatusOK, gin.H{
						"success": true,
						"message": "用户更新成功",
						"data":    updateUser,
					})
					return
				}
			}

			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "用户不存在",
			})
		})

		// 删除用户
		api.DELETE("/users/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "无效的用户ID",
				})
				return
			}

			for i, user := range users {
				if user.ID == id {
					users = append(users[:i], users[i+1:]...)
					c.JSON(http.StatusOK, gin.H{
						"success": true,
						"message": "用户删除成功",
					})
					return
				}
			}

			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "用户不存在",
			})
		})
	}

	r.Run(":8080")
}
