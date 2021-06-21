package bet

import (
	"context"
	"github.com/clarke94/roulette-service/internal/pkg/bet"
	"github.com/google/go-cmp/cmp"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var db *gorm.DB

func TestStorage_Create(t *testing.T) {
	tests := []struct {
		name    string
		model   bet.Bet
		ctx     context.Context
		want    string
		wantErr bool
	}{
		{
			name: "expect success given valid bet",
			model: bet.Bet{
				ID:       "8117bb87-148c-4fb1-8971-a2d4373b3f19",
				TableID:  "8117bb87-148c-4fb1-8971-a2d4373b3f19",
				Bet:      "foo",
				Type:     "bar",
				Amount:   10,
				Currency: "GBP",
			},
			ctx:     context.Background(),
			want:    "8117bb87-148c-4fb1-8971-a2d4373b3f19",
			wantErr: false,
		},
		{
			name: "expect fail given id already exists",
			model: bet.Bet{
				ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			},
			ctx:     nil,
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(db)

			got, err := s.Create(tt.ctx, tt.model)
			if !cmp.Equal(err != nil, tt.wantErr) {
				t.Fatal(cmp.Diff(err != nil, tt.wantErr))
			}

			if !cmp.Equal(got, tt.want) {
				t.Fatal(cmp.Diff(got, tt.want))
			}
		})
	}
}

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

	err = db.AutoMigrate(Bet{})
	if err != nil {
		log.Fatalf("Could not migrate data: %s", err)
	}

	db.Create([]bet.Bet{
		{
			ID:       "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			TableID:  "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
			Bet:      "foo",
			Type:     "bar",
			Amount:   10,
			Currency: "GBP",
		},
	})

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
