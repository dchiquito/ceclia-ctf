package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dchiquit/ceclia-ctf/phase3/assets"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
)

var (
	Info  *log.Logger
	Error *log.Logger

	loginTpl  *template.Template
	appTpl    *template.Template
	robotsTpl *template.Template
)

func init() {
	Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime)

	loginTpl = loadHtmlTemplate("login.html")
	appTpl = loadHtmlTemplate("app.html")
	robotsTpl = loadHtmlTemplate("robots.txt") // not an HTML file, but whatever
}

func initLoggers(r *http.Request) {
	remote := r.RemoteAddr
	Info = log.New(os.Stdout, "[INFO]["+remote+"] ", log.Ldate|log.Ltime)
	Error = log.New(os.Stderr, "[ERROR]["+remote+"] ", log.Ldate|log.Ltime)
}

// Loads an HTML template.
// [[ and ]] are used as delimiters instead of the default {{ and }} to avoid conflicts with vue.js.
func loadHtmlTemplate(fileName string) *template.Template {
	templateHtml := string(assets.MustAsset("templates/" + fileName))
	return template.Must(template.New(fileName).Delims("[[", "]]").Parse(templateHtml))
}

// Render a template, or server error.
func render(w http.ResponseWriter, r *http.Request, tpl *template.Template, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, data); err != nil {
		Error.Printf("\nRender Error: %v\n", err)
		return
	}
	w.Write(buf.Bytes())
}

// Renders the login message with the given error message.
// If message is "", no message will be rendered.
func renderLoginPage(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		Info.Printf("Rendering login page\n")
	} else {
		Info.Printf("Rendering login page with error message \"%v\"\n", message)
	}
	fullData := map[string]interface{}{
		"Message": message,
		"IsError": true,
	}
	render(w, r, loginTpl, fullData)
}

// Renders the given app page with the given message.
// If isError, the message will be rendered as an error. Otherwise, it is rendered as a success message
// There is a flag hidden on undefined pages
func renderAppPage(w http.ResponseWriter, r *http.Request, page string, username string, message string, isError bool) {
	if message == "" {
		Info.Printf("Rendering app page %v\n", page)
	} else if isError {
		Info.Printf("Rendering app page %v with error message \"%v\"\n", page, message)
	} else {
		Info.Printf("Rendering app page %v with info message \"%v\"\n", page, message)
	}
	admin := UserIsAdmin(username)
	progress := ProgressForUser(username)
	usernames := ListUsers()
	leaderboard := make(map[string]string)
	for _, username := range usernames {
		solved := 0
		userProgress := ProgressForUser(username)
		for _, p := range userProgress {
			if p.Solved {
				solved += 1
			}
		}
		percent := 100 * float64(solved) / float64(len(userProgress))
		leaderboard[username] = fmt.Sprintf("%.1f%%", percent)
	}
	fullData := map[string]interface{}{
		"Username":    username,
		"Admin":       admin,
		"Page":        page,
		"Progress":    progress,
		"Leaderboard": leaderboard,
		"Message":     message,
		"IsError":     isError,
	}
	render(w, r, appTpl, fullData)
}

// GET,POST /login
// If the auth token is valid and specifies authorized==true, no further validation is done.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	initLoggers(r)
	Info.Printf("\n")
	Info.Printf("/login\n")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if cookie, err := r.Cookie("auth"); err == nil {
		username, _, authorized := ParseToken(cookie.Value)
		if authorized {
			Info.Printf("Username %v is already authorized!\n", username)
			http.Redirect(w, r, "/app", 302)
			return
		}
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" && password == "" {
		renderLoginPage(w, r, "")
		return
	}

	Info.Printf("Attempting login for (%v,%v)\n", username, password)

	if username == "" {
		renderLoginPage(w, r, "Please specify a username")
		return
	}
	if password == "" {
		renderLoginPage(w, r, "Please specify a password")
		return
	}
	// This is the admin password from phase 4
	if password == "CTF{d4mn_u_sm4rt_gurl}" {
		renderLoginPage(w, r, "LOL my real password isn't actually a flag :/ I'm only using it for debugging")
		return
	}

	authToken := GenerateToken(username, password)
	authorized := UserIsAuthorized(username, password)
	cookie := &http.Cookie{
		Name:  "auth",
		Value: authToken,
	}
	http.SetCookie(w, cookie)
	if authorized {
		Info.Printf("Login successful!\n")
		http.Redirect(w, r, "/app", 302)
		return
	}
	renderLoginPage(w, r, "Login Failed!")
}

