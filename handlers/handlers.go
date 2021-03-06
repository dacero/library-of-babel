package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/dacero/labyrinth-of-babel/repository"
	"github.com/dacero/labyrinth-of-babel/models"
)

func ViewHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cellId := r.URL.Path[len("/view/"):]
		cell, err := lob.GetCell(cellId)
		if err != nil {
			log.Printf("Error when returning card: %s", err)
			notFound, err := ioutil.ReadFile("./templates/card_not_found.html")
			if err != nil {
				log.Printf("Error when returning card: %s", err)
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, string(notFound))
		} else {
			t, err := template.ParseFiles("./templates/card.gohtml")
			if err != nil {
				log.Printf("Error when returning card: %s", err)
			}
			err = t.Execute(w, cell)
			if err != nil {
				log.Printf("Error when returning card: %s", err)
			}
		}
	})
}

func CreateHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//parse the form and create the cell
		log.Printf("New cell title: %s", r.PostFormValue("title"))
		newCell := models.Cell{Title: r.PostFormValue("title"),
			Body: r.PostFormValue("body"),
			Room: r.PostFormValue("room"),
			Sources: []models.Source{ models.Source{Source:r.PostFormValue("source")} } }
		log.Printf("Cell being created: %s", newCell)
		//call repository to create it
		newCellId, err := lob.NewCell(newCell)
		log.Printf("NewCellId: %s", newCellId)
		if err != nil {
			log.Printf("Error when creating card: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error when creating card: %s", err)
		} 
		//redirect to view the new cell card
		http.Redirect(w, r, "/view/"+newCellId, http.StatusFound)
	})
}

func PageHandler() func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageName := r.URL.Path[len("/page/"):]
		page, err := ioutil.ReadFile("./templates/" + pageName)
		if err != nil {
			notFound, err := ioutil.ReadFile("./templates/card_not_found.html")
			if err != nil {
				log.Printf("Error when returning page: %s", err)
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, string(notFound))
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(page))
	})
}