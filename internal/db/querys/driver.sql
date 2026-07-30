--name CreateDriver : one
INSERT INTO drivers (
        id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id),
)VALUES(
    $1 ,$2
)
RETURNING *;