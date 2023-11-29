package helpers

type Response struct {
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    Meta        `json:"meta,omitempty"`
}

type Meta struct {
	Page      int    `json:"page,omitempty"`
	NextPage  string `json:"next,omitempty"`
	PrevPage  string `json:"prev,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	TotalData int    `json:"total_data,omitempty"`
	TotalPage int    `json:"total_page,omitempty"`
}
