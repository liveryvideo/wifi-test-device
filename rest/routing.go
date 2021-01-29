package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"wifi-test-device/wifi"

	"github.com/gorilla/mux"
)

type HostDeviceStatus struct {
	HostName                  string
	OperatingSystem           string
	TestDeviceSoftwareVersion string
	NetworkStatus             []wifi.NetworkInterface
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
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(headerAdapter)

	apiRouter.HandleFunc("/devices", fetchDevices).Methods("GET")
	apiRouter.HandleFunc("/settings", handleSettings).Methods("GET", "POST")
	apiRouter.HandleFunc("/status", fetchStatus).Methods("GET")
	apiRouter.HandleFunc("/logs", fetchLogs).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	router.Use(mux.CORSMethodMiddleware(router))
}

func headerAdapter(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		nextHandler.ServeHTTP(responseWriter, request)
	})
}

func fetchDevices(responseWriter http.ResponseWriter, request *http.Request) {
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
	status, interfaceError := wifi.GetNetworkInterfaces()
	if interfaceError != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte("Could not get network status."))
		return
	}

	hostName, hostNameError := os.Hostname()
	if hostNameError != nil {
		log.Println("Failed to fetch host-device name:", hostNameError)
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
	start, err := strconv.Atoi(request.URL.Query().Get("start"))

	if err != nil {
		start = 0
	}

	end, err := strconv.Atoi(request.URL.Query().Get("end"))

	if err != nil {
		end = -1
	}

	logs := wifi.FetchLogs(start, end)
	json, err := json.Marshal(logs)

	if err != nil {
		log.Println("Failed to marshal logs:", err)
		responseWriter.WriteHeader(500)
		responseWriter.Write([]byte("Failed to marshal logs."))
		return
	}
	responseWriter.Write(json)
}
