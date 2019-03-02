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

func loadHtmlTemplate(fileName string) *template.Template {
	templateHtml := string(assets.MustAsset("templates/" + fileName))
	return template.Must(template.New(fileName).Delims("[[", "]]").Parse(templateHtml))
}

var (
	Info  *log.Logger
	Error *log.Logger

	loginTpl  *template.Template
	appTpl    *template.Template
	robotsTpl *template.Template
)

func init() {
	Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)

	loginTpl = loadHtmlTemplate("login.html")
	appTpl = loadHtmlTemplate("app.html")
	robotsTpl = loadHtmlTemplate("robots.txt")
}

// Render a template, or server error.
func render(w http.ResponseWriter, r *http.Request, tpl *template.Template, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, data); err != nil {
		Info.Printf("\nRender Error: %v\n", err)
		return
	}
	w.Write(buf.Bytes())
}

func renderLoginPage(w http.ResponseWriter, r *http.Request, message string) {
	fullData := map[string]interface{}{
		"Message": message,
		"IsError": true,
	}
	render(w, r, loginTpl, fullData)
}

func renderAppPage(w http.ResponseWriter, r *http.Request, page string, message string, isError bool) {
	cookie, err := r.Cookie("auth")
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	username, _, _ := ParseToken(cookie.Value)
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

// LoginHandler renders the login.html template
func LoginHandler(w http.ResponseWriter, r *http.Request) {
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

	Info.Printf("Attempting login for username %v password %v\n", username, password)

	if username == "" && password == "" {
		renderLoginPage(w, r, "")
		return
	}
	if username == "" {
		renderLoginPage(w, r, "Please specify a username")
		return
	}
	if password == "" {
		renderLoginPage(w, r, "Please specify a password")
		return
	}
	if password == "CTF{d4mn_u_sm4rt_gurl}" {
		renderLoginPage(w, r, "LOL my real password isn't actually a flag :/ I'm only using it for debugging")
		return
	}

	authToken := GenerateToken(username, password)
	username, password, authorized := ParseToken(authToken)
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

// AppHandler renders the app.html template
func AppHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	cookie, err := r.Cookie("auth")
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	_, _, authorized := ParseToken(cookie.Value)
	if !authorized {
		http.Redirect(w, r, "/login", 302)
		return
	}

	pageValues := r.URL.Query()["page"]
	page := "progress"
	if len(pageValues) > 0 {
		page = pageValues[0]
	}
	renderAppPage(w, r, page, "", false)
}

func HintHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	cookie, err := r.Cookie("auth")
	username, password, authorized := ParseToken(cookie.Value)
	if !authorized {
		http.Redirect(w, r, "/login", 302)
		return
	}

	indexStr := r.FormValue("index")
	if indexStr == "" {
		Info.Printf("No index specified\n")
		renderAppPage(w, r, "progress", "No index specified!", true)
		return
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		Info.Printf("Incorrectly formatted index %v\n", indexStr)
		renderAppPage(w, r, "progress", "Incorrectly formatted index!", true)
		return
	}
	err = RequestHint(username, password, index)
	if err != nil {
		renderAppPage(w, r, "progress", err.Error(), true)
		return
	}
	renderAppPage(w, r, "progress", "Fine. Here ya go", false)
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	cookie, err := r.Cookie("auth")
	username, password, authorized := ParseToken(cookie.Value)
	if !authorized {
		http.Redirect(w, r, "/login", 302)
		return
	}

	flag := r.FormValue("flag")
	if flag == "" {
		Info.Printf("No flag specified\n")
		renderAppPage(w, r, "progress", "No flag specified", true)
		return
	}
	err = Submit(username, password, flag)
	if err != nil {
		renderAppPage(w, r, "progress", err.Error(), true)
		return
	}
	renderAppPage(w, r, "progress", "You got it!", false)
}

func ResetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	cookie, err := r.Cookie("auth")
	username, password, authorized := ParseToken(cookie.Value)
	if !authorized {
		http.Redirect(w, r, "/login", 302)
		return
	}

	err = ResetUsers(username, password)
	if err != nil {
		renderAppPage(w, r, "admin", err.Error(), true)
		return
	}
	renderAppPage(w, r, "admin", "All user progress reset", false)
}

func RobotsHandler(w http.ResponseWriter, r *http.Request) {
	render(w, r, robotsTpl, nil)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/login", LoginHandler).Methods("GET", "POST")
	router.HandleFunc("/app", AppHandler).Methods("GET")
	router.HandleFunc("/app/hint", HintHandler).Methods("GET", "POST")
	router.HandleFunc("/app/submit", SubmitHandler).Methods("GET", "POST")
	router.HandleFunc("/app/reset", ResetHandler).Methods("GET")
	router.HandleFunc("/robots.txt", RobotsHandler).Methods("GET")
	router.PathPrefix("/static/").Handler(http.FileServer(&assetfs.AssetFS{Asset: assets.Asset, AssetDir: assets.AssetDir, AssetInfo: assets.AssetInfo, Prefix: ""}))
	Info.Printf("Up and running!\n")
	log.Fatal(http.ListenAndServe(":9596", router))
}
