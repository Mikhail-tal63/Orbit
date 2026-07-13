
DB_URL=host=localhost port=5432 user=postgres password=postgres dbname=orbit sslmode=disable

migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down

migrate-status:
	goose -dir migrations postgres "$(DB_URL)" status

migration:
	goose -dir migrations create $(name) sql