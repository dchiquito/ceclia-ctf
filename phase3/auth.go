package main

import (
    "net/http"
    "github.com/dgrijalva/jwt-go"
)

// TODO maybe this should be a real key
var hmacKey []byte

func init() {
    hmacKey = []byte{}
}

func keyFunction(token *jwt.Token) (interface{}, error) {
    return hmacKey, nil
}

func GenerateToken(username string, password string) string {
    Info.Printf("Generating JWT token for %v %v", username, password)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username,
        "password": password,
        "authorized": UserIsAuthorized(username, password),
    })

    tokenString, err := token.SignedString(hmacKey)
    Info.Println(tokenString, err)

    if err != nil {
        return ""
    }
    return tokenString
}

func ParseToken(tokenString string) (string, string, bool) {
    token, err := jwt.Parse(tokenString, keyFunction)
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        Info.Println("Error parsing JWT token: %v", err)
    }
    return claims["username"].(string), claims["password"].(string), claims["authorized"].(bool)
}

func GetUsernameFromCookie(r *http.Request) string {
    cookie, err := r.Cookie("username")
    if err != nil {
        return ""
    }
    return cookie.Value
}

func GetPasswordFromCookie(r *http.Request) string {
    cookie, err := r.Cookie("password")
    if err != nil {
        return ""
    }
    return cookie.Value
}

func GetAuthFromCookie(r *http.Request) string {
    cookie, err := r.Cookie("auth")
    if err != nil {
        return ""
    }
    return cookie.Value
}

// Tests if the given request has the "auth" cookie set 
func HasBeenAuthorized(r *http.Request) bool {
    auth := GetAuthFromCookie(r)
    if auth != "true" {
        Info.Printf("auth cookie not present, not authorized\n")
        return false
    }
    Info.Printf("Authorized using auth cookie\n")
    return true
}

// Tests if the given request has a valid username/password
func IsAuthorized(r *http.Request) bool {
    username := GetUsernameFromCookie(r)
    password := GetPasswordFromCookie(r)
    if !UserIsAuthorized(username, password) {
        Info.Printf("Incorrect password, not authorized\n")
        return false
    }
    Info.Printf("Authorized username %v password %v\n", username, password)
    return true
}

// Tests if the given request has a valid username/password and sets the authorized cookie appropriately
func Authenticate(w http.ResponseWriter, r *http.Request) bool {
    if IsAuthorized(r) {
        http.SetCookie(w, &http.Cookie{
            Name:  "auth",
            Value: "true",
        })
        return true
    } else {
        http.SetCookie(w, &http.Cookie{
            Name:  "auth",
            Value: "false",
        })
        return false
    }
}


