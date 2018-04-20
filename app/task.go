package app

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nzlov/gorm"
)

type TaskHistory struct {
	gorm.Model

	TaskID uint

	Action int `description:"动作 1 创建 2 修改 3 完成 4 重新开始 5 删除"`

	Items HistoryItems `gorm:"type:jsonb"`

	UserID uint

	User User `gorm:"save_associations:false"`
}
type HistoryItem struct {
	Field string
	Old   interface{}
	New   interface{}
}
type HistoryItems []HistoryItem

func (j HistoryItems) Value() (driver.Value, error) {
	return json.Marshal(&j)
}

// Scan scan value into Jsonb
func (j *HistoryItems) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	return json.Unmarshal(bytes, j)
}

type TaskEasy struct {
	ID   uint
	Name string
}

//任务间隔
type TaskTimeRect struct {
	Start int64
	End   int64
}

type TaskTimeRects []TaskTimeRect

func (j TaskTimeRects) Value() (driver.Value, error) {
	return json.Marshal(&j)
}

// Scan scan value into Jsonb
func (j *TaskTimeRects) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	return json.Unmarshal(bytes, j)
}

type Task struct {
	gorm.Model

	Name         string
	UserGroupID  uint
	UserID       uint
	ParentTaskID uint
	PTask        bool   `description:"是否为父级任务"`
	Description  string `description:"任务描述"`

	CreateUserID uint `gorm:"createuserid" description:"创建人ID"`

	Start   int64 `description:"开始时间"`
	End     int64 `description:"结束时间"`
	RealEnd int64 `description:"真正结束时间"`

	Status int `description:"状态 1 创建 2 完成 3 重新开始"`

	TimeRect TaskTimeRect `gorm:"type:jsonb" description:"任务时间段"`

	TaskHistory []TaskHistory `gorm:"ASSOCIATION_AUTOUPDATE:false"`

	User       User      `gorm:"save_associations:false"`
	CreateUser User      `gorm:"save_associations:false"`
	UserGroup  UserGroup `gorm:"save_associations:false"`
}

func TaskAdd(c *gin.Context) {
	tx := NewTx(c)
	obj := struct {
		Name        string `form:"name" binding:"required"`
		UserGroupID uint   `form:"usergroupid"`
		TaskID      *uint  `form:"taskid" description:"上级任务ID"`
		UserID      uint   `gorm:"userid"`
		Description string `form:"description" description:"任务描述"`
		Start       int64  `form:"start" description:"开始时间"`
		End         int64  `form:"end" description:"结束时间"`
	}{}

	if err := c.Bind(&obj); err == nil {
		exit := Task{}
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
		if obj.Start > obj.End {
			tx.Error(http.StatusBadRequest, CodeParamsError, "start > end")
			return
		}
		session, _ := c.Get("Session")
		task := Task{
			Name:         obj.Name,
			UserGroupID:  obj.UserGroupID,
			UserID:       obj.UserID,
			CreateUserID: session.(*AuthSession).User.ID,
			Description:  obj.Description,
			Status:       1,
			Start:        obj.Start,
			End:          obj.End,

			TaskHistory: []TaskHistory{
				{
					Action: 1,
					UserID: session.(*AuthSession).User.ID,
					Items:  HistoryItems{},
				},
			},
		}
		need := false
		if obj.TaskID != nil && *obj.TaskID > 0 {
			parent := Task{}
			err = tx.DB.First(&parent, *obj.TaskID).Error
			if err != nil {
				tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
				return
			}
			need = true
			task.ParentTaskID = parent.ID
			task.TaskHistory[0].Items = append(task.TaskHistory[0].Items, HistoryItem{
				Field: "父级任务",
				Old:   "",
				New:   parent.Name,
			})
		}
		err = DB.Create(&task).Error
		if err != nil {
			tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		} else {
			if need {
				updatesTask(tx, task, task.ParentTaskID, false)
			}
			tx.Ok(CodeOK, obj)
		}
	} else {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
	}
}

