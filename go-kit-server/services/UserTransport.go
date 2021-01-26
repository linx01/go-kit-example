package services

import (
	"context"
	"encoding/json"
	"errors"
	"example/utils"
	"net/http"
	"strconv"

	mymux "github.com/gorilla/mux"
)

// DecodeUserRequest ...
func DecodeUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	// if r.URL.Query().Get("userid") != "" {
	// 	userID, _ := strconv.Atoi(r.URL.Query().Get("userid"))
	// 	return UserRequest{ID: userID}, nil
	// }
	// return nil, errors.New("参数错误")
	// 使用第三方路由
	args := mymux.Vars(r)
	uid, ok := args["uid"]
	if ok {
		uidInt, _ := strconv.Atoi(uid)
		return UserRequest{ID: uidInt, Method: r.Method}, nil
	} else {
		return nil, errors.New("参数错误")
	}
}

// EncodeUserResponse ...
func EncodeUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// MyErrorEncoder 自定义错误处理
func MyErrorEncoder(c context.Context, err error, w http.ResponseWriter) {
	contentType := "text/plain;charset=utf-8"
	w.Header().Set("Content-Type", contentType)
	result := err.(*utils.MyError) // error接口断言，转为MyError类型
	w.WriteHeader(result.Code)     // 错误状态码写入
	body := result.Message
	w.Write([]byte(body))
}
