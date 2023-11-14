package auth

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	GenerateToken(EmployeeID int, UserID int) (string, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
}
type jwtService struct {
}

var SECRET_KEY = []byte(os.Getenv("SECRET_KEY")) 

func NewService() *jwtService {
	return &jwtService{}
}

// Generate Token
func (s *jwtService) GenerateToken(EmployeeID int, UserID int) (string, error) {
	
	jwtClaim := jwt.MapClaims{
		"employee_id": 	EmployeeID,
		"user_id":		UserID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)

	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

// Validate Token
func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		// if token is not valid
		if !ok {
			return nil, errors.New("Invalid token")
		}

		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return token, err
	}

	return token, nil
}
