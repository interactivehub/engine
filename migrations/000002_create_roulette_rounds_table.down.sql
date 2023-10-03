ALTER TABLE IF EXISTS roulette_round_entries
DROP CONSTRAINT IF EXISTS fk_roulette_round_user_id;

ALTER TABLE IF EXISTS roulette_round_entries
DROP CONSTRAINT IF EXISTS fk_roulette_round_id;

DROP TABLE IF EXISTS roulette_round_entries;

DROP TABLE IF EXISTS roulette_rounds;