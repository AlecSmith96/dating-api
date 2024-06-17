package adapters

import (
	"database/sql"
	"fmt"
	"github.com/AlecSmith96/dating-api/internal/entities"
	_ "github.com/lib/pq"
	. "github.com/onsi/gomega"
	"github.com/pressly/goose/v3"
	"testing"
)

const (
	databaseConnectionString = "postgres://postgres:postgres@localhost:5432/users?sslmode=disable"
)

// SetUpMigrationTestDB is a function that creates a new database to run a migration test in
func SetUpMigrationTestDB(dbName string) (*sql.DB, error) {
	setupDb, err := sql.Open("postgres", databaseConnectionString)
	if err != nil {
		return nil, err
	}

	// Create a new database for this test
	_, err = setupDb.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	if err != nil {
		return nil, err
	}
	_, err = setupDb.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		return nil, err
	}
	err = setupDb.Close()
	if err != nil {
		return nil, err

	}
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", dbName))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func setupAdapter(g *WithT, db *sql.DB) (*PostgresAdapter, *sql.DB) {
	setupDb, err := sql.Open("postgres", databaseConnectionString)
	g.Expect(err).ToNot(HaveOccurred())

	return NewPostgresAdapter(setupDb, 3000000, "something-secret"), setupDb
}

func TestPostgresAdapter_PerformDataMigration_HappyPath(t *testing.T) {
	g := NewGomegaWithT(t)

	db, err := SetUpMigrationTestDB("perform_data_migration_test")
	g.Expect(err).ToNot(HaveOccurred())

	adapter, db := setupAdapter(g, db)

	initialVersion, err := goose.GetDBVersion(db)
	g.Expect(err).ToNot(HaveOccurred())

	err = adapter.PerformDataMigration("../../db/goose")
	g.Expect(err).ToNot(HaveOccurred())
	latestVersion, _ := goose.GetDBVersion(db)
	g.Expect(latestVersion).To(BeNumerically(">=", initialVersion))
}

func TestSchema(t *testing.T) {
	g := NewGomegaWithT(t)
	db, err := SetUpMigrationTestDB("schema")
	g.Expect(err).ToNot(HaveOccurred())

	err = goose.UpTo(db, "../../db/goose", 20240614195956) // initial migration
	g.Expect(err).ToNot(HaveOccurred())

	_, err = db.Exec("SELECT * FROM platform_user;")
	g.Expect(err).ToNot(HaveOccurred())

	_, err = db.Exec("SELECT * FROM token;")
	g.Expect(err).ToNot(HaveOccurred())
}

func TestAddUserSwipeAndUserMatchTables(t *testing.T) {
	g := NewGomegaWithT(t)
	db, err := SetUpMigrationTestDB("add_user_swipe_table")
	g.Expect(err).ToNot(HaveOccurred())

	err = goose.UpTo(db, "../../db/goose", 20240614195956) // initial migration
	g.Expect(err).ToNot(HaveOccurred())

	_, err = db.Exec("SELECT * FROM user_swipe;")
	g.Expect(err).To(MatchError("pq: relation \"user_swipe\" does not exist"))

	_, err = db.Exec("SELECT * FROM user_match;")
	g.Expect(err).To(MatchError("pq: relation \"user_match\" does not exist"))

	err = goose.UpTo(db, "../../db/goose", 20240615135847) // current migration
	g.Expect(err).ToNot(HaveOccurred())

	_, err = db.Exec("SELECT * FROM user_swipe;")
	g.Expect(err).ToNot(HaveOccurred())

	_, err = db.Exec("SELECT * FROM user_match;")
	g.Expect(err).ToNot(HaveOccurred())
}

func TestAddLocationFields(t *testing.T) {
	g := NewGomegaWithT(t)
	db, err := SetUpMigrationTestDB("add_user_swipe_table")
	g.Expect(err).ToNot(HaveOccurred())

	err = goose.UpTo(db, "../../db/goose", 20240615135847) // initial migration
	g.Expect(err).ToNot(HaveOccurred())

	row, err := db.Query("SELECT * FROM platform_user WHERE email = 'admin';")
	g.Expect(err).ToNot(HaveOccurred())

	if row.Next() {
		var user entities.User
		err = row.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.Gender,
			&user.DateOfBirth,
			&user.Location.Latitude,
			&user.Location.Longitude,
		)
		g.Expect(err).To(MatchError("sql: expected 6 destination arguments in Scan, not 8"))
	}

	err = goose.UpTo(db, "../../db/goose", 20240617164832) // current migration
	g.Expect(err).ToNot(HaveOccurred())

	row, err = db.Query("SELECT * FROM platform_user WHERE email = 'admin';")
	g.Expect(err).ToNot(HaveOccurred())

	if row.Next() {
		var user entities.User
		err = row.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.Gender,
			&user.DateOfBirth,
			&user.Location.Latitude,
			&user.Location.Longitude,
		)
		g.Expect(err).ToNot(HaveOccurred())
	}
}
