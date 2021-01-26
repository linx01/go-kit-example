package services

// UserRequest ...
type UserRequest struct {
	ID     int `json:"id"`
	Method string
}

// UserResponse ...
type UserResponse struct {
	Result string `json:"result"`
}
