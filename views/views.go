package views

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

const (
	AlertMsgGeneric = "Something went wrong. Please try again, and contact us if the problem persists."
)

var (
	LayoutDir   = "views/layouts/"
	TemplateDir = "views/templates/"
	TemplateExt = ".gohtml"
)

func NewView(layout string, files ...string) *View {
	files = addTemplatePath(files)
	files = addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.New("").ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

// Render is used to render the view with predefined layout
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch data.(type) {
	case Data:
		vd = data.(Data)
	default:
		vd = Data{
			Yield: data,
		}
	}
	var buf bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		log.Println(err)
		http.Error(w, AlertMsgGeneric, http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

// layoutFiles return a slice of strings representing
// the layout files used in our application
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// addTemplatePath takes in a slice of strings
// representing file paths for templates, and it prepends
// the TemplateDir directory to each string in the slice
func addTemplatePath(files []string) []string {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
	return files
}

// addTemplateExt takes in a slice of strings
// representing file paths for templates, and it appends
// the TemplateExt extension to each string in the slice
func addTemplateExt(files []string) []string {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
	return files
}
