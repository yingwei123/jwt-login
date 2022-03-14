package server

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type templateHandlerFactory struct {
	tmpl          *template.Template
	MongoDBClient mongoDBClient
	Authenticator authenticator
}

func (rt Router) NewTemplateHandlerFactory(templateDirPath string) templateHandlerFactory {
	var tmpl *template.Template
	tmpl = template.Must(tmpl.ParseGlob(filepath.Join(templateDirPath, "*.gohtml")))
	template.Must(tmpl.ParseGlob(filepath.Join(templateDirPath, "partials/*.gohtml")))

	return templateHandlerFactory{tmpl, rt.MongoDBClient, rt.Authenticator}
}

type templateData struct {
	PageName string
}

func (t templateHandlerFactory) Handler(templateName string, pageName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}

		userID, err := t.Authenticator.ValidateRequest(r)
		if err != nil && err.Error() != "no access token" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if pageName == "Login Page" {
			if userID != "" {
				http.Redirect(w, r, "/default", http.StatusSeeOther)
				return
			}
		}

		if pageName == "Signup Page" {
			if userID != "" {
				http.Redirect(w, r, "/default", http.StatusSeeOther)
				return
			}
		}

		data := &templateData{
			PageName: pageName,
		}

		err = t.tmpl.ExecuteTemplate(w, templateName, data)
		if err != nil {
			println(err.Error())
			http.NotFound(w, r)
			return
		}

	}
}

func (t templateHandlerFactory) DefaultHandler(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}

		_, err := t.Authenticator.ValidateRequest(r)
		if err != nil {
			println(err.Error())
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		pageName := "Default Page"

		data := &templateData{
			PageName: pageName,
		}

		err = t.tmpl.ExecuteTemplate(w, templateName, data)
		if err != nil {
			println(err.Error())
			http.NotFound(w, r)
			return
		}

	}
}
