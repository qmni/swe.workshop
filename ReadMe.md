# Programmierworkshop am 19.6.2026

## Namen

Ngoc Hien Do 94148
Quang Nguyen 90863

## Link zum Git-Repository

https://github.com/qmni/swe.workshop

## KI-Werkzeuge

### Agenten

OpenAI Codex

### Chat-URLs, z.B. https://chatgpt.com

https://chatgpt.com

## Frameworks und Bibliotheken

## Projektstruktur

Die aktive Go-Anwendung liegt in diesen Ordnern:

* `cmd/server`: Einstiegspunkt der API
* `internal/app`: Fiber-Server und Routing
* `internal/httpapi`: HTTP-Handler und Request-Tests
* `internal/database`: DB-Konfiguration, Verbindung und Migrationen
* `internal/model`: Datenmodelle
* `integration`: End-to-End-Test gegen die Testdatenbank
* `scripts`: Hilfsskripte fuer lokale Checks, Demo-Requests und Keycloak-Token

Ergaenzende bzw. uebernommene Artefakte liegen getrennt davon:

* `src`: Konfigurations- und Ressourcenstruktur des Teampartner-Projekts als fachliche Referenz
* `prisma`: Prisma-Schema als Referenz fuer die Datenbankstruktur
* `extras/compose`: Zusaetzliche Compose-Stacks fuer Infrastruktur wie Keycloak, Postgres, Monitoring und SonarQube
* `extras/doc`: Architektur- und Projektdokumentation
* `bruno`: API-Requests fuer manuelle Tests mit Bruno
* `requests.http`: Alternative HTTP-Beispielrequests fuer IDEs

Build-Artefakte werden nicht versioniert. Das lokale Verzeichnis `bin/` ist nur fuer erzeugte Binaries aus `make build` gedacht.

### REST-Schnittstelle (Lesen und Neuanlegen)

Go mit Fiber (`github.com/gofiber/fiber/v2`)

Implementierte Endpunkte:

* `GET /health`
* `GET /players`
* `GET /players/:id`
* `POST /players`
* Zusatz: `PUT /players/:id`
* Zusatz: `DELETE /players/:id`

Beispiel zum Neuanlegen:

```bash
curl -X POST http://localhost:8080/players \
  -H "Content-Type: application/json" \
  -d '{"username":"testplayer","email":"testplayer@example.com","level":10,"experience":500,"playerClass":"MAGE"}'
```

Beispiel zum Lesen:

```bash
curl http://localhost:8080/players
```

### Validierung (nur Neuanlegen)

`github.com/go-playground/validator/v10`

Validiert wird beim `POST /players`:

* `username`: Pflichtfeld, 3 bis 60 Zeichen
* `email`: Pflichtfeld, gueltige E-Mail-Adresse, maximal 120 Zeichen
* `level`: optional, 1 bis 100
* `experience`: optional, mindestens 0
* `playerClass`: Pflichtfeld, einer von `WARRIOR`, `MAGE`, `ROGUE`, `PRIEST`, `HUNTER`
* `guildId`: optional, positive ID

### OR-Mapping (für PostgreSQL)

GORM:

* `gorm.io/gorm`
* `gorm.io/driver/postgres`

Die Tabellen `player` und `guild` sowie die PostgreSQL-Enums `PlayerClass` und `PlayerStatus` werden beim Start angelegt, falls sie noch nicht existieren. Die Struktur orientiert sich an den Partner-Dateien `prisma/schema.prisma` und `src/config/resources/postgresql/create-table.sql`.

Die Datenbankverbindung ist ueber Umgebungsvariablen konfigurierbar, damit auch der DB-Server des Teampartners verwendet werden kann. Beispielwerte stehen in `.env.example`.

### Optional: OIDC mit Keycloak

Keycloak ist optional und wird in der Go-API nicht verpflichtend erzwungen. Die Infrastruktur ist vorbereitet:

* `docker-compose.keycloak.yml`
* `keycloak/swe-workshop-realm.json`
* `scripts/keycloak-token.sh`

Starten:

```bash
make keycloak-run
```

Admin-Konsole:

```text
http://localhost:8880
```

Admin-Login:

```text
Username: tmp
Password: p
```

Importierter Realm:

```text
swe-workshop
```

Importierter Client:

```text
swe-workshop-client
```

Demo-User:

```text
user / p
admin / p
```

Token abrufen:

```bash
make keycloak-token
```

### Einfacher Integrationstest

Integrationstest in `integration/players_test.go`.

Der Test startet die API gegen eine PostgreSQL-Testdatenbank, legt einen Player per REST an und liest die Playerliste wieder aus.

Ausführen:

```bash
docker compose -f docker-compose.test.yml up --build --abort-on-container-exit
```

Anwendung lokal starten:

```bash
docker compose up --build
```

Danach ist die API unter `http://localhost:8080` erreichbar.

Alternativ koennen die haeufigen Befehle ueber das `Makefile` gestartet werden:

```bash
make test
make build
make run
make integration-test
```

Lokale Pruefung ohne Docker:

```bash
./scripts/check-local.sh
```

Demo-Requests nach dem Start der Anwendung:

```bash
./scripts/demo-requests.sh
```

Zusaetzlich liegen Beispiel-Requests in `requests.http`, damit die API direkt aus einer IDE oder mit einem REST-Client getestet werden kann.

Fuer manuelle Tests mit Bruno liegt eine Collection im Ordner `bruno/`.

### Demo-Ablauf

1. Anwendung starten: `docker compose up --build`
2. Health-Check aufrufen: `GET /health`
3. Player anlegen: `POST /players`
4. Playerliste lesen: `GET /players`
5. Validierung testen, indem ein Player ohne Username, mit falscher E-Mail oder ungueltiger `playerClass` gesendet wird
6. Optional: Player aktualisieren mit `PUT /players/:id`
7. Optional: Player loeschen mit `DELETE /players/:id`

Der Demo-Ablauf wurde lokal mit Docker ausgefuehrt. Health-Check, Lesen, Neuanlegen und Validierungsfehler waren erfolgreich. Der Integrationstest mit Docker wurde ebenfalls erfolgreich ausgefuehrt.

## Prompts/Requests an KI-Agent/en

Initialer Request:

> swe workshop so heisst das projekt was ich jetzt machen muss ich hab 4h stunden zeit ... richte alles für mich ein und das ist der github https://github.com/qmni/swe.workshop

Umgesetzte KI-Aufgaben:

* Projektstruktur für Go erstellt
* REST-API mit Lesen und Neuanlegen umgesetzt
* Validierung für Neuanlegen ergänzt
* PostgreSQL-Anbindung mit GORM umgesetzt
* Docker-Compose für Anwendung und Datenbank erstellt
* Einfachen Integrationstest vorbereitet
* Diese `ReadMe.md` für die Abgabe ausgefüllt
