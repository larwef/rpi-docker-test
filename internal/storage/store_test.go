package storage

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/larwef/rpi-docker-test/pkg/enemy"
	postgresdocker "github.com/larwef/rpi-docker-test/test/postgres-docker"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
	gotestAssert "gotest.tools/v3/assert"
)

func TestEnemyStore_AddEnemy(t *testing.T) {
	pg, err := postgresdocker.New()
	assert.NoError(t, err)
	defer pg.Shutdown()

	db, err := sql.Open("pgx", pg.ConnectionString())
	assert.NoError(t, err)

	es, err := NewEnemyStore(db)
	assert.NoError(t, err)

	now = func() time.Time { return time.Date(2021, time.December, 31, 14, 59, 5, 0, time.UTC) }
	id = func() string { return "someID" }
	res, err := es.AddEnemy(context.Background(), &enemy.AddEnemyRequest{
		Name:   "Voldemort",
		Email:  "voldemort@bar.com",
		Rating: 10.0,
	})
	assert.NoError(t, err)
	gotestAssert.DeepEqual(t, &enemy.AddEnemyResponse{
		Enemy: &enemy.Enemy{
			Id:          "someID",
			Name:        "Voldemort",
			Email:       "voldemort@bar.com",
			Rating:      10.0,
			LastUpdated: timestamppb.New(time.Date(2021, time.December, 31, 14, 59, 5, 0, time.UTC)),
		},
	}, res, protocmp.Transform())
}

func TestEnemyStore_GetEnemy(t *testing.T) {
	pg, err := postgresdocker.New()
	assert.NoError(t, err)
	defer pg.Shutdown()

	db, err := sql.Open("pgx", pg.ConnectionString())
	assert.NoError(t, err)

	es, err := NewEnemyStore(db)
	assert.NoError(t, err)

	_, err = db.Exec("INSERT INTO enemies (enemy_id, full_name, email, rating, last_updated) VALUES ($1, $2, $3, $4, $5);",
		"enemyID", "Voldemort", "voldemort@bar.com", 9.9, time.Date(2021, time.December, 30, 14, 59, 45, 0, time.UTC))
	assert.NoError(t, err)

	res, err := es.GetEnemy(context.Background(), &enemy.GetEnemyRequest{Id: "enemyID"})
	assert.NoError(t, err)
	gotestAssert.DeepEqual(t, &enemy.GetEnemyResponse{
		Enemy: &enemy.Enemy{
			Id:          "enemyID",
			Name:        "Voldemort",
			Email:       "voldemort@bar.com",
			Rating:      9.9,
			LastUpdated: timestamppb.New(time.Date(2021, time.December, 30, 14, 59, 45, 0, time.UTC)),
		},
	}, res, protocmp.Transform())
}

func TestEnemyStore_UpdateEnemy(t *testing.T) {
	pg, err := postgresdocker.New()
	assert.NoError(t, err)
	defer pg.Shutdown()

	db, err := sql.Open("pgx", pg.ConnectionString())
	assert.NoError(t, err)

	es, err := NewEnemyStore(db)
	assert.NoError(t, err)

	_, err = db.Exec("INSERT INTO enemies (enemy_id, full_name, email, rating, last_updated) VALUES ($1, $2, $3, $4, $5);",
		"enemyID", "Voldemort", "voldemort@bar.com", 9.9, time.Date(2021, time.December, 30, 14, 59, 45, 0, time.UTC))
	assert.NoError(t, err)

	now = func() time.Time { return time.Date(2021, time.December, 31, 14, 59, 5, 0, time.UTC) }
	res, err := es.UpdateEnemy(context.Background(), &enemy.UpdateEnemyRequest{
		Id:     "enemyID",
		Email:  "voldemort@foo.com",
		Rating: 11.0,
	})
	assert.NoError(t, err)
	gotestAssert.DeepEqual(t, &enemy.UpdateEnemyResponse{
		Enemy: &enemy.Enemy{
			Id:          "enemyID",
			Name:        "Voldemort",
			Email:       "voldemort@foo.com",
			Rating:      11.0,
			LastUpdated: timestamppb.New(time.Date(2021, time.December, 31, 14, 59, 5, 0, time.UTC)),
		},
	}, res, protocmp.Transform())
}

func TestEnemyStore_ListEnemy(t *testing.T) {
	pg, err := postgresdocker.New()
	assert.NoError(t, err)
	defer pg.Shutdown()

	db, err := sql.Open("pgx", pg.ConnectionString())
	assert.NoError(t, err)

	es, err := NewEnemyStore(db)
	assert.NoError(t, err)

	q := "INSERT INTO enemies (enemy_id, full_name, email, rating, last_updated) VALUES ($1, $2, $3, $4, $5);"
	_, err = db.Exec(q, "enemy1", "Enemy One", "enemy1@bar.com", 1.1, time.Date(2021, time.December, 1, 11, 59, 5, 0, time.UTC))
	assert.NoError(t, err)
	_, err = db.Exec(q, "enemy2", "Enemy Two", "enemy2@bar.com", 2.2, time.Date(2021, time.December, 2, 12, 59, 5, 0, time.UTC))
	assert.NoError(t, err)
	_, err = db.Exec(q, "enemy3", "Enemy Three", "enemy3@bar.com", 3.3, time.Date(2021, time.December, 3, 13, 59, 5, 0, time.UTC))
	assert.NoError(t, err)
	_, err = db.Exec(q, "enemy4", "Enemy Four", "enemy4@bar.com", 4.4, time.Date(2021, time.December, 4, 14, 59, 5, 0, time.UTC))
	assert.NoError(t, err)
	_, err = db.Exec(q, "enemy5", "Enemy Five", "enemy5@bar.com", 5.5, time.Date(2021, time.December, 5, 15, 59, 5, 0, time.UTC))
	assert.NoError(t, err)

	res, err := es.ListEnemy(context.Background(), &enemy.ListEnemiesRequest{})
	assert.NoError(t, err)
	gotestAssert.DeepEqual(t, &enemy.ListEnemiesResponse{
		Enemies: []*enemy.Enemy{
			{
				Id:          "enemy1",
				Name:        "Enemy One",
				Email:       "enemy1@bar.com",
				Rating:      1.1,
				LastUpdated: timestamppb.New(time.Date(2021, time.December, 1, 11, 59, 5, 0, time.UTC)),
			},
			{
				Id:          "enemy2",
				Name:        "Enemy Two",
				Email:       "enemy2@bar.com",
				Rating:      2.2,
				LastUpdated: timestamppb.New(time.Date(2021, time.December, 2, 12, 59, 5, 0, time.UTC)),
			},
			{
				Id:          "enemy3",
				Name:        "Enemy Three",
				Email:       "enemy3@bar.com",
				Rating:      3.3,
				LastUpdated: timestamppb.New(time.Date(2021, time.December, 3, 13, 59, 5, 0, time.UTC)),
			},
			{
				Id:          "enemy4",
				Name:        "Enemy Four",
				Email:       "enemy4@bar.com",
				Rating:      4.4,
				LastUpdated: timestamppb.New(time.Date(2021, time.December, 4, 14, 59, 5, 0, time.UTC)),
			},
			{
				Id:          "enemy5",
				Name:        "Enemy Five",
				Email:       "enemy5@bar.com",
				Rating:      5.5,
				LastUpdated: timestamppb.New(time.Date(2021, time.December, 5, 15, 59, 5, 0, time.UTC)),
			},
		},
	}, res, protocmp.Transform())
}
