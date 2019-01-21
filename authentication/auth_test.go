package authentication

import (
	. "github.com/alex1ai/ugr-master-cc/data"
	"github.com/mongodb/mongo-go-driver/bson"
	"os"
	"strconv"
	"testing"
)

var db *DB

func getEnv(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		value = def
	}
	return value
}

func setupDB(t *testing.T) {
	if db == nil {
		Database = "testing"
		Collection = DbCollection

		data := DB{}

		mPort := getEnv("MONGO_PORT", "27017")
		mIP := getEnv("MONGO_IP", "localhost")

		portI, err := strconv.Atoi(mPort)
		err = data.Connect(mIP, portI)
		if err != nil {
			t.Error(err)
		}
		if err != nil {
			t.Fatal(err)
		}
		db = &data
	}
}

func TestCreateAdmin(t *testing.T) {
	setupDB(t)
	db.Reset()
	registered := userNameIsRegistered("admin", db)
	if registered {
		t.Error("Admin is already registered")
	}
	created, err := RegisterAdmin(db)
	if !created || err != nil {
		t.Errorf("Admin is already registered or error occured %s", err.Error())
	}
	registered = userNameIsRegistered("admin", db)
	if !registered {
		t.Error("Admin was not saved in DB")
	}

}

func TestAddUser(t *testing.T) {
	created, err := AddUserIfNotThere("test", "test123", db)
	if err != nil {
		t.Error(err)
	}
	if !created {
		t.Error("User should not be in DB already")
	}
	_, err = db.Delete(bson.M{"name": "test"})
	if err != nil {
		t.Error("Could not delete User after test")
	}
}

func TestCheckPassword(t *testing.T) {
	_, err := AddUserIfNotThere("test", "test123", db)
	if err != nil {
		t.Error(err)
	}
	correct := checkPassword("test", "test123", db)
	if !correct {
		t.Error("Passwords should have matched but didn't")
	}

	correct = checkPassword("test", "test12", db)
	if correct {
		t.Error("Passwords should NOT have matched but did")
	}
}

func TestTokenCreation(t *testing.T) {

	_, err := CreateToken("test", "test123", db)
	if err != nil {
		t.Error(err)
	}

}

func TestValidateToken(t *testing.T) {
	tokenString, err := CreateToken("test", "test123", db)
	if err != nil {
		t.Error(err)
	}
	claims, ok := ValidateToken(tokenString)
	if !ok {
		t.Error()
	}
	if claims["user"] != "test" {
		t.Error("Could not read correct user test from token")
	}
}
