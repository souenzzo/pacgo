package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Origin struct {
	Src string `json:"src"`
}

func addSrcHandler(w http.ResponseWriter, r *http.Request) {
	var O Origin
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &O)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = addSrc(O.Src)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func listSrcHandler(w http.ResponseWriter, r *http.Request) {
	names := listSrc()
	n, err := json.Marshal(&names)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%s", n)
}

func makePkgHandler(w http.ResponseWriter, r *http.Request) {
	pkgName := mux.Vars(r)["pkg"]
	err := makePkg(pkgName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Fprintln(w, "Ok!")
}

func web(addr string) error {
	r := mux.NewRouter()
	r.HandleFunc("/src", addSrcHandler).Methods("POST")
	r.HandleFunc("/src", listSrcHandler).Methods("GET")
	r.HandleFunc("/make/{pkg}", makePkgHandler).Methods("POST")
	s := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	return s.ListenAndServe()
}
