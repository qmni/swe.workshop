package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/qmni/swe.workshop/internal/app"
	"github.com/qmni/swe.workshop/internal/database"
)

func TestPlayerAPI(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION") != "1" {
		t.Skip("set RUN_INTEGRATION=1 to run integration tests")
	}

	db, err := database.Open(database.ConfigFromEnv())
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("database handle: %v", err)
	}
	defer sqlDB.Close()
	if _, err := sqlDB.Exec(`TRUNCATE player, guild RESTART IDENTITY CASCADE`); err != nil {
		t.Fatalf("truncate player tables: %v", err)
	}

	srv := app.New(db)
	go func() {
		_ = srv.Listen(":18080")
	}()
	defer func() {
		_ = srv.Shutdown()
	}()

	client := http.Client{Timeout: 5 * time.Second}
	waitForHealth(t, client)

	createBody := []byte(`{"username":"testplayer","email":"testplayer@example.com","level":10,"experience":500,"playerClass":"MAGE"}`)
	createResp, err := client.Post("http://localhost:18080/players", "application/json", bytes.NewReader(createBody))
	if err != nil {
		t.Fatalf("post player: %v", err)
	}
	defer createResp.Body.Close()
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", createResp.StatusCode)
	}

	listResp, err := client.Get("http://localhost:18080/players")
	if err != nil {
		t.Fatalf("get players: %v", err)
	}
	defer listResp.Body.Close()
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", listResp.StatusCode)
	}

	var players []map[string]any
	if err := json.NewDecoder(listResp.Body).Decode(&players); err != nil {
		t.Fatalf("decode players: %v", err)
	}
	if len(players) != 1 || players[0]["username"] != "testplayer" {
		t.Fatalf("unexpected players response: %#v", players)
	}
}

func waitForHealth(t *testing.T, client http.Client) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := client.Get("http://localhost:18080/health")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatal("server did not become healthy")
}
