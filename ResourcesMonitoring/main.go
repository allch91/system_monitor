package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	common "./common"
)

const (
    STATIC_DIR = "/static/"
    PORT       = "3001"
)

var router = NewRouter()

func NewRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
		router.PathPrefix("/css/").Handler(http.StripPrefix("/css/",
		http.FileServer(http.Dir("static/css/"))))
    return router
}

func main() {

	//====== GET METHODS
	router.HandleFunc("/", common.LoginPageHandler)
	router.HandleFunc("/adminPage",common.HomePageHandler)
	router.HandleFunc("/cpu_graph", common.CpuProcessHandler)
	router.HandleFunc("/ram_graph", common.RamProcessHandler)

	//====== POST METHODS
	router.HandleFunc("/login", common.LoginHandler).Methods("POST")
	router.HandleFunc("/ram_data",common.RamData).Methods("POST")
	router.HandleFunc("/cpu_data",common.CpuData).Methods("POST")
	router.HandleFunc("/data",common.AdminHandler).Methods("POST")

	//====== CONFIG
	http.Handle("/",router)
	log.Println("Server running on http://localhost:8080")
  log.Fatal(http.ListenAndServe(":3001", nil))
}
