CREATE INDEX IF NOT EXISTS tasks_title_idx ON tasks USING GIN (to_tsvector('simple', title));
