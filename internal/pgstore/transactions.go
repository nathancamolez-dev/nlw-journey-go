package pgstore

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nathancamolez-dev/nlw-journey-go/internal/api/spec"
)

func (q *Queries) CreateTrip(
	ctx context.Context,
	pool *pgxpool.Pool,
	params spec.CreateTripRequest,
) (uuid.UUID, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf(
			"pgstore: failed to begin a transaction to create a new trip: %w",
			err,
		)
	}

	defer func() { _ = tx.Rollback(ctx) }() // making explicit ignored error

	qtx := q.WithTx(tx)
	tripID, err := qtx.InsertTrip(ctx, InsertTripParams{
		Destination: params.Destination,
		OwnerEmail:  string(params.OwnerEmail),
		OwnerName:   params.OwnerName,
		StartsAt:    pgtype.Timestamp{Valid: true, Time: params.StartsAt},
		EndsAt:      pgtype.Timestamp{Valid: true, Time: params.EndsAt},
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("pgstore: failed to a new trip: %w", err)
	}

	participants := make([]InviteParticipantsToTripParams, len(params.EmailsToInvite))
	for i, eti := range params.EmailsToInvite {
		participants[i] = InviteParticipantsToTripParams{
			TripID: tripID,
			Email:  string(eti),
		}
	}

	if _, err := qtx.InviteParticipantsToTrip(ctx, participants); err != nil {
		return uuid.UUID{}, fmt.Errorf("pgstore: failed to invite participants: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.UUID{}, fmt.Errorf(
			"pgstore: failed to commit a transaction to create a new trip: %w",
			err,
		)
	}

	return tripID, nil

}
