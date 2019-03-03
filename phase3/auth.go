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
		Error.Printf("Error generating token: %v\n", err.Error())
		return ""
	}
	return tokenString
}

// Parses a token into the username, password, and authorized value
func ParseToken(tokenString string) (string, string, bool) {
	Info.Printf("Parsing JWT token %v\n", tokenString)
	token, err := jwt.Parse(tokenString, keyFunction)
	if err != nil {
		Error.Printf("Error parsing JWT token: %v\n", err)
		return "", "", false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		Error.Printf("Error mapping claims from the JWT token\n")
		return "", "", false
	}
	Info.Printf("Token claims: %v\n", claims)
	username := ""
	if claim, ok := claims["username"].(string); ok {
		username = claim
	}
	password := ""
	if claim, ok := claims["password"].(string); ok {
		password = claim
	}
	authorized := false
	if claim, ok := claims["authorized"].(bool); ok {
		authorized = claim
	}
	Info.Printf("username=%v password=%v authorized=%v\n", username, password, authorized)
	return username, password, authorized
}
