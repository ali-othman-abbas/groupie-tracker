//TODO: ADD PAGES FOR INDIVIDUAL ARTISTS
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"project/env"
	"slices"
	"strings"
	"unicode"
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
	name  string
	dates []string
}

type ArtistResponse struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type RelationResponse struct {
	Index []Relation `json:"index"`
}

type Response interface {
	RelationResponse | []ArtistResponse
}
type Relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

const (
	ARTSTSURL  string = "https://groupietrackers.herokuapp.com/api/artists"
	RELATIONURL string = "https://groupietrackers.herokuapp.com/api/relation"
	HOMEPAGEPATH string = "./templates/index.html"
	ARTISTPAGEPATH string = "./templates/artist.html"
)

var (
	nameToArtists = make(map[string]*Artist)
	artistsArr    = []*Artist{}
)

func main() {
	initalizeData()
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/page/", artistPageHandler)

	fmt.Println("Starting server at port", env.PORT)
	err := http.ListenAndServe(env.IP+env.PORT, nil)
	if err != nil {
		panic(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		writeError(w, "Invalid Request Path", 400)
		return
	}
	data := struct {
		Artists []*Artist
	} {
		Artists: artistsArr,
	}
	err := writeTemplate(w, HOMEPAGEPATH, data)
	if err != nil {
		fmt.Println("Inside rootHandler:- ")
		fmt.Println(err)
		writeError(w, "Internal Server Error", 500)
	}
}

func artistPageHandler(w http.ResponseWriter, r *http.Request) {
	artistName := strings.TrimPrefix(r.URL.Path, "/page/")
	artist, ok := nameToArtists[artistName]
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
		fmt.Println("inside artistPageHandler")
		fmt.Println(err)
		writeError(w, "Internal Server Error", 500)
		return
	}
}

func initalizeData() {
	artistsResp, err := get[[]ArtistResponse](ARTSTSURL)
	artists := *artistsResp
	if err != nil {
		fmt.Println("inside initalizeData:-\n" + err.Error())
	}
	relationResp, err := get[RelationResponse](RELATIONURL)
	if err != nil {
		fmt.Println("inside initalizeData:-\n" + err.Error())
	}
	relations := relationResp.Index

	for i := range artists {
		artistObj := artists[i]
		relationObj := relations[i]
		if artistObj.Id == 21 {
			artistObj.Image = fmt.Sprintf("http://%s%s/static/mamonas_assas.webp", env.IP, env.PORT)
		}
		artist := Artist{
			Id:               artistObj.Id,
			Name:             artistObj.Name,
			Image:            artistObj.Image,
			Members:          artistObj.Members,
			CreationDate:     artistObj.CreationDate,
			FirstAlbum:       artistObj.FirstAlbum,
			LocationDatesArr: getLocationDates(relationObj.DatesLocations),
		}
		nameToArtists[artist.Name] = &artist
		artistsArr = append(artistsArr, &artist)
	}
	
	slices.SortFunc(artistsArr, func(x *Artist, y *Artist) int {
		if x.Name < y.Name {
			return -1
		} else if x.Name > y.Name {
			return 1
		}
		
		return 0
	})
}

func getLocationDates(dateLocations map[string][]string) []LocationDates {
	res := []LocationDates{}
	for location, dates := range dateLocations {
		res = append(res, LocationDates{
			name:  correctlyFormatLocation(location),
			dates: dates,
		})
	}

	slices.SortFunc(res, func(x LocationDates, y LocationDates) int {
		if x.name < y.name {
			return -1
		} else if x.name > y.name {
			return 1
		}
		return 0
	})

	return res
}

func correctlyFormatLocation(location string) string {
	stateAndCountry := strings.Split(location, "-")
	state := strings.Split(stateAndCountry[0], "_")
	country := strings.Split(stateAndCountry[1], "_")
	state = capitalizeStrings(state)
	country = capitalizeStrings(country)
	stateAndCountry[0] = strings.Join(state, " ")
	stateAndCountry[1] = strings.Join(country, " ")
	if stateAndCountry[1] == "Usa" || stateAndCountry[1] == "Uk" {
		stateAndCountry[1] = strings.ToUpper(stateAndCountry[1])
	}
	newLocation := strings.Join(stateAndCountry, ", ")
	return newLocation
}

func capitalizeStrings(strArr []string) []string {
	for i, _ := range strArr {
		runeArr := []rune(strArr[i])
		runeArr[0] = unicode.ToUpper(runeArr[0])
		strArr[i] = string(runeArr)
	}

	return strArr
}

func get[T Response](url string) (*T, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var response T
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
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
