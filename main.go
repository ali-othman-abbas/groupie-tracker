// TODO: ADD PAGES FOR INDIVIDUAL ARTISTS
package main

import (
	"fmt"
	"net/http"
	"project/requests"
	"project/env"
	"project/routes"
)

func main() {
	requests.InitalizeData()
	http.HandleFunc("/", routes.RootHandler)
	http.HandleFunc("/page/", routes.ArtistPageHandler)

	fmt.Println("Starting server at port", env.PORT)
	err := http.ListenAndServe(env.IP+env.PORT, nil)
	if err != nil {
		panic(err)
	}
}
