CREATE TABLE IF NOT EXISTS posts (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  createdAt TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  content TEXT NOT NULL,
  title TEXT NOT NULL,
  updatedAt TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  userId uuid NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS posts;
