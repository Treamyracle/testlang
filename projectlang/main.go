package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   int    `json:"id"`
	Nama string `json:"nama"`
	Role string `json:"role"`
}
type Rapat struct {
	ID     int    `json:"id"`
	Judul  string `json:"judul"`
	Tempat string `json:"tempat"`
}

type Absensi struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	RapatID int    `json:"rapat_id"`
	Status  string `json:"status"`
}

var db *sql.DB

// func initDatabase() {
// 	var err error
// 	db, err = sql.Open("mysql", "root:ngok@tcp(localhost:3306)/database")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func main() {
	 // Load connection string from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env", err)
	}

	// Open a connection to PlanetScale
	db, err = sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	rows, err := db.Query("SHOW TABLES")
	 if err != nil {
		 log.Fatalf("failed to query: %v", err)
	 }
	 defer rows.Close()
 
	 var tableName string
	 for rows.Next() {
		 if err := rows.Scan(&tableName); err != nil {
			 log.Fatalf("failed to scan row: %v", err)
		 }
		 log.Println(tableName)
	 }
 
	 defer db.Close()

	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/rapat", getRapats)
	http.HandleFunc("/rapat/", getRapat)
	http.HandleFunc("/absensi", getAbsensis)
	http.HandleFunc("/absensi/", getAbsensi)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	rows, err := db.Query("SELECT id, nama, role FROM user")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Nama, &user.Role)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getRapats(w http.ResponseWriter, r *http.Request) {
	var rapats []Rapat
	rows, err := db.Query("SELECT id, judul, tempat FROM rapat")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var rapat Rapat
		err := rows.Scan(&rapat.ID, &rapat.Judul, &rapat.Tempat)
		if err != nil {
			log.Fatal(err)
		}
		rapats = append(rapats, rapat)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rapats)
}

func getAbsensis(w http.ResponseWriter, r *http.Request) {
	var absensis []Absensi
	rows, err := db.Query("SELECT id, user_id, rapat_id, status FROM absensi	")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var absensi Absensi
		err := rows.Scan(&absensi.ID, &absensi.UserID, &absensi.RapatID, &absensi.Status)
		if err != nil {
			log.Fatal(err)
		}
		absensis = append(absensis, absensi)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(absensis)
}

func getAbsensi(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/absensi/"):]
	if idStr == "" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var absensi Absensi
	err = db.QueryRow("SELECT id, user_id, rapat_id, status FROM absensi WHERE id = ?", id).
		Scan(&absensi.ID, &absensi.UserID, &absensi.RapatID, &absensi.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(absensi)
}

func getRapat(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/rapat/"):])
	if err != nil {
		log.Fatal(err)
	}

	var absensis []Absensi
	rows, err := db.Query("SELECT id, user_id, rapat_id, status FROM absensi WHERE rapat_id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var absensi Absensi
		err := rows.Scan(&absensi.ID, &absensi.UserID, &absensi.RapatID, &absensi.Status)
		if err != nil {
			log.Fatal(err)
		}
		absensis = append(absensis, absensi)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(absensis)
}
