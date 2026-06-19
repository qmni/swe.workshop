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

### REST-Schnittstelle (Lesen und Neuanlegen)

Go mit Fiber (`github.com/gofiber/fiber/v2`)

Implementierte Endpunkte:

* `GET /health`
* `GET /products`
* `GET /products/:id`
* `POST /products`

Beispiel zum Neuanlegen:

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Notebook","description":"Workshop product","priceCents":1299}'
```

Beispiel zum Lesen:

```bash
curl http://localhost:8080/products
```

### Validierung (nur Neuanlegen)

`github.com/go-playground/validator/v10`

Validiert wird beim `POST /products`:

* `name`: Pflichtfeld, 2 bis 120 Zeichen
* `description`: maximal 500 Zeichen
* `priceCents`: Pflichtfeld, mindestens 1, maximal 10000000

### OR-Mapping (für PostgreSQL)

GORM:

* `gorm.io/gorm`
* `gorm.io/driver/postgres`

Die Tabelle `products` wird beim Start per `AutoMigrate` angelegt.

Die Datenbankverbindung ist ueber Umgebungsvariablen konfigurierbar, damit auch der DB-Server des Teampartners verwendet werden kann. Beispielwerte stehen in `.env.example`.

### Optional: OIDC mit Keycloak

Nicht implementiert. Keycloak war laut Aufgabenstellung optional.

### Einfacher Integrationstest

Integrationstest in `integration/products_test.go`.

Der Test startet die API gegen eine PostgreSQL-Testdatenbank, legt ein Produkt per REST an und liest die Produktliste wieder aus.

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

Zusaetzlich liegen Beispiel-Requests in `requests.http`, damit die API direkt aus einer IDE oder mit einem REST-Client getestet werden kann.

### Demo-Ablauf

1. Anwendung starten: `docker compose up --build`
2. Health-Check aufrufen: `GET /health`
3. Produkt anlegen: `POST /products`
4. Produktliste lesen: `GET /products`
5. Validierung testen, indem ein Produkt ohne Namen oder mit `priceCents: 0` gesendet wird

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
