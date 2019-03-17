package models

var users map[int]*UserModel
var posts map[int]*PostModel
var autoIncrement map[string]int

func init() {
	users = make(map[int]*UserModel)
	posts = make(map[int]*PostModel)
	autoIncrement = make(map[string]int)
}

func GetNextId(name string) int {
	i, ok := autoIncrement[name]
	if ok {
		return i
	} else {
		autoIncrement[name] = 0
		return 0
	}
}