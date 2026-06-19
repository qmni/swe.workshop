#!/usr/bin/env sh
set -eu

BASE_URL="${BASE_URL:-http://localhost:8080}"
RUN_ID="$(date +%s)"
USERNAME="demo_player_${RUN_ID}"
EMAIL="demo_player_${RUN_ID}@example.com"

echo "Health check"
curl -sS "$BASE_URL/health"
echo

echo "List players before create"
curl -sS "$BASE_URL/players"
echo

echo "Create player"
CREATE_RESPONSE="$(curl -sS -X POST "$BASE_URL/players" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"${USERNAME}\",\"email\":\"${EMAIL}\",\"level\":10,\"experience\":500,\"playerClass\":\"MAGE\"}")"
echo "$CREATE_RESPONSE"
PLAYER_ID="$(printf '%s' "$CREATE_RESPONSE" | sed -n 's/.*"id":\([0-9][0-9]*\).*/\1/p')"
if [ -z "$PLAYER_ID" ]; then
  echo "Could not read player id from create response" >&2
  exit 1
fi

echo "List players after create"
curl -sS "$BASE_URL/players"
echo

echo "Update player $PLAYER_ID"
curl -sS -X PUT "$BASE_URL/players/$PLAYER_ID" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"updated_${USERNAME}\",\"email\":\"updated_${EMAIL}\",\"level\":20,\"experience\":1000,\"playerClass\":\"ROGUE\",\"status\":\"ACTIVE\"}"
echo

echo "Validation demo"
curl -sS -X POST "$BASE_URL/players" \
  -H "Content-Type: application/json" \
  -d '{"username":"","email":"not-an-email","playerClass":"PALADIN"}'
echo

echo "Delete player $PLAYER_ID"
curl -sS -X DELETE "$BASE_URL/players/$PLAYER_ID" -w "HTTP %{http_code}"
echo
