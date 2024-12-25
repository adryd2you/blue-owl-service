DROP TABLE IF EXISTS api_services CASCADE;

CREATE TABLE api_services (
    id SERIAL PRIMARY KEY,
    api_project_id INT NOT NULL,
    name TEXT NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'active',
    FOREIGN KEY (api_project_id) REFERENCES api_projects(id) ON DELETE CASCADE
);