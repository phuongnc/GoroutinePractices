package src

const (
	SUCCESS        = 200
	FAIL           = 500
	INVALID_PARAMS = 400
)

var MsgFlags = map[int]string{
	SUCCESS:        "Success",
	FAIL:           "fail",
	INVALID_PARAMS: "Invalid param request",
}

func GetMsg(code int) string {
	msg, _ := MsgFlags[code]
	return msg
}
