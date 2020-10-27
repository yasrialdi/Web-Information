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

type mahasiswa struct {
	NoBp     int    `json:"NoBp"`
	Nama     string `json:"Nama"`
	Fakultas string `json:"Fakultas"`
	Jurusan  string `json:"Jurusan"`
	Alamat   struct {
		Jalan     string `json:"Jalan"`
		Kelurahan string `json:"Kelurahan"`
		Kecamatan string `json:"Kecamatan"`
		Kabupaten string `json:"Kabupaten"`
		Provinsi  string `json:"Provinsi"`
	} `json:"Alamat"`
	Nilai []nilai `json:"Nilai"`
}

type nilai struct {
	NoBp       int     `json:"NoBp"`
	IDMatkul   int     `json:"IdMatkul"`
	NamaMatkul string  `json:"NamaMatkul"`
	Nilai      float64 `json:"Nilai"`
	Semester   string  `json:"Semester"`
}

func getNilai(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var mhsx []mahasiswa

	params := mux.Vars(r)

	sql := `SELECT
				nobp,
				IFNULL(nama,'') nama,
				IFNULL(jalan,'') jalan,
				IFNULL(kelurahan,'') kelurahan,
				IFNULL(kecamatan,'') kecamatan,
				IFNULL(kabupaten,'') kabupaten,
				IFNULL(provinsi,'') provinsi,
				IFNULL(fakultas,'') fakultas,
				IFNULL(jurusan,'') jurusan				
			FROM mahasiswa WHERE nobp IN (?)`

	result, err := db.Query(sql, params["NoBp"])

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		var mhs mahasiswa
		err := result.Scan(&mhs.NoBp, &mhs.Nama, &mhs.Alamat.Jalan, &mhs.Alamat.Kelurahan, &mhs.Alamat.Kecamatan, &mhs.Alamat.Kabupaten, &mhs.Alamat.Provinsi, &mhs.Fakultas, &mhs.Jurusan)

		if err != nil {
			panic(err.Error())
		}

		sqlNilai := `SELECT
						nobp		
						, matkul.id_matkul
						, matkul.nama
						, nilai
						, semester
					FROM
						nilai INNER JOIN matkul 
							ON (nilai.id_matkul = matkul.id_matkul)
					WHERE nobp = ?`

		resultNilai, errNilai := db.Query(sqlNilai, mhs.NoBp)

		defer resultNilai.Close()

		if errNilai != nil {
			panic(err.Error())
		}

		for resultNilai.Next() {
			var nilaix nilai
			err := resultNilai.Scan(&nilaix.NoBp, &nilaix.IDMatkul, &nilaix.NamaMatkul, &nilaix.Nilai, &nilaix.Semester)
			if err != nil {
				panic(err.Error())
			}
			mhs.Nilai = append(mhs.Nilai, nilaix)
		}
		mhsx = append(mhsx, mhs)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mhsx)
}

func getNilaiAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var mhsx []mahasiswa

	sql := `SELECT
				nobp,
				IFNULL(nama,'') nama,
				IFNULL(jalan,'') jalan,
				IFNULL(kelurahan,'') kelurahan,
				IFNULL(kecamatan,'') kecamatan,
				IFNULL(kabupaten,'') kabupaten,
				IFNULL(provinsi,'') provinsi,
				IFNULL(fakultas,'') fakultas,
				IFNULL(jurusan,'') jurusan				
			FROM mahasiswa`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		var mhs mahasiswa
		err := result.Scan(&mhs.NoBp, &mhs.Nama, &mhs.Alamat.Jalan, &mhs.Alamat.Kelurahan, &mhs.Alamat.Kecamatan, &mhs.Alamat.Kabupaten, &mhs.Alamat.Provinsi, &mhs.Fakultas, &mhs.Jurusan)

		if err != nil {
			panic(err.Error())
		}

		sqlNilai := `SELECT
						nobp		
						, matkul.id_matkul
						, matkul.nama
						, nilai
						, semester
					FROM
						nilai INNER JOIN matkul 
							ON (nilai.id_matkul = matkul.id_matkul)
					WHERE nobp = ?`

		resultNilai, errNilai := db.Query(sqlNilai, mhs.NoBp)

		defer resultNilai.Close()

		if errNilai != nil {
			panic(err.Error())
		}

		for resultNilai.Next() {
			var nilaix nilai
			err := resultNilai.Scan(&nilaix.NoBp, &nilaix.IDMatkul, &nilaix.NamaMatkul, &nilaix.Nilai, &nilaix.Semester)
			if err != nil {
				panic(err.Error())
			}
			mhs.Nilai = append(mhs.Nilai, nilaix)
		}
		mhsx = append(mhsx, mhs)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mhsx)
}

func main() {

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/kemahasiswaanpnp")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	r := mux.NewRouter()
	r.HandleFunc("/nilai/{NoBp}", getNilai).Methods("GET")
	r.HandleFunc("/nilai", getNilaiAll).Methods("GET")

	fmt.Println("Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
