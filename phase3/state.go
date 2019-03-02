package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dchiquit/ceclia-ctf/phase3/assets"
)

type Challenge struct {
	Name     string
	Flag     string
	Hint     string
	Solved   bool
	HintUsed bool
}

type User struct {
	Username string
	Password string
	Admin    bool
	Progress []Challenge
}

var (
	challengesTemplate []Challenge
	usersTemplate      []User

	users []User
)

func init() {
	rawChallengesTemplate := assets.MustAsset("json/challenges.json")
	rawUsersTemplate := assets.MustAsset("json/users.json")

	if err := json.Unmarshal(rawChallengesTemplate, &challengesTemplate); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(rawUsersTemplate, &usersTemplate); err != nil {
		panic(err)
	}
	Info.Printf("Initial values loaded from binary:\n")
	Info.Printf("Users: %v\n", usersTemplate)
	Info.Printf("Challenges: %v\n", challengesTemplate)

	for i, _ := range usersTemplate {
		usersTemplate[i].Progress = make([]Challenge, len(challengesTemplate))
		for j, _ := range usersTemplate[i].Progress {
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
		for i, user := range usersTemplate {
			// shallow copy, they will still share Progress
			users[i] = user
			// deep copy Progress
			users[i].Progress = make([]Challenge, len(user.Progress))
			copy(users[i].Progress, user.Progress)
		}
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
	for i, user := range users {
		userlist[i] = user.Username
	}
	return userlist
}

func FindUser(username string) *User {
	for _, user := range users {
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
	if !UserIsAuthorized(username, password) {
		return errors.New(fmt.Sprintf("%v is not authorized to request hints", username))
	}
	user := FindUser(username)
	if user.Progress[index].HintUsed {
		return errors.New(fmt.Sprintf("Hint already requested for challenge #%v", index))
	}
	user.Progress[index].HintUsed = true
	SaveUsers()
	return nil
}

func Submit(username string, password string, flag string) error {
	Info.Printf("Attempting submission for %v %v %v\n", username, password, flag)
	if !UserIsAuthorized(username, password) {
		Info.Printf("%v is not authorized to submit flags without a valid password", username)
		return errors.New(fmt.Sprintf("%v is not authorized to submit flags without a valid password", username))
	}
	user := FindUser(username)
	for index, challenge := range user.Progress {
		if challenge.Flag == flag {
			if challenge.Solved {
				Info.Printf("Flag has already been submitted for challenge #%v", index)
				return errors.New(fmt.Sprintf("Flag has already been submitted for challenge #%v", index))
			}
			Info.Printf("Flag %v accepted!", flag)
			user.Progress[index].Solved = true
			Info.Printf("challenge.Solved %v\n", challenge.Solved)
			Info.Printf("user.Progress[index].Solved %v\n", user.Progress[index].Solved)
			Info.Printf("users[0].Progress[index].Solved %v\n", users[0].Progress[index].Solved)
			SaveUsers()
			return nil
		}
	}
	Info.Printf("flag %v was incorrect\n", flag)
	return errors.New(fmt.Sprintf("Flag %v was incorrect", flag))
}

func ResetUsers(username string, password string) error {
	if !UserIsAdmin(username) {
		Info.Printf("%v is not an admin, not permitted to reset users", username)
		return errors.New(fmt.Sprintf("%v is not an admin, not permitted to reset users", username))
	}
	if !UserIsAuthorized(username, password) {
		Info.Printf("Wrong password for %v! How dare you hack my site!", username)
		return errors.New(fmt.Sprintf("Wrong password for %v! How dare you hack my site!", username))
	}
	Info.Printf("Resetting users")
	users = make([]User, len(usersTemplate))
	for i, user := range usersTemplate {
		// shallow copy, they will still share Progress
		users[i] = user
		// deep copy Progress
		users[i].Progress = make([]Challenge, len(user.Progress))
		copy(users[i].Progress, user.Progress)
	}
	SaveUsers()
	return nil
}
