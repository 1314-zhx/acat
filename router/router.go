package router

import (
	"acat/controller"
	"acat/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// 返回gin的引擎
func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.GinRecovery(true))
	r.LoadHTMLGlob("templates/*/*.html")
	// 加载静态文件
	r.Static("/static", "./static")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "PONG")
		zap.L().Info("PING-PONG")
	})
	// 跟路由，跳转到/index路由下
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/index")
	})
	// 渲染初始页面
	r.GET("/index", controller.IndexHandler)
	r.GET("/introduce", func(c *gin.Context) {
		c.HTML(200, "introduce.html", nil)
	})
	// 方向介绍模板
	r.GET("/tech/go", controller.TechGoHandler)
	r.GET("/tech/java", controller.TechJavaHandler)
	r.GET("/tech/frontend", controller.TechFrontendHandler)
	// 查看预约面试时间表
	r.GET("/check", controller.CheckHandler)
	// 用来渲染前端数据
	r.GET("/schedule", func(c *gin.Context) {
		c.HTML(http.StatusOK, "check.html", nil)
	})
	// 用户存在路由组
	userRouter := r.Group("/user")
	{
		// 用户中心首页
		userRouter.GET("/center", func(c *gin.Context) {
			c.HTML(200, "center.html", nil)
		})
		// 用户注册
		userRouter.POST("/register", controller.RegisterHandler)
		userRouter.GET("/register", func(c *gin.Context) {
			c.HTML(200, "register.html", nil)
		})
		userRouter.POST("/forget", controller.ForgetHandler)
		// 用户登录
		userRouter.POST("/login", controller.LoginHandler)
		userRouter.GET("/login", controller.ShowLoginHandler)
		authed := userRouter.Group("/auth")
		authed.Use(middleware.AuthUserHTML())
		{
			// 用户查询面试结果
			authed.POST("/result", controller.ResultHandler)
			authed.GET("/result", func(c *gin.Context) {
				c.HTML(200, "result.html", nil)
			})
			// 用户报名面试
			authed.POST("/signup", controller.PostHandler)
			authed.GET("/signup", func(c *gin.Context) {
				c.HTML(200, "booking.html", nil)
			})
			// 用户联系管理员，提交问题，参看邮箱
			authed.POST("/conversation", controller.LetterHandler)
			authed.GET("/conversation", func(c *gin.Context) {
				c.HTML(200, "message.html", nil)
			})
			// 展示管理员列表
			authed.GET("/show_admin", controller.ShowAdminHandler)

			// 用户提交更新面试预约的请求
			authed.POST("/update", controller.UpdateHandler)
			authed.GET("/my_slot", controller.ShowSlotHandler)
			authed.GET("/update", func(c *gin.Context) {
				c.HTML(200, "update.html", nil)
			})
			// 用户个人中心
			authed.GET("/account", controller.AccountHandler)
		}
	}
	// 管理员登录
	r.POST("/admin_login", controller.AdminLoginHandler)
	r.GET("/admin_login", func(c *gin.Context) {
		c.HTML(200, "admin_login.html", nil)
	})
	// 管理员路由
	adminRouter := r.Group("/admin")
	adminRouter.Use(middleware.AuthUserHTML())
	{
		// 管理员中心
		adminRouter.GET("/center", func(c *gin.Context) {
			c.HTML(200, "admin_center.html", nil)
		})
		// 管理员设置面试表
		adminRouter.POST("/settimetable", controller.AdminSetScheduleHandler)
		// 管理员设置面试结果
		adminRouter.POST("/setresult", controller.AdminSetInterviewResultHandler)
		// 管理员发邮件(短信和QQ邮箱和一块)
		adminRouter.POST("/postemail", controller.AdminPostEmailHandler)
		// 管理员查看收件箱并回信
		adminRouter.POST("/letter", controller.LetterHandler)
	}
	r.NoRoute(controller.Norouter)
	return r
}
