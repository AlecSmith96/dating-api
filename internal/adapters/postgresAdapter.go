package adapters

import (
	"database/sql"
	"errors"
	"github.com/AlecSmith96/dating-api/internal/entities"
	"github.com/AlecSmith96/dating-api/internal/usecases"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"log/slog"
	"time"
)

type PostgresAdapter struct {
	db              *sql.DB
	jwtExpiryMillis int
	jwtSecretKey    string
}

var _ usecases.UserCreator = &PostgresAdapter{}
var _ usecases.UserAuthenticator = &PostgresAdapter{}

func NewPostgresAdapter(db *sql.DB, jwtExpiryMillis int, jwtSecretKey string) *PostgresAdapter {
	return &PostgresAdapter{
		db:              db,
		jwtExpiryMillis: jwtExpiryMillis,
		jwtSecretKey:    jwtSecretKey,
	}
}

// PerformDataMigration is a function that ensure that the database has had all migration ran against it on startup
func (p *PostgresAdapter) PerformDataMigration(gooseDir string) error {
	return goose.Up(p.db, gooseDir)
}

func (p *PostgresAdapter) CreateUser(user *entities.User) (*entities.User, error) {
	var returnedUser entities.User
	err := p.db.QueryRow("INSERT INTO platform_user(email, password, name, gender, date_of_birth) VALUES ($1, $2, $3, $4, $5) RETURNING *;",
		user.Email,
		user.Password,
		user.Name,
		user.Gender,
		user.DateOfBirth).
		Scan(
			&returnedUser.ID,
			&returnedUser.Email,
			&returnedUser.Password,
			&returnedUser.Name,
			&returnedUser.Gender,
			&returnedUser.DateOfBirth,
		)
	if err != nil {
		slog.Debug("creating new user", "err", err)
		return nil, err
	}

	return &returnedUser, nil
}

func (p *PostgresAdapter) LoginUser(email string, password string) (*entities.User, error) {
	var returnedUser entities.User
	err := p.db.QueryRow("SELECT * FROM platform_user WHERE email = $1 AND password = $2 LIMIT 1;", email, password).
		Scan(
			&returnedUser.ID,
			&returnedUser.Email,
			&returnedUser.Password,
			&returnedUser.Name,
			&returnedUser.Gender,
			&returnedUser.DateOfBirth,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Debug("incorrect details for user", "email", email)
			return nil, entities.ErrUserNotFound
		}

		slog.Debug("authenticating user", "err", err)
		return nil, err
	}

	return &returnedUser, nil
}

type MyCustomClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func (p *PostgresAdapter) IssueJWT(userID uuid.UUID) (*entities.Token, error) {
	var returnedUser entities.User
	err := p.db.QueryRow("SELECT * FROM platform_user WHERE id = $1;", userID).
		Scan(
			&returnedUser.ID,
			&returnedUser.Email,
			&returnedUser.Password,
			&returnedUser.Name,
			&returnedUser.Gender,
			&returnedUser.DateOfBirth,
		)
	if err != nil {
		slog.Debug("getting user to issue jwt", "err", err)
		return nil, err
	}

	expirationTime := time.Now().Add(300000 * time.Millisecond)

	// Create the claims
	claims := MyCustomClaims{
		UserID: 12345,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "dating-api",
			Subject:   "user authentication",
			Audience:  jwt.ClaimStrings{"dating-api users"},
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "unique-id-12345",
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(p.jwtSecretKey)
	if err != nil {
		slog.Debug("failed to sign jwt", "err", err)
		return nil, err
	}

	// TODO: finish dis

	return &entities.Token{
		ID:       uuid.New().String(),
		UserID:   uuid.UUID{},
		Value:    tokenString,
		IssuedAt: time.Time{},
	}, nil
}
