package grpc

import (
	"context"

	"github.com/tomakado/kvevri/internal/pb"
	"github.com/tomakado/kvevri/store"
)

type Server struct {
	pb.UnimplementedStoreServer
	store *store.Store
}

func NewServer(store *store.Store) *Server {
	return &Server{store: store}
}

func (s *Server) Get(_ context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	values := make([]*pb.ValuePair, 0, len(req.Keys))

	for _, key := range req.Keys {
		value, ok := s.store.Get(key)
		if !ok {
			continue
		}

		values = append(
			values,
			&pb.ValuePair{
				Key:   key,
				Value: value,
			},
		)
	}

	return &pb.GetResponse{Values: values}, nil
}

func (s *Server) Set(_ context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	for _, pair := range req.Values {
		s.store.Set(pair.Key, pair.Value)
	}

	return &pb.SetResponse{}, nil
}

func (s *Server) Delete(_ context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	for _, key := range req.Keys {
		s.store.Delete(key)
	}

	return &pb.DeleteResponse{}, nil
}

func (s *Server) Keys(_ context.Context, req *pb.KeysRequest) (*pb.KeysResponse, error) {
	return &pb.KeysResponse{Keys :s.store.Keys()}, nil
}
