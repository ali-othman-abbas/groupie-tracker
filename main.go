package main

import (
	"encoding/json"
	"fmt"
	_ "html/template"
	"net/http"
	"project/env"
	"slices"
	"strings"
	"unicode"
)

type Artist struct {
	id               int
	name             string
	image            string
	members          []string
	creationDate     int
	firstAlbum       string
	locationDatesArr []LocationDates
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
	artistsUrl  string = "https://groupietrackers.herokuapp.com/api/artists"
	relationUrl string = "https://groupietrackers.herokuapp.com/api/relation"
)

var (
	nameToArtists = make(map[string]*Artist)
	artistsArr    = []Artist{}
)

func main() {
	initalizeData()
	http.HandleFunc("/", rootHandler)

	fmt.Println("Starting server at port", env.PORT)
	err := http.ListenAndServe(env.IP+env.PORT, nil)
	if err != nil {
		panic(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	for _, artist := range artistsArr {
		fmt.Fprintf(w, "%v, %v, %v, %v, %v, %v, %v\n", artist.id, artist.image, artist.name, artist.members, artist.creationDate, artist.firstAlbum, artist.locationDatesArr)
		fmt.Fprint(w, "----------------------\n")
	}
}

func initalizeData() {
	artistsResp, err := get[[]ArtistResponse](artistsUrl)
	artists := *artistsResp
	if err != nil {
		fmt.Println("inside initalizeData:-\n" + err.Error())
	}
	relationResp, err := get[RelationResponse](relationUrl)
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
			id:               artistObj.Id,
			name:             artistObj.Name,
			image:            artistObj.Image,
			members:          artistObj.Members,
			creationDate:     artistObj.CreationDate,
			firstAlbum:       artistObj.FirstAlbum,
			locationDatesArr: getLocationDates(relationObj.DatesLocations),
		}
		nameToArtists[artist.name] = &artist
		artistsArr = append(artistsArr, artist)
	}

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
