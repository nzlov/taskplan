package app

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nzlov/gorm"
)

type Leave struct {
	gorm.Model

	UserID uint  `form:"userid" binding:"required"`
	Start  int64 `form:"start" description:"开始时间"`
	End    int64 `form:"end" description:"结束时间"`

	User User `gorm:"save_associations:false" validate:"-" form:"-" binding:"-"`
}

func LeaveAdd(c *gin.Context) {
	tx := NewTx(c)
	var obj Leave
	if err := c.Bind(&obj); err == nil {
		exit := Leave{}
		err = DB.LogMode(true).Where("\"start\" between ? and ?", obj.Start, obj.End).Or("\"end\" between ? and ?", obj.Start, obj.End).First(&exit).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
				return
			}
		} else {
			tx.Error(http.StatusOK, CodeKeyMany, "假期重复")
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

func LeaveList(c *gin.Context) {
	tx := NewTx(c)
	id := c.Param("id")
	query := map[string]interface{}{}
	if id != "" {
		query["id = ?"] = id
	}
	offsets := c.Query("offset")
	limits := c.Query("limit")

	offset := int64(-1)
	if strings.TrimSpace(offsets) != "" {
		offset, _ = strconv.ParseInt(offsets, 10, 64)
	}
	limit := int64(-1)
	if strings.TrimSpace(limits) != "" {
		limit, _ = strconv.ParseInt(limits, 10, 64)
	}

	// TODO 请假列表过滤 用户、时间段

	session, _ := c.Get("Session")
	if c.Query("list") == "self" {
		query["user_id = ?"] = session.(*AuthSession).User.ID
	}

	objs := []Leave{}
	total, err := DBFind(tx.DB, new(Leave), &objs, query, nil, c.Query("order"), offset, limit, true)
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
	} else {
		tx.Ok(CodeOK, map[string]interface{}{
			"total": total,
			"data":  objs,
		})
	}
}
func LeaveDel(c *gin.Context) {
	tx := NewTx(c)
	ids := c.Param("id")
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
		return
	}
	form := Leave{}
	form.ID = uint(id)
	err = tx.DB.Delete(&form).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
	}

	tx.Ok(CodeOK, nil)
}
