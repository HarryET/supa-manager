-- name: CreateAccount :exec
INSERT INTO public.accounts (email, password_hash, username)
VALUES ($1, $2, $3);

-- name: SetAccountName :exec
UPDATE public.accounts
SET first_name = $2,
    last_name  = $3
WHERE id = $1;