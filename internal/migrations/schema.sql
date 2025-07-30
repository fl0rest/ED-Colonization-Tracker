CREATE TABLE events (
  id        BIGSERIAL PRIMARY KEY,
  raw_text  TEXT      NOT NULL,
  time      INTEGER   NOT NULL
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
