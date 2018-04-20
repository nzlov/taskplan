package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nzlov/gorm"
	_ "github.com/nzlov/gorm/dialects/postgres"
)

var (
	DB    *gorm.DB
	DBLog *gorm.DB
)

func InitDB(url string) error {
	var err error
	DB, err = gorm.Open("postgres", url)
	if err != nil {
		return err
	}
	return nil
}

func InitDBLog(url string) error {
	var err error
	DBLog, err = gorm.Open("postgres", url)
	if err != nil {
		return err
	}
	return nil
}

type Tx struct {
	c  *gin.Context
	DB *gorm.DB
}

func NewTx(c *gin.Context) *Tx {
	return &Tx{
		c:  c,
		DB: DB.Begin(),
	}
}

func (tx *Tx) Ok(code int, obj interface{}) {
	tx.DB.Commit()
	tx.c.JSON(http.StatusOK, RespData(code, obj))
}
func (tx *Tx) Error(httpcode, code int, obj interface{}) {
	tx.DB.Rollback()
	tx.c.JSON(httpcode, RespData(code, obj))
}

func InitDBModel(models ...interface{}) {
	DB.AutoMigrate(models...)
}

func InitDBLogModel(models ...interface{}) {
	DBLog.AutoMigrate(models...)
}

func DBCount(db *gorm.DB, m interface{}, reqtotal bool) int {
	total := 0
	if reqtotal {
		db.Model(m).Offset(-1).Limit(-1).Count(&total)
	}
	return total
}

func DBFind(db *gorm.DB, obj, objs interface{}, where map[string]interface{}, or []map[string]interface{}, order string, offset, limit int64, needTotal bool) (int, error) {
	db = GenQueryDB(db, where, or)
	return DBCount(db, obj, needTotal), DBOLO(db, order, int(offset), int(limit)).Find(objs).Error
}

func DBOLO(db *gorm.DB, order string, offset, limit int) *gorm.DB {
	order = strings.Trim(order, ",")
	if order != "" {
		orders := strings.Split(order, ",")
		for _, o := range orders {
			n := strings.TrimSpace(o)
			if strings.HasPrefix(n, "-") {
				db = db.Order(strings.Trim(n, "-") + " desc")
			} else {

				db = db.Order(o)
			}
		}
	}
	if offset > -1 {
		db = db.Offset(offset)
	}
	if limit > -1 {
		db = db.Limit(limit)
	}
	return db
}

func GenQueryDB(db *gorm.DB, Where map[string]interface{}, Or []map[string]interface{}) *gorm.DB {
	sqlWhere := ""
	sqlV := []interface{}{}
	for k, v := range Where {
		sqlWhere += "and " + k
		sqlV = append(sqlV, v)
	}
	sqlWhere = strings.TrimLeft(sqlWhere, "and")
	if Or != nil && len(Or) > 0 {
		fmt.Println("OR:", Or, len(Or))
		for _, or := range Or {
			sqls := sqlWhere
			values := append([]interface{}{}, sqlV...)
			for sql, value := range or {
				sqls += "and " + sql
				values = append(values, value)
			}
			db = db.Or(strings.TrimLeft(sqls, "and"), values...)
		}
	} else {
		db = db.Where(sqlWhere, sqlV...)
	}
	return db
}