// GET /app
// Renders the app page specified by the page URL query parameter.
// There is a flag hidden when page is undefined.
func AppHandler(w http.ResponseWriter, r *http.Request) {
	initLoggers(r)
	Info.Printf("\n")
	Info.Printf("/app\n")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	cookie, err := r.Cookie("auth")
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	username, _, authorized := ParseToken(cookie.Value)
	if !authorized {
		http.Redirect(w, r, "/login", 302)
		return
	}

	pageValues := r.URL.Query()["page"]
	page := "progress"
	if len(pageValues) > 0 {
		page = pageValues[0]
	}
	renderAppPage(w, r, page, username, "", false)
}

// POST /app/hint
// Attempts to request a hint.
// This is a write operation, so an authorization check is performed.
func HintHandler(w http.ResponseWriter, r *http.Request) {
	initLoggers(r)
	Info.Printf("\n")
	Info.Printf("/app/hint\n")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	cookie, err := r.Cookie("auth")
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	username, password, authorized := ParseToken(cookie.Value)
	if !authorized {
		http.Redirect(w, r, "/login", 302)
		return
	}

	indexStr := r.FormValue("index")
	if indexStr == "" {
		renderAppPage(w, r, "progress", username, "No index specified!", true)
		return
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		Error.Printf("Incorrectly formatted index %v\n", indexStr)
		renderAppPage(w, r, "progress", username, "Incorrectly formatted index!", true)
		return
	}
	err = RequestHint(username, password, index)
	if err != nil {
		renderAppPage(w, r, "progress", username, err.Error(), true)
		return
	}
	renderAppPage(w, r, "progress", username, "Fine. Here ya go", false)
}

// POST /app/submit
// Attempts to submit a flag.
// This is a write operation, so an authorization check is performed.
func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	initLoggers(r)
	Info.Printf("\n")
	Info.Printf("/app/submit\n")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	cookie, err := r.Cookie("auth")
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	username, password, authorized := ParseToken(cookie.Value)
	if !authorized {
		http.Redirect(w, r, "/login", 302)
		return
	}

	flag := r.FormValue("flag")
	if flag == "" {
		renderAppPage(w, r, "progress", username, "No flag specified", true)
		return
	}
	err = Submit(username, password, flag)
	if err != nil {
		renderAppPage(w, r, "progress", username, err.Error(), true)
		return
	}
	renderAppPage(w, r, "progress", username, "You got it!", false)
}

// GET /app/reset
// Resets the users.json file, reverting all the progress of all users.
// This is an admin only write operation, so an authorization check is performed.
func ResetHandler(w http.ResponseWriter, r *http.Request) {
	initLoggers(r)
	Info.Printf("\n")
	Info.Printf("/app/reset\n")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	cookie, err := r.Cookie("auth")
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	username, password, authorized := ParseToken(cookie.Value)
	if !authorized {
		http.Redirect(w, r, "/login", 302)
		return
	}

	err = ResetUsers(username, password)
	if err != nil {
		renderAppPage(w, r, "admin", username, err.Error(), true)
		return
	}
	renderAppPage(w, r, "admin", username, "All user progress reset", false)
}

// GET /robots.txt
// a flag is concealed in robots.txt
func RobotsHandler(w http.ResponseWriter, r *http.Request) {
	initLoggers(r)
	Info.Printf("\n")
	Info.Printf("/robots.txt\n")
	render(w, r, robotsTpl, nil)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/login", LoginHandler).Methods("GET", "POST")
	router.HandleFunc("/app", AppHandler).Methods("GET")
	router.HandleFunc("/app/hint", HintHandler).Methods("POST")
	router.HandleFunc("/app/submit", SubmitHandler).Methods("POST")
	router.HandleFunc("/app/reset", ResetHandler).Methods("GET")
	router.HandleFunc("/robots.txt", RobotsHandler).Methods("GET")
	router.PathPrefix("/static/").Handler(http.FileServer(&assetfs.AssetFS{Asset: assets.Asset, AssetDir: assets.AssetDir, AssetInfo: assets.AssetInfo, Prefix: ""}))
	Info.Printf("Up and running!\n")
	log.Fatal(http.ListenAndServe(":9596", router))
}
