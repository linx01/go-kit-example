package services

import (
	// "bufio"

	"context"
	"encoding/json"

	// "fmt"
	// "io"
	"net/http"
	"strconv"
)

var buf [128]byte
var ret []byte

// GetUserInfoEncodeReqFunc ...
func GetUserInfoEncodeReqFunc(c context.Context, req *http.Request, r interface{}) error {
	// 断言
	usereq := r.(UserRequest)
	req.URL.Path += "/user/" + strconv.Itoa(usereq.ID)
	return nil
}

// GetUserInfoDecodeResFunc ...
func GetUserInfoDecodeResFunc(c context.Context, res *http.Response) (response interface{}, err error) {

	userRes := UserResponse{}
	err = json.NewDecoder(res.Body).Decode(&userRes)
	if err != nil {
		return userRes, err
	}

	return userRes, nil

}
