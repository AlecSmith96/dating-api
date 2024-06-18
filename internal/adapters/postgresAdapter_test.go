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
)

func TestNewPostgresAdapter(t *testing.T) {
	g := NewWithT(t)
	db, _, err := sqlmock.New()
	g.Expect(err).ToNot(HaveOccurred())

	adapter := adapters.NewPostgresAdapter(db, 0, "something-secret")
	g.Expect(adapter).To(BeAssignableToTypeOf(&adapters.PostgresAdapter{}))

	defer db.Close()
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
