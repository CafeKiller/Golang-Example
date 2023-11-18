package err

const (
	SUCCESS = 2000

	INVALID_PARAMS = 4000

	ERROR = 5000

	ERROR_EXIST_TAG         = 5101
	ERROR_NOT_EXIST_TAG     = 5102
	ERROR_NOT_EXIST_ARTICLE = 5103

	ERROR_AUTH_CHECK_TOKEN_FAIL    = 5201
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 5202
	ERROR_AUTH_TOKEN               = 5203
	ERROR_AUTH                     = 5204
)
