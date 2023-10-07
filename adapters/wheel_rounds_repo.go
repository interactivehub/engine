package adapters

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/interactivehub/engine/domain/games/wheel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type sqlWheelRound struct {
	ID                uuid.UUID              `db:"id"`
	OutcomeIdx        int                    `db:"outcome_idx"`
	RoundStartTime    time.Time              `db:"round_start_time"`
	SpinStartTime     time.Time              `db:"spin_start_time"`
	RoundEndTime      time.Time              `db:"round_end_time"`
	ServerSeed        string                 `db:"server_seed"`
	ClientSeed        string                 `db:"client_seed"`
	BlindedServerSeed string                 `db:"blinded_server_seed"`
	Nonce             int                    `db:"nonce"`
	Status            wheel.WheelRoundStatus `db:"status"`
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
		RoundStartTime:    round.RoundStartTime,
		SpinStartTime:     round.SpinStartTime,
		RoundEndTime:      round.RoundEndTime,
		ServerSeed:        round.StringServerSeed(),
		ClientSeed:        round.StringClientSeed(),
		BlindedServerSeed: round.StringBlindedServerSeed(),
		Nonce:             int(round.Nonce),
		Status:            round.Status,
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
		Select("id", "outcome_idx", "round_start_time", "spin_start_time", "round_end_time", "server_seed", "client_seed", "blinded_server_seed", "nonce", "status").
		From("wheel_rounds").
		ToSql()
	if err != nil {
		return wheel.WheelRound{}, errors.Wrap(err, "failed to retrieve latest wheel round")
	}

	var sqlRound sqlWheelRound

	if err = r.db.GetContext(ctx, &sqlRound, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return wheel.WheelRound{}, nil
		}

		return wheel.WheelRound{}, err
	}

	return wheel.WheelRound{
		ID: sqlRound.ID,
		ProvablyFair: wheel.ProvablyFair{
			ServerSeed:        []byte(sqlRound.ServerSeed),
			ClientSeed:        []byte(sqlRound.ClientSeed),
			BlindedServerSeed: []byte(sqlRound.BlindedServerSeed),
			Nonce:             uint64(sqlRound.Nonce),
		},
		Status: sqlRound.Status,
		// Entries: TODO: Gotta join with wheel_round_entries table
		Outcome:        wheel.WheelItems()[sqlRound.OutcomeIdx],
		RoundStartTime: sqlRound.RoundStartTime,
		SpinStartTime:  sqlRound.SpinStartTime,
		RoundEndTime:   sqlRound.RoundEndTime,
	}, nil
}

func (r WheelRoundsRepo) PersistWheelRound(ctx context.Context, round wheel.WheelRound) error {
	sqlR := newFromWheelRound(round)
	// var sqlRe []sqlWheelRoundEntry

	// for _, e := range round.Entries {
	// 	sqlRe = append(sqlRe, *newFromWheelRoundEntry(round.ID, e))
	// }

	sql := "INSERT INTO wheel_rounds (id, outcome_idx, round_start_time, spin_start_time, round_end_time, server_seed, client_seed, blinded_server_seed, nonce, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);"

	// tx, err := r.db.BeginTx(ctx, nil)
	// if err != nil {
	// 	return err
	// }

	// defer tx.Rollback()

	_, err := r.db.ExecContext(ctx, sql, sqlR.ID, sqlR.OutcomeIdx, sqlR.RoundStartTime, sqlR.SpinStartTime, sqlR.RoundEndTime, sqlR.ServerSeed, sqlR.ClientSeed, sqlR.BlindedServerSeed, sqlR.Nonce, sqlR.Status)
	if err != nil {
		return errors.Wrap(err, "failed to persist wheel round")
	}

	// if len(sqlRe) > 0 {
	// 	baseReSql := sq.
	// 		Insert("wheel_round_entries").
	// 		Columns("round_id", "user_id", "wager", "pick")

	// 	for _, re := range sqlRe {
	// 		baseReSql.Values(re.RoundID, re.UserID, re.Wager, re.Pick)
	// 	}

	// 	reSql, reArgs, err := baseReSql.ToSql()
	// 	if err != nil {
	// 		return errors.Wrap(err, "failed to persist wheel round entries")
	// 	}

	// 	_, err = tx.ExecContext(ctx, reSql, reArgs...)
	// 	if err != nil {
	// 		return errors.Wrap(err, "failed to persist wheel round entries")
	// 	}
	// }

	// err = tx.Commit()
	// if err != nil {
	// 	return errors.Wrap(err, "failed to commit transaction")
	// }

	return nil
}
