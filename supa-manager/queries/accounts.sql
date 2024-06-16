-- name: CreateAccount :one
INSERT INTO public.accounts (email, password_hash, username)
VALUES ($1, $2, $3)
RETURNING *;

-- name: SetAccountName :exec
UPDATE public.accounts
SET first_name = $2,
    last_name  = $3
WHERE id = $1;

-- name: GetAccountByEmail :one
SELECT * FROM public.accounts WHERE email = $1;

-- name: GetAccountByID :one
SELECT * FROM public.accounts WHERE id = $1;

-- name: GetAccountByGoTrueID :one
SELECT * FROM public.accounts WHERE gotrue_id = $1;