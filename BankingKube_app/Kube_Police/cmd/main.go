package main

import (
	"log"
	"net/http"

	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/server"
)

func main() {
	log.Println("Starting Webhook Server...")

	certFile := getSecret("WEBHOOK_CERT_SECRET_NAME", "eu-central-1")
	keyFile := getSecret("WEBHOOK_KEY_SECRET_NAME", "eu-central-1")

	if certFile == "" || keyFile == "" {
		log.Fatal("Failed to retrieve certificate or key file from Secrets Manager")
	}

	// Register the admission handler
	http.HandleFunc("/validate", admission.HandleAdmissionRequest)

	// Create and start server
	srv := server.NewServer(certFile, keyFile)
	server.StartServer(srv)
}
