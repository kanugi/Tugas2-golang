package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

type waterAndWind struct {
	Status struct {
		Water int `json:"water"`
		Wind  int `json:"wind"`
	} `json:"status"`
}

var WaterAndWind *waterAndWind

func handlers() {
	r := mux.NewRouter()
	fmt.Println("starting...")

	r.HandleFunc("/up-to-date", upToDate).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "OK")
		if err != nil {
			log.Println(err)
		}
	})

	http.Handle("/", r)
}

func upToDate(w http.ResponseWriter, r *http.Request) {
	readJSON()
	fmt.Println(WaterAndWind)
	response := make(map[string]interface{})

	// wind status
	wind := WaterAndWind.Status.Wind
	switch {
	case wind <= 6:
		response["windStatus"] = "aman"
	case wind <= 15:
		response["windStatus"] = "siaga"
	default:
		response["windStatus"] = "bahaya"
	}

	// water status
	water := WaterAndWind.Status.Water
	switch {
	case water <= 5:
		response["waterStatus"] = "aman"
	case water <= 8:
		response["waterStatus"] = "siaga"
	default:
		response["waterStatus"] = "bahaya"
	}

	response["water"] = water
	response["wind"] = wind
	//res, err := json.Marshal(response)
	//if err != nil {
	//	w.Write([]byte(err.Error()))
	//}

	randJSON()


	t, err := template.ParseFiles("template.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = t.Execute(w, response)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//w.Write([]byte(res))
}

func readJSON() {
	open, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Println("can't open data.json", err)
	}
	err = json.Unmarshal([]byte(open), &WaterAndWind)
	if err != nil {
		log.Println("err marshal")
	}
}

func randJSON()  {
	wind := rand.Intn(100)
	water := rand.Intn(100)

	WaterAndWind.Status.Wind = wind
	WaterAndWind.Status.Water = water

	res, err := json.Marshal(WaterAndWind)
	if err != nil {
		log.Println(err)
	}

	err = ioutil.WriteFile("data.json", res, 0644)
	log.Println(err)
}

func main() {
	handlers()
	fmt.Println("running port:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("Err handle", err)
	}
}
