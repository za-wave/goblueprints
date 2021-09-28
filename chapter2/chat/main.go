package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/za-wave/goblueprints/chapter1/trace"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

func main() {
	// err := godotenv.Load("../../.env")
	// // err := godotenv.Load(fmt.Sprintf("../%s.env", os.Getenv("GO_ENV")))

	// //もし err がnilではないなら、"読み込み出来ませんでした"が出力されます。
	// if err != nil {
	// 	fmt.Printf("読み込み出来ませんでした: %v", err)
	// }
	// AUTH_SECURITY_KEY := os.Getenv("AUTH_SECURITY_KEY")
	// GOOGLE_CLIENT_ID := os.Getenv("GOOGLE_CLIENT_ID")
	// GOOGLE_CLIENT_SECRET := os.Getenv("GOOGLE_CLIENT_SECRET")
	// log.Println(AUTH_SECURITY_KEY)
	// fmt.Println(GOOGLE_CLIENT_SECRET)

	// 	AUTH_SECURITY_KEY=79Uhwc4GwG2fquNQmljpNqyswvPDbg6sms3F3WtsOZhCK0IVYryCcw7TfGtZYiJ8
	// GOOGLE_CLIENT_ID=594946738580-roi12cdupdvji6l168qp09u852mccjcp.apps.googleusercontent.com
	// GOOGLE_CLIENT_SECRET_KEY=JTlECG7BnT-OZOf7EG4FW4KH

	var host = flag.String("host", ":8080", "The host of the application.")

	flag.Parse()

	// setup for gomniauth
	// gomniauth.SetSecurityKey("OGLLaiLEn5ljt7bT9YgUqk85iDLQqWQnK0C6RILdFtxIfbBBdK9BumymjXbVjXJP")
	// gomniauth.SetSecurityKey(AUTH_SECURITY_KEY)
	// gomniauth.WithProviders(
	// 	// facebook.New("", "", "http://localhost:8080/auth/callback/facebook"),
	// 	// github.New("", "", "http://localhost:8080/auth/callback/github"),
	// 	google.New(GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, "http://localhost:8080/auth/callback/google"),
	// )

	gomniauth.SetSecurityKey("79Uhwc4GwG2fquNQmljpNqyswvPDbg6sms3F3WtsOZhCK0IVYryCcw7TfGtZYiJ8")
	gomniauth.WithProviders(
		// facebook.New("", "", "http://localhost:8080/auth/callback/facebook"),
		// github.New("", "", "http://localhost:8080/auth/callback/github"),
		google.New("594946738580-roi12cdupdvji6l168qp09u852mccjcp.apps.googleusercontent.com", "JTlECG7BnT-OZOf7EG4FW4KH", "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	// get the room going
	go r.run()

	// start the web server
	log.Println("Starting web server on", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
