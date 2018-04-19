package app

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"
)

func RandomString(lens int) string {
	now := time.Now()
	return MakeMd5(strconv.FormatInt(now.UnixNano(), 10))[:lens]
}

func MakeMd5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}
