package authentication

import (
	"context"
	"errors"
	"fmt"
	"github.com/alex1ai/ugr-master-cc/data"
	"github.com/dgrijalva/jwt-go"
	"github.com/mongodb/mongo-go-driver/bson"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Name     string `json:"name" bson:"name"`
	Password string `json:"pass" bson:"pass"`
}


const (
	// TODO: Secret should be env variable
	AUTHSECRET     = "akskdjfkölkjaksdASDFAWERkmdlaksdfajdfi;HDKzuiwehrjahljhfaiwulezrualihds"
	DbCollection = "users"
	tokenValidTime = time.Minute * time.Duration(1)
)

func CreateToken(userName string, password string, db *data.DB) (string, error) {
	if !checkPassword(userName, password, db) {
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
	// useful if you use multiple keys for your application.
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

func userNameIsRegistered(userName string, db *data.DB) bool {
	collection := db.Client.Database(data.Database).Collection(DbCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	res := collection.FindOne(ctx, bson.M{"name": userName})

	var user User
	if err := res.Decode(&user); err != nil {
		return false
	}

	return true
}

func checkPassword(userName string, password string, db *data.DB) bool {
	collection := db.Client.Database(data.Database).Collection(DbCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	res := collection.FindOne(ctx, bson.M{"name": userName})

	var user User
	if err := res.Decode(&user); err != nil {
		log.Debug(err)
		return false
	}
	valid := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if valid != nil {
		return false
	}
	return true
}

func RegisterAdmin(db *data.DB) (created bool, err error) {
	if !userNameIsRegistered("admin", db) {
		collection := db.Client.Database(data.Database).Collection(DbCollection)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		adminUser := User{"admin", "$2a$10$Is68fj8SktejE4mPz8AII.xvU.kTw2.7JgAfrvMmRrD4.Lku1Xngq"}

		_, err := collection.InsertOne(ctx, adminUser)

		return true, err
	}
	return false, nil

}

func AddUserIfNotThere(userName string, password string, db *data.DB) (created bool, err error) {
	if !userNameIsRegistered(userName, db) {
		collection := db.Client.Database(data.Database).Collection(DbCollection)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return false, err
		}

		user := User{userName, string(hashPassword)}

		_, err = collection.InsertOne(ctx, user)

		return true, err
	}
	return false, nil
}
