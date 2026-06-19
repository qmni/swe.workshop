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

	duplicateResp, err := client.Post("http://localhost:18080/players", "application/json", bytes.NewReader(createBody))
	if err != nil {
		t.Fatalf("post duplicate player: %v", err)
	}
	defer duplicateResp.Body.Close()
	if duplicateResp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 Conflict for duplicate player, got %d", duplicateResp.StatusCode)
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

	updateBody := []byte(`{"username":"updatedplayer","email":"updated@example.com","level":20,"experience":1000,"playerClass":"ROGUE","status":"ACTIVE"}`)
	updateReq, err := http.NewRequest(http.MethodPut, "http://localhost:18080/players/1", bytes.NewReader(updateBody))
	if err != nil {
		t.Fatalf("create update request: %v", err)
	}
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp, err := client.Do(updateReq)
	if err != nil {
		t.Fatalf("put player: %v", err)
	}
	defer updateResp.Body.Close()
	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for update, got %d", updateResp.StatusCode)
	}

	deleteReq, err := http.NewRequest(http.MethodDelete, "http://localhost:18080/players/1", nil)
	if err != nil {
		t.Fatalf("create delete request: %v", err)
	}
	deleteResp, err := client.Do(deleteReq)
	if err != nil {
		t.Fatalf("delete player: %v", err)
	}
	defer deleteResp.Body.Close()
	if deleteResp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204 No Content for delete, got %d", deleteResp.StatusCode)
	}

	getDeletedResp, err := client.Get("http://localhost:18080/players/1")
	if err != nil {
		t.Fatalf("get deleted player: %v", err)
	}
	defer getDeletedResp.Body.Close()
	if getDeletedResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found for deleted player, got %d", getDeletedResp.StatusCode)
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
