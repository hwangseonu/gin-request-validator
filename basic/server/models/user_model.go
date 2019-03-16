package models

type UserModel struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *UserModel) Save() {
	users[u.Id] = u
}

func NewUserModel(username, password string) *UserModel {
	return &UserModel{
		Id:       GetNextId("users"),
		Username: username,
		Password: password,
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