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
	getByIdQuery = `
        SELECT id, outcome_idx, round_start_time, spin_start_time, round_end_time, server_seed, client_seed, blinded_server_seed, nonce, status
        FROM wheel_rounds
        WHERE id=$1;
    `
	getLatestQuery = `
        SELECT id, outcome_idx, round_start_time, spin_start_time, round_end_time, server_seed, client_seed, blinded_server_seed, nonce, status
        FROM wheel_rounds
        ORDER BY round_end_time DESC;
    `
	upsertQuery = `
        INSERT INTO wheel_rounds (id, outcome_idx, round_start_time, spin_start_time, round_end_time, server_seed, client_seed, blinded_server_seed, nonce, status)
        VALUES (:id, :outcome_idx, :round_start_time, :spin_start_time, :round_end_time, :server_seed, :client_seed, :blinded_server_seed, :nonce, :status)
        ON CONFLICT (id) DO UPDATE 
        SET 
            outcome_idx = :outcome_idx,
            round_start_time = :round_start_time,
            spin_start_time = :spin_start_time,
            round_end_time = :round_end_time,
            server_seed = :server_seed,
            client_seed = :client_seed,
            blinded_server_seed = :blinded_server_seed,
            nonce = :nonce,
            status = :status;
    `
	createEntryQuery = `
        INSERT INTO wheel_round_entries (round_id, user_id, bet, pick, entry_time)
        VALUES (:round_id, :user_id, :bet, :pick, :entry_time); 
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
	Entries           []*sqlWheelRoundEntry  `db:"-"`
}

type sqlWheelRoundEntry struct {
	RoundID   uuid.UUID `db:"round_id"`
	UserID    string    `db:"user_id"`
	Bet       float64   `db:"bet"`
	Pick      string    `db:"pick"`
	EntryTime time.Time `db:"entry_time"`
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

func (r WheelRoundsRepo) GetByID(ctx context.Context, id uuid.UUID) (*wheel.WheelRound, error) {
	sqlRound := &sqlWheelRound{}

	if err := r.db.GetContext(ctx, sqlRound, getByIdQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to get wheel round by id")
	}

	return sqlRound.toWheelRound(), nil
}

func (r WheelRoundsRepo) GetLatest(ctx context.Context) (*wheel.WheelRound, error) {
	sqlRound := &sqlWheelRound{}

	if err := r.db.GetContext(ctx, sqlRound, getLatestQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to get latest wheel round")
	}

	return sqlRound.toWheelRound(), nil
}

func (r WheelRoundsRepo) Upsert(ctx context.Context, round *wheel.WheelRound) error {
	sqlRound := &sqlWheelRound{}
	sqlRound.fromWheelRound(round)

	_, err := r.db.NamedExecContext(ctx, upsertQuery, sqlRound)
	if err != nil {
		return errors.Wrap(err, "failed to upsert wheel round")
	}

	return nil
}

func (r WheelRoundsRepo) CreateEntry(ctx context.Context, roundEntry *wheel.WheelRoundEntry) error {
	round, err := r.GetByID(ctx, roundEntry.RoundID)
	if err != nil {
		return errors.Wrap(err, "failed to create wheel round entry")
	}

	if round == nil {
		return errors.New("failed to create wheel round entry: unknown round id")
	}

	sqlRoundEntry := &sqlWheelRoundEntry{}
	sqlRoundEntry.fromWheelRoundEntry(roundEntry)

	_, err = r.db.NamedExecContext(ctx, createEntryQuery, sqlRoundEntry)
	if err != nil {
		return errors.Wrap(err, "failed to create wheel round entry")
	}

	return nil
}

func (r *sqlWheelRound) fromWheelRound(round *wheel.WheelRound) {
	if round == nil {
		return
	}

	r.ID = round.ID
	r.OutcomeIdx = round.GetOutcomeIdx()
	r.RoundStartTime = round.RoundStartTime
	r.SpinStartTime = round.SpinStartTime
	r.RoundEndTime = round.RoundEndTime
	r.ServerSeed = round.StringServerSeed()
	r.ClientSeed = round.StringClientSeed()
	r.BlindedServerSeed = round.StringBlindedServerSeed()
	r.Nonce = int(round.Nonce)
	r.Status = round.Status
}

func (r *sqlWheelRound) toWheelRound() *wheel.WheelRound {
	return &wheel.WheelRound{
		ID: r.ID,
		ProvablyFair: wheel.ProvablyFair{
			ServerSeed:        []byte(r.ServerSeed),
			ClientSeed:        []byte(r.ClientSeed),
			BlindedServerSeed: []byte(r.BlindedServerSeed),
			Nonce:             uint64(r.Nonce),
		},
		Status:         r.Status,
		Outcome:        wheel.WheelItems()[r.OutcomeIdx],
		RoundStartTime: r.RoundStartTime,
		SpinStartTime:  r.SpinStartTime,
		RoundEndTime:   r.RoundEndTime,
	}
}

func (e *sqlWheelRoundEntry) fromWheelRoundEntry(entry *wheel.WheelRoundEntry) {
	if entry == nil {
		return
	}

	e.RoundID = entry.RoundID
	e.UserID = entry.UserID
	e.Bet = entry.Bet
	e.Pick = string(entry.Pick)
	e.EntryTime = entry.EntryTime
}
