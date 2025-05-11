package server

import (
	cerr "VKTest/pkg/customErrors"
	pb "VKTest/pkg/grpc/pb/subpubservice"
	"VKTest/pkg/pubsub"
	"context"
	"errors"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc/codes"
)

type Config struct {
	Port string `yaml:"port"`
}

type server struct {
	pb.UnimplementedPubSubServer
	hub pubsub.SubPub
}

type Server struct {
	PBServer *grpc.Server
	Listener *net.Listener
	cfg      Config
}

func New(cfg Config, hub pubsub.SubPub) (*Server, error) {
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return nil, err
	}
	s := grpc.NewServer()
	pb.RegisterPubSubServer(s, &server{})
	log.Printf("Auth server listening at %v\n", lis.Addr())
	return &Server{PBServer: s, Listener: &lis, cfg: cfg}, nil
}

func (s *server) Subscribe(req *pb.SubscribeRequest, sub pb.PubSub_SubscribeServer) error {
	subscriber, err := s.hub.Subscribe(req.Key, func(msg interface{}) {})
	if err != nil {
		if errors.Is(err, cerr.ErrSubClosed) {
			return status.Error(codes.Unavailable, err.Error())
		}
		return status.Error(codes.Internal, err.Error())
	}
	msgs := subscriber.GetMessages()
	for msg := range msgs {
		if err := sub.SendMsg(pb.Event{Data: msg.(string)}); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
	return nil
}

func (s *server) Publish(ctx context.Context, req *pb.PublishRequest) (*emptypb.Empty, error) {
	err := s.hub.Publish(req.Key, req.Data)
	if err != nil {
		if errors.Is(err, cerr.ErrSubClosed) {
			return nil, status.Error(codes.Unavailable, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
