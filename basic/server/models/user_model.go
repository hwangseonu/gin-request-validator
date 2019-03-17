package models

import "github.com/gin-gonic/gin"

type UserModel struct {
	Id       int      `json:"id"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Nickname string   `json:"nickname"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

func (u *UserModel) Save() {
	users[u.Id] = u
}

func NewUserModel(username, password, nickname, email string, roles ...string) *UserModel {
	return &UserModel{
		Id:       GetNextId("users"),
		Username: username,
		Password: password,
		Nickname: nickname,
		Email:    email,
		Roles:    roles,
	}
}

func FindUserById(id int) *UserModel {
	return users[id]
}

func FindUserByUsername(username string) *UserModel {
	for _, u := range users {
		if u.Username == username {
			return u
		}
	}
	return nil
}

func ExistsUserByUsername(username string) bool {
	return FindUserByUsername(username) != nil
}

func ExistsUserByUsernameOrNicknameOrEmail(username, nickname, email string) bool {
	for _, u := range users {
		if u.Username == username || u.Nickname == nickname || u.Email == email {
			return true
		}
	}
	return false
}

func GetAccounts() gin.Accounts {
	accounts := make(map[string]string)
	for _, u := range users {
		accounts[u.Username] = u.Password
	}
	return accounts
}
