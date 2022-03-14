package server

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Router struct {
	ResourcesPath string
	ServerURL     string
	MongoDBClient mongoDBClient
	Authenticator authenticator
}

type authenticator interface {
	GenerateValidToken(userName string) (string, time.Time, error)
	ValidateRequest(r *http.Request) (string, error)
}

type Response struct {
	Status  int
	Message string
}

type mongoDBClient interface {
	CreateItem(item Item) (string, error)
	CreateNewUser(user NewUser) (string, error)
	FindUserByEmail(email string) (NewUser, error)
}

func (rt Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	template := rt.NewTemplateHandlerFactory(filepath.Join(rt.ResourcesPath, "templates"))

	router.Handle("/default", template.DefaultHandler("default.gohtml"))
	router.Handle("/mongo-db", rt.AddItem())
	router.Handle("/login", template.Handler("login.gohtml", "Login Page"))
	router.Handle("/signup", template.Handler("signup.gohtml", "Signup Page"))
	router.Handle("/new-user", rt.SignUpNewUser())
	router.Handle("/auth", rt.LoginHandler())

	n := negroni.New()
	n.Use(negroni.NewLogger())

	faviconMiddleware := negroni.NewStatic(http.Dir(filepath.Join(rt.ResourcesPath, "favicon.ico")))
	faviconMiddleware.Prefix = "/favicon.ico"
	n.Use(faviconMiddleware)

	publicMiddleware := negroni.NewStatic(http.Dir(filepath.Join(rt.ResourcesPath, "public")))
	publicMiddleware.Prefix = "/public"

	n.Use(publicMiddleware)

	n.UseHandler(router)
	n.ServeHTTP(w, r)
}
