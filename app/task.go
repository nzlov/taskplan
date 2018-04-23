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

func (t TaskTimeRect) toString() string {
	return fmt.Sprintf("%d-%d", t.Start, t.End)
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

	TimeRect TaskTimeRects `gorm:"type:jsonb" description:"任务时间段"`

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
		resizeTaskTimeRect(tx, &task)
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
		Name        *string        `form:"name"`
		UserGroupID *uint          `form:"usergroupid"`
		UserID      *uint          `form:"userid"`
		TaskID      *uint          `form:"taskid"`
		Description *string        `form:"description" description:"任务描述"`
		Start       *int64         `form:"start" description:"开始时间"`
		End         *int64         `form:"end" description:"结束时间"`
		Remark      *string        `form:"remark" description:"备注"`
		TimeRect    *TaskTimeRects `form:"timerect" description:"任务时间段"`
	}{}
	otaskid := task.ParentTaskID

	if err = c.Bind(&obj); err == nil {
		exit := Task{}

		has := false
		need := false
		needTime := false
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
			needTime = true
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
			needTime = true
		}
		if obj.TimeRect != nil {
			items = append(items, HistoryItem{
				Field: "结束时间",
				Old:   task.TimeRect,
				New:   *obj.TimeRect,
			})
			task.TimeRect = *obj.TimeRect
			has = true
		}
		if task.Start > task.End {
			tx.Error(http.StatusBadRequest, CodeParamsError, "start > end")
			return
		}
		if has {
			task.TaskHistory = append(task.TaskHistory, TaskHistory{
				Action: 2,
				UserID: session.(*AuthSession).User.ID,
				Items:  items,
			})
			if needTime && !task.PTask {
				resizeTaskTimeRect(tx, &task)
			}
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
		task.TimeRect = TaskTimeRects{}
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
	query := map[string]interface{}{}
	all := c.Query("all")
	objs := []Task{}

	id := c.Param("id")
	if id != "" {
		query["id = ?"] = id
	}
	or := []map[string]interface{}{}
	filter := strings.TrimSpace(c.Query("filter"))
	if filter != "" {
		or = append(or, map[string]interface{}{
			"name like ?": "%" + filter + "%",
		})
		or = append(or, map[string]interface{}{
			"description like ?": "%" + filter + "%",
		})
	}
	if all != "t" {
		total, err := DBFind(tx.DB, new(Task), &objs, query, or, c.Query("order"), -1, -1, true)
		if err != nil {
			tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
		} else {
			o := []TaskEasy{}
			d, _ := json.Marshal(objs)
			json.Unmarshal(d, &o)

			tx.Ok(CodeOK, map[string]interface{}{
				"total": total,
				"data":  o,
			})
		}
		return
	}
	offsets := c.Query("offset")
	limits := c.Query("limit")
	pid := strings.TrimSpace(c.Query("pid"))

	offset := int64(-1)
	if strings.TrimSpace(offsets) != "" {
		offset, _ = strconv.ParseInt(offsets, 10, 64)
	}
	limit := int64(-1)
	if strings.TrimSpace(limits) != "" {
		limit, _ = strconv.ParseInt(limits, 10, 64)
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

	total, err := DBFind(tx.DB, new(Task), &objs, query, or, c.Query("order")+",-created_at", offset, limit, true)
	if err != nil {
		tx.Error(http.StatusInternalServerError, CodeDBError, err.Error())
	} else {
		tx.Ok(CodeOK, map[string]interface{}{
			"total": total,
			"data":  objs,
		})
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

// 任务时间间隔缩放
func resizeTaskTimeRect(tx *Tx, task *Task) {
	//获取假日安排
	holidays := []Holiday{}
	err := tx.DB.Where("day between ? and ?", task.Start, task.End).Find(&holidays).Order("day").Error
	if err != nil {
		fmt.Println("genTaskTimeRect:", err.Error())
		return
	}
	holidaymap := map[string]struct{}{}
	for _, v := range holidays {
		holidaymap[time.Unix(v.Day, 0).Format("20006-01-02")] = struct{}{}
	}

	if len(task.TimeRect) == 0 {
		task.TimeRect = genTaskTimeRect(holidaymap, task.Start, task.End)
		return
	}
	newTimeRects := []TaskTimeRect{}
	ostart := task.TimeRect[0].Start
	oend := task.TimeRect[len(task.TimeRect)-1].End

	fmt.Println("resizeTaskTimeRect:", ostart, oend, time.Unix(ostart, 0).Format("2006-01-02 15:04:05"), time.Unix(oend, 0).Format("2006-01-02 15:04:05"))
	fmt.Println("=> resizeTaskTimeRect:", task.Start, task.End, time.Unix(task.Start, 0).Format("2006-01-02 15:04:05"), time.Unix(task.End, 0).Format("2006-01-02 15:04:05"))

	// 扩展开始部分
	if task.Start < ostart {
		newTimeRects = append(newTimeRects, genTaskTimeRect(holidaymap, task.Start, ostart)...)
	}

	for _, v := range task.TimeRect {
		if v.Start >= task.Start {
			newTimeRects = append(newTimeRects, v)
		}
	}

	// 扩展结束部分
	if task.End > oend {
		newTimeRects = append(newTimeRects, genTaskTimeRect(holidaymap, oend, task.End)...)
	}

	tr := []TaskTimeRect{}
	trl := len(newTimeRects)
	for i := 0; i < trl; i++ {
		//如果两个时间段差1s则连接时间段
		n := newTimeRects[i]

		if n.Start > task.End {
			continue
		}
		if i+1 < trl {
			if n.End+1 == newTimeRects[i+1].Start {
				n.End = newTimeRects[i+1].End
				i++
			}
		}
		if n.End > task.End {
			n.End = task.End
		}
		if n.Start == n.End {
			continue
		}
		tr = append(tr, n)
	}

	task.TimeRect = TaskTimeRects(tr)
}

// 生成任务的时间间隔
func genTaskTimeRect(holidaymap map[string]struct{}, start, end int64) (ts []TaskTimeRect) {
	if start == 0 || end == 0 {
		return
	}
	if start > end {
		return
	}
	fmt.Println("genTaskTimeRect:", time.Unix(start, 0).Format("2006-01-02 15:04:05"), time.Unix(end, 0).Format("2006-01-02 15:04:05"))

	ts = append(ts, gentimerect(start, end)...)

	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)
	//判断当前日期加一天后是否大于结束日期
	//如果大于结束日期开始计算last日期的时间段
	//循环
	last := startTime.AddDate(0, 0, 1)
	last, err := time.ParseInLocation("2006-01-02 15:04:05", last.Format("2006-01-02")+" 08:29:59", time.Local)
	if err != nil {
		fmt.Println("genTaskTimeRect:", err.Error())
		return
	}
	v := last.Sub(endTime)

	for v < 0 {
		//获取last的年月日
		lastYMD := last.Format("2006-01-02")
		if _, ok := holidaymap[lastYMD]; !ok {
			ts = append(ts, gentimerect(last.Unix(), end)...)
		}
		last, _ = time.ParseInLocation("2006-01-02 15:04:05", last.AddDate(0, 0, 1).Format("2006-01-02")+" 08:30:00", time.Local)
		v = last.Sub(endTime)
	}
	return
}
func gentimerect(start, end int64) TaskTimeRects {
	ts := TaskTimeRects{}

	startTime := time.Unix(start, 0)
	ymd := startTime.Format("2006-01-02")
	rgCheck := time.Unix(end, 0).Format("2006-01-02") == ymd
	//	ams, err := time.ParseInLocation("2006-01-02 15:04:05", ymd+" 08:30:00", loc)
	//	if err != nil {
	//		panic(err)
	//	}
	//	ames := ams.Unix()
	ame, err := time.ParseInLocation("2006-01-02 15:04:05", ymd+" 12:00:00", time.Local)
	if err != nil {
		panic(err)
	}
	amet := ame.Unix()

	//start时间在早上8点半之前则从8:30开始，否则用start时间
	//	if ames < start {
	//		ames = start
	//	} else if !rgCheck {
	//		ames = start
	//	}
	// nstart 是否在上午
	if start < amet {
		// end时间是否超过上午下班点
		if amet < end {
			ts = append(ts, TaskTimeRect{
				Start: start,
				End:   amet,
			})
		} else {
			ts = append(ts, TaskTimeRect{
				Start: start,
				End:   end,
			})
			return ts
		}
	}

	pms, err := time.ParseInLocation("2006-01-02 15:04:05", ymd+PM.start, time.Local)
	if err != nil {
		panic(err)
	}
	pmes := pms.Unix()
	pme, err := time.ParseInLocation("2006-01-02 15:04:05", ymd+PM.end, time.Local)
	if err != nil {
		panic(err)
	}
	pmet := pme.Unix()

	if pmes < start {
		pmes = start
	}
	if pmes < pmet {
		if pmet < end {
			if rgCheck {
				pmet = end
			}
			ts = append(ts, TaskTimeRect{
				Start: pmes,
				End:   pmet,
			})
		} else {
			ts = append(ts, TaskTimeRect{
				Start: pmes,
				End:   end,
			})
			return ts
		}
	}

	return ts
}
