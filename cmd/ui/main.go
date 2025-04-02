package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"ui-service/internal/server"
	pb "ui-service/proto/ui"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func handleGetRecommendation(w http.ResponseWriter, r *http.Request) {
	deviceID, recommendation := server.GetLatestRecommendation()

	response := struct {
		DeviceID       string `json:"deviceId"`
		Recommendation string `json:"recommendation"`
	}{
		DeviceID:       deviceID,
		Recommendation: recommendation,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func startGRPCServer(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	uiServer := server.NewUIServer()
	pb.RegisterUIServiceServer(grpcServer, uiServer)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on port %d", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func main() {
	grpcPort := flag.Int("grpc_port", 50053, "Port for gRPC server")
	httpPort := flag.Int("http_port", 8080, "Port for HTTP server")
	flag.Parse()

	if _, err := os.Stat("templates"); os.IsNotExist(err) {
		os.MkdirAll("templates", 0o755)
	}

	if _, err := os.Stat("public"); os.IsNotExist(err) {
		os.MkdirAll("public", 0o755)
	}

	go startGRPCServer(*grpcPort)

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/recommendation", handleGetRecommendation)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Printf("HTTP server listening on port %d", *httpPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
