package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

var templates *template.Template

// InitTemplates loads all HTML templates from the web/templates directory
func InitTemplates() error {
	var err error
	templates, err = template.ParseGlob("web/templates/*.html")
	if err != nil {
		return fmt.Errorf("failed to parse templates: %v", err)
	}
	return nil
}

// Render helper to execute templates cleanly inside handlers
func Render(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if templates == nil {
		http.Error(w, "Template Render Error: templates not initialized", http.StatusInternalServerError)
		return
	}
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template Render Error: %v", err), http.StatusInternalServerError)
	}
}
