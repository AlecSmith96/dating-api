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
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
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

	issuedAt := time.Now()
	expirationTime := issuedAt.Add(300000 * time.Millisecond)

	claims := MyCustomClaims{
		UserID: returnedUser.ID.String(),
		Email:  returnedUser.Email,
		Name:   returnedUser.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "dating-api",
			Subject:   "user authentication",
			Audience:  jwt.ClaimStrings{"dating-api users"},
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "unique-id-12345",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(p.jwtSecretKey))
	if err != nil {
		slog.Debug("failed to sign jwt", "err", err)
		return nil, err
	}

	var returnedToken entities.Token
	err = p.db.QueryRow("INSERT INTO token (user_id, value, issued_at) VALUES ($1, $2, $3) RETURNING *;", returnedUser.ID, tokenString, issuedAt).
		Scan(&returnedToken.ID, &returnedToken.UserID, &returnedToken.Value, &returnedToken.IssuedAt)
	if err != nil {
		slog.Debug("writing token to storage", "err", err)
		return nil, err
	}

	return &entities.Token{
		ID:       returnedToken.ID,
		UserID:   returnedToken.UserID,
		Value:    returnedToken.Value,
		IssuedAt: returnedToken.IssuedAt,
	}, nil
}
