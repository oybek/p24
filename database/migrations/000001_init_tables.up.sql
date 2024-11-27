
CREATE EXTENSION fuzzystrmatch;

CREATE TABLE IF NOT EXISTS apteka (
  id SERIAL PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "phone" VARCHAR NOT NULL,
  "address" VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS medicine (
  id SERIAL PRIMARY KEY,
  "name" VARCHAR UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS apteka_medicine (
  "apteka_id" INTEGER NOT NULL,
  "medicine_id" INTEGER NOT NULL,
  "amount" INTEGER NOT NULL,
  "updated" TIMESTAMPTZ NOT NULL,
  PRIMARY KEY ("apteka_id", "medicine_id")
);

CREATE TABLE IF NOT EXISTS users (
  "chat_id" BIGINT PRIMARY KEY,
  "apteka_id" INTEGER NOT NULL,
  "reader" VARCHAR NOT NULL
);