func TaskUpdate(c *gin.Context) {
	tx := NewTx(c)

	ids := c.Param("id")
	var task Task
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
		return
	}
	err = tx.DB.First(&task, id).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		return
	}
	session, _ := c.Get("Session")
	working := task.Start > 0 && task.Start <= time.Now().Unix()

	//任务开始时间小于当前时间 则任务任务已经开始 只有拥有 `task.expire` 延期权限才可以修改任务信息
	if _, ok := session.(*AuthSession).Permission["task.expire"]; !ok && working {
		tx.Error(http.StatusBadRequest, CodeParamsError, "task working")
		return
	}

	// 父级任务不允许编辑
	childs := 0
	err = tx.DB.Model(new(Task)).Where("parent_task_id = ?", id).Count(&childs).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		return
	}
	if childs > 0 {
		tx.Error(http.StatusBadRequest, CodeParamsError, "parent task don't edit")
		return
	}

	obj := struct {
		Name        *string `form:"name"`
		UserGroupID *uint   `form:"usergroupid"`
		UserID      *uint   `form:"userid"`
		TaskID      *uint   `form:"taskid"`
		Description *string `form:"description" description:"任务描述"`
		Start       *int64  `form:"start" description:"开始时间"`
		End         *int64  `form:"end" description:"结束时间"`
		Remark      *string `form:"remark" description:"备注"`
	}{}
	otaskid := task.ParentTaskID
	has := false
	need := false

	if err = c.Bind(&obj); err == nil {
		exit := Task{}

		items := HistoryItems{}
		if working {
			fmt.Println("WORK:", working, obj.Remark)
			if obj.Remark == nil {
				tx.Error(http.StatusBadRequest, CodeParamsError, "任务已经开始，如果去修改请填写备注！")
				return
			}
			items = append(items, HistoryItem{
				Field: "Remark",
				New:   *obj.Remark,
			})
			has = true
		}

		if obj.Name != nil && *obj.Name != task.Name {
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
			items = append(items, HistoryItem{
				Field: "Name",
				Old:   task.Name,
				New:   *obj.Name,
			})
			task.Name = *obj.Name
			has = true
		}
		if obj.UserGroupID != nil && *obj.UserGroupID != task.UserGroupID {
			if *obj.UserGroupID == 0 {
				items = append(items, HistoryItem{
					Field: "资源组",
					Old:   task.UserGroup.Name,
					New:   "",
				})
			} else {
				usergroup := UserGroup{}
				err = tx.DB.First(&usergroup, *obj.UserGroupID).Error
				if err != nil {
					tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
					return
				}
				items = append(items, HistoryItem{
					Field: "资源组",
					Old:   task.UserGroup.Name,
					New:   usergroup.Name,
				})
			}
			task.UserGroupID = *obj.UserGroupID
			task.UserGroup.ID = *obj.UserGroupID
			has = true
		}
		if obj.UserID != nil && *obj.UserID != task.UserID {
			if *obj.UserID == 0 {
				items = append(items, HistoryItem{
					Field: "执行人",
					Old:   task.User.Name,
					New:   "",
				})
			} else {
				user := User{}
				err = tx.DB.First(&user, *obj.UserID).Error
				if err != nil {
					tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
					return
				}
				items = append(items, HistoryItem{
					Field: "执行人",
					Old:   task.User.Name,
					New:   user.RealName,
				})
			}
			task.UserID = *obj.UserID
			task.User.ID = *obj.UserID
			has = true
		}
		if obj.TaskID != nil && *obj.TaskID != task.ParentTaskID {
			if *obj.TaskID == 0 {
				items = append(items, HistoryItem{
					Field: "父级任务：",
					Old:   task.Name,
					New:   "",
				})
			} else {
				parent := Task{}
				err = tx.DB.First(&parent, *obj.TaskID).Error
				if err != nil {
					tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
					return
				}
				items = append(items, HistoryItem{
					Field: "父级任务：",
					Old:   task.Name,
					New:   parent.Name,
				})
			}
			task.ParentTaskID = *obj.TaskID
			has = true
			need = true
		}
		if obj.Description != nil && *obj.Description != task.Description {
			items = append(items, HistoryItem{
				Field: "简介",
				Old:   task.Description,
				New:   *obj.Description,
			})
			task.Description = *obj.Description
			has = true
		}
		if obj.Start != nil && *obj.Start != task.Start {
			items = append(items, HistoryItem{
				Field: "开始时间",
				Old:   formatTime(task.Start),
				New:   formatTime(*obj.Start),
			})
			task.Start = *obj.Start
			has = true
			need = true
		}
		if obj.End != nil && *obj.End != task.End {
			items = append(items, HistoryItem{
				Field: "结束时间",
				Old:   formatTime(task.End),
				New:   formatTime(*obj.End),
			})
			task.End = *obj.End
			has = true
			need = true
		}
		if has {
			task.TaskHistory = append(task.TaskHistory, TaskHistory{
				Action: 2,
				UserID: session.(*AuthSession).User.ID,
				Items:  items,
			})
			err = tx.DB.Save(&task).Error
			if err != nil {
				tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
			} else {
				if need && task.ParentTaskID > 0 {
					if task.ParentTaskID != otaskid {
						fmt.Println("更新老父级")
						updatesTask(tx, task, otaskid, true)
					}
					fmt.Println("更新父级")
					updatesTask(tx, task, task.ParentTaskID, false)
				}
			}
		}
		tx.Ok(CodeOK, task)
	} else {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
	}
}

//递归更新父级任务
// action 0 更新 1 删除
func updatesTask(tx *Tx, child Task, id uint, del bool) {
	fmt.Printf("updatesTask init:%+v\n%v----%v\n", child, id, del)
	if child.ParentTaskID == 0 {
		return
	}
	task := Task{}
	err := tx.DB.First(&task, id).Error
	if err != nil {
		fmt.Println("updatesTask has Error:", err)
		return
	}
	need := false

	start, end := int64(9999999999999999), int64(0)
	allok := true
	tasks := []Task{}
	err = tx.DB.Where("parent_task_id=? and id != ?", task.ID, child.ID).Find(&tasks).Error
	if err != nil {
		fmt.Println("updatesTask has Error:", err)
		return
	}
	tasks = append(tasks, child)
	fmt.Println("updatesTask", "childs count:", len(tasks))

	ptask := false
	for _, v := range tasks {
		if del && v.ID == task.ID {
			fmt.Println("updatesTask", "抛弃当前")
			continue
		}
		ptask = true
		if v.End > end {
			fmt.Println("updatesTask", v.End, end, "替换 end")
			end = v.End
		}
		if v.Start < start {
			fmt.Println("updatesTask", v.Start, start, "替换 start")
			start = v.Start
		}
		if allok && v.Status != 2 {
			fmt.Println("updatesTask", "allok")
			allok = false
		}
	}

	if task.Start != start {
		need = true
		task.Start = start
	}
	if task.End != end {
		need = true
		task.End = end
	}

	fmt.Println("updatesTask", "allok", allok)
	if allok {
		if task.Status != 2 {
			fmt.Println("updatesTask", "allok", allok, task.Status, "399")
			task.Status = 2
			task.RealEnd = time.Now().Unix()
			need = true
		}
	} else {
		fmt.Println("updatesTask", "allok", allok, task.Status, "404")
		if task.Status == 2 {
			fmt.Println("updatesTask", "allok", allok, task.Status, "406")
			task.Status = 3
			task.RealEnd = 0
			need = true
		}
	}

	fmt.Println("updatesTask update", "need", need)
	//判断是否需要更新
	if need || task.PTask != ptask {
		task.PTask = ptask
		fmt.Println("updatesTask update")
		err = tx.DB.Save(&task).Error
		if err != nil {
			return
		}
		if task.ParentTaskID > 0 {
			updatesTask(tx, task, task.ParentTaskID, false)
		}
	}
}

