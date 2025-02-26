-- name: InsertTrip :one
INSERT INTO trips("destination","owner_email","owner_name","starts_at","ends_at") VALUES ($1, $2, $3, $4, $5) RETURNING "id";

-- name: GetTrip :one
select *
from trips
where id = $1
;

-- name: UpdateTrip :exec
UPDATE trips SET "destination" = $1, "ends_at" = $2, "starts_at" = $3, "is_confirmed" =$4  WHERE id = $5 ;

-- name: GetParticipant :one
select "id", "trip_id", "email", "is_confirmed"
from participants
where id = $1
;

-- name: ConfirmParticipant :exec
UPDATE participants SET "is_confirmed" = true WHERE id = $1;

-- name: GetParticipants :many
select *
from participants
where trip_id = $1
;

-- name: InviteParticipantToTrip :one
INSERT INTO participants("trip_id","email") VALUES ($1,$2) RETURNING "id";

-- name: InviteParticipantsToTrip :copyfrom
INSERT INTO participants("trip_id","email") VALUES ($1 , $2);


-- name: CreateActivity :one
INSERT INTO activities ("trip_id","title","occurs_at") VALUES ($1, $2, $3) RETURNING "id";

-- name: GetActivities :many
select *
from activities
where trip_id = $1
;

-- name: CreateTripLink :one
INSERT INTO links ("trip_id","title","url") VALUES ($1, $2, $3) RETURNING "id";

-- name: GetLinks :many
select *
from links
where trip_id = $1
;

