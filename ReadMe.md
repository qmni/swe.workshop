# Programmierworkshop am 19.6.2026

## Namen

Ngoc Hien Do 94148
Quang Nguyen 90863

## Link zum Git-Repository

[https://github.com/qmni/swe.workshop](https://github.com/qmni/swe.workshop)

## KI-Werkzeuge

### Agenten

* Claude Code (Anthropic)
* OpenAI Codex

### Chat-URLs, z.B. ChatGPT

[https://chatgpt.com](https://chatgpt.com)
<https://claude.ai>

## Voraussetzungen

Fuer den lokalen Start werden benoetigt:

* Go 1.22 oder neuer
* Docker und Docker Compose
* Optional: `make` fuer die kuerzeren Hilfsbefehle

## Frameworks und Bibliotheken

Backend und API:

* Go 1.22
* Fiber (`github.com/gofiber/fiber/v2`) als Web-Framework
* Validator (`github.com/go-playground/validator/v10`) fuer Request-Validierung

Datenbank und Persistenz:

* PostgreSQL
* GORM (`gorm.io/gorm`) als ORM
* PostgreSQL-Treiber fuer GORM (`gorm.io/driver/postgres`)

Authentifizierung (optional):

* Keycloak als OIDC-Provider
* JWT-Pruefung mit `github.com/golang-jwt/jwt/v5`
* JWKS-Abruf mit `github.com/MicahParks/keyfunc/v3`

Tests und Entwicklungswerkzeuge:

* Go Testing (`go test`) fuer Unit- und Integrationstests
* Docker / Docker Compose fuer lokale Laufzeit und Testumgebung
* Bruno und `requests.http` fuer manuelle API-Tests

## Projektstruktur

Die aktive Go-Anwendung liegt in diesen Ordnern:

* `cmd/server`: Einstiegspunkt der API
* `internal/app`: Fiber-Server und Routing
* `internal/httpapi`: HTTP-Handler und Request-Tests
* `internal/middleware`: Keycloak-JWT-Pruefung (optionale Absicherung der Player-Routen)
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

### REST-Schnittstelle (CRUD)

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

Die Antwort von `GET /players` ist aktuell eine JSON-Liste von Player-Objekten.

### Validierung (Neuanlegen und Aktualisieren)

`github.com/go-playground/validator/v10`

Validiert wird beim `POST /players`:

* `username`: Pflichtfeld, 3 bis 60 Zeichen
* `email`: Pflichtfeld, gueltige E-Mail-Adresse, maximal 120 Zeichen
* `level`: optional, 1 bis 100
* `experience`: optional, mindestens 0
* `playerClass`: Pflichtfeld, einer von `WARRIOR`, `MAGE`, `ROGUE`, `PRIEST`, `HUNTER`
* `guildId`: optional, positive ID

Zusätzlich wird beim `PUT /players/:id` validiert:

* `username`: Pflichtfeld, 3 bis 60 Zeichen
* `email`: Pflichtfeld, gueltige E-Mail-Adresse, maximal 120 Zeichen
* `level`: Pflichtfeld, 1 bis 100
* `experience`: mindestens 0
* `playerClass`: Pflichtfeld, einer von `WARRIOR`, `MAGE`, `ROGUE`, `PRIEST`, `HUNTER`
* `status`: Pflichtfeld, einer von `ACTIVE`, `BANNED`, `DELETED`
* `guildId`: optional, positive ID

Bei doppeltem `username` oder `email` wird ein `409 Conflict` zurueckgegeben.

### OR-Mapping (für PostgreSQL)

GORM:

* `gorm.io/gorm`
* `gorm.io/driver/postgres`

Die Tabellen `player` und `guild` sowie die PostgreSQL-Enums `PlayerClass` und `PlayerStatus` werden beim Start angelegt, falls sie noch nicht existieren. Die Struktur orientiert sich an den Partner-Dateien `prisma/schema.prisma` und `src/config/resources/postgresql/create-table.sql`.

Die Datenbankverbindung ist ueber Umgebungsvariablen konfigurierbar, damit auch der DB-Server des Teampartners verwendet werden kann. Beispielwerte stehen in `.env.example`.

Wichtige Umgebungsvariablen sind:

* `DB_HOST`: Host der PostgreSQL-Datenbank
* `DB_PORT`: Port der PostgreSQL-Datenbank
* `DB_USER`: Datenbankbenutzer
* `DB_PASSWORD`: Datenbankpasswort
* `DB_NAME`: Datenbankname
* `DB_SSLMODE`: SSL-Modus fuer die Verbindung
* `PORT`: Port fuer die Go-API
* `KEYCLOAK_JWKS_URL`: optionaler JWKS-Endpunkt fuer Keycloak

### Optional: OIDC mit Keycloak

Keycloak ist optional. Die Absicherung wird ueber die Umgebungsvariable `KEYCLOAK_JWKS_URL` gesteuert:

* Ist `KEYCLOAK_JWKS_URL` gesetzt, prueft die Go-API jeden Zugriff auf die `/players`-Routen gegen ein gueltiges Keycloak-JWT (Bearer-Token). Ohne gueltiges Token wird `401 Unauthorized` zurueckgegeben. `GET /health` bleibt immer ohne Token erreichbar.
* Ist die Variable nicht gesetzt, laeuft die API ohne Token-Pruefung durch.

Im mitgelieferten `docker-compose.yml` ist `KEYCLOAK_JWKS_URL` gesetzt, die Absicherung ist beim Start ueber Docker also aktiv. Die Pruefung ist in `internal/middleware/auth.go` umgesetzt. Die Infrastruktur ist vorbereitet:

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

### Schnellstart

1. Anwendung starten: `docker compose up --build`
2. Health-Check pruefen: `GET /health` auf `http://localhost:8080/health`
3. Ersten Player anlegen oder die Liste mit `GET /players` pruefen

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

* Repository analysieren: aktueller Stand, fehlende Punkte fuer die Abgabe, klare Prioritaeten
* Go-API lauffaehig umsetzen mit `GET /health`, `GET /players`, `GET /players/:id`, `POST /players`, `PUT /players/:id`, `DELETE /players/:id`
* Validierung fuer Player-Requests umsetzen (`username` Pflicht, `email` gueltig, `playerClass` nur erlaubte Werte, Fehlerstatus `400`)
* PostgreSQL mit GORM anbinden, Tabellen beim Start anlegen/migrieren, Konfiguration ueber Env-Variablen (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`)
* Docker und Docker Compose so konfigurieren, dass API und DB gemeinsam starten und die API unter `localhost:8080` erreichbar ist
* Integrationstest fuer zentrale Faelle erstellen: Create, List, Update, Delete, Duplicate mit `409`
* Nach jeder Code-Aenderung kurz dokumentieren, welche Dateien angepasst wurden und warum
* Nach jedem groesseren Schritt ausfuehren und pruefen (`go test` oder Compose-Testlauf), Fehler iterativ nachziehen
* README kompakt halten: Start, wichtigste Befehle, Testlauf, Demo-Ablauf
* Offene Aufgaben priorisieren in Muss-Punkte fuer die Abgabe und optionale Ergaenzungen
