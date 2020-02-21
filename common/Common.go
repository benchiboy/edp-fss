package common

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	ERR_CODE_SUCCESS = "0000"
	ERR_CODE_DBERROR = "1001"
	ERR_CODE_TOKENER = "1003"
	ERR_CODE_PARTOEN = "1005"
	ERR_CODE_JSONERR = "2001"
	ERR_CODE_URLERR  = "2005"
	ERR_CODE_NOTFIND = "3000"
	ERR_CODE_EXPIRED = "6000"
	ERR_CODE_TYPEERR = "4000"
	ERR_CODE_STATUS  = "5000"
	ERR_CODE_FAILED  = "9000"
	ERR_CODE_TOOBUSY = "6010"
	ERR_CODE_PDFERR  = "4040"

	STATUS_DISABLED = "B"
	STATUS_ENABLED  = "E"
	STATUS_SUCC     = "S"
	STATUS_INIT     = "I"
	STATUS_FAIL     = "F"
	STATUS_DOING    = "D"

	DEFAULT_PATH      = "/tmp/pdfile/"
	DEFAULT_FONT_SIZE = 14

	FIELD_LOGIN_PASS  = "user_pwd"
	FIELD_ERRORS      = "pwderr_cnt"
	FIELD_UPDATE_TIME = "update_date"
	FIELD_PROC_STATUS = "status"

	EMPTY_STRING = ""
)

var (
	ERROR_MAP map[string]string = map[string]string{
		ERR_CODE_SUCCESS: "执行成功:",
		ERR_CODE_DBERROR: "DB执行错误:",
		ERR_CODE_JSONERR: "JSON格式错误:",
		ERR_CODE_EXPIRED: "时效已经到期:",
		ERR_CODE_TYPEERR: "类型转换错误:",
		ERR_CODE_STATUS:  "状态不正确:",
		ERR_CODE_TOKENER: "获取TOKEN失败:",
		ERR_CODE_PARTOEN: "解析TOKEN错误:",
		ERR_CODE_NOTFIND: "查询没发现提示:",
		ERR_CODE_TOOBUSY: "短信发送太频繁:",
		ERR_CODE_PDFERR:  "创建PDF文件出错:",
	}
)

type ErrorResp struct {
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

const (
	USER_CHARGE = "用户充值"
	FLOW_CHARGE = "charge"
	FLOW_INIT   = "i"
	FLOW_SUCC   = "s"
	FLOW_FAIL   = "f"

	NOW_TIME_FORMAT    = "2006-01-02 15:04:05"
	FIELD_ACCOUNT_BAL  = "Account_bal"
	FIELD_UPDATED_TIME = "Updated_time"

	CODE_SUCC    = "0000"
	CODE_NOEXIST = "1000"

	CODE_FAIL = "2000"

	RESP_SUCC = "0000"
	RESP_FAIL = "1000"

	CODE_TYPE_EDU       = "EDU"
	CODE_TYPE_POSITION  = "POSITION"
	CODE_TYPE_SALARY    = "SALARY"
	CODE_TYPE_WORKYEARS = "WORKYEARS"
	CODE_TYPE_POSICLASS = "POSICLASS"
	CODE_TYPE_REWARDS   = "REWARDS"

	TOKEN_KEY = "u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4"
)

func PrintHead(a ...interface{}) {
	log.Println("========》", a)
}

func PrintTail(a ...interface{}) {
	log.Println("《========", a)
}

func Write_Response(response interface{}, w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8087")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "1728000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "content-type,Action, Module,Authorization")
	fmt.Fprintf(w, string(json))
}
