package app

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nzlov/gorm"
)

type Holiday struct {
	gorm.Model

	Day int64 `form:"day" binding:"required"`
}

func HolidayAdd(c *gin.Context) {
	tx := NewTx(c)
	var obj Holiday
	if err := c.Bind(&obj); err == nil {
		exit := Holiday{}
		err = DB.Where("day = ?", obj.Day).First(&exit).Error
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

func HolidayUpdate(c *gin.Context) {
	tx := NewTx(c)

	ids := c.Param("id")
	var obj Holiday
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
		exit := Holiday{}
		err = tx.DB.Where("id <> ? and day = ?", ids, obj.Day).First(&exit).Error
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
func HolidayList(c *gin.Context) {
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

	objs := []Holiday{}
	total, err := DBFind(tx.DB, new(Holiday), &objs, query, nil, c.Query("order"), offset, limit, true)
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
	} else {
		tx.Ok(CodeOK, map[string]interface{}{
			"total": total,
			"data":  objs,
		})
	}
}
func HolidayDel(c *gin.Context) {
	tx := NewTx(c)
	ids := c.Param("id")
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
		return
	}
	form := Holiday{}
	form.ID = uint(id)
	err = tx.DB.Delete(&form).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
	}

	tx.Ok(CodeOK, nil)
}
