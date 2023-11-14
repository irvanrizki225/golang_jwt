package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"net/http"
	"strings"
	"log"

	// "github.com/irvanrizki225/golang_jwt/utilities"
	// "github.com/irvanrizki225/golang_jwt/models"
	"github.com/irvanrizki225/golang_jwt/auth"
	"github.com/irvanrizki225/golang_jwt/controllers"
	"github.com/irvanrizki225/golang_jwt/helpers"
)

//jika mau push ke repo harus ganti pass db nya

func main() {
	// //connec DB
	// utilities.ConnecDB()

	// // //migrate table
	// models.Migrate(db)

	authServices := auth.NewService()

	router := gin.Default()
	api := router.Group("/api/v1")

	//register user
	api.POST("/register", controllers.RegisterUser)

	//get data employee
	api.GET("/employee", authMiddleware(authServices), controllers.GetEmployee)
	//run server
	log.Fatal(http.ListenAndServe(":8080", router))

	log.Println("Server running at port 8080")
}

// auth middleware
func authMiddleware(authService auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helpers.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")

		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helpers.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helpers.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		employeeID := int(claim["employee_id"].(float64))

		c.Set("currentUser", employeeID)
	}
}
