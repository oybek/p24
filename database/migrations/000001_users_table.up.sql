
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
  "chat_id" BIGINT PRIMARY KEY,
  "uuid" uuid DEFAULT uuid_generate_v4 ()
);
