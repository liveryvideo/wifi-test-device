package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"wifi-test-device/wifi"

	"github.com/gorilla/mux"
)

var router *mux.Router

func getRouter() *mux.Router {
	if router == nil {
		router = mux.NewRouter().StrictSlash(true)
		initRouter()
	}
	return router
}

func initRouter() {
	router.HandleFunc("/api/devices", fetchDevices).Methods("GET")
	router.HandleFunc("/api/settings", handleSettings).Methods("GET", "POST")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	router.Use(mux.CORSMethodMiddleware(router))
}

func fetchDevices(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")

	json, err := json.Marshal(wifi.GetLeasedDevices())
	if err != nil {
		fmt.Printf("Could not marshal devices: %s\r\n", err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte("Could not marshal devices."))
		return
	}

	responseWriter.Write([]byte(json))
}

func handleSettings(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		fetchSettings(responseWriter, request)
		break
	case "POST":
		updateSettings(responseWriter, request)
		break
	}
}

func fetchSettings(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")

	json, err := json.Marshal(wifi.GetGlobalRules())

	if err != nil {
		fmt.Printf("Could not marshal devices: %s\r\n", err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte("Could not marshal devices."))
		return
	}

	responseWriter.Write([]byte(json))
}

func updateSettings(responseWriter http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	globalRules := wifi.GetGlobalRules()
	data, readError := ioutil.ReadAll(request.Body)

	if readError != nil {
		fmt.Printf("Could not read request: %s\r\n", readError)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte("Could not read request."))
		responseWriter.WriteHeader(500)
		return
	}

	unmarshalError := json.Unmarshal(data, globalRules)

	if unmarshalError != nil {
		fmt.Printf("Could not unmarshal request: %s\r\n", unmarshalError)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte("Could not unmarshal request."))
		responseWriter.WriteHeader(400)
		return
	}

	wifi.UpdateGlobalRules()
	responseWriter.WriteHeader(200)
}
