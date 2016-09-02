package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/urfave/negroni"
)

var db *leveldb.DB

func init() {
	var file = flag.String("F", "level", "DB path")
	flag.Parse()
	db, _ = leveldb.OpenFile(*file, nil)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/get/{key}", GetHandler)
	r.HandleFunc("/set/{key}/{value}", SetHandler)
	r.HandleFunc("/delete/{key}", DeleteHandler)

	n := negroni.Classic()
	n.UseHandler(r)

	http.ListenAndServe(":8080", n)
}
func GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	data, err := db.Get([]byte(key), nil)
	if err == leveldb.ErrNotFound {
		fmt.Fprintf(w, "key not found")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Fprintf(w, byteString(data))
	}
}
func SetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value := vars["value"]
	err := db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		panic(err)
	} else {
		fmt.Fprintf(w, "set OK!")
	}
}
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	err := db.Delete([]byte(key), nil)
	if err != nil {
		panic(err)
	} else {
		fmt.Fprintf(w, "delete OK!")
	}
}
func byteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}
