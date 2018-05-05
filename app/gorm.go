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
func DBFind2(db *gorm.DB, obj, objs interface{}, where map[string]interface{}, order string, offset, limit int64, needTotal bool, ors ...[]map[string]interface{}) (int, error) {
	db = GenQueryDB(db, where, ors...)
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

func GenQueryDB(db *gorm.DB, Where map[string]interface{}, Ors ...[]map[string]interface{}) *gorm.DB {
	for k, v := range Where {
		if len(k) == 0 {
			continue
		}
		db = db.Where(k, v)
	}
	if len(Ors) > 0 {
		s, v := GenOr(Ors...)
		fmt.Printf("GENOR OK:%s:%+v\n", s, v)
		if len(s) > 0 {
			db = db.Where(s, v...)
		}
	}
	return db
}

func GenOr(Ors ...[]map[string]interface{}) (string, []interface{}) {
	fmt.Printf("GENOR:%+v\n", Ors)
	sqlstr := ""
	sqlV := []interface{}{}
	for _, or := range Ors {
		if len(or) == 0 {
			continue
		}
		orstr := ""
		for _, o := range or {
			if len(o) == 0 {
				continue
			}
			ostr := ""
			for sql, value := range o {
				ostr += fmt.Sprintf("AND %s ", sql)
				sqlV = append(sqlV, value)
			}
			orstr = orstr + " OR (" + strings.Trim(ostr, "AND") + ") "
		}
		sqlstr = sqlstr + "AND (" + strings.Trim(orstr, " OR") + ") "
	}

	return strings.Trim(sqlstr, "AND"), sqlV
}
