package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/larwef/rpi-docker-test/pkg/enemy"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
	gotestAssert "gotest.tools/v3/assert"
)

type storageMock struct {
	addEnemy    func(ctx context.Context, req *enemy.AddEnemyRequest) (*enemy.AddEnemyResponse, error)
	getEnemy    func(ctx context.Context, req *enemy.GetEnemyRequest) (*enemy.GetEnemyResponse, error)
	updateEnemy func(ctx context.Context, req *enemy.UpdateEnemyRequest) (*enemy.UpdateEnemyResponse, error)
	listEnemies func(ctx context.Context, req *enemy.ListEnemiesRequest) (*enemy.ListEnemiesResponse, error)
}

func (s *storageMock) AddEnemy(ctx context.Context, req *enemy.AddEnemyRequest) (*enemy.AddEnemyResponse, error) {
	return s.addEnemy(ctx, req)
}

func (s *storageMock) GetEnemy(ctx context.Context, req *enemy.GetEnemyRequest) (*enemy.GetEnemyResponse, error) {
	return s.getEnemy(ctx, req)
}

func (s *storageMock) UpdateEnemy(ctx context.Context, req *enemy.UpdateEnemyRequest) (*enemy.UpdateEnemyResponse, error) {
	return s.updateEnemy(ctx, req)
}

func (s *storageMock) ListEnemies(ctx context.Context, req *enemy.ListEnemiesRequest) (*enemy.ListEnemiesResponse, error) {
	return s.listEnemies(ctx, req)
}

func TestServer_AddEnemy(t *testing.T) {
	tests := []struct {
		name    string
		give    *enemy.AddEnemyRequest
		storage *storageMock
		want    *enemy.AddEnemyResponse
		wantErr error
	}{
		{
			name: "Test error from storage",
			give: &enemy.AddEnemyRequest{
				Name:   "Enemy One",
				Email:  "enemy1@bar.com",
				Rating: 1.1,
			},
			storage: &storageMock{
				addEnemy: func(ctx context.Context, req *enemy.AddEnemyRequest) (*enemy.AddEnemyResponse, error) {
					return nil, errors.New("some error")
				},
			},
			wantErr: errors.New("some error"),
		},
		{
			name:    "Test empty name",
			give:    &enemy.AddEnemyRequest{},
			wantErr: status.Error(codes.InvalidArgument, "enemy name can't be empty"),
		},
		{
			name: "Test empty email",
			give: &enemy.AddEnemyRequest{
				Name: "Some Enemy",
			},
			wantErr: status.Error(codes.InvalidArgument, "enemy email can't be empty"),
		},
		{
			name: "Test empty rating",
			give: &enemy.AddEnemyRequest{
				Name:  "Some Enemy",
				Email: "someenemy@bar.com",
			},
			wantErr: status.Error(codes.InvalidArgument, "rating must be > 0"),
		},
		{
			name: "Test successful",
			give: &enemy.AddEnemyRequest{
				Name:   "Enemy One",
				Email:  "enemy1@bar.com",
				Rating: 1.1,
			},
			storage: &storageMock{
				addEnemy: func(ctx context.Context, req *enemy.AddEnemyRequest) (*enemy.AddEnemyResponse, error) {
					return &enemy.AddEnemyResponse{
						Enemy: &enemy.Enemy{
							Id:          "someEnemy",
							Name:        "Enemy One",
							Email:       "enemy1@bar.com",
							Rating:      1.1,
							LastUpdated: timestamppb.New(time.Date(2021, time.December, 30, 20, 57, 35, 0, time.UTC)),
						},
					}, nil
				},
			},
			want: &enemy.AddEnemyResponse{
				Enemy: &enemy.Enemy{
					Id:          "someEnemy",
					Name:        "Enemy One",
					Email:       "enemy1@bar.com",
					Rating:      1.1,
					LastUpdated: timestamppb.New(time.Date(2021, time.December, 30, 20, 57, 35, 0, time.UTC)),
				},
			},
		},
	}

	for _, test := range tests {
		srv := New(test.storage)
		res, err := srv.AddEnemy(context.Background(), test.give)
		gotestAssert.DeepEqual(t, test.want, res, protocmp.Transform())
		assert.Equal(t, test.wantErr, err)
	}
}

func TestServer_GetEnemy(t *testing.T) {
	tests := []struct {
		name    string
		give    *enemy.GetEnemyRequest
		storage *storageMock
		want    *enemy.GetEnemyResponse
		wantErr error
	}{
		{
			name: "Test error from storage",
			give: &enemy.GetEnemyRequest{Id: "enemy1"},
			storage: &storageMock{
				getEnemy: func(ctx context.Context, req *enemy.GetEnemyRequest) (*enemy.GetEnemyResponse, error) {
					return nil, errors.New("some error")
				},
			},
			wantErr: errors.New("some error"),
		},
		{
			name:    "Test empty id",
			give:    &enemy.GetEnemyRequest{},
			wantErr: status.Error(codes.InvalidArgument, "id can't be empty"),
		},
		{
			name: "Test successful",
			give: &enemy.GetEnemyRequest{Id: "enemy1"},
			storage: &storageMock{
				getEnemy: func(ctx context.Context, req *enemy.GetEnemyRequest) (*enemy.GetEnemyResponse, error) {
					return &enemy.GetEnemyResponse{
						Enemy: &enemy.Enemy{
							Id:          "enemy1",
							Name:        "Enemy One",
							Email:       "enemy1@bar.com",
							Rating:      1.1,
							LastUpdated: timestamppb.New(time.Date(2021, time.December, 30, 20, 57, 35, 0, time.UTC)),
						},
					}, nil
				},
			},
			want: &enemy.GetEnemyResponse{
				Enemy: &enemy.Enemy{
					Id:          "enemy1",
					Name:        "Enemy One",
					Email:       "enemy1@bar.com",
					Rating:      1.1,
					LastUpdated: timestamppb.New(time.Date(2021, time.December, 30, 20, 57, 35, 0, time.UTC)),
				},
			},
		},
	}

	for _, test := range tests {
		srv := New(test.storage)
		res, err := srv.GetEnemy(context.Background(), test.give)
		gotestAssert.DeepEqual(t, test.want, res, protocmp.Transform())
		assert.Equal(t, test.wantErr, err)
	}
}

