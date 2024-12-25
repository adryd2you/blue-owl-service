DROP TABLE IF EXISTS api_type CASCADE;
DROP TABLE IF EXISTS http_method CASCADE;

CREATE TYPE http_method AS ENUM ('GET', 'POST', 'PUT', 'DELETE', 'PATCH');

CREATE TYPE api_type  AS ENUM ('REST', 'GRPC');

DROP TABLE IF EXISTS api_endpoints CASCADE;

CREATE TABLE api_endpoints (
  id SERIAL PRIMARY KEY,
  api_service_id INT NOT NULL,
  api_type api_type NOT NULL,
  name VARCHAR(255) NOT NULL,
  path VARCHAR(255) NOT NULL,
  method http_method NOT NULL,
  description TEXT,
  FOREIGN KEY (api_service_id) REFERENCES api_services(id) ON DELETE CASCADE
);