package server

import (
	"context"

	"github.com/larwef/rpi-docker-test/pkg/enemy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Storage interface {
	AddEnemy(ctx context.Context, req *enemy.AddEnemyRequest) (*enemy.AddEnemyResponse, error)
	GetEnemy(ctx context.Context, req *enemy.GetEnemyRequest) (*enemy.GetEnemyResponse, error)
	UpdateEnemy(ctx context.Context, req *enemy.UpdateEnemyRequest) (*enemy.UpdateEnemyResponse, error)
	ListEnemies(ctx context.Context, req *enemy.ListEnemiesRequest) (*enemy.ListEnemiesResponse, error)
}

type Server struct {
	enemy.UnimplementedEnemyServiceServer
	storage Storage
}

func New(s Storage) *Server {
	return &Server{storage: s}
}

func (s *Server) AddEnemy(ctx context.Context, req *enemy.AddEnemyRequest) (*enemy.AddEnemyResponse, error) {
	switch {
	case req.GetName() == "":
		return nil, status.Error(codes.InvalidArgument, "enemy name can't be empty")
	case req.GetEmail() == "":
		return nil, status.Error(codes.InvalidArgument, "enemy email can't be empty")
	case req.GetRating() == 0.0:
		return nil, status.Error(codes.InvalidArgument, "rating must be > 0")
	}
	return s.storage.AddEnemy(ctx, req)
}

func (s *Server) GetEnemy(ctx context.Context, req *enemy.GetEnemyRequest) (*enemy.GetEnemyResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id can't be empty")
	}
	return s.storage.GetEnemy(ctx, req)
}

func (s *Server) UpdateEnemy(ctx context.Context, req *enemy.UpdateEnemyRequest) (*enemy.UpdateEnemyResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id can't be empty")
	}
	return s.storage.UpdateEnemy(ctx, req)
}

func (s *Server) ListEnemies(ctx context.Context, req *enemy.ListEnemiesRequest) (*enemy.ListEnemiesResponse, error) {
	return s.storage.ListEnemies(ctx, req)
}
