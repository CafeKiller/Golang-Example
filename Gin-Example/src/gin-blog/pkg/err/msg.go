package err

var MsgFlags = map[int]string{
	SUCCESS: "[Ok] 成功",

	ERROR: "[Fail] 错误",

	INVALID_PARAMS: "[Fail] 请求参数异常",

	ERROR_EXIST_TAG:         "[Fail] 该标签名已存在",
	ERROR_NOT_EXIST_TAG:     "[Fail] 该标签不存在",
	ERROR_NOT_EXIST_ARTICLE: "[Fail] 该文章不存在",

	ERROR_AUTH_CHECK_TOKEN_FAIL:    "[Fail] Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "[Fail] Token已超时",
	ERROR_AUTH_TOKEN:               "[Fail] Token 生成失败",
	ERROR_AUTH:                     "[Fail] Token错误",
}

// GetMsg 获取异常/错误信息提示
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
