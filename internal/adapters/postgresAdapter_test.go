package adapters_test

import (
	"database/sql"
	"errors"
	"github.com/AlecSmith96/dating-api/internal/adapters"
	"github.com/AlecSmith96/dating-api/internal/entities"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

const (
	mockJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYWI4ZjFiZTAtOWJjMC00YjAxLTg1NzgtMzlkNWM0MDlmMWFlIiwiZW1haWwiOiJhZG1pbiIsIm5hbWUiOiJhZG1pbiIsImlzcyI6ImRhdGluZy1hcGkiLCJzdWIiOiJ1c2VyIGF1dGhlbnRpY2F0aW9uIiwiYXVkIjpbImRhdGluZy1hcGkgdXNlcnMiXSwiZXhwIjoxNzE4NjMyOTI5LCJpYXQiOjE3MTg2MzI2MjksImp0aSI6InVuaXF1ZS1pZC0xMjM0NSJ9.wWukfV2Pc_qFCAxyK6v6g4U8daODLZqaF4wXU4XHGOE"
)

func TestNewPostgresAdapter(t *testing.T) {
	g := NewWithT(t)
	db, _, err := sqlmock.New()
	g.Expect(err).ToNot(HaveOccurred())

	adapter := adapters.NewPostgresAdapter(db, 0, "something-secret")
	g.Expect(adapter).To(BeAssignableToTypeOf(&adapters.PostgresAdapter{}))

	defer db.Close()
}

func TestPostgresAdapter_LoginUser(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	g.Expect(err).ToNot(HaveOccurred())

	adapter := adapters.NewPostgresAdapter(db, 0, "something-secret")

	user := &entities.User{
		ID:          uuid.New(),
		Email:       gofakeit.Email(),
		Password:    gofakeit.Password(true, true, true, true, true, 15),
		Name:        gofakeit.Name(),
		Gender:      gofakeit.Gender(),
		DateOfBirth: gofakeit.Date(),
		Location: entities.Location{
			Latitude:  gofakeit.Address().Latitude,
			Longitude: gofakeit.Address().Longitude,
		},
	}

	mock.ExpectQuery(`SELECT \* FROM platform_user WHERE email = \$1 AND password = \$2 LIMIT 1;`).WithArgs(user.Email, user.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "name", "gender", "date_of_birth", "location_latitude", "location_longitude"}).
			AddRow(user.ID, user.Email, user.Password, user.Name, user.Gender, user.DateOfBirth, user.Location.Latitude, user.Location.Longitude))

	_, err = adapter.LoginUser(user.Email, user.Password)
	g.Expect(err).ToNot(HaveOccurred())
}

func TestPostgresAdapter_ErrNoRows(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	g.Expect(err).ToNot(HaveOccurred())

	adapter := adapters.NewPostgresAdapter(db, 0, "something-secret")

	user := &entities.User{
		ID:          uuid.New(),
		Email:       gofakeit.Email(),
		Password:    gofakeit.Password(true, true, true, true, true, 15),
		Name:        gofakeit.Name(),
		Gender:      gofakeit.Gender(),
		DateOfBirth: gofakeit.Date(),
		Location: entities.Location{
			Latitude:  gofakeit.Address().Latitude,
			Longitude: gofakeit.Address().Longitude,
		},
	}

	mock.ExpectQuery(`SELECT \* FROM platform_user WHERE email = \$1 AND password = \$2 LIMIT 1;`).WithArgs(user.Email, user.Password).
		WillReturnError(sql.ErrNoRows)

	_, err = adapter.LoginUser(user.Email, user.Password)
	g.Expect(err).To(MatchError(entities.ErrUserNotFound))
}

func TestPostgresAdapter_GenericErr(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	g.Expect(err).ToNot(HaveOccurred())

	adapter := adapters.NewPostgresAdapter(db, 0, "something-secret")

	user := &entities.User{
		ID:          uuid.New(),
		Email:       gofakeit.Email(),
		Password:    gofakeit.Password(true, true, true, true, true, 15),
		Name:        gofakeit.Name(),
		Gender:      gofakeit.Gender(),
		DateOfBirth: gofakeit.Date(),
		Location: entities.Location{
			Latitude:  gofakeit.Address().Latitude,
			Longitude: gofakeit.Address().Longitude,
		},
	}

	mock.ExpectQuery(`SELECT \* FROM platform_user WHERE email = \$1 AND password = \$2 LIMIT 1;`).WithArgs(user.Email, user.Password).
		WillReturnError(errors.New("an error occurred"))

	_, err = adapter.LoginUser(user.Email, user.Password)
	g.Expect(err).To(MatchError("an error occurred"))
}

