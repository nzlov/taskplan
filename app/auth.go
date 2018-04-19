package app

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nzlov/cache2go"
)

type AuthSession struct {
	User       User
	Token      string              `json:"token"`
	Menu       []string            `json:"menu"`
	Permission map[string]struct{} `json:"permission"`
}

var Cache = cache2go.Cache("session")

func Auth(p string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// time.Sleep(time.Second / 2)
		user := strings.TrimSpace(c.Request.Header.Get("X-AppUser"))
		if user == "" {
			c.JSON(403, map[string]interface{}{
				"code": CodeNoLogin,
			})
			c.Abort()
			return
		}
		sign := strings.TrimSpace(c.Request.Header.Get("X-AppSign"))
		if sign == "" {
			c.JSON(403, map[string]interface{}{
				"code": CodeNoLogin,
			})
			c.Abort()
			return
		}
		r, err := Cache.Value(user)
		if err != nil {
			c.JSON(403, map[string]interface{}{
				"code": CodeNoLogin,
			})
			c.Abort()
			return
		}
		session := r.Data().(*AuthSession)
		if sign == genSign(user, session.Token, c.Request.URL.String()) {
			if _, ok := session.Permission[p]; ok {
				c.Set("Session", session)
				c.Next()
				return
			} else {
				c.JSON(403, map[string]interface{}{
					"code": CodeNoPermission,
				})
			}
		} else {
			c.JSON(403, map[string]interface{}{
				"code": CodeNoLogin,
			})
		}
		c.Abort()
	}
}

func CheckLogin(c *gin.Context) {
	user := strings.TrimSpace(c.Request.Header.Get("X-AppUser"))
	if user == "" {
		c.JSON(403, map[string]interface{}{
			"code": CodeNoLogin,
		})
		return
	}
	sign := strings.TrimSpace(c.Request.Header.Get("X-AppSign"))
	if sign == "" {
		c.JSON(403, map[string]interface{}{
			"code": CodeNoLogin,
		})
		return
	}
	r, err := Cache.Value(user)
	if err != nil {
		c.JSON(403, map[string]interface{}{
			"code": CodeNoLogin,
		})
		return
	}
	session := r.Data().(*AuthSession)
	if sign == genSign(user, session.Token, "/api/login") {
		c.JSON(200, map[string]interface{}{
			"code": CodeLoginOk,
		})
	} else {
		c.JSON(403, map[string]interface{}{
			"code": CodeNoLogin,
		})
	}
}

func genSign(user, token, u string) string {
	u, _ = url.QueryUnescape(u)
	sign := MakeMd5(user + u + token)
	fmt.Println(user, token, u, sign)
	return sign
}
