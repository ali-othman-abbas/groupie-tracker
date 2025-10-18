package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/env"
	"project/routes"
	"slices"
	"strings"
	"unicode"
)


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
)

func InitalizeData() {
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
		artist := routes.Artist{
			Id:               artistObj.Id,
			Name:             artistObj.Name,
			Image:            artistObj.Image,
			Members:          artistObj.Members,
			CreationDate:     artistObj.CreationDate,
			FirstAlbum:       artistObj.FirstAlbum,
			LocationDatesArr: getLocationDates(relationObj.DatesLocations),
		}
		routes.NameToArtists[artist.Name] = &artist
		routes.ArtistsArr = append(routes.ArtistsArr, &artist)
	}

	slices.SortFunc(routes.ArtistsArr, func(x *routes.Artist, y *routes.Artist) int {
		return strings.Compare(x.Name, y.Name)
	})
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

func getLocationDates(dateLocations map[string][]string) []routes.LocationDates {
	res := []routes.LocationDates{}
	for location, dates := range dateLocations {
		res = append(res, routes.LocationDates{
			Name:  correctlyFormatLocation(location),
			Dates: dates,
		})
	}

	slices.SortFunc(res, func(x routes.LocationDates, y routes.LocationDates) int {
		return strings.Compare(x.Name, y.Name)
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
