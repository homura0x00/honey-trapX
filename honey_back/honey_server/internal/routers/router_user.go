package routers

import (
	"honey_back/honey_server/internal/controllers"
	"honey_back/honey_server/internal/middleware"
)

func UserRouter() {
	userRouter := apiRouterGroup.Group("/user")
	{
		/**
		register: 用户注册
		login: 用户登陆
		logout: 用户注销（前端退出）
		*/
		userRouter.POST("/register", controllers.RegisterUser)
		userRouter.POST("/login", controllers.LoginUser)
		userRouter.POST("/logout", controllers.UserLogout)
	}
	{
		/** 普通权限
		update: 更改用户昵称
		*/
		userRouter.GET("/get/login", controllers.GetLoginUser)
	}
	userRouter.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		/** 管理员权限
		update: 更新用户信息
		delete: 删除用户数据（软删除）
		*/
		userRouter.POST("/update", controllers.UpdateUser)
		userRouter.POST("/delete", controllers.DeleteUser)
	}
}
