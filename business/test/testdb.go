package test

import (
	"context"
	"fmt"
	"log"
	"sales-api/business/data/dbmigrate"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabase struct {
	DB        *sqlx.DB
	container testcontainers.Container
}

func (d *TestDatabase) TearDown() {
	d.container.Terminate(context.Background())
}

func SetUpTestDatabase() *TestDatabase {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	container, db, err := createPostgresContainer(ctx)
	if err != nil {
		log.Fatalf("createPostgresContainer failed %v", err)
	}
	return &TestDatabase{
		container: container,
		DB:        db,
	}
}

// ====================================================================================

func createPostgresContainer(ctx context.Context) (container testcontainers.Container, db *sqlx.DB, err error) {
	containerPort := "5432"
	req := testcontainers.ContainerRequest{
		Image: "postgres:16.4",
		Env: map[string]string{
			"POSTGRES_USER":     "root",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "testdb",
		},
		ExposedPorts: []string{containerPort + "/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
			wait.ForListeningPort(nat.Port(containerPort)),
		).WithDeadline(5 * time.Minute),
	}

	container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req, Started: true,
	})

	if err != nil {
		return container, nil, fmt.Errorf("failed to start container: %v", err)
	}
	host, err := container.Host(ctx)
	if err != nil {
		return container, nil, fmt.Errorf("failed to get container external host: %v", err)
	}
	p, err := container.MappedPort(ctx, nat.Port(containerPort))
	if err != nil {
		return container, nil, fmt.Errorf("failed to get container external port: %v", err)
	}
	log.Println("container ready and running at port: ", p.Port())
	connStr := fmt.Sprintf("postgresql://root:password@%s:%s/testdb?sslmode=disable", host, p.Port())

	db, err = sqlx.Open("pgx", connStr)
	if err != nil {
		return container, db, fmt.Errorf("failed to establish database connection: %v", err)
	}
	defer db.Close()
	err = dbmigrate.Migration(ctx, db)
	return
}