func TaskList(c *gin.Context) {
	tx := NewTx(c)
	id := c.Param("id")
	query := map[string]interface{}{}
	if id != "" {
		query["id = ?"] = id
	}
	offsets := c.Query("offset")
	limits := c.Query("limit")
	all := c.Query("all")
	pid := strings.TrimSpace(c.Query("pid"))
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
		query["description like ?"] = "%" + filter + "%"
	}

	if pid != "" {
		query["parent_task_id = ? "] = pid
	} else {
		query["parent_task_id = ? "] = 0
	}

	session, _ := c.Get("Session")
	list := c.DefaultQuery("list", "")
	switch list {
	case "group":
		query["user_group_id = ?"] = session.(*AuthSession).User.UserGroupID
	case "self":
		query["user_id = ?"] = session.(*AuthSession).User.ID

	}

	objs := []Task{}
	if all == "t" {
		total, err := DBFind(tx.DB.LogMode(true), new(Task), &objs, query, nil, c.Query("order")+",-created_at", offset, limit, true)
		if err != nil {
			tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		} else {
			tx.Ok(CodeOK, map[string]interface{}{
				"total": total,
				"data":  objs,
			})
		}
	} else {
		total, err := DBFind(tx.DB, new(Task), &objs, query, nil, c.Query("order"), -1, -1, true)
		if err != nil {
			tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		} else {
			o := []TaskEasy{}
			d, _ := json.Marshal(objs)
			fmt.Println(string(d))
			json.Unmarshal(d, &o)

			tx.Ok(CodeOK, map[string]interface{}{
				"total": total,
				"data":  o,
			})
		}
	}
}
func TaskDel(c *gin.Context) {
	tx := NewTx(c)
	ids := c.Param("id")
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
		return
	}
	task := Task{}
	err = tx.DB.First(&task, id).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		return
	}
	session, _ := c.Get("Session")
	task.TaskHistory = append(task.TaskHistory, TaskHistory{
		Action: 5,
		Items:  HistoryItems{},
		UserID: session.(*AuthSession).User.ID,
	})
	err = tx.DB.Save(&task).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
	}
	if task.PTask {
		tasks := []Task{}
		err = tx.DB.Where("parent_task_id=?", task.ID).Find(&tasks).Error
		if err != nil {
			return
		}
		for _, t := range tasks {
			t.TaskHistory = append(t.TaskHistory, TaskHistory{
				Action: 5,
				Items:  HistoryItems{},
				UserID: session.(*AuthSession).User.ID,
			})
			err = tx.DB.Save(&t).Error
			if err != nil {
				tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
			}
			err = tx.DB.Delete(&t).Error
			if err != nil {
				tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
			}
		}
	}

	err = tx.DB.Delete(&task).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
	}
	if task.ParentTaskID > 0 {
		updatesTask(tx, task, task.ParentTaskID, false)
	}
	tx.Ok(CodeOK, nil)
}

func TaskOpen(c *gin.Context) {
	tx := NewTx(c)

	ids := c.Param("id")
	var task Task
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
		return
	}
	err = tx.DB.First(&task, id).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		return
	}

	if task.Status != 2 {
		tx.Error(http.StatusBadRequest, CodeParamsError, "任务为结束")
		return
	}
	task.Status = 3
	session, _ := c.Get("Session")

	task.TaskHistory = append(task.TaskHistory, TaskHistory{
		Action: 4,
		Items:  HistoryItems{},
		UserID: session.(*AuthSession).User.ID,
	})

	err = tx.DB.Save(&task).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		return
	}
	if task.ParentTaskID > 0 {
		updatesTask(tx, task, task.ParentTaskID, false)
	}
	tx.Ok(CodeOK, task)
}

func TaskDone(c *gin.Context) {
	tx := NewTx(c)

	ids := c.Param("id")
	var task Task
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		tx.Error(http.StatusBadRequest, CodeParamsError, err.Error())
		return
	}
	err = tx.DB.First(&task, id).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		return
	}

	if task.Status == 2 {
		tx.Error(http.StatusBadRequest, CodeParamsError, "任务已结束")
		return
	}
	task.Status = 2
	task.RealEnd = time.Now().Unix()
	session, _ := c.Get("Session")

	task.TaskHistory = append(task.TaskHistory, TaskHistory{
		Action: 3,
		Items:  HistoryItems{},
		UserID: session.(*AuthSession).User.ID,
	})

	err = tx.DB.Save(&task).Error
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		return
	}
	if task.ParentTaskID > 0 {
		updatesTask(tx, task, task.ParentTaskID, false)
	}
	tx.Ok(CodeOK, task)
}

func formatTime(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}