func TestPostgresAdapter_IssueJWT(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	g.Expect(err).ToNot(HaveOccurred())

	adapter := adapters.NewPostgresAdapter(db, 0, "something-secret")

	user := &entities.User{
		ID:          uuid.New(),
		Email:       gofakeit.Email(),
		Password:    gofakeit.Password(true, true, true, true, true, 15),
		Name:        gofakeit.Name(),
		Gender:      gofakeit.Gender(),
		DateOfBirth: gofakeit.Date(),
		Location: entities.Location{
			Latitude:  gofakeit.Address().Latitude,
			Longitude: gofakeit.Address().Longitude,
		},
	}

	token := &entities.Token{
		ID:       uuid.New().String(),
		UserID:   user.ID,
		Value:    mockJWT,
		IssuedAt: time.Now(),
	}

	mock.ExpectQuery(`SELECT \* FROM platform_user WHERE platform_user.id = \$1;`).WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "name", "gender", "date_of_birth", "location_latitude", "location_longitude"}).
			AddRow(user.ID, user.Email, user.Password, user.Name, user.Gender, user.DateOfBirth, user.Location.Latitude, user.Location.Longitude))
	mock.ExpectQuery(`INSERT INTO token \(user_id, value, issued_at\) VALUES \(\$1, \$2, \$3\) RETURNING \*;`).WithArgs(user.ID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "value", "issued_at"}).
			AddRow(token.ID, token.UserID, token.Value, token.IssuedAt))

	returnedToken, err := adapter.IssueJWT(user.ID)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(returnedToken).To(Equal(token))
}

func TestPostgresAdapter_IssueJWT_GetUserErr(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	g.Expect(err).ToNot(HaveOccurred())

	adapter := adapters.NewPostgresAdapter(db, 0, "something-secret")

	user := &entities.User{
		ID:          uuid.New(),
		Email:       gofakeit.Email(),
		Password:    gofakeit.Password(true, true, true, true, true, 15),
		Name:        gofakeit.Name(),
		Gender:      gofakeit.Gender(),
		DateOfBirth: gofakeit.Date(),
		Location: entities.Location{
			Latitude:  gofakeit.Address().Latitude,
			Longitude: gofakeit.Address().Longitude,
		},
	}

	mock.ExpectQuery(`SELECT \* FROM platform_user WHERE platform_user.id = \$1;`).WithArgs(user.ID).
		WillReturnError(errors.New("an error occurred"))

	returnedToken, err := adapter.IssueJWT(user.ID)
	g.Expect(err).To(MatchError("an error occurred"))
	g.Expect(returnedToken).To(BeNil())
}

func TestPostgresAdapter_IssueJWT_CreatingTokenErr(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	g.Expect(err).ToNot(HaveOccurred())

	adapter := adapters.NewPostgresAdapter(db, 0, "something-secret")

	user := &entities.User{
		ID:          uuid.New(),
		Email:       gofakeit.Email(),
		Password:    gofakeit.Password(true, true, true, true, true, 15),
		Name:        gofakeit.Name(),
		Gender:      gofakeit.Gender(),
		DateOfBirth: gofakeit.Date(),
		Location: entities.Location{
			Latitude:  gofakeit.Address().Latitude,
			Longitude: gofakeit.Address().Longitude,
		},
	}

	mock.ExpectQuery(`SELECT \* FROM platform_user WHERE platform_user.id = \$1;`).WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "name", "gender", "date_of_birth", "location_latitude", "location_longitude"}).
			AddRow(user.ID, user.Email, user.Password, user.Name, user.Gender, user.DateOfBirth, user.Location.Latitude, user.Location.Longitude))
	mock.ExpectQuery(`INSERT INTO token \(user_id, value, issued_at\) VALUES \(\$1, \$2, \$3\) RETURNING \*;`).WithArgs(user.ID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("an error occurred"))

	returnedToken, err := adapter.IssueJWT(user.ID)
	g.Expect(err).To(MatchError("an error occurred"))
	g.Expect(returnedToken).To(BeNil())
}

func TestPostgresAdapter_CreateUser(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	g.Expect(err).ToNot(HaveOccurred())

	adapter := adapters.NewPostgresAdapter(db, 0, "something-secret")

	user := &entities.User{
		ID:          uuid.New(),
		Email:       gofakeit.Email(),
		Password:    gofakeit.Password(true, true, true, true, true, 15),
		Name:        gofakeit.Name(),
		Gender:      gofakeit.Gender(),
		DateOfBirth: gofakeit.Date(),
		Location: entities.Location{
			Latitude:  gofakeit.Address().Latitude,
			Longitude: gofakeit.Address().Longitude,
		},
	}

	mock.ExpectQuery(`INSERT INTO platform_user\(email, password, name, gender, date_of_birth, location_latitude, location_longitude\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\) RETURNING \*;`).
		WithArgs(user.Email, user.Password, user.Name, user.Gender, user.DateOfBirth, user.Location.Latitude, user.Location.Longitude).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "name", "gender", "date_of_birth", "location_latitude", "location_longitude"}).
		AddRow(user.ID, user.Email, user.Password, user.Name, user.Gender, user.DateOfBirth, user.Location.Latitude, user.Location.Longitude))

	userResp, err := adapter.CreateUser(user)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(userResp).To(Equal(user))
}

func TestPostgresAdapter_CreateUser_ReturnsErr(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	g.Expect(err).ToNot(HaveOccurred())

	adapter := adapters.NewPostgresAdapter(db, 0, "something-secret")

	user := &entities.User{
		ID:          uuid.New(),
		Email:       gofakeit.Email(),
		Password:    gofakeit.Password(true, true, true, true, true, 15),
		Name:        gofakeit.Name(),
		Gender:      gofakeit.Gender(),
		DateOfBirth: gofakeit.Date(),
		Location: entities.Location{
			Latitude:  gofakeit.Address().Latitude,
			Longitude: gofakeit.Address().Longitude,
		},
	}

	mock.ExpectQuery(`INSERT INTO platform_user\(email, password, name, gender, date_of_birth, location_latitude, location_longitude\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\) RETURNING \*;`).
		WithArgs(user.Email, user.Password, user.Name, user.Gender, user.DateOfBirth, user.Location.Latitude, user.Location.Longitude).WillReturnError(errors.New("an error occurred"))

	userResp, err := adapter.CreateUser(user)
	g.Expect(err).To(MatchError("an error occurred"))
	g.Expect(userResp).To(BeNil())
}

func TestPostgresAdapter_DiscoverNewUsers(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	g.Expect(err).ToNot(HaveOccurred())

	adapter := adapters.NewPostgresAdapter(db, 0, "something-secret")

	ownerUserID := uuid.New()
	pageInfo := entities.PageInfo{
		MinAge:           18,
		MaxAge:           80,
		PreferredGenders: []string{"female"},
	}

	users := []entities.UserDiscovery{
		{
			ID:          uuid.New(),
			Email:       gofakeit.Email(),
			Password:    gofakeit.Password(true, true, true, true, true, 15),
			Name:        gofakeit.Name(),
			Gender:      gofakeit.Gender(),
			DateOfBirth: gofakeit.Date(),
			Age:         23,
			Location: entities.Location{
				Latitude:  gofakeit.Address().Latitude,
				Longitude: gofakeit.Address().Longitude,
			},
		},
		{
			ID:          uuid.New(),
			Email:       gofakeit.Email(),
			Password:    gofakeit.Password(true, true, true, true, true, 15),
			Name:        gofakeit.Name(),
			Gender:      gofakeit.Gender(),
			DateOfBirth: gofakeit.Date(),
			Age:         23,
			Location: entities.Location{
				Latitude:  gofakeit.Address().Latitude,
				Longitude: gofakeit.Address().Longitude,
			},
		},
	}

	mock.ExpectQuery("SELECT pu\\.\\* FROM \\( SELECT pu\\.\\*, DATE_PART\\('year', AGE\\(pu\\.date_of_birth\\)\\) AS age FROM platform_user pu \\) pu LEFT JOIN user_swipe us ON pu\\.id = us\\.swiped_user_id AND us\\.owner_user_id = \\$1 WHERE pu\\.id != \\$1 AND us\\.id IS NULL AND pu\\.age >= \\$2 AND pu\\.age <= \\$3 AND pu\\.gender IN \\(\\$4\\);").
		WithArgs(ownerUserID, pageInfo.MinAge, pageInfo.MaxAge, pageInfo.PreferredGenders[0]).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "name", "gender", "date_of_birth", "location_latitude", "location_longitude", "age"}).
			AddRow(users[0].ID, users[0].Email, users[0].Password, users[0].Name, users[0].Gender, users[0].DateOfBirth, users[0].Location.Latitude, users[0].Location.Longitude, users[0].Age).
			AddRow(users[1].ID, users[1].Email, users[1].Password, users[1].Name, users[1].Gender, users[1].DateOfBirth, users[1].Location.Latitude, users[1].Location.Longitude, users[0].Age))

	returnedUsers, err := adapter.DiscoverNewUsers(ownerUserID, pageInfo)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(returnedUsers).To(HaveLen(2))
}

func TestPostgresAdapter_DiscoverNewUsers_ErrNoRows(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	g.Expect(err).ToNot(HaveOccurred())

	adapter := adapters.NewPostgresAdapter(db, 0, "something-secret")

	ownerUserID := uuid.New()
	pageInfo := entities.PageInfo{
		MinAge:           18,
		MaxAge:           80,
		PreferredGenders: []string{"female"},
	}

	mock.ExpectQuery("SELECT pu\\.\\* FROM \\( SELECT pu\\.\\*, DATE_PART\\('year', AGE\\(pu\\.date_of_birth\\)\\) AS age FROM platform_user pu \\) pu LEFT JOIN user_swipe us ON pu\\.id = us\\.swiped_user_id AND us\\.owner_user_id = \\$1 WHERE pu\\.id != \\$1 AND us\\.id IS NULL AND pu\\.age >= \\$2 AND pu\\.age <= \\$3 AND pu\\.gender IN \\(\\$4\\);").
		WithArgs(ownerUserID, pageInfo.MinAge, pageInfo.MaxAge, pageInfo.PreferredGenders[0]).
		WillReturnError(sql.ErrNoRows)

	returnedUsers, err := adapter.DiscoverNewUsers(ownerUserID, pageInfo)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(returnedUsers).To(HaveLen(0))
}

func TestPostgresAdapter_DiscoverNewUsers_GenericErr(t *testing.T) {
	g := NewWithT(t)
	db, mock, err := sqlmock.New()
	g.Expect(err).ToNot(HaveOccurred())

	adapter := adapters.NewPostgresAdapter(db, 0, "something-secret")

	ownerUserID := uuid.New()
	pageInfo := entities.PageInfo{
		MinAge:           18,
		MaxAge:           80,
		PreferredGenders: []string{"female"},
	}

	mock.ExpectQuery("SELECT pu\\.\\* FROM \\( SELECT pu\\.\\*, DATE_PART\\('year', AGE\\(pu\\.date_of_birth\\)\\) AS age FROM platform_user pu \\) pu LEFT JOIN user_swipe us ON pu\\.id = us\\.swiped_user_id AND us\\.owner_user_id = \\$1 WHERE pu\\.id != \\$1 AND us\\.id IS NULL AND pu\\.age >= \\$2 AND pu\\.age <= \\$3 AND pu\\.gender IN \\(\\$4\\);").
		WithArgs(ownerUserID, pageInfo.MinAge, pageInfo.MaxAge, pageInfo.PreferredGenders[0]).
		WillReturnError(errors.New("an error occurred"))

	returnedUsers, err := adapter.DiscoverNewUsers(ownerUserID, pageInfo)
	g.Expect(err).To(MatchError("an error occurred"))
	g.Expect(returnedUsers).To(BeNil())
}
