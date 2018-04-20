package app

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {

	err := InitDB("postgres://:@localhost/taskplan?sslmode=disable")
	if err != nil {
		panic(err)
	}
	//DB.LogMode(true)

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
