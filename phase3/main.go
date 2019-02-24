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

func renderAppPage(w http.ResponseWriter, r *http.Request, page string, message string) {
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
    }
    fullData := map[string]interface{}{
        "Username":    GetUsernameFromCookie(r),
        "Admin":       admin,
        "Page":        page,
        "Progress":    progress,
        "Leaderboard": leaderboard,
        "Message":     message,
    }
    render(w, r, appTpl, fullData)
}

// LoginHandler renders the login.html template
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    if HasBeenAuthorized(r) || Authenticate(w,r) {
        http.Redirect(w, r, "/app", 302)
        return
    }

    render(w, r, loginTpl, nil)
}

// AppHandler renders the app.html template
func AppHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    pageValues := r.URL.Query()["page"]
    page := "progress"
    if len(pageValues) > 0 {
        page = pageValues[0]
    }
    if !HasBeenAuthorized(r) {
        http.Redirect(w, r, "/login", 302)
        return
    }
    renderAppPage(w, r, page, "")
}

func HintHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    indexValues := r.URL.Query()["index"]
    if len(indexValues) == 0 {
        Info.Printf("No index specified\n")
        renderAppPage(w, r, "progress", "No index specified!")
        return
    }
    index, err := strconv.Atoi(indexValues[0])
    if err != nil {
        Info.Printf("Incorrectly formatted index %v\n", indexValues[0])
        renderAppPage(w, r, "progress", "Incorrectly formatted index!")
        return
    }
    err = RequestHint(GetUsernameFromCookie(r), GetPasswordFromCookie(r), index)
    if err != nil {
        renderAppPage(w, r, "progress", err.Error())
        return
    }
    renderAppPage(w, r, "progress", "Fine. Here ya go")
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    flagValues := r.URL.Query()["flag"]
    if len(flagValues) == 0 {
        Info.Printf("No flag specified\n")
        renderAppPage(w, r, "progress", "No flag specified")
        return
    }
    flag := flagValues[0]
    err := Submit(GetUsernameFromCookie(r), GetPasswordFromCookie(r), flag)
    if err != nil {
        renderAppPage(w, r, "progress", err.Error())
        return
    }
    renderAppPage(w, r, "progress", "You got it!")
}

func ResetHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    err := ResetUsers(GetUsernameFromCookie(r), GetPasswordFromCookie(r))
    if err != nil {
        renderAppPage(w, r, "admin", err.Error())
        return
    }
    renderAppPage(w, r, "admin", "All user progress reset")
}

func RobotsHandler(w http.ResponseWriter, r *http.Request) {
    render(w, r, robotsTpl, nil)
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/login", LoginHandler).Methods("GET")
    router.HandleFunc("/app", AppHandler).Methods("GET")
    router.HandleFunc("/app/hint", HintHandler).Methods("GET")
    router.HandleFunc("/app/submit", SubmitHandler).Methods("GET")
    router.HandleFunc("/app/reset", ResetHandler).Methods("GET")
    router.HandleFunc("/robots.txt", RobotsHandler).Methods("GET")
    router.PathPrefix("/static/").Handler(http.FileServer(&assetfs.AssetFS{Asset: assets.Asset, AssetDir: assets.AssetDir, AssetInfo: assets.AssetInfo, Prefix: ""}))
    Info.Printf("Up and running!\n")
    log.Fatal(http.ListenAndServe(":9596", router))
}
