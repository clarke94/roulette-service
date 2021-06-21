package test

import (
	"github.com/clarke94/roulette-service/storage/bet"
	"github.com/clarke94/roulette-service/storage/table"
	"github.com/clarke94/roulette-service/test/data"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var db *gorm.DB

// TestMain connects to a database with docker to start integration testing.
func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"listen_addresses = '*'",
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: "5432"},
			},
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	err = resource.Expire(60)
	if err != nil {
		log.Fatalf("Could not expire resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err = pool.Retry(func() error {
		db, err = gorm.Open(
			postgres.Open("postgres://postgres:postgres@localhost:5432?sslmode=disable"),
		)
		if err != nil {
			return err
		}

		d, err := db.DB()
		if err != nil {
			return err
		}

		return d.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	err = db.AutoMigrate(table.Table{}, bet.Bet{})
	if err != nil {
		log.Fatalf("Could not migrate data: %s", err)
	}

	db.Create(data.TableData)
	db.Create(data.BetData)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
