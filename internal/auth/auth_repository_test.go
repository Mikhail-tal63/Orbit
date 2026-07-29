package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Mikhail-Tal63/Orbit/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, *AuthRepositoryImpl) {
	t.Helper()
	ctx := context.Background()

	var connStr string
	envConnStr := os.Getenv("TEST_DATABASE_URL")
	if envConnStr == "" {
		envConnStr = os.Getenv("DATABASE_URL")
	}

	if envConnStr != "" {
		connStr = envConnStr
	} else {
	
		pgContainer, err := postgrescontainer.Run(ctx,
			"postgres:16-alpine",
			postgrescontainer.WithDatabase("orbit_test"),
			postgrescontainer.WithUsername("postgres"),
			postgrescontainer.WithPassword("postgres"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(30*time.Second),
			),
		)
		if err != nil {
			t.Skipf("Skipping PostgreSQL repository integration test: docker container start failed: %v", err)
			return nil, nil
		}

		t.Cleanup(func() {
			if err := pgContainer.Terminate(context.Background()); err != nil {
				t.Logf("failed to terminate container: %v", err)
			}
		})

		connStr, err = pgContainer.ConnectionString(ctx, "sslmode=disable")
		require.NoError(t, err, "failed to build container connection string")
	}

	sqlDB, err := sql.Open("pgx", connStr)
	require.NoError(t, err, "failed to open database/sql connection for goose migrations")
	defer sqlDB.Close()

	err = goose.SetDialect("postgres")
	require.NoError(t, err, "failed to set goose dialect")

	migrationsDir := findMigrationsDir(t)
	err = goose.Up(sqlDB, migrationsDir)
	require.NoError(t, err, "failed to run goose migrations")


	poolConfig, err := pgxpool.ParseConfig(connStr)
	require.NoError(t, err, "failed to parse database config")

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	require.NoError(t, err, "failed to connect pgxpool")

	t.Cleanup(func() {
		pool.Close()
	})

	queries := db.New(pool)
	repo := NewAuthRepository(queries)

	return pool, repo
}

func findMigrationsDir(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	require.NoError(t, err)

	for {
		candidate := filepath.Join(dir, "migrations")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatalf("migrations directory not found from working directory")
	return ""
}

func truncateTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(), "TRUNCATE TABLE users CASCADE;")
	require.NoError(t, err, "failed to truncate users table")
}

func randomCreateUserParams(t *testing.T) db.CreateUserParams {
	t.Helper()
	uid := uuid.New()
	shortID := uid.String()[:8]
	return db.CreateUserParams{
		ID:           uid,
		FirstName:    "Integration",
		LastName:     "Tester",
		Username:     fmt.Sprintf("user_%s", shortID),
		Email:        fmt.Sprintf("user_%s@orbit.test", shortID),
		PasswordHash: "$2a$10$abcdefghijklmnopqrstuvwxyz123456",
	}
}

func TestAuthRepository_CreateUser(t *testing.T) {
	pool, repo := setupTestDB(t)
	if pool == nil {
		return
	}

	t.Run("Success", func(t *testing.T) {
		truncateTables(t, pool)
		ctx := context.Background()

		params := randomCreateUserParams(t)
		user, err := repo.CreateUser(ctx, params)

		require.NoError(t, err)
		require.NotNil(t, user)

		assert.Equal(t, params.ID, user.ID)
		assert.Equal(t, params.FirstName, user.FirstName)
		assert.Equal(t, params.LastName, user.LastName)
		assert.Equal(t, params.Username, user.Username)
		assert.Equal(t, params.Email, user.Email)
		assert.Equal(t, params.PasswordHash, user.PasswordHash)
		assert.Equal(t, "", user.Phone)
		assert.Equal(t, "passenger", user.Role)
		assert.True(t, user.IsActive)
		assert.False(t, user.CreatedAt.Time.IsZero())
		assert.False(t, user.UpdatedAt.Time.IsZero())
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		truncateTables(t, pool)
		ctx := context.Background()

		params1 := randomCreateUserParams(t)
		_, err := repo.CreateUser(ctx, params1)
		require.NoError(t, err)

		params2 := randomCreateUserParams(t)
		params2.Email = params1.Email // Duplicate email

		user, err := repo.CreateUser(ctx, params2)
		require.Error(t, err)
		require.Nil(t, user)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			assert.Equal(t, "23505", pgErr.Code, "should fail with PostgreSQL unique constraint violation")
		}
	})

	t.Run("DuplicateUsername", func(t *testing.T) {
		truncateTables(t, pool)
		ctx := context.Background()

		params1 := randomCreateUserParams(t)
		_, err := repo.CreateUser(ctx, params1)
		require.NoError(t, err)

		params2 := randomCreateUserParams(t)
		params2.Username = params1.Username // Duplicate username

		user, err := repo.CreateUser(ctx, params2)
		require.Error(t, err)
		require.Nil(t, user)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			assert.Equal(t, "23505", pgErr.Code, "should fail with PostgreSQL unique constraint violation")
		}
	})
}

func TestAuthRepository_GetUserByEmail(t *testing.T) {
	pool, repo := setupTestDB(t)
	if pool == nil {
		return
	}

	t.Run("Found", func(t *testing.T) {
		truncateTables(t, pool)
		ctx := context.Background()

		params := randomCreateUserParams(t)
		createdUser, err := repo.CreateUser(ctx, params)
		require.NoError(t, err)

		foundUser, err := repo.GetUserByEmail(ctx, params.Email)
		require.NoError(t, err)
		require.NotNil(t, foundUser)

		assert.Equal(t, createdUser.ID, foundUser.ID)
		assert.Equal(t, createdUser.Email, foundUser.Email)
		assert.Equal(t, createdUser.Username, foundUser.Username)
	})

	t.Run("NotFound", func(t *testing.T) {
		truncateTables(t, pool)
		ctx := context.Background()

		foundUser, err := repo.GetUserByEmail(ctx, "nonexistent@orbit.test")
		require.Error(t, err)
		require.True(t, errors.Is(err, pgx.ErrNoRows), "expected pgx.ErrNoRows error")
		require.Nil(t, foundUser)
	})
}

func TestAuthRepository_GetUserByUsername(t *testing.T) {
	pool, repo := setupTestDB(t)
	if pool == nil {
		return
	}

	t.Run("Found", func(t *testing.T) {
		truncateTables(t, pool)
		ctx := context.Background()

		params := randomCreateUserParams(t)
		createdUser, err := repo.CreateUser(ctx, params)
		require.NoError(t, err)

		foundUser, err := repo.GetUserByUsername(ctx, params.Username)
		require.NoError(t, err)
		require.NotNil(t, foundUser)

		assert.Equal(t, createdUser.ID, foundUser.ID)
		assert.Equal(t, createdUser.Username, foundUser.Username)
	})

	t.Run("NotFound", func(t *testing.T) {
		truncateTables(t, pool)
		ctx := context.Background()

		foundUser, err := repo.GetUserByUsername(ctx, "unknown_user_123")
		require.Error(t, err)
		require.True(t, errors.Is(err, pgx.ErrNoRows), "expected pgx.ErrNoRows error")
		require.Nil(t, foundUser)
	})
}

func TestAuthRepository_GetUserByID(t *testing.T) {
	pool, repo := setupTestDB(t)
	if pool == nil {
		return
	}

	t.Run("Found", func(t *testing.T) {
		truncateTables(t, pool)
		ctx := context.Background()

		params := randomCreateUserParams(t)
		createdUser, err := repo.CreateUser(ctx, params)
		require.NoError(t, err)

		foundUser, err := repo.GetUserByID(ctx, params.ID)
		require.NoError(t, err)
		require.NotNil(t, foundUser)

		assert.Equal(t, createdUser.ID, foundUser.ID)
		assert.Equal(t, createdUser.Email, foundUser.Email)
	})

	t.Run("NotFound", func(t *testing.T) {
		truncateTables(t, pool)
		ctx := context.Background()

		randomID := uuid.New()
		foundUser, err := repo.GetUserByID(ctx, randomID)
		require.Error(t, err)
		require.True(t, errors.Is(err, pgx.ErrNoRows), "expected pgx.ErrNoRows error")
		require.Nil(t, foundUser)
	})
}
