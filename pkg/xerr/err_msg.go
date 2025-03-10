package xerr

var codeText = map[int]string{
	SERVER_COMMON_ERROR: "服务器异常，请稍后再试",
	REQUEST_PARAM_ERROR: "请求参数错误",
	DATABASE_ERROR:      "数据库繁忙，请稍后再试",
}

func ErrMsg(errCode int) string {
	if text, ok := codeText[errCode]; ok {
		return text
	}
	return codeText[SERVER_COMMON_ERROR]
}
