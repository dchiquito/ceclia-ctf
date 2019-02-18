package main

import (
    "encoding/json"
    "io/ioutil"
    "os"

    "github.com/dchiquit/ceclia-ctf-go/assets"
)

type Challenge struct {
    Name string
    Flag string
    Hint string
    Solved bool
    HintUsed bool
}

type User struct {
    Username string
    Password string
    Admin bool
    Progress []Challenge
}

var (
    challengesTemplate []Challenge
    usersTemplate []User

    users []User
)

func init() {
    Info.Printf("ehehehehe initd")

    rawChallengesTemplate := assets.MustAsset("json/challenges.json")
    rawUsersTemplate := assets.MustAsset("json/users.json")

    if err := json.Unmarshal(rawChallengesTemplate, &challengesTemplate); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(rawUsersTemplate, &usersTemplate); err != nil {
        panic(err)
    }

    for i,_ := range usersTemplate {
        usersTemplate[i].Progress = make([]Challenge, len(challengesTemplate))
        for j,_ := range usersTemplate[i].Progress {
            usersTemplate[i].Progress[j] = challengesTemplate[j]
            if usersTemplate[i].Admin {
                usersTemplate[i].Progress[j].Solved = true
            }
        }
    }

    os.MkdirAll("ceclia-ctf", os.ModePerm)
    if _, err := os.Stat("ceclia-ctf/users.json"); err == nil {
        LoadUsers()
    } else if os.IsNotExist(err) {
        Info.Printf("users.json is missing, creating it\n")
        users = make([]User, len(usersTemplate))
        copy(users, usersTemplate)

        SaveUsers()
    }
}

func LoadUsers() {
        Info.Printf("Loading existing users.json\n")
        rawUsers, err := ioutil.ReadFile("ceclia-ctf/users.json")
        if err != nil {
            panic(err)
        }
        if err := json.Unmarshal(rawUsers, &users); err != nil {
            panic(err)
        }
        Info.Printf("Loaded users.json %v\n", users)
}

func SaveUsers() {
    Info.Printf("Saving users.json")
    rawUsers, err := json.MarshalIndent(users, "", "    ")
    if err != nil {
        panic(err)
    }
    err = ioutil.WriteFile("ceclia-ctf/users.json", rawUsers, 0644)
    if err != nil {
        panic(err)
    }
    Info.Printf("users.json saved %v\n", users)
}

func ListUsers() []string {
    userlist := make([]string, len(users))
    for i,user := range users {
        userlist[i] = user.Username
    }
    return userlist
}

func FindUser(username string) *User {
    for _,user := range users {
        if user.Username == username {
            return &user
        }
    }
    return nil
}

func UserExists(username string) bool {
    return FindUser(username) != nil
}

func UserIsAdmin(username string) bool {
    user := FindUser(username)
    return user != nil && user.Admin
}

func UserIsAuthorized(username string, password string) bool {
    user := FindUser(username)
    return username != "" && password != "" && user != nil && user.Password == password
}

func ProgressForUser(username string) []Challenge {
    user := FindUser(username)
    if user == nil {
        return make([]Challenge, 0)
    }
    return user.Progress
}

func RequestHint(username string, password string, index int) error {
    Info.Printf("Requesting hint for %v %v %v\n", username, password, index)
    return nil
}

func Submit(username string, password string, flag string) error {
    Info.Printf("Attempting submission for %v %v %v\n", username, password, flag)
    return nil
}
