package authentication

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"reflect"
	"time"
)

type User struct {
	Name     string `json:"name"`
	Password string `json:"pass"`
}

var users = []User{{"test", "test123"},}

const (
	// TODO: Secret should be env variable
	AUTHSECRET     = "akskdjfkölkjaksdASDFAWERkmdlaksdfajdfi;HDKzuiwehrjahljhfaiwulezrualihds"
	tokenValidTime = time.Minute * time.Duration(1)
)

func CreateToken(userName string, password string) (string, error) {
	if !isRegistered(userName, password) {
		return "", errors.New(fmt.Sprintf("User %s is not registered or wrong password", userName))
	}
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":    userName,
		"expires": time.Now().Add(tokenValidTime).Unix(),
	})

	// Sign and get the complete encoded token as a string using the AUTHSECRET
	return token.SignedString([]byte(AUTHSECRET))
}

func ValidateToken(tokenString string) (jwt.MapClaims, bool) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// AUTHSECRET is a []byte containing your AUTHSECRET, e.g. []byte("my_secret_key")
		return []byte(AUTHSECRET), nil
	})
	if err != nil {
		return nil, false
	}
	if token.Valid {
		// TODO: Token only valid as long as expiration time is before now!
		return claims, true
	}
	return nil, false

}

// TODO: Database lookup instead of local array; password hashing
func isRegistered(userName string, password string) bool {
	user := User{userName, password}
	for _, u := range users {
		if reflect.DeepEqual(user, u) {
			return true
		}
	}
	return false
}
