package vo

type RespCode uint

const (
	Success RespCode = iota
	Failed
)

type BaseResponse struct {
	Code    RespCode    `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int         `json:"total"`
}
