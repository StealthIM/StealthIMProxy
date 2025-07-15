package msap

type sendRequest struct {
	Msg  string `json:"msg" validate:"required"`
	Type int32  `json:"type" validate:"required"`
}

type recallRequest struct {
	MsgID string `json:"msgid" validate:"required"` // int64
}
