package api

type ApiResponse struct {
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg,omitempty"`
	Data   any    `json:"data,omitempty"`
}
