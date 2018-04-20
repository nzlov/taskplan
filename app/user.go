package app

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nzlov/gorm"
)

type User struct {
	gorm.Model

	Name     string
	RealName string `form:"realname" binding:"required"`
	Password string `json:"-" form:"password"`
	Status   int    `description:"状态 0 正常 1 禁用"`

	UserGroupID uint `form:"usergroupid"`
	RoleID      uint `form:"roleid"`

	UserGroup UserGroup `gorm:"save_associations:false" validate:"-" form:"-" binding:"-"`
	Role      Role      `gorm:"save_associations:false" validate:"-" form:"-" binding:"-"`
}

type UserRegist struct {
	Name     string `form:"name" binding:"required"`
	RealName string `form:"realname" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func UserAdd(c *gin.Context) {
	tx := NewTx(c)
	var form UserRegist
	if err := c.Bind(&form); err == nil {
		exit := User{}
		err = tx.DB.Where("name = ?", form.Name).First(&exit).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
				return
			}
		} else {
			tx.Error(http.StatusOK, CodeKeyMany, nil)
			return
		}
		user := User{
			Name:     form.Name,
			RealName: form.RealName,
			Password: MakeMd5(form.Password),
			Status:   0,
		}
		err = tx.DB.Create(&user).Error
		if err != nil {
			tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		} else {
			tx.Ok(CodeOK, user)
		}
	} else {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
	}
}

func UserUpdate(c *gin.Context) {
	tx := NewTx(c)

	ids := c.Param("id")
	var form User
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
		return
	}
	err = tx.DB.First(&form, id).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		return
	}

	if err = c.Bind(&form); err == nil {
		exit := User{}
		err = tx.DB.Where("id <> ? and name = ?", ids, form.Name).First(&exit).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
				return
			}
		} else {
			tx.Error(http.StatusOK, CodeKeyMany, nil)
			return
		}
		form.Role.ID = form.RoleID
		form.UserGroup.ID = form.UserGroupID
		// TODO 修改密码判断权限
		err = tx.DB.Save(&form).Error
		if err != nil {
			tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		} else {
			tx.Ok(CodeOK, form)
		}
	} else {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
	}

}

func UserChange(c *gin.Context) {
	tx := NewTx(c)

	ids := c.Param("id")
	var form User
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
		return
	}
	err = tx.DB.First(&form, id).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		return
	}

	req := struct {
		Status *int `form:"status"`
	}{}

	if err = c.Bind(&req); err == nil && req.Status != nil {
		form.Status = *req.Status
		err = tx.DB.Save(&form).Error
		if err != nil {
			tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		} else {
			tx.Ok(CodeOK, form)
		}
	} else {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
	}
}

func UserResetPassword(c *gin.Context) {
	tx := NewTx(c)

	ids := c.Param("id")
	var form User
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
		return
	}
	err = tx.DB.First(&form, id).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		return
	}
	np := RandomString(6)
	form.Password = MakeMd5(np)
	err = tx.DB.Save(&form).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
	} else {
		tx.Ok(CodeOK, np)
	}
}

func UserList(c *gin.Context) {
	tx := NewTx(c)
	id := c.Param("id")
	query := map[string]interface{}{}
	if id != "" {
		query["id = ?"] = id
	}
	offsets := c.Query("offset")
	limits := c.Query("limit")
	filter := strings.TrimSpace(c.Query("filter"))

	offset := int64(-1)
	if strings.TrimSpace(offsets) != "" {
		offset, _ = strconv.ParseInt(offsets, 10, 64)
	}
	limit := int64(-1)
	if strings.TrimSpace(limits) != "" {
		limit, _ = strconv.ParseInt(limits, 10, 64)
	}

	if filter != "" {
		query["name like ?"] = "%" + filter + "%"
		query["real_name like ?"] = "%" + filter + "%"
	}

	cs, _ := c.Get("Session")
	session := cs.(*AuthSession)
	if _, ok := session.Permission["list.all"]; !ok {
		if _, ok = session.Permission["list.group"]; ok {
			query["user_group_id = ?"] = session.User.UserGroupID
		} else {
			query["id = ?"] = session.User.ID
		}
	}

	Users := []User{}
	total, err := DBFind(tx.DB, new(User), &Users, query, nil, c.Query("order"), offset, limit, true)
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
	} else {
		tx.Ok(CodeOK, map[string]interface{}{
			"total": total,
			"data":  Users,
		})
	}
}
func UserDel(c *gin.Context) {
	tx := NewTx(c)
	ids := c.Param("id")
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
		return
	}
	form := User{}
	err = tx.DB.First(&form, id).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		return
	}
	err = tx.DB.Delete(&form).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
	}
	Cache.Delete(form.Name)
	tx.Ok(CodeOK, nil)
}
