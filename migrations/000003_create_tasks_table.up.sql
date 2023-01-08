CREATE TABLE IF NOT EXISTS tasks (
    id bigserial PRIMARY KEY,  
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    content text NOT NULL,
    done boolean NOT NULL DEFAULT false,
    version integer NOT NULL DEFAULT 1
);
