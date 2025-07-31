CREATE TABLE events (
  id          INTEGER   PRIMARY KEY AUTOINCREMENT,
  raw_text    TEXT      NOT NULL,
  completion  REAL      NOT NULL,
  time        INTEGER   NOT NULL,
  marketId    INTEGER   NOT NULL
);

CREATE TABLE resources (
  id        BIGSERIAL PRIMARY KEY,
  name      TEXT      NOT NULL,
  required  INTEGER   NOT NULL,
  provided  INTEGER   NOT NULL,
  diff      INTEGER   NOT NULL,
  payment   INTEGER   NOT NULL,
  time      INTEGER   NOT NULL
);

CREATE TABLE resourceIds (
  id    INTEGER NOT NULL,
  name  TEXT    NOT NULL
);
