package main

import (
	"html/template"
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

func main() {
	var err error

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

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
