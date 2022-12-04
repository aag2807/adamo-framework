package render

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
)

type Render struct {
	Renderer   string
	RootPath   string
	Secure     bool
	Port       string
	ServerName string
	JetViews   *jet.Set
}

type TemplateData[T any] struct {
	IsAuthenticated bool
	Secure          bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]T
	CSRFToken       string
	Port            string
	ServerName      string
}

func (ren *Render) Page(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {
	engineType := strings.ToLower(ren.Renderer)

	switch engineType {
	case "go":
		return ren.GoPage(w, r, view, data)
	case "jet":
		return ren.JetPage(w, r, view, variables, data)
	default:
		return errors.New(" no engine selected")
	}
}

// GoPage renders a regular html/template page
func (ren *Render) GoPage(w http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	filePath := fmt.Sprintf("%s/views/%s.go.html", ren.RootPath, view)
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		return err
	}

	tmplData := &TemplateData[interface{}]{}
	if data != nil {
		tmplData = data.(*TemplateData[interface{}])
	}

	err = tmpl.Execute(w, &tmplData)
	if err != nil {
		return err
	}

	return nil
}

// JetPage renders a jet page template
func (ren *Render) JetPage(w http.ResponseWriter, r *http.Request, templateName string, variables, data interface{}) error {
	var vars jet.VarMap

	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	tmplData := &TemplateData[interface{}]{}
	if data != nil {
		tmplData = data.(*TemplateData[interface{}])
	}

	tmplName := fmt.Sprintf("%s.jet", templateName)
	t, err := ren.JetViews.GetTemplate(tmplName)
	if err != nil {
		log.Println(err)
		return err
	}

	if err = t.Execute(w, vars, tmplData); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
