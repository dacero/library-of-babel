package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"encoding/json"

	"github.com/dacero/labyrinth-of-babel/repository"
	"github.com/dacero/labyrinth-of-babel/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func ViewHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cellId := mux.Vars(r)["id"]
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

func EditHandler(lob repository.LobRepository, store *sessions.CookieStore) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth, _ := checkAuthorization(w, r, store); !auth {
			http.Redirect(w, r, "/page/auth.html", http.StatusFound)
		}
		
		cellId := mux.Vars(r)["id"]
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
			t, err := template.ParseFiles("./templates/edit_card.gohtml")
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

func SourcesHandler(lob repository.LobRepository, store *sessions.CookieStore) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth, _ := checkAuthorization(w, r, store); !auth {
			http.Redirect(w, r, "/page/auth.html", http.StatusFound)
		}
		
		cellId := mux.Vars(r)["id"]
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
			t, err := template.ParseFiles("./templates/edit_sources.gohtml")
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

func AddSourceHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cellId := mux.Vars(r)["id"]
		newSource := models.Source{Source: r.PostFormValue("source") }
		_, err := lob.AddSourceToCell(cellId, newSource)
		if err != nil {
			log.Printf("Error when adding source: %s", err)
			notFound, err := ioutil.ReadFile("./templates/card_not_found.html")
			if err != nil {
				log.Printf("Error when returning card: %s", err)
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, string(notFound))
		} else {
			http.Redirect(w, r, "/cell/"+cellId+"/sources", http.StatusFound)
		}
	})
}

func RemoveSourceHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cellId := mux.Vars(r)["id"]
		source := models.Source{Source: r.PostFormValue("source") }
		_, err := lob.RemoveSourceFromCell(cellId, source)
		if err != nil {
			log.Printf("Error when removing source: %s", err)
			notFound, err := ioutil.ReadFile("./templates/card_not_found.html")
			if err != nil {
				log.Printf("Error when removing source: %s", err)
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, string(notFound))
		} else {
			http.Redirect(w, r, "/cell/"+cellId+"/sources", http.StatusFound)
		}
	})
}

func LinksHandler(lob repository.LobRepository, store *sessions.CookieStore) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth, _ := checkAuthorization(w, r, store); !auth {
			http.Redirect(w, r, "/page/auth.html", http.StatusFound)
		}
		
		cellId := mux.Vars(r)["id"]
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
			t, err := template.ParseFiles("./templates/edit_links.gohtml")
			if err != nil {
				log.Printf("Error when displaying links: %s", err)
			}
			err = t.Execute(w, cell)
			if err != nil {
				log.Printf("Error when displaying links: %s", err)
			}
		}
	})
}

func LinkCellsHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cellA := mux.Vars(r)["id"]
		cellB := r.PostFormValue("cellToLink")
		err := lob.LinkCells(cellA, cellB)
		if err != nil {
			log.Printf("Error when linking cells: %s", err)
		}
		http.Redirect(w, r, "/cell/"+cellA+"/links", http.StatusFound)
	})
}

func UnlinkCellsHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cellA := mux.Vars(r)["id"]
		cellB := r.PostFormValue("cellToUnlink")
		err := lob.UnlinkCells(cellA, cellB)
		if err != nil {
			log.Printf("Error when linking cells: %s", err)
		}
		http.Redirect(w, r, "/cell/"+cellA+"/links", http.StatusFound)
	})
}

