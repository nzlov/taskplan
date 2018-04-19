package app

func RespData(code int, obj interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code": code,
		"data": obj,
	}
}
