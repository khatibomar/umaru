CREATE TABLE books (
  id     BIGSERIAL PRIMARY KEY,
  name   text      NOT NULL,
  link   text      NOT NULL,
  image  text      NOT NULL,
  status integer   NOT NULL
);