package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dchiquit/ceclia-ctf/phase3/assets"
)

// A single challenge in the CTF
type Challenge struct {
	Name     string
	Flag     string
	Hint     string
	Solved   bool
	HintUsed bool
}

// A user in the app. Only ceclia and d4ni3l are planned for.
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
	// Load the default values for the challenge and user JSON files
	rawChallengesTemplate := assets.MustAsset("json/challenges.json")
	rawUsersTemplate := assets.MustAsset("json/users.json")
	if err := json.Unmarshal(rawChallengesTemplate, &challengesTemplate); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(rawUsersTemplate, &usersTemplate); err != nil {
		panic(err)
	}

	// The default users.json file does not contain challenge data, so we populate it here
	for i, _ := range usersTemplate {
		usersTemplate[i].Progress = make([]Challenge, len(challengesTemplate))
		for j, _ := range usersTemplate[i].Progress {
			usersTemplate[i].Progress[j] = challengesTemplate[j]
			if usersTemplate[i].Admin {
				usersTemplate[i].Progress[j].Solved = true
			}
		}
	}

	// Load from the saved state directory, or create and initialize it if it doesn't exist
	os.MkdirAll("ceclia-ctf-data", os.ModePerm)
	if _, err := os.Stat("ceclia-ctf-data/users.json"); err == nil {
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

// Loads the users.json file into the users object
func LoadUsers() {
	Info.Printf("Loading existing users.json\n")
	rawUsers, err := ioutil.ReadFile("ceclia-ctf-data/users.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(rawUsers, &users); err != nil {
		panic(err)
	}
	Info.Printf("Loaded users.json\n")
}

// Saves the current state of users to the users.json file
func SaveUsers() {
	Info.Printf("Saving users.json\n")
	rawUsers, err := json.MarshalIndent(users, "", "    ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("ceclia-ctf-data/users.json", rawUsers, 0644)
	if err != nil {
		panic(err)
	}
	Info.Printf("users.json saved\n")
}

// Lists all users in the users object. Should only return [ceclia, d4ni3l]
func ListUsers() []string {
	userlist := make([]string, len(users))
	for i, user := range users {
		userlist[i] = user.Username
	}
	return userlist
}

// Looks up the User object for a given username
func findUser(username string) *User {
	for _, user := range users {
		if user.Username == username {
			return &user
		}
	}
	return nil
}

// Tests if a user is an admin. Only d4ni3l should be an admin
func UserIsAdmin(username string) bool {
	user := findUser(username)
	return user != nil && user.Admin
}

// Tests a user's login credentials
func UserIsAuthorized(username string, password string) bool {
	user := findUser(username)
	return username != "" && password != "" && user != nil && user.Password == password
}

// Gets the challenge progress for a user
func ProgressForUser(username string) []Challenge {
	user := findUser(username)
	if user == nil {
		return make([]Challenge, 0)
	}
	return user.Progress
}

// Sets HintUsed to true for the given user and challenge and updates users.json.
// The user needs to be authorized and HintUsed must currently be false.
func RequestHint(username string, password string, index int) error {
	Info.Printf("Requesting hint for challenge #%v for (%v,%v)\n", index, username, password)
	if !UserIsAuthorized(username, password) {
		return errors.New(fmt.Sprintf("%v has the wrong password! Write access is not permitted", username))
	}
	user := findUser(username)
	if user.Progress[index].HintUsed {
		return errors.New(fmt.Sprintf("Hint already requested for challenge #%v", index))
	}
	user.Progress[index].HintUsed = true
	SaveUsers()
	return nil
}

// Attempts to submit a flag.
// All challenges will be searched to see if any match the given flag.
// The user needs to be authorized and Solved must currently be false.
// If a valid challenge is found, Solved will be set to true and the users.json file will be updated.
func Submit(username string, password string, flag string) error {
	Info.Printf("Attempting submission of %v for %v %v\n", flag, username, password)
	if !UserIsAuthorized(username, password) {
		return errors.New(fmt.Sprintf("%v has the wrong password! Write access is not permitted", username))
	}
	user := findUser(username)
	for index, challenge := range user.Progress {
		if challenge.Flag == flag {
			if challenge.Solved {
				return errors.New(fmt.Sprintf("Flag has already been submitted for challenge #%v", index))
			}
			Info.Printf("Flag %v accepted!\n", flag)
			user.Progress[index].Solved = true
			SaveUsers()
			return nil
		}
	}
	Info.Printf("flag %v was incorrect\n", flag)
	return errors.New(fmt.Sprintf("Flag %v was incorrect", flag))
}

// Resets the users.json file.
// The given credentials must be to an authorized admin user.
func ResetUsers(username string, password string) error {
	if !UserIsAdmin(username) {
		return errors.New(fmt.Sprintf("%v is not an admin! How rude!", username))
	}
	if !UserIsAuthorized(username, password) {
		return errors.New(fmt.Sprintf("Wrong password for %v! How dare you hack my site!", username))
	}
	Info.Printf("Resetting users on behalf of %v\n", username)
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
