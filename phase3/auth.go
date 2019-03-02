package main

import "github.com/dgrijalva/jwt-go"

// TODO maybe this should be a real key
var hmacKey []byte

func init() {
    hmacKey = []byte{}
}

func keyFunction(token *jwt.Token) (interface{}, error) {
    return hmacKey, nil
}

func GenerateToken(username string, password string) string {
    Info.Printf("Generating JWT token for (%v,%v)\n", username, password)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username,
        "password": password,
        "authorized": UserIsAuthorized(username, password),
    })

    tokenString, err := token.SignedString(hmacKey)
    Info.Println(tokenString)

    if err != nil {
        return ""
    }
    return tokenString
}

func ParseToken(tokenString string) (string, string, bool) {
    token, err := jwt.Parse(tokenString, keyFunction)
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        Info.Printf("Error parsing JWT token: %v\n", err)
    }
    return claims["username"].(string), claims["password"].(string), claims["authorized"].(bool)
}

