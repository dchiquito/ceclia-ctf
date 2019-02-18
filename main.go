package main

import (
    "bytes"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"

    "github.com/dchiquit/ceclia-ctf-go/assets"
    "github.com/gorilla/mux"
)

func loadHtmlTemplate(fileName string) *template.Template {
    templateHtml := string(assets.MustAsset("templates/" + fileName))
    return template.Must(template.New(fileName).Delims("[[", "]]").Parse(templateHtml))
}

var (
    Info *log.Logger
    Error *log.Logger

    loginTpl  *template.Template
    appTpl    *template.Template
    robotsTpl *template.Template
)

func init() {
    Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
    Error = log.New(os.Stderr, "[ERROR] ",log.Ldate|log.Ltime|log.Lshortfile)

    loginTpl  = loadHtmlTemplate("login.html")
    appTpl    = loadHtmlTemplate("app.html")
    robotsTpl = loadHtmlTemplate("robots.txt")
}

// Render a template, or server error.
func render(w http.ResponseWriter, r *http.Request, tpl *template.Template, data interface{}) {
    buf := new(bytes.Buffer)
    if err := tpl.Execute(buf, data); err != nil {
        fmt.Printf("\nRender Error: %v\n", err)
        return
    }
    w.Write(buf.Bytes())
}

// LoginHandler renders the login.html template
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    if HasBeenAuthorized(r) || Authenticate(w,r) {
        http.Redirect(w, r, "/app?page=progress", 302)
        return
    }

    render(w, r, loginTpl, nil)
}

// AppHandler renders the app.html template
func AppHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    if !HasBeenAuthorized(r) {
        http.Redirect(w, r, "/login", 302)
        return
    }
    values := r.URL.Query()["page"]
    Info.Printf("All them values! %v %v\n", values, len(values))
    page := "progress"
    if len(values) > 0 {
        page = values[0]
    }
    username := GetUsernameFromCookie(r)
    admin := UserIsAdmin(username)
    progress := ProgressForUser(username)
    usernames := ListUsers()
    leaderboard := make(map[string]string)
    for _,username := range usernames {
        solved := 0
        userProgress := ProgressForUser(username)
        for _,p := range userProgress {
            if p.Solved {
                solved += 1
            }
        }
        percent := 100 * float64(solved) / float64(len(userProgress))
        leaderboard[username] = fmt.Sprintf("%.1f%%", percent)
        Info.Printf("HEEEEEY %v %v %v\n", leaderboard[username], solved, percent)
    }
    fullData := map[string]interface{}{
        "Username":    GetUsernameFromCookie(r),
        "Admin":       admin,
        "Page":        page,
        "Progress":    progress,
        "Leaderboard": leaderboard,
    }
    render(w, r, appTpl, fullData)
}

func HintHandler(w http.ResponseWriter, r *http.Request) {
    values := r.URL.Query()["index"]
    if len(values) == 0 {
        Info.Printf("No index specified\n")
        r.WriteHeader(400)
    }
    index, err := strconv.ParseInt(values[0])
    if err != nil {
        Info.Printf("Incorrectly formatted index %v\n", values[0])
        r.WriteHeader(400)
    }
    err := RequestHint(GetUsernameFromCookie(), GetPasswordFromCookie(), index)
    if err != nil {
        r.WriteHeader(400)
    }
    r.WriteHeader(200)
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
    values := r.URL.Query()["flag"]
    if len(values) == 0 {
        Info.Printf("No flag specified\n")
        r.WriteHeader(400)
    }
    flag := values[0]
    err := Submit(GetUsernameFromCookie(), GetPasswordFromCookie(), flag)
    if err != nil {
        r.WriteHeader(400)
    }
    r.WriteHeader(200)
}

func RobotsHandler(w http.ResponseWriter, r *http.Request) {
    render(w, r, robotsTpl, nil)
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/login", LoginHandler).Methods("GET")
    router.HandleFunc("/app", AppHandler).Methods("GET")
    router.HandleFunc("/app/hint", HintHandler).Methods("POST")
    router.HAndleFunc("/app/submit", SubmitHandler).Methods("POST")
    router.HandleFunc("/robots.txt", RobotsHandler).Methods("GET")
    router.PathPrefix("/app/js/").Handler(http.StripPrefix("/app/js/", http.FileServer(http.Dir("./js"))))
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
    Info.Printf("Up and running!\n")
    log.Fatal(http.ListenAndServe(":9596", router))
}
