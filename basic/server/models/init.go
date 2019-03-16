package models

var db map[string]interface{}

func init() {
	db = make(map[string]interface{})
}