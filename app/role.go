package app

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/nzlov/gorm"
)

type Role struct {
	gorm.Model

	Name       string         `form:"name" binding:"required"`
	Permission pq.StringArray `gorm:"type:varchar(100)[]" form:"permission" binding:"required"`
	Menu       pq.StringArray `gorm:"type:varchar(100)[]" form:"menu" binding:"required"`
}

func RoleAdd(c *gin.Context) {
	tx := NewTx(c)
	var obj Role
	if err := c.Bind(&obj); err == nil {
		exit := Role{}
		err = DB.Where("name = ?", obj.Name).First(&exit).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
				return
			}
		} else {
			tx.Error(http.StatusOK, CodeKeyMany, nil)
			return
		}

		err = DB.Create(&obj).Error
		if err != nil {
			tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		} else {
			tx.Ok(CodeOK, obj)
		}
	} else {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
	}
}

func RoleUpdate(c *gin.Context) {
	tx := NewTx(c)

	ids := c.Param("id")
	var obj Role
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
		return
	}
	err = tx.DB.First(&obj, id).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		return
	}

	if err = c.Bind(&obj); err == nil {
		exit := Role{}
		err = tx.DB.Where("id <> ? and name = ?", ids, obj.Name).First(&exit).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
				return
			}
		} else {
			tx.Error(http.StatusOK, CodeKeyMany, nil)
			return
		}

		err = tx.DB.Save(&obj).Error
		if err != nil {
			tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		} else {
			tx.Ok(CodeOK, obj)
		}
	} else {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
	}
}
func RoleList(c *gin.Context) {
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
	}

	objs := []Role{}
	total, err := DBFind(tx.DB, new(Role), &objs, query, nil, c.Query("order"), offset, limit, true)
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
	} else {
		tx.Ok(CodeOK, map[string]interface{}{
			"total": total,
			"data":  objs,
		})
	}
}
func RoleDel(c *gin.Context) {
	tx := NewTx(c)
	ids := c.Param("id")
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
		return
	}
	form := Role{}
	form.ID = uint(id)
	err = tx.DB.Delete(&form).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
	}
	tx.Ok(CodeOK, nil)
}
