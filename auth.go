package main

import (
    "net/http"
)

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


