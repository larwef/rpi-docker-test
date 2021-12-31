-- name: AddEnemy :one
INSERT INTO enemies (enemy_id, full_name, email, rating, last_updated)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetEnemy :one
SELECT * FROM enemies
WHERE enemy_id = $1;

-- name: UpdateEnemy :one
UPDATE enemies
SET
    full_name = COALESCE(NULLIF(@full_name::text, ''), full_name),
    email = COALESCE(NULLIF(@email::text, ''), email),
    rating = COALESCE(NULLIF(@rating::real, 0.0), rating),
    last_updated = @last_updated::timestamp
WHERE enemy_id = @enemy_id::text
RETURNING *;

-- name: ListEnemies :many
SELECT * FROM enemies;
