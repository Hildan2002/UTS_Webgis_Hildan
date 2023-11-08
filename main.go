package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Marker struct {
	ID      int
	Lat     float64
	Lng     float64
	Title   string
	Content string
}

var db *sql.DB

func main() {
	db, _ = sql.Open("sqlite3", "markers.db")
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS markers (id INTEGER PRIMARY KEY, lat REAL, lng REAL, title TEXT, content TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	r := mux.NewRouter()
	r.HandleFunc("/", ListMarkers).Methods("GET")
	r.HandleFunc("/add", AddMarker).Methods("POST")
	r.HandleFunc("/edit/{id}", EditMarker).Methods("GET")
	r.HandleFunc("/update/{id}", UpdateMarker).Methods("POST")
	r.HandleFunc("/delete/{id}", DeleteMarker).Methods("GET")

	fmt.Printf("Memulai server pada port 8080\n")
	fmt.Printf("server dapat diakses pada http://localhost:8080 \n")

	http.Handle("/", r)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func ListMarkers(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.Query("SELECT id, lat, lng, title, content FROM markers")
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var markers []Marker
	for rows.Next() {
		marker := Marker{}
		err := rows.Scan(&marker.ID, &marker.Lat, &marker.Lng, &marker.Title, &marker.Content)
		if err != nil {
			return
		}
		markers = append(markers, marker)
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	err = tmpl.Execute(w, markers)
	if err != nil {
		return
	}
}

func AddMarker(w http.ResponseWriter, r *http.Request) {
	lat := r.FormValue("lat")
	lng := r.FormValue("lng")
	title := r.FormValue("title")
	content := r.FormValue("content")

	_, err := db.Exec("INSERT INTO markers (lat, lng, title, content) VALUES (?, ?, ?, ?)", lat, lng, title, content)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditMarker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	row := db.QueryRow("SELECT id, lat, lng, title, content FROM markers WHERE id = ?", id)
	marker := Marker{}
	err := row.Scan(&marker.ID, &marker.Lat, &marker.Lng, &marker.Title, &marker.Content)
	if err != nil {
		return
	}

	tmpl, _ := template.ParseFiles("templates/edit.html")
	err = tmpl.Execute(w, marker)
	if err != nil {
		return
	}
}

func UpdateMarker(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")
	latStr := r.FormValue("lat")
	lngStr := r.FormValue("lng")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE markers SET title=?, content=?, lat=?, lng=? WHERE id=?", title, content, lat, lng, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteMarker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec("DELETE FROM markers WHERE id=?", id)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
