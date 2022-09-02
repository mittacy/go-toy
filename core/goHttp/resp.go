package goHttp

type IReply interface {
	GetCode() int
	GetUnknownCode() int
	GetSuccessCode() int
	GetMsg() string
	IsSuccess() bool
	UnmarshalData(dataResult interface{}) error
}

func DefaultReply() IReply {
	return &Reply{}
}

type Reply struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func (r *Reply) GetCode() int {
	return r.Code
}

func (r *Reply) GetUnknownCode() int {
	return 500
}

func (r *Reply) GetSuccessCode() int {
	return 0
}

func (r *Reply) GetMsg() string {
	return r.Msg
}

func (r *Reply) IsSuccess() bool {
	return r.GetCode() == r.GetSuccessCode()
}

func (r *Reply) UnmarshalData(dataResult interface{}) error {
	b, err := jsonMarshal(r.Data)
	if err != nil {
		return err
	}

	return jsonUnmarshal(b, dataResult)
}
