CREATE TYPE TASK_STATUS AS ENUM ('todo', 'in-progress', 'done');

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(50),
    description TEXT NOT NULL,
    status TASK_STATUS DEFAULT 'todo',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);