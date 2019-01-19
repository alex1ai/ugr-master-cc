package authentication

import (
	"testing"
)

func TestTokenCreation(t *testing.T) {

	_, err := CreateToken("test", "test123")
	if err != nil {
		t.Error(err)
	}

}

func TestValidateToken(t *testing.T) {
	tokenString, err := CreateToken("test", "test123")
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
