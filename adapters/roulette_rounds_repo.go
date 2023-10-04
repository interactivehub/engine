package adapters

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/interactivehub/engine/domain/games/roulette"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type sqlRouletteRound struct {
	ID                uuid.UUID `db:"id"`
	OutcomeIdx        int       `db:"outcome_idx"`
	StartTime         time.Time `db:"start_time"`
	EndTime           time.Time `db:"end_time"`
	ServerSeed        string    `db:"server_seed"`
	ClientSeed        string    `db:"client_seed"`
	BlindedServerSeed string    `db:"blinded_server_seed"`
	Nonce             int       `db:"nonce"`
}

type sqlRouletteRoundEntry struct {
	RoundID   uuid.UUID `db:"round_id"`
	UserID    string    `db:"user_id"`
	Wager     float64   `db:"wager"`
	Pick      string    `db:"pick"`
	EnteredAt time.Time `db:"entered_at"`
}

func newFromRouletteRoundEntry(roundId uuid.UUID, entry roulette.RouletteRoundEntry) *sqlRouletteRoundEntry {
	return &sqlRouletteRoundEntry{
		RoundID: roundId,
		UserID:  entry.UserID,
		Wager:   entry.Wager,
		Pick:    string(entry.Pick),
	}
}

func newFromRouletteRound(round roulette.RouletteRound) *sqlRouletteRound {
	return &sqlRouletteRound{
		ID:                round.ID,
		OutcomeIdx:        round.GetOutcomeIdx(),
		StartTime:         round.StartTime,
		EndTime:           round.EndTime,
		ServerSeed:        round.StringServerSeed(),
		ClientSeed:        round.StringClientSeed(),
		BlindedServerSeed: round.StringBlindedServerSeed(),
		Nonce:             int(round.Nonce),
	}
}

type RouletteRoundsRepo struct {
	db *sqlx.DB
}

func NewRouletteRoundsRepo(db *sqlx.DB) *RouletteRoundsRepo {
	if db == nil {
		panic("missing db")
	}

	return &RouletteRoundsRepo{db}
}

func (r RouletteRoundsRepo) PersistRouletteRound(ctx context.Context, round roulette.RouletteRound) error {
	sqlR := newFromRouletteRound(round)
	var sqlRe []sqlRouletteRoundEntry

	for _, e := range round.Entries {
		sqlRe = append(sqlRe, *newFromRouletteRoundEntry(round.ID, e))
	}

	rSql, rArgs, err := sq.
		Insert("roulette_rounds").
		Columns("id", "outcome_idx", "start_time", "end_time", "server_seed", "client_seed", "blinded_server_seed", "nonce").
		Values(sqlR.ID, sqlR.OutcomeIdx, sqlR.StartTime, sqlR.EndTime, sqlR.ServerSeed, sqlR.ClientSeed, sqlR.BlindedServerSeed, sqlR.Nonce).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to persist roulette round")
	}

	baseReSql := sq.
		Insert("roulette_round_entries").
		Columns("round_id", "user_id", "wager", "pick")

	for _, re := range sqlRe {
		baseReSql.Values(re.RoundID, re.UserID, re.Wager, re.Pick)
	}

	reSql, reArgs, err := baseReSql.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to persist roulette round entries")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, rSql, rArgs...)
	if err != nil {
		return errors.Wrap(err, "failed to persist roulette round")
	}

	_, err = tx.ExecContext(ctx, reSql, reArgs...)
	if err != nil {
		return errors.Wrap(err, "failed to persist roulette round entries")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
