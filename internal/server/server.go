package server

import (
	"context"
	"log"
	"sync"

	pb "ui-service/proto/ui"
)

var (
	latestRecommendation string
	latestDeviceID       string
	mutex                sync.RWMutex
)

type UIServer struct {
	pb.UnimplementedUIServiceServer
}

func NewUIServer() *UIServer {
	return &UIServer{}
}

func (s *UIServer) DisplayRecommendation(ctx context.Context, req *pb.RecommendationRequest) (*pb.RecommendationResponse, error) {
	log.Printf("Received recommendation for device %s", req.DeviceId)

	mutex.Lock()
	latestDeviceID = req.DeviceId
	latestRecommendation = req.RecommendationData
	mutex.Unlock()

	return &pb.RecommendationResponse{Received: true}, nil
}

func GetLatestRecommendation() (string, string) {
	mutex.RLock()
	defer mutex.RUnlock()
	return latestDeviceID, latestRecommendation
}