func checkAuthorization(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore) (bool, error) {
	//store == nil is for testing purposes
	if store == nil {
		return true, nil
	}
	
	session, err := store.Get(r, "lob-session")
	if err != nil {
		return false, err
	}
	auth := session.Values["authenticated"]
	if auth == nil {
		return false, nil
	}
	if !auth.(bool) {
		session.AddFlash("You don't have access!")
		err = session.Save(r, w)
		if err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func SaveHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {		
		//parse the form and create the cell
		updateCell := models.Cell{Id: r.PostFormValue("cellId"),
			Title: r.PostFormValue("title"),
			Body: r.PostFormValue("body"),
			Room: r.PostFormValue("room")}
		//call repository to create it
		_, err := lob.UpdateCell(updateCell)
		if err != nil {
			log.Printf("Error when updating card: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error when updating card: %s", err)
		} 
		//redirect to view the new cell card
		http.Redirect(w, r, "/cell/"+r.PostFormValue("cellId"), http.StatusFound)
	})
}

func CreateHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//parse the form and create the cell
		log.Printf("New cell title: %s", r.PostFormValue("title"))
		newCell := models.Cell{Title: r.PostFormValue("title"),
			Body: r.PostFormValue("body"),
			Room: r.PostFormValue("room")}
		//call repository to create it
		newCellId, err := lob.NewCell(newCell)
		if err != nil {
			log.Printf("Error when creating card: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error when creating card: %s", err)
		} 
		//redirect to view the new cell card
		http.Redirect(w, r, "/cell/"+newCellId, http.StatusFound)
	})
}

func PageHandler() func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageName := mux.Vars(r)["page"]
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

func SearchSourcesHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		term := r.FormValue("term")
		sources := lob.SearchSources(term)
		returnString := "["
		for _, source := range sources {
			returnString += `"` + source.String() + `",`
		}
		returnString = returnString[:len(returnString)-1] + "]"
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, returnString)
	})
}

func SearchRoomsHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		term := r.FormValue("term")
		rooms := lob.SearchRooms(term)
		returnString := "["
		for _, room := range rooms {
			returnString += `"` + room + `",`
		}
		returnString = returnString[:len(returnString)-1] + "]"
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, returnString)
	})
}

func SearchCellsHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		term := r.FormValue("term")
		log.Printf("Searching for cells with %s", term)
		cells := lob.SearchCells(term)
		type CellLinkAlias struct {
			Id   string `json:"value"`
			Text string `json:"label"`
		}
		var alias []CellLinkAlias
		for _, cell := range cells {
			cellAlias := CellLinkAlias{Id: cell.Id, Text: cell.Summary()}
			alias = append(alias, cellAlias)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(alias)
		if err != nil {
			log.Printf("Error when returning search result: %s", err)
		}
	})
}


func RoomListHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rooms, err := lob.ListRooms()
		if err != nil {
			log.Printf("Error when obtaining the list of rooms: %s", err)
			notFound, err := ioutil.ReadFile("./templates/card_not_found.html")
			if err != nil {
				log.Printf("Error when obtaining the list of rooms: %s", err)
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, string(notFound))
		} else {
			t, err := template.ParseFiles("./templates/rooms.gohtml")
			if err != nil {
				log.Printf("Error when parsing the rooms template: %s", err)
			}
			type data struct {
				Rooms []models.CollectionOfCells
			}
			roomsData := data{Rooms: rooms}
			err = t.Execute(w, roomsData)
			if err != nil {
				log.Printf("Error when returning card: %s", err)
			}
		}
	})
}

func RoomHandler(lob repository.LobRepository) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		room := mux.Vars(r)["room"]
		cells, err := lob.ListCellsInRoom(room)
		if err != nil {
			log.Printf("Error when entering room: %s", err)
			notFound, err := ioutil.ReadFile("./templates/card_not_found.html")
			if err != nil {
				log.Printf("Error when entering room: %s", err)
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, string(notFound))
		} else {
			t, err := template.ParseFiles("./templates/cells_collection.gohtml")
			if err != nil {
				log.Printf("Error when parsing the room template: %s", err)
			}
			type data struct {
				Name	string
				Cells 	[]models.Cell
			}
			roomCells := data{Name: room, Cells: cells}
			err = t.Execute(w, roomCells)
			if err != nil {
				log.Printf("Error when returning card: %s", err)
			}
		}
	})
}