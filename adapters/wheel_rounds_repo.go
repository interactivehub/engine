package adapters

import (
	"context"
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/interactivehub/engine/domain/games/wheel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type sqlWheelRound struct {
	ID                uuid.UUID `db:"id"`
	OutcomeIdx        int       `db:"outcome_idx"`
	StartTime         time.Time `db:"start_time"`
	EndTime           time.Time `db:"end_time"`
	ServerSeed        string    `db:"server_seed"`
	ClientSeed        string    `db:"client_seed"`
	BlindedServerSeed string    `db:"blinded_server_seed"`
	Nonce             int       `db:"nonce"`
}

type sqlWheelRoundEntry struct {
	RoundID   uuid.UUID `db:"round_id"`
	UserID    string    `db:"user_id"`
	Wager     float64   `db:"wager"`
	Pick      string    `db:"pick"`
	EnteredAt time.Time `db:"entered_at"`
}

func newFromWheelRoundEntry(roundId uuid.UUID, entry wheel.WheelRoundEntry) *sqlWheelRoundEntry {
	return &sqlWheelRoundEntry{
		RoundID: roundId,
		UserID:  entry.UserID,
		Wager:   entry.Wager,
		Pick:    string(entry.Pick),
	}
}

func newFromWheelRound(round wheel.WheelRound) *sqlWheelRound {
	return &sqlWheelRound{
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

type WheelRoundsRepo struct {
	db *sqlx.DB
}

func NewWheelRoundsRepo(db *sqlx.DB) *WheelRoundsRepo {
	if db == nil {
		panic("missing db")
	}

	return &WheelRoundsRepo{db}
}

func (r WheelRoundsRepo) GetLatest(ctx context.Context) (wheel.WheelRound, error) {
	query, args, err := sq.
		Select("id", "outcome_idx", "start_time", "end_time", "server_seed", "client_seed", "blinded_server_seed", "nonce").
		From("wheel_rounds").
		ToSql()
	if err != nil {
		return wheel.WheelRound{}, errors.Wrap(err, "failed to retrieve latest wheel round")
	}

	var round wheel.WheelRound

	if err = r.db.GetContext(ctx, &round, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return wheel.WheelRound{}, nil
		}

		return wheel.WheelRound{}, err
	}

	return round, nil
}

func (r WheelRoundsRepo) PersistWheelRound(ctx context.Context, round wheel.WheelRound) error {
	sqlR := newFromWheelRound(round)
	var sqlRe []sqlWheelRoundEntry

	for _, e := range round.Entries {
		sqlRe = append(sqlRe, *newFromWheelRoundEntry(round.ID, e))
	}

	rSql, rArgs, err := sq.
		Insert("wheel_rounds").
		Columns("id", "outcome_idx", "start_time", "end_time", "server_seed", "client_seed", "blinded_server_seed", "nonce").
		Values(sqlR.ID, sqlR.OutcomeIdx, sqlR.StartTime, sqlR.EndTime, sqlR.ServerSeed, sqlR.ClientSeed, sqlR.BlindedServerSeed, sqlR.Nonce).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to persist wheel round")
	}

	log.Println(rSql)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, rSql, rArgs...)
	if err != nil {
		return errors.Wrap(err, "failed to persist wheel round")
	}

	if len(sqlRe) > 0 {
		baseReSql := sq.
			Insert("wheel_round_entries").
			Columns("round_id", "user_id", "wager", "pick")

		for _, re := range sqlRe {
			baseReSql.Values(re.RoundID, re.UserID, re.Wager, re.Pick)
		}

		reSql, reArgs, err := baseReSql.ToSql()
		if err != nil {
			return errors.Wrap(err, "failed to persist wheel round entries")
		}

		_, err = tx.ExecContext(ctx, reSql, reArgs...)
		if err != nil {
			return errors.Wrap(err, "failed to persist wheel round entries")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
