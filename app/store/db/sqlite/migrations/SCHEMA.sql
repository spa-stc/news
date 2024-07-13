-- Days Table.
CREATE TABLE IF NOT EXISTS days (
  date TEXT PRIMARY KEY,
  lunch TEXT NOT NULL,
  x_period TEXT NOT NULL,
  rotation_day TEXT NOT NULL,
  location TEXT NOT NULL,
  notes TEXT NOT NULL,
  ap_info TEXT NOT NULL,
  cc_info TEXT NOT NULL,
  grade_9 TEXT NOT NULL,
  grade_10 TEXT NOT NULL,
  grade_11 TEXT NOT NULL,
  grade_12 TEXT NOT NULL,

  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  updated_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))
);
