package main

import (
	"context"
	pb "github.com/CuteHat/SmartGabby/generated/discovery-service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

var (
	peers = make(map[string]*net.TCPAddr)
)

type DiscoveryServiceServerImpl struct {
	pb.UnimplementedDiscoveryServiceServer
}

func (server *DiscoveryServiceServerImpl) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.DataLoss, "Could not extract peer info")
	}

	existingPeer := peers[request.PreferredName]
	if existingPeer != nil {
		return nil, status.Error(codes.InvalidArgument, "Peer with given name is already registered")
	}

	tcpAddr, ok := peerInfo.Addr.(*net.TCPAddr)
	if !ok {
		return nil, status.Error(codes.DataLoss, "Could not extract peer address")
	}
	peers[request.PreferredName] = tcpAddr
	return &pb.RegisterResponse{Message: "You have been successfully registered"}, nil
}

func (server *DiscoveryServiceServerImpl) FetchPeers(ctx context.Context, req *emptypb.Empty) (*pb.FetchPeersResponse, error) {
	peerInfos := make([]*pb.PeerInfo, len(peers))

	for preferredName, peer := range peers {
		peerInfos = append(peerInfos, &pb.PeerInfo{PeerPreferredName: preferredName, PeerIp: peer.IP.String()})
	}

	return &pb.FetchPeersResponse{Peers: peerInfos}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error listening on port 8080", err)
	}

	server := grpc.NewServer()
	pb.RegisterDiscoveryServiceServer(server, &DiscoveryServiceServerImpl{})

	err = server.Serve(listener)
	if err != nil {
		log.Fatal("Error starting grpc server on 8080", err)
	}
}
