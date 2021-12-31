package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/larwef/rpi-docker-test/pkg/enemy"
	"github.com/rs/xid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Simplify testing.
var now = func() time.Time {
	return time.Now()
}
var id = func() string {
	return xid.New().String()
}

type EnemyStore struct {
	queries *Queries
}

func NewEnemyStore(db *sql.DB) (*EnemyStore, error) {
	if err := migrateUp(db); err != nil {
		return nil, err
	}
	return &EnemyStore{
		queries: New(db),
	}, nil
}

func (e *EnemyStore) AddEnemy(ctx context.Context, req *enemy.AddEnemyRequest) (*enemy.AddEnemyResponse, error) {
	enmy, err := e.queries.AddEnemy(ctx, AddEnemyParams{
		EnemyID:     id(),
		FullName:    req.GetName(),
		Email:       req.GetEmail(),
		Rating:      req.GetRating(),
		LastUpdated: now(),
	})
	if err != nil {
		return nil, err
	}
	return &enemy.AddEnemyResponse{
		Enemy: &enemy.Enemy{
			Id:          enmy.EnemyID,
			Name:        enmy.FullName,
			Email:       enmy.Email,
			Rating:      enmy.Rating,
			LastUpdated: timestamppb.New(enmy.LastUpdated),
		},
	}, nil
}

func (e *EnemyStore) GetEnemy(ctx context.Context, req *enemy.GetEnemyRequest) (*enemy.GetEnemyResponse, error) {
	enmy, err := e.queries.GetEnemy(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &enemy.GetEnemyResponse{
		Enemy: &enemy.Enemy{
			Id:          enmy.EnemyID,
			Name:        enmy.FullName,
			Email:       enmy.Email,
			Rating:      enmy.Rating,
			LastUpdated: timestamppb.New(enmy.LastUpdated),
		},
	}, nil
}

func (e *EnemyStore) UpdateEnemy(ctx context.Context, req *enemy.UpdateEnemyRequest) (*enemy.UpdateEnemyResponse, error) {
	enmy, err := e.queries.UpdateEnemy(ctx, UpdateEnemyParams{
		FullName:    req.Name,
		Email:       req.Email,
		Rating:      req.Rating,
		LastUpdated: now(),
		EnemyID:     req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &enemy.UpdateEnemyResponse{
		Enemy: &enemy.Enemy{
			Id:          enmy.EnemyID,
			Name:        enmy.FullName,
			Email:       enmy.Email,
			Rating:      enmy.Rating,
			LastUpdated: timestamppb.New(enmy.LastUpdated),
		},
	}, nil
}

func (e *EnemyStore) ListEnemy(ctx context.Context, req *enemy.ListEnemiesRequest) (*enemy.ListEnemiesResponse, error) {
	enemies, err := e.queries.ListEnemies(ctx)
	if err != nil {
		return nil, err
	}
	var res []*enemy.Enemy
	for _, enmy := range enemies {
		res = append(res, &enemy.Enemy{
			Id:          enmy.EnemyID,
			Name:        enmy.FullName,
			Email:       enmy.Email,
			Rating:      enmy.Rating,
			LastUpdated: timestamppb.New(enmy.LastUpdated),
		})
	}
	return &enemy.ListEnemiesResponse{
		Enemies: res,
	}, nil
}
