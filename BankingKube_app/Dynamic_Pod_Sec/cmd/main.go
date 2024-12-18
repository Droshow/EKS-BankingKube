package main

import (
	"log"
	"net/http"

	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/server"
)

func main() {
	log.Println("Starting Webhook Server...")

	certFile := "/tls/tls.crt"
	keyFile := "/tls/tls.key"

	// Register the admission handlers
	http.HandleFunc("/validate/context", admission.HandleAdmissionRequest)
	http.HandleFunc("/validate/volumes", admission.HandleAdmissionRequest)
	http.HandleFunc("/validate/network", admission.HandleAdmissionRequest)
	http.HandleFunc("/validate/api", admission.HandleAdmissionRequest)
	http.HandleFunc("/validate/image", admission.HandleAdmissionRequest)
	http.HandleFunc("/validate/rbac", admission.HandleAdmissionRequest)
	http.HandleFunc("/validate/resources", admission.HandleAdmissionRequest)
	http.HandleFunc("/mutate/pod", admission.HandleAdmissionRequest)

	// Create and start the server
	srv := server.NewServer(certFile, keyFile)
	server.StartServer(srv)
}
