package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type ExampleData struct {
	Title string
	Items []string
}

type PageData struct {
	DataStructs []ExampleData
}

var (
	exampleDataArr []ExampleData = []ExampleData{
		{
			Title: "Struct One",
			Items: []string{
				"Item One",
				"Item Two",
				"Item Three",
			},
		},
		{
			Title: "Struct Two",
			Items: []string{
				"Item Alpha",
				"Item Beta",
				"Item Gamma",
			},
		},
	}
)

var (
	//go:embed static/css/output.css
	embedCSSFile []byte

	//go:embed static/templates/*.html
	templatesFS embed.FS
)

func main() {
	var err error

	templatesFS, err := fs.Sub(templatesFS, "static/templates")
	if err != nil {
		log.Fatalf("error during embedded file system: %v", err)
	}

	http.HandleFunc("/static/templates/index", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("static/templates/index.html")
		if err != nil {
			log.Fatalf("error parsing template: %v", err)
		}
		err = tmpl.Execute(w, PageData{
			DataStructs: exampleDataArr,
		})
		if err != nil {
			log.Fatalf("error during template execute: %v", err)
		}
	})

	log.Printf("Serving template at http://localhost:8080/static/templates/index")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("error during http serving: %v", err)
	}
}
