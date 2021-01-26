package services

import (
	"context"
	"example/utils"
	"fmt"
	"strconv"

	"golang.org/x/time/rate"

	"github.com/go-kit/kit/endpoint"
)

// UserRequest ...
type UserRequest struct {
	ID     int `json:"id"`
	Method string
}

// UserResponse ...
type UserResponse struct {
	Result string `json:"result"`
}

// RateLimitMiddleWare ...
func RateLimitMiddleWare(r *rate.Limiter) endpoint.Middleware { // endpoint.Middleware func(endpoint)endpoint
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if r.Allow() {
				return next(ctx, request)
			}
			return nil, utils.GenError(400, "rate limit.")

		}
	}
}

// GenerateEndpoint ...
func GenerateEndpoint(s IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if r, ok := request.(UserRequest); ok {
			userID := r.ID
			switch r.Method {
			case "GET":
				userName := s.GetUserName(userID) + strconv.Itoa(utils.ServicePort)
				return UserResponse{Result: userName}, nil
			case "DELETE":
				result := s.DeleteUserName(userID)
				if result != nil {
					return UserResponse{Result: result.Error()}, result
				} else {
					return UserResponse{Result: fmt.Sprintf("%d删除成功!", userID)}, result
				}

			}

		}
		return UserResponse{Result: ""}, nil
	}
}
