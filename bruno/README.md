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
05-get-player-1
06-update-player-1
05-get-player-1
04-validation-error
07-delete-player-1
```

Expected results:

* `01-health`: HTTP 200 and `{"status":"ok"}`
* `02-list-players`: HTTP 200 and a JSON list
* `03-create-player`: HTTP 201 and a new player
* `05-get-player-1`: HTTP 200 and player with ID 1
* `06-update-player-1`: HTTP 200 and updated player data
* `04-validation-error`: HTTP 400 and validation details
* `07-delete-player-1`: HTTP 204

If your created player has another ID, replace `/players/1` in requests 05, 06 and 07 with that ID.
