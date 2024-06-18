package adapters_test

import (
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
