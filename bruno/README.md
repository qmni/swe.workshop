# Bruno Tests

1. Start the application:

```bash
make run
```

2. Open this folder in Bruno:

```text
/Users/hien/Desktop/swe-workshop/bruno
```

3. Run the requests in this order:

```text
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

* `01-health`: HTTP 200 and `{"status":"ok"}`
* `02-list-players`: HTTP 200 and a JSON list
* `03-create-player`: HTTP 201, a new player and stored `playerId`
* `05-get-created-player`: HTTP 200 and the player created in step 03
* `06-update-created-player`: HTTP 200 and updated player data
* `04-validation-error`: HTTP 400 and validation details
* `07-delete-created-player`: HTTP 204 and deletes the player created in step 03
