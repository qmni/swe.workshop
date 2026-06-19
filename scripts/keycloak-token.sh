#!/usr/bin/env sh
set -eu

KEYCLOAK_URL="${KEYCLOAK_URL:-http://localhost:8880}"

curl -sS -X POST "$KEYCLOAK_URL/realms/swe-workshop/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=swe-workshop-client" \
  -d "username=user" \
  -d "password=p"
