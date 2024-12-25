DROP TABLE IF EXISTS api_projects CASCADE;

CREATE TABLE api_projects (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  status VARCHAR(50) DEFAULT 'active'
);