package rest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

const (
	HOST = "localhost"
)

func StartServer(port int) {
	certificateManager := getCertificateManager()
	handler := getRouter()

	openHttp(certificateManager, handler, port)
}

func getCertificateManager() *autocert.Manager {
	hostPolicy := func(ctx context.Context, host string) error {
		if host == HOST {
			return nil
		}
		return fmt.Errorf("acme/autocert: only %s host is allowed", HOST)
	}
	return &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache("."),
	}
}

func openHttp(certificateManager *autocert.Manager, handler http.Handler, port int) {
	fmt.Printf("Running REST server on port: %d\r\n", port)
	httpServer := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler,
		Addr:         ":" + strconv.Itoa(port),
	}
	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
