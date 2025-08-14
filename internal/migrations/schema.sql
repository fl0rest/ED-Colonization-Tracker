CREATE TABLE depotEvents (
  id          INTEGER   PRIMARY KEY AUTOINCREMENT,
  completion  REAL      NOT NULL,
  time        INTEGER   NOT NULL,
  marketId    INTEGER   NOT NULL,
  raw_text    TEXT      NOT NULL
);

CREATE TABLE dockEvents (
  id          INTEGER   PRIMARY KEY AUTOINCREMENT,
  time        INTEGER   NOT NULL,
  marketId    TEXT      NOT NULL,
  systemName  TEXT      NOT NULL,
  stationName TEXT      NOT NULL
);

CREATE TABLE events (
  id            INTEGER   PRIMARY KEY AUTOINCREMENT,
  time          INTEGER   NOT NULL,
  completion    REAL      NOT NULL,
  marketId      INTEGER   NOT NULL,
  stationId     INTEGER   NOT NULL,
  raw_resources TEXT      NOT NULL
);

CREATE TABLE stations (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  systemName  TEXT    NOT NULL,
  stationName TEXT    NOT NULL,
  marketId    INTEGER NOT NULL
);

CREATE TABLE resources (
  id        INTEGER   NOT NULL,
  eventId   INTEGER   NOT NULL,
  name      TEXT      NOT NULL,
  required  INTEGER   NOT NULL,
  provided  INTEGER   NOT NULL,
  diff      INTEGER   NOT NULL,
  payment   INTEGER   NOT NULL,
  time      INTEGER   NOT NULL,
  stationId INTEGER   NOT NULL,
  PRIMARY KEY (id, stationId)
);

CREATE TABLE resourceIds (
  id    INTEGER NOT NULL,
  name  TEXT    NOT NULL
);
