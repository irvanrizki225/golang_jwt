package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"github.com/irvanrizki225/golang_jwt/auth"
	"github.com/irvanrizki225/golang_jwt/helpers"
	"github.com/irvanrizki225/golang_jwt/models"
	"github.com/irvanrizki225/golang_jwt/utilities"
)

//conec DB
var db = utilities.ConnecDB()

// input register user
type inputRegisterUser struct {
	Name 			string	`json:"name" binding:"required"`
	Occupation		string	`json:"occupation" binding:"required"`
	Username		string	`json:"username" binding:"required"`
	Password		string	`json:"password" binding:"required"`
}

//func register user
func RegisterUser(c *gin.Context) {
	var input inputRegisterUser

	//validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		//validate error
		var errors []string

		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, e.Error())
		}

		errorMessage := gin.H{"errors": errors}

		response := helpers.APIResponse("Register Failed", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}

	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		response := helpers.APIResponse("Create Encripted Password Failed", 500, "error", nil)
		c.JSON(500, response)
		return
	}

	//create user
	newUser := models.User{
		Username: 	input.Username,
		Password: 	string(hashedPassword),
	}

	if err := db.Create(&newUser).Error; err != nil {
		response := helpers.APIResponse("Create User Failed", 500, "error", nil)
		c.JSON(500, response)
		return
	}

	//create employee
	newEmployee := models.Employee{
		UserID: 	newUser.ID,
		Name: 		input.Name,
		Occupation: input.Occupation,
	}

	if err := db.Create(&newEmployee).Error; err != nil {
		response := helpers.APIResponse("Create Employee Failed", 500, "error", nil)
		c.JSON(500, response)
		return
	}

	//generate token
	token, err := auth.NewService().GenerateToken(newEmployee.ID, newUser.ID)
	if err != nil {
		response := helpers.APIResponse("Create Token Failed", 500, "error", nil)
		c.JSON(500, response)
		return
	}

	

	//save token
	newEmployee.Token = token
	if err := db.Save(&newEmployee).Error; err != nil {
		response := helpers.APIResponse("Update Token Failed", 500, "error", nil)
		c.JSON(500, response)
		return
	}

	//response
	response := helpers.APIResponse("Register Success", 200, "success", newEmployee)
	c.JSON(200, response)
}


// Get all Employee
func GetEmployee(c *gin.Context) {
	var employee []models.Employee
	
	// Get employee ID from token
	id := c.MustGet("currentUser").(int)
	fmt.Println("ID : ", id)

	// Get employee data
	if err := db.Where("id = ?", id).Find(&employee).Error; err != nil {
		response := helpers.APIResponse("Get Employee Failed", 400, "error", nil)
		c.JSON(400, response)
		return
	}

	response := helpers.APIResponse("Get Employee Success", 200, "success", employee)
	c.JSON(200, response)
}