CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  createdAt TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  email TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  password TEXT NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS users;