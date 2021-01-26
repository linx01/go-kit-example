package services

import "errors"

// IUserService ...
type IUserService interface {
	GetUserName(userID int) (userName string)
	DeleteUserName(userID int) (err error)
}

// UserService ...
type UserService struct{}

// GetUserName ...
func (u *UserService) GetUserName(userID int) (userName string) {
	if userID == 100 {
		return "admin"
	}
	return "guest"
}

// DeleteUserName ...
func (u *UserService) DeleteUserName(userID int) (err error) {
	if userID == 100 {
		return errors.New("管理人员100不可删除！")
	}
	return nil
}
