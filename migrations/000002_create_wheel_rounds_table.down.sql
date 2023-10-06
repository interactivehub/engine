ALTER TABLE IF EXISTS wheel_round_entries
DROP CONSTRAINT IF EXISTS fk_wheel_round_user_id;

ALTER TABLE IF EXISTS wheel_round_entries
DROP CONSTRAINT IF EXISTS fk_wheel_round_id;

DROP TABLE IF EXISTS wheel_round_entries;

DROP TABLE IF EXISTS wheel_rounds;