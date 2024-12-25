DROP TABLE IF EXISTS api_examples CASCADE;

CREATE TABLE api_examples (
  id SERIAL PRIMARY KEY,
  api_endpoint_id INT NOT NULL,
  request_header TEXT,
  request_body TEXT,
  request_parameter TEXT,
  response_status_code INT NOT NULL,
  response_body TEXT,
  response_header TEXT,
  test_state INT NOT NULL DEFAULT 0,
  FOREIGN KEY (api_endpoint_id) REFERENCES api_endpoints(id) ON DELETE CASCADE
);