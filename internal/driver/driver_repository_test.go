package driver

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Mikhail-Tal63/Orbit/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/magiconair/properties/assert"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, *DriverRepositoryImpl) {
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
	repo := NewDriverRepository(queries)

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
	_, err := pool.Exec(context.Background(), "TRUNCATE TABLE drivers CASCADE;")
	require.NoError(t, err, "failed to truncate users table")
}

func createTestUser(t *testing.T, pool *pgxpool.Pool) db.User {
	t.Helper()

	queries := db.New(pool)

	user, err := queries.CreateUser(context.Background(), db.CreateUserParams{
		ID:           uuid.New(),
		FirstName:    "John",
		LastName:     "Doe",
		Username:     "john_" + uuid.New().String()[:8],
		Email:        uuid.New().String()[:8] + "@test.com",
		PasswordHash: "hashed-password",
	})

	require.NoError(t, err)

	return user
}
func TestDriverRepository_CreateDriver(t *testing.T) {
	pool, repo := setupTestDB(t)
	if pool == nil {
		return
	}

	t.Run("success", func(t *testing.T) {
		truncateTables(t, pool)

		ctx := context.Background()

		user := createTestUser(t, pool)

		params := db.CreateDriverParams{
			ID:     uuid.New(),
			UserID: user.ID,
		}

		driver, err := repo.CreateDriver(ctx, params)

		require.NoError(t, err)
		require.NotNil(t, driver)

		assert.Equal(t, params.ID, driver.ID)
		assert.Equal(t, params.UserID, driver.UserID)
	})

	t.Run("duplicate user", func(t *testing.T) {
		truncateTables(t, pool)

		ctx := context.Background()

		user := createTestUser(t, pool)

		params := db.CreateDriverParams{
			ID:     uuid.New(),
			UserID: user.ID,
		}

		_, err := repo.CreateDriver(ctx, params)
		require.NoError(t, err)

		_, err = repo.CreateDriver(ctx, db.CreateDriverParams{
			ID:     uuid.New(),
			UserID: user.ID,
		})

		require.Error(t, err)
	})

	t.Run("user does not exist", func(t *testing.T) {
		truncateTables(t, pool)

		ctx := context.Background()

		params := db.CreateDriverParams{
			ID:     uuid.New(),
			UserID: uuid.New(),
		}

		_, err := repo.CreateDriver(ctx, params)

		require.Error(t, err)
	})
}

//fuck github actions
