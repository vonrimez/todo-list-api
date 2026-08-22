DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS users;

CREATE TYPE TASK_STATUS AS ENUM ('todo', 'in-progress', 'done');

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(50),
    description TEXT NOT NULL,
    status TASK_STATUS,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    jwt TEXT
);