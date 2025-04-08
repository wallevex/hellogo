package hello

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

var JwtSecret = []byte("my-secret-key")

type MyCustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func GenToken(name string) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "hellogo",
			Subject:   "0x1",
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			NotBefore: nil,
			IssuedAt:  nil,
			ID:        "JWT ID",
		},
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(JwtSecret)
	return tokenString, err
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return JwtSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	return token, err
}

func TestJWT(t *testing.T) {
	tokenString, err := GenToken("alice")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tokenString)

	token, err := ParseToken(tokenString)
	if err != nil {
		t.Fatal(err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["name"], claims["exp"])
	}
}
