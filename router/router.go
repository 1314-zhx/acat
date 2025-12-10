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
		userRouter.GET("/forget", func(c *gin.Context) {
			c.HTML(200, "forget.html", nil)
		})
		userRouter.POST("/reset-password", controller.ReSetPasswordHandler)
		// 用户登录
		userRouter.POST("/login", controller.LoginHandler)
		userRouter.GET("/login", controller.ShowLoginHandler)
		// 用户下载电子版试题
		userRouter.GET("/download/file", controller.DownloadHandler)
		userRouter.GET("/download", func(c *gin.Context) {
			c.HTML(200, "download.html", nil)
		})
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
			// 用户查看回信
			authed.POST("/check_reply", controller.CheckReplyHandler)
			authed.GET("/check_reply", func(c *gin.Context) {
				c.HTML(200, "check_reply.html", nil)
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
			authed.GET("/logout", controller.LoginOutHandler)

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
		adminRouter.POST("/set_schedule", controller.AdminSetScheduleHandler)
		adminRouter.GET("/set_schedule", func(c *gin.Context) {
			c.HTML(200, "set_schedule.html", nil)
		})
		// 管理员设置面试结果
		adminRouter.POST("/set_result", controller.AdminSetInterviewResultHandler)
		adminRouter.GET("/set_result", func(c *gin.Context) {
			c.HTML(200, "set_result.html", nil)
		})
		adminRouter.POST("/set_pass", controller.SetPassHandler)
		// 管理员发邮件，公布面试结果(QQ邮箱和一块)
		adminRouter.POST("/post_email", controller.AdminPublishHandler)
		adminRouter.POST("/get_pass_user", controller.GetPassUserHandler)
		adminRouter.GET("/post_email", func(c *gin.Context) {
			c.HTML(200, "publish_email.html", nil)
		})
		// 管理员查看收件箱并回信
		adminRouter.POST("/mailbox", controller.AdminLetterHandler)
		adminRouter.POST("/reply", controller.AdminReplyHandler)
		adminRouter.GET("/mailbox", func(c *gin.Context) {
			c.HTML(200, "admin_check_reply_letter.html", nil)
		})
		// 管理员上传电子版试题
		adminRouter.POST("/upload", controller.UploadHandler)
		adminRouter.GET("/upload", func(c *gin.Context) {
			c.HTML(200, "upload.html", nil)
		})
	}
	r.NoRoute(controller.Norouter)
	return r
}
