package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)
const (
	HOMEPAGEPATH string = "./templates/index.html"
	ARTISTPAGEPATH string = "./templates/artist.html"
)
type Artist struct {
	Id               int
	Name             string
	Image            string
	Members          []string
	CreationDate     int
	FirstAlbum       string
	LocationDatesArr []LocationDates
}
type LocationDates struct {
	Name  string
	Dates []string
}

var (
	NameToArtists = make(map[string]*Artist)
	ArtistsArr    = []*Artist{}
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		writeError(w, "Invalid Request Path", 400)
		return
	}
	data := struct {
		Artists []*Artist
	} {
		Artists: ArtistsArr,
	}
	err := writeTemplate(w, HOMEPAGEPATH, data)
	if err != nil {
		fmt.Println("Inside RootHandler:- ")
		fmt.Println(err)
		writeError(w, "Internal Server Error", 500)
	}
}

func ArtistPageHandler(w http.ResponseWriter, r *http.Request) {
	artistName := strings.TrimPrefix(r.URL.Path, "/page/")
	artist, ok := NameToArtists[artistName]
	if !ok {
		writeError(w, "Artist Page Doesn't Exist", 400)
		return
	} 
	
	data := struct {
		ArtistPage *Artist
	} {
		ArtistPage: artist,
	}
	
	err := writeTemplate(w, ARTISTPAGEPATH, data)
	
	if err != nil {
		fmt.Println("inside ArtistPageHandler")
		fmt.Println(err)
		writeError(w, "Internal Server Error", 500)
		return
	}
}

func writeTemplate(w http.ResponseWriter, path string, data any) error {
	templ, err := template.ParseFiles(path)
	if err != nil {
		fmt.Println("Inside writeTemplate")
		return err
	}
	
	err = templ.Execute(w, data)
	if err != nil {
		fmt.Println("Inside writeTemplate")
		return err
	}
	return nil
}

func writeError(w http.ResponseWriter, errMsg string, statusCode int) {
	http.Error(w, errMsg, statusCode)
}