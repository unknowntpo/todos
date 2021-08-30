CREATE TABLE IF NOT EXISTS tasks (
    id bigserial PRIMARY KEY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    user_id integer NOT NULL,
    title text NOT NULL,
    content text NOT NULL,
    done boolean NOT NULL DEFAULT false,
    version integer NOT NULL DEFAULT 1
);
