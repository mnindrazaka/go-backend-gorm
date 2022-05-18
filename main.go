package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Album struct {
	ID     int
	Title  string
	Artist string
	Price  float32
}

type AlbumRequest struct {
	Title  string
	Artist string
	Price  float32
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome"))
}

func handleAlbums(w http.ResponseWriter, r *http.Request) {
	var albums []Album
	result := db.Table("album").Find(&albums)

	if result.Error != nil {
		w.Write([]byte("Failed to get data"))
	}

	json.NewEncoder(w).Encode(albums)
}

func handleAlbumDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var formattedId, err = strconv.Atoi(id)

	if err != nil {
		w.Write([]byte("Failed to get data"))
	}

	var album = Album{ID: formattedId}
	result := db.Table("album").First(&album)

	if result.Error != nil {
		w.Write([]byte("Failed to get data"))
	}

	json.NewEncoder(w).Encode(album)
}

func handleAlbumCreate(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var albumRequest AlbumRequest
	json.Unmarshal(reqBody, &albumRequest)

	result := db.Table("album").Create(&albumRequest)

	if result.Error != nil {
		w.Write([]byte("Failed to create data"))
	}

	w.Write([]byte("Success to create data"))
}

func handleAlbumUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)

	var albumRequest AlbumRequest
	json.Unmarshal(reqBody, &albumRequest)

	var formattedId, err = strconv.Atoi(id)

	if err != nil {
		w.Write([]byte("Failed to get data"))
	}

	album := Album{ID: formattedId, Title: albumRequest.Title, Artist: albumRequest.Artist, Price: albumRequest.Price}

	result := db.Table("album").Save(&album)

	if result.Error != nil {
		w.Write([]byte("Failed to update data"))
	}

	w.Write([]byte("Success to update data"))
}

func handleAlbumDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var formattedId, err = strconv.Atoi(id)

	if err != nil {
		w.Write([]byte("Failed to get data"))
	}

	album := Album{ID: formattedId}

	result := db.Table("album").Delete(&album)

	if result.Error != nil {
		w.Write([]byte("Failed to delete data"))
	}

	w.Write([]byte("Success to delete data"))
}

var db *gorm.DB

func main() {
	dsn := "root:roottoor@tcp(127.0.0.1:3306)/recordings?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = gorm.Open(mysql.Open(dsn))

	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", handleHome)
	router.HandleFunc("/albums", handleAlbumCreate).Methods("POST")
	router.HandleFunc("/albums", handleAlbums)
	router.HandleFunc("/albums/{id}", handleAlbumUpdate).Methods("PUT")
	router.HandleFunc("/albums/{id}", handleAlbumDelete).Methods("DELETE")
	router.HandleFunc("/albums/{id}", handleAlbumDetail)

	http.ListenAndServe(":3000", router)
}
