-- Aufruf: psql --dbname=player --username=player --file=/init/player/sql/copy-csv.sql

SET search_path TO public;

INSERT INTO guild (id, name, description, "foundedAt", version)
VALUES
  (10, 'Warriors', 'Eine Guild', CURRENT_TIMESTAMP, 0),
  (20, 'Mages', 'Magische Guild', CURRENT_TIMESTAMP, 0),
  (30, 'Hunters', 'Jaeger Guild', CURRENT_TIMESTAMP, 0);

INSERT INTO player (
  id,
  username,
  email,
  level,
  experience,
  "playerClass",
  status,
  "guildId",
  "createdAt",
  "updatedAt",
  version
)
VALUES
  (1, 'player1', 'player1@example.com', 1, 100, 'WARRIOR', 'ACTIVE', 10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0),
  (20, 'player2', 'player2@example.com', 10, 500, 'MAGE', 'ACTIVE', 20, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0),
  (30, 'testplayer3', 'player3@example.com', 30, 1500, 'ROGUE', 'ACTIVE', 10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0),
  (40, 'delta', 'delta@example.com', 40, 2500, 'PRIEST', 'ACTIVE', 30, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0),
  (50, 'echo', 'echo@example.com', 50, 3500, 'HUNTER', 'BANNED', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0),
  (60, 'remove', 'remove@example.com', 60, 4500, 'WARRIOR', 'ACTIVE', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0),
  (70, 'updateTarget', 'update-target@example.com', 70, 5500, 'PRIEST', 'ACTIVE', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0),
  (80, 'bannedTarget', 'banned-target@example.com', 80, 6500, 'HUNTER', 'BANNED', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0);

SELECT setval(
  pg_get_serial_sequence('guild', 'id'),
  (SELECT max(id) FROM guild)
);

SELECT setval(
  pg_get_serial_sequence('player', 'id'),
  (SELECT max(id) FROM player)
);
