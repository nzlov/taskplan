package app

import (
	"github.com/gin-gonic/gin"
)

func Router(app *gin.Engine) {
	app.POST("/api/login", LoginEndpoint)
	app.GET("/api/login", CheckLogin)
	app.POST("/api/user", UserAdd)

	app.GET("/api/permission", Auth("permission.list"), PermissionList)
	app.GET("/api/permission/:id", Auth("permission.info"), PermissionList)
	app.POST("/api/permission", Auth("permission.add"), PermissionAdd)
	app.POST("/api/permission/:id", Auth("permission.update"), PermissionUpdate)
	app.DELETE("/api/permission/:id", Auth("permission.del"), PermissionDel)

	app.GET("/api/role", Auth("role.list"), RoleList)
	app.GET("/api/role/:id", Auth("role.info"), RoleList)
	app.POST("/api/role", Auth("role.add"), RoleAdd)
	app.POST("/api/role/:id", Auth("role.update"), RoleUpdate)
	app.DELETE("/api/role/:id", Auth("role.del"), RoleDel)

	app.GET("/api/user", Auth("user.list"), UserList)
	app.GET("/api/user/:id", Auth("user.info"), UserList)
	app.POST("/api/user/:id", Auth("user.update"), UserUpdate)
	app.POST("/api/user/:id/change", Auth("user.change"), UserChange)
	app.POST("/api/user/:id/resetpassword", Auth("user.resetpassword"), UserResetPassword)
	app.DELETE("/api/user/:id", Auth("user.del"), UserDel)

	app.GET("/api/usergroup", Auth("usergroup.list"), UserGroupList)
	app.GET("/api/usergroup/:id", Auth("usergroup.info"), UserGroupList)
	app.POST("/api/usergroup", Auth("usergroup.add"), UserGroupAdd)
	app.POST("/api/usergroup/:id", Auth("usergroup.update"), UserGroupUpdate)
	app.DELETE("/api/usergroup/:id", Auth("usergroup.del"), UserGroupDel)

	app.GET("/api/task", Auth("task.list"), TaskList)
	app.GET("/api/task/:id", Auth("task.info"), TaskList)
	app.POST("/api/task", Auth("task.add"), TaskAdd)
	app.POST("/api/task/:id", Auth("task.update"), TaskUpdate)
	app.POST("/api/task/:id/open", Auth("task.open"), TaskOpen)
	app.POST("/api/task/:id/done", Auth("task.done"), TaskDone)
	app.DELETE("/api/task/:id", Auth("task.del"), TaskDel)
}
