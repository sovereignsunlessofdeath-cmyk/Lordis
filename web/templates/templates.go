package templates

import (
	"fmt"
	"html/template"
	"net/http"
)

var Templ *template.Template

// InitTemplates loads all HTML templates from the web/templates directory
func InitTemplates() error {
	var err error
	// Parses all .html files inside the web/templates folder
	Templ, err = template.ParseGlob("web/templates/*.html")
	if err != nil {
		return fmt.Errorf("failed to parse templates: %v", err)
	}
	return nil
}

// Render helper to cleanly execute templates and handle errors
func Render(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := Templ.ExecuteTemplate(w, name, data)
	if err != nil {
		// Fallback error if a template fails to render
		http.Error(w, fmt.Sprintf("Template Render Error: %v", err), http.StatusInternalServerError)
	}
}
