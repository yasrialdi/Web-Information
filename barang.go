package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

type Barang struct {
	KodeBarang   string `json:"KodeBarang"`
	NamaBarang   string `json:"NamaBarang"`
	JumlahBarang string `json:"JumlahBarang"`
	HargaBarang  string `json:"HargaBarang"`
	StokBarang   string `json:"StokBarang"`
}

// Get all orders

func getBarangs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var barangs []Barang

	sql := `SELECT
				KodeBarang,
				IFNULL(NamaBarang,''),
				IFNULL(JumlahBarang,'') JumlahBarang,
				IFNULL(HargaBarang,'') HargaBarang,
				IFNULL(StokBarang,'') StokBarang
			FROM barangs`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var barang Barang
		err := result.Scan(&barang.KodeBarang, &barang.NamaBarang, &barang.JumlahBarang,
			&barang.HargaBarang, &barang.StokBarang)

		if err != nil {
			panic(err.Error())
		}
		barangs = append(barangs, barang)
	}

	json.NewEncoder(w).Encode(barangs)
}

func createBarang(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		KodeBarang := r.FormValue("KodeBarang")
		NamaBarang := r.FormValue("NamaBarang")
		JumlahBarang := r.FormValue("JumlahBarang")
		HargaBarang := r.FormValue("HargaBarang")
		StokBarang := r.FormValue("StokBarang")
		stmt, err := db.Prepare("INSERT INTO barangs (KodeBarang,NamaBarang,JumlahBarang,HargaBarang,StokBarang) VALUES (?,?,?,?,?)")

		_, err = stmt.Exec(KodeBarang, NamaBarang, JumlahBarang, HargaBarang, StokBarang)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}

func getBarang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var barangs []Barang
	params := mux.Vars(r)

	sql := `SELECT
			KodeBarang,
			IFNULL(NamaBarang,''),
			IFNULL(JumlahBarang,'') JumlahBarang,
			IFNULL(HargaBarang,'') HargaBarang,
			IFNULL(StokBarang,'') StokBarang
			FROM barangs WHERE KodeBarang = ?`

	result, err := db.Query(sql, params["kode"])

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var barang Barang

	for result.Next() {

		err := result.Scan(&barang.KodeBarang, &barang.NamaBarang, &barang.JumlahBarang,
			&barang.HargaBarang, &barang.StokBarang)

		if err != nil {
			panic(err.Error())
		}

		barangs = append(barangs, barang)
	}

	json.NewEncoder(w).Encode(barangs)
}

func updateBarang(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newNamaBarang := r.FormValue("NamaBarang")

		stmt, err := db.Prepare("UPDATE barangs SET NamaBarang = ? WHERE KodeBarang = ?")

		_, err = stmt.Exec(newNamaBarang, params["kode"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Barang with KodeBarang = %s was updated", params["kode"])
	}
}

func deleteBarang(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM barangs WHERE KodeBarang = ?")

	_, err = stmt.Exec(params["kode"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Barang with KodeBarang = %s was deleted", params["kode"])
}

func delBarang(w http.ResponseWriter, r *http.Request) {

	KodeBarang := r.FormValue("KodeBarang")
	NamaBarang := r.FormValue("NamaBarang")

	stmt, err := db.Prepare("DELETE FROM barangs WHERE KodeBarang = ? AND NamaBarang = ?")

	_, err = stmt.Exec(KodeBarang, NamaBarang)

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Barang with KodeBarang = %s was deleted", KodeBarang)
}

func getPost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var barangs []Barang

	KodeBarang := r.FormValue("KodeBarang")
	NamaBarang := r.FormValue("NamaBarang")

	sql := `SELECT
			KodeBarang,
			IFNULL(NamaBarang,''),
			IFNULL(JumlahBarang,'') JumlahBarang,
			IFNULL(HargaBarang,'') HargaBarang,
			IFNULL(StokBarang,'') StokBarang
			FROM barangs WHERE KodeBarang = ? AND NamaBarang = ?`

	result, err := db.Query(sql, KodeBarang, NamaBarang)

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var barang Barang

	for result.Next() {

		err := result.Scan(&barang.KodeBarang, &barang.NamaBarang, &barang.JumlahBarang,
			&barang.HargaBarang, &barang.StokBarang)

		if err != nil {
			panic(err.Error())
		}

		barangs = append(barangs, barang)
	}

	json.NewEncoder(w).Encode(barang)

}

// Main function
func main() {

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/northwind")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/barangs", getBarangs).Methods("GET")
	r.HandleFunc("/barangs/{kode}", getBarang).Methods("GET")
	r.HandleFunc("/barangs", createBarang).Methods("POST")
	r.HandleFunc("/barangs/{kode}", updateBarang).Methods("PUT")
	r.HandleFunc("/barangs/{kode}", deleteBarang).Methods("DELETE")

	//New
	r.HandleFunc("/getbarang", getPost).Methods("POST")

	//DelBarang
	r.HandleFunc("/delbarang", delBarang).Methods("POST")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