func TestServer_UpdateEnemy(t *testing.T) {
	tests := []struct {
		name    string
		give    *enemy.UpdateEnemyRequest
		storage *storageMock
		want    *enemy.UpdateEnemyResponse
		wantErr error
	}{
		{
			name: "Test error from storage",
			give: &enemy.UpdateEnemyRequest{Id: "enemy1"},
			storage: &storageMock{
				updateEnemy: func(ctx context.Context, req *enemy.UpdateEnemyRequest) (*enemy.UpdateEnemyResponse, error) {
					return nil, errors.New("some error")
				},
			},
			wantErr: errors.New("some error"),
		},
		{
			name:    "Test empty id",
			give:    &enemy.UpdateEnemyRequest{},
			wantErr: status.Error(codes.InvalidArgument, "id can't be empty"),
		},
		{
			name: "Test successful",
			give: &enemy.UpdateEnemyRequest{
				Id:     "enemy1",
				Name:   "New Name",
				Email:  "new@bar.com",
				Rating: 10.0,
			},
			storage: &storageMock{
				updateEnemy: func(ctx context.Context, req *enemy.UpdateEnemyRequest) (*enemy.UpdateEnemyResponse, error) {
					return &enemy.UpdateEnemyResponse{
						Enemy: &enemy.Enemy{
							Id:          "enemy1",
							Name:        "New Name",
							Email:       "new@bar.com",
							Rating:      10.0,
							LastUpdated: timestamppb.New(time.Date(2021, time.December, 30, 20, 57, 35, 0, time.UTC)),
						},
					}, nil
				},
			},
			want: &enemy.UpdateEnemyResponse{
				Enemy: &enemy.Enemy{
					Id:          "enemy1",
					Name:        "New Name",
					Email:       "new@bar.com",
					Rating:      10.0,
					LastUpdated: timestamppb.New(time.Date(2021, time.December, 30, 20, 57, 35, 0, time.UTC)),
				},
			},
		},
	}

	for _, test := range tests {
		srv := New(test.storage)
		res, err := srv.UpdateEnemy(context.Background(), test.give)
		gotestAssert.DeepEqual(t, test.want, res, protocmp.Transform())
		assert.Equal(t, test.wantErr, err)
	}
}

func TestServer_ListEnemies(t *testing.T) {
	tests := []struct {
		name    string
		give    *enemy.ListEnemiesRequest
		storage *storageMock
		want    *enemy.ListEnemiesResponse
		wantErr error
	}{
		{
			name: "Test some error",
			give: &enemy.ListEnemiesRequest{},
			storage: &storageMock{
				listEnemies: func(ctx context.Context, req *enemy.ListEnemiesRequest) (*enemy.ListEnemiesResponse, error) {
					return nil, errors.New("some error")
				},
			},
			wantErr: errors.New("some error"),
		},
		{
			name: "Test successful",
			give: &enemy.ListEnemiesRequest{},
			storage: &storageMock{
				listEnemies: func(ctx context.Context, req *enemy.ListEnemiesRequest) (*enemy.ListEnemiesResponse, error) {
					return &enemy.ListEnemiesResponse{
						Enemies: []*enemy.Enemy{
							{
								Id:          "enemy1",
								Name:        "Enemy One",
								Email:       "enemy1@bar.com",
								Rating:      1.1,
								LastUpdated: timestamppb.New(time.Date(2021, time.December, 30, 20, 57, 35, 0, time.UTC)),
							},
							{
								Id:          "enemy2",
								Name:        "Enemy Two",
								Email:       "enemy2@bar.com",
								Rating:      2.2,
								LastUpdated: timestamppb.New(time.Date(2021, time.December, 29, 19, 34, 10, 0, time.UTC)),
							},
							{
								Id:          "enemy3",
								Name:        "Enemy Three",
								Email:       "enemy3@bar.com",
								Rating:      3.3,
								LastUpdated: timestamppb.New(time.Date(2021, time.December, 28, 15, 15, 12, 0, time.UTC)),
							},
						},
					}, nil
				},
			},
			want: &enemy.ListEnemiesResponse{
				Enemies: []*enemy.Enemy{
					{
						Id:          "enemy1",
						Name:        "Enemy One",
						Email:       "enemy1@bar.com",
						Rating:      1.1,
						LastUpdated: timestamppb.New(time.Date(2021, time.December, 30, 20, 57, 35, 0, time.UTC)),
					},
					{
						Id:          "enemy2",
						Name:        "Enemy Two",
						Email:       "enemy2@bar.com",
						Rating:      2.2,
						LastUpdated: timestamppb.New(time.Date(2021, time.December, 29, 19, 34, 10, 0, time.UTC)),
					},
					{
						Id:          "enemy3",
						Name:        "Enemy Three",
						Email:       "enemy3@bar.com",
						Rating:      3.3,
						LastUpdated: timestamppb.New(time.Date(2021, time.December, 28, 15, 15, 12, 0, time.UTC)),
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		srv := New(test.storage)
		res, err := srv.ListEnemies(context.Background(), test.give)
		gotestAssert.DeepEqual(t, test.want, res, protocmp.Transform())
		assert.Equal(t, test.wantErr, err)
	}
}
