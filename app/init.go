package app

import (
	"time"

	"github.com/gin-gonic/gin"
)

var (
	PM = struct {
		start string
		end   string
	}{}
)

func Init(r *gin.Engine) {

	u := "postgres://:@localhost/taskplan?sslmode=disable"
	switch gin.Mode() {
	case "release":
		u = "postgres://taskplan:shengyun123@localhost/taskplan?sslmode=disable"
	}

	err := InitDB(u)
	if err != nil {
		panic(err)
	}
	if gin.Mode() != "release" {
		DB.LogMode(true)
	}

	genPM()

	InitDBModel(
		new(Permission),
		new(Role),
		new(User),
		new(UserGroup),
		new(Task),
		new(TaskHistory),
		new(Holiday),
		new(Leave),
	)
	DB.InstantSet("gorm:auto_preload", true)

	Router(r)

}

//计算下午上下班时间
func genPM() {
	d := time.Now().Month()
	if d > 4 && d < 10 {
		PM.start = " 13:30:00"
		PM.end = " 18:00:00"
	} else {
		PM.start = " 13:00:00"
		PM.end = " 17:30:00"
	}
}
