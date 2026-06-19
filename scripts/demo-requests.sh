#!/usr/bin/env sh
set -eu

BASE_URL="${BASE_URL:-http://localhost:8080}"

echo "Health check"
curl -sS "$BASE_URL/health"
echo

echo "List players before create"
curl -sS "$BASE_URL/players"
echo

echo "Create player"
curl -sS -X POST "$BASE_URL/players" \
  -H "Content-Type: application/json" \
  -d '{"username":"demo_player","email":"demo_player@example.com","level":10,"experience":500,"playerClass":"MAGE"}'
echo

echo "List players after create"
curl -sS "$BASE_URL/players"
echo

echo "Update player 1"
curl -sS -X PUT "$BASE_URL/players/1" \
  -H "Content-Type: application/json" \
  -d '{"username":"updated_demo_player","email":"updated_demo_player@example.com","level":20,"experience":1000,"playerClass":"ROGUE","status":"ACTIVE"}'
echo

echo "Validation demo"
curl -sS -X POST "$BASE_URL/players" \
  -H "Content-Type: application/json" \
  -d '{"username":"","email":"not-an-email","playerClass":"PALADIN"}'
echo

echo "Delete player 1"
curl -sS -X DELETE "$BASE_URL/players/1" -w "HTTP %{http_code}"
echo
