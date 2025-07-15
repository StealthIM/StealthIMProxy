package errorcode

const (
	// ProxyFailed 服务器出现无法确定类型的自动捕获的错误
	ProxyFailed int32 = 1500 + iota
	// ProxyBadRequest 错误的请求
	ProxyBadRequest
	// ProxyAuthFailed 认证失败
	ProxyAuthFailed
)
