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

echo "Validation demo"
curl -sS -X POST "$BASE_URL/players" \
  -H "Content-Type: application/json" \
  -d '{"username":"","email":"not-an-email","playerClass":"PALADIN"}'
echo
