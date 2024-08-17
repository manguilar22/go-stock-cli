CREATE TABLE stocks
(
    id       SERIAL PRIMARY KEY,
    symbol   VARCHAR(10),
    period1  BIGINT,
    period2  BIGINT,
    interval VARCHAR(5),
    date     DATE,
    open     NUMERIC(10, 6),
    high     NUMERIC(10, 6),
    low      NUMERIC(10, 6),
    close    NUMERIC(10, 6),
    volume   BIGINT
);