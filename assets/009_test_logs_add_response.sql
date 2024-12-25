ALTER TABLE
  test_logs
ADD
  COLUMN response_status_code INT,
ADD
  COLUMN response_body TEXT,
ADD
  COLUMN response_header TEXT,
ADD
  COLUMN test_message TEXT;