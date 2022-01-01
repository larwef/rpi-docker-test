package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/larwef/rpi-docker-test/pkg/enemy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	serviceURL = "10.0.0.18:8080"
)

func main() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial(serviceURL, opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := enemy.NewEnemyServiceClient(conn)

	// for i := 0; i < 10; i++ {
	// 	addEnemy(client, &enemy.AddEnemyRequest{
	// 		Name:   fmt.Sprintf("Enemy %d", i),
	// 		Email:  fmt.Sprintf("enemy%d@bar.com", i),
	// 		Rating: float32(i + 1),
	// 	})
	// }

	listEnemies(client)
}

func addEnemy(client enemy.EnemyServiceClient, req *enemy.AddEnemyRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := client.AddEnemy(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully added:\n%s\n", protojson.Format(res))
}

func listEnemies(client enemy.EnemyServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := client.ListEnemies(ctx, &enemy.ListEnemiesRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(protojson.Format(res))
}

func updateEnemy(client enemy.EnemyServiceClient, req *enemy.UpdateEnemyRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := client.UpdateEnemy(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(protojson.Format(res))
}
