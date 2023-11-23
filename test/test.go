package main

import "gorm.io/gorm"

func main() {
	var db *gorm.DB
	key := "用户可控" // 可控将导致注入

	db.Order(key)

}
