package adapters

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/interactivehub/engine/domain/games/wheel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	getLatestRoundQuery = `
        SELECT id, outcome_idx, round_start_time, spin_start_time, round_end_time, server_seed, client_seed, blinded_server_seed, nonce, status
        FROM wheel_rounds
        ORDER BY round_end_time DESC;
    `
	createWheelRoundQuery = `
        INSERT INTO wheel_rounds (id, outcome_idx, round_start_time, spin_start_time, round_end_time, server_seed, client_seed, blinded_server_seed, nonce, status)
        VALUES (:id, :outcome_idx, :round_start_time, :spin_start_time, :round_end_time, :server_seed, :client_seed, :blinded_server_seed, :nonce, :status);
    `
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
	Entries           []sqlWheelRoundEntry   `db:"-"`
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
	var sqlRound sqlWheelRound

	if err := r.db.GetContext(ctx, &sqlRound, getLatestRoundQuery); err != nil {
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
		Status:         sqlRound.Status,
		Outcome:        wheel.WheelItems()[sqlRound.OutcomeIdx],
		RoundStartTime: sqlRound.RoundStartTime,
		SpinStartTime:  sqlRound.SpinStartTime,
		RoundEndTime:   sqlRound.RoundEndTime,
	}, nil
}

func (r WheelRoundsRepo) CreateWheelRound(ctx context.Context, round wheel.WheelRound) error {
	sqlRound := newFromWheelRound(round)

	_, err := r.db.NamedExecContext(ctx, createWheelRoundQuery, sqlRound)
	if err != nil {
		return errors.Wrap(err, "failed to create wheel round")
	}

	return nil
}
