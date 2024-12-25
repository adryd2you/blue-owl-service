DROP TABLE IF EXISTS api_responses CASCADE;

CREATE TABLE api_responses (
  id SERIAL PRIMARY KEY,
  api_endpoint_id INT NOT NULL,
  parameter_name VARCHAR(255) NOT NULL,
  data_type VARCHAR(50) NOT NULL,
  description TEXT,
  is_required BOOLEAN NOT NULL DEFAULT TRUE,
  default_value VARCHAR(255),
  FOREIGN KEY (api_endpoint_id) REFERENCES api_endpoints(id) ON DELETE CASCADE
);