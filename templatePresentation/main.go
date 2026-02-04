package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	templates, err := loadTemplates()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/slides/", slideHandler(templates))

	log.Println("serving at http://localhost:8080")
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func slideHandler(templates *Templates) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slideStr := strings.TrimPrefix(r.URL.Path, "/slides/")
		slideNum, err := strconv.Atoi(slideStr)
		if err != nil || slideNum < 1 || slideNum > len(templates.Slides) {
			http.NotFound(w, r)
			return
		}

		data := struct {
			Slide int
			Next  int
			Prev  int
			Total int
		}{
			Slide: slideNum,
			Next:  slideNum + 1,
			Prev:  slideNum - 1,
			Total: len(templates.Slides),
		}

		tmpl := templates.Slides[slideNum-1]
		if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type Templates struct {
	Slides []*template.Template
}

func loadTemplates() (*Templates, error) {
	layout := "templates/layout.html"
	navigation := "templates/navbar.html"

	slideFiles := []string{
		"templates/titleSlide.html",
		"templates/whatIsTemplating.html",
		"templates/whyTemplating.html",
		"templates/templateExampleWithNoData.html",
		//"templates/addingDataIntoTemplating.html",
		//"templates/passingFunctionsThroughTemplating.html",
		//"templates/AConvenientWayToOrganizeTemplates.html",
	}

	slides := make([]*template.Template, 0, len(slideFiles))

	for _, slide := range slideFiles {
		t, err := template.ParseFiles(layout, navigation, slide)
		if err != nil {
			return nil, err
		}
		slides = append(slides, t)
	}

	return &Templates{Slides: slides}, nil
}
