package models

import (
	"fmt"

	"gorm.io/gorm"
	"github.com/irvanrizki225/golang_jwt/utilities"
)

var db = utilities.ConnecDB()

type User struct {
	ID 			int    	`json:"id" gorm:"primary_key"`
	Username 	string 	`json:"username" gorm:"unique"`
	Password 	string 	`json:"password"`
	gorm.Model
}

type Employee struct {
	ID 			int    	`json:"id" gorm:"primary_key"`
	UserID 		int 	`json:"user_id" gorm:"foreignkey:UserID"`
	Name 		string 	`json:"name"`
	Occupation 	string 	`json:"occupation"`
	Token 		string 	`json:"token"`
	gorm.Model
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Employee{})
	fmt.Println("Successfully Migrate User Table!")
}


