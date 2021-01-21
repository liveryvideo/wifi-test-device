package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"wifi-test-device/wifi"

	"github.com/gorilla/mux"
)

type HostDeviceStatus struct {
	HostName                  string
	OperatingSystem           string
	TestDeviceSoftwareVersion string
	NetworkStatus             []wifi.NetworkStatus
}

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
	router.HandleFunc("/api/status", fetchStatus).Methods("GET")
	router.HandleFunc("/api/logs", fetchLogs).Methods("GET")

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
		return
	}

	unmarshalError := json.Unmarshal(data, globalRules)

	if unmarshalError != nil {
		fmt.Printf("Could not unmarshal request: %s\r\n", unmarshalError)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte("Could not unmarshal request."))
		return
	}

	wifi.UpdateGlobalRules()
	responseWriter.WriteHeader(200)
}

func fetchStatus(responseWriter http.ResponseWriter, request *http.Request) {
	status, err := wifi.GetNetworkStatus()

	responseWriter.Header().Set("Content-Type", "application/json")

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte("Could not get network status."))
		return
	}

	hostName, err := os.Hostname()
	if err != nil {
		log.Println("Failed to fetch host-device name:", err)
		hostName = "Unknown"
	}

	hostDeviceStatus := HostDeviceStatus{
		HostName:                  hostName,
		OperatingSystem:           runtime.GOOS,
		TestDeviceSoftwareVersion: "1.0",
		NetworkStatus:             status,
	}
	json, marshalError := json.Marshal(hostDeviceStatus)

	if marshalError != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte("Failed to marshal network status."))
		return
	}

	responseWriter.Write([]byte(json))
}

func fetchLogs(responseWriter http.ResponseWriter, request *http.Request) {
	logs := wifi.FetchLogs()
	json, err := json.Marshal(logs)

	if err != nil {
		log.Println("Failed to marshal logs:", err)
		responseWriter.WriteHeader(500)
		responseWriter.Write([]byte("Failed to marshal logs."))
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(json)
}
