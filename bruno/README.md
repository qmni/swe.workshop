# Bruno Tests

1. Start the application:

```bash
make run
```

2. Open this folder in Bruno:

```text
/Users/hien/Desktop/swe-workshop/bruno
```

3. Select the `local` environment (top right corner in Bruno).

4. If Keycloak auth is enabled (`KEYCLOAK_JWKS_URL` set), run the token
   request first. It stores the access token for all player requests:

```text
00-keycloak-token
```

The token expires after a few minutes. Re-run `00-keycloak-token` when a
request returns HTTP 401.

5. Run the requests in this order:

```text
00-keycloak-token
01-health
02-list-players
03-create-player
02-list-players
05-get-created-player
06-update-created-player
05-get-created-player
04-validation-error
07-delete-created-player
```

Expected results:

* `00-keycloak-token`: HTTP 200 and an `access_token` stored as `accessToken`
* `01-health`: HTTP 200 and `{"status":"ok"}`
* `02-list-players`: HTTP 200 and a JSON list
* `03-create-player`: HTTP 201, a new player and stored `playerId`
* `05-get-created-player`: HTTP 200 and the player created in step 03
* `06-update-created-player`: HTTP 200 and updated player data
* `04-validation-error`: HTTP 400 and validation details
* `07-delete-created-player`: HTTP 204 and deletes the player created in step 03
