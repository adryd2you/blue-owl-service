ALTER TABLE
  api_examples
ADD
  COLUMN IF NOT EXISTS test_state INT NOT NULL DEFAULT 0;