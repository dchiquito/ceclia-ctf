package main

import "github.com/dgrijalva/jwt-go"

var hmacKey []byte

func init() {
	// TODO maybe this should be a real key
	hmacKey = []byte{}
}

// Simple key generator function that wraps hmacKey
func keyFunction(token *jwt.Token) (interface{}, error) {
	return hmacKey, nil
}

// Generates an auth token for a given username and password.
// It will automatically determine whether or not the user is authorized.
func GenerateToken(username string, password string) string {
	Info.Printf("Generating JWT token for (%v,%v)\n", username, password)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":   username,
		"password":   password,
		"authorized": UserIsAuthorized(username, password),
	})

	tokenString, err := token.SignedString(hmacKey)
	Info.Println(tokenString)

	if err != nil {
		return ""
	}
	return tokenString
}

// Parses a token into the username, password, and authorized value
func ParseToken(tokenString string) (string, string, bool) {
	Info.Printf("Parsing JWT token %v", tokenString)
	token, err := jwt.Parse(tokenString, keyFunction)
	if err != nil {
		Info.Printf("Error parsing JWT token: %v\n", err)
		return "", "", false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		Info.Printf("Error mapping claims from the JWT token\n")
		return "", "", false
	}
	username := claims["username"].(string)
	password := claims["password"].(string)
	authorized := claims["authorized"].(bool)
	Info.Printf("username=%v password=%v authorized=%v\n", username, password, authorized)
	return username, password, authorized
}
