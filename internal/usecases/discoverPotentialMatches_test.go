package usecases_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlecSmith96/dating-api/internal/entities"
	"github.com/AlecSmith96/dating-api/internal/usecases"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("discovering potential matches", func() {
	var w *httptest.ResponseRecorder
	var requestBody *usecases.DiscoverPotentialMatchesRequestBody
	var requestBodyJSON []byte

	var validateJwtForUserUUID uuid.UUID
	var validateJwtForUserErr error
	var validateJwtForUserCallCount int

	var discoverNewUsersResponse []entities.UserDiscovery
	var discoverNewUsersErr error
	var discoverNewUsersCallCount int

	var getUsersLocationResponse *entities.Location
	var getUsersLocationErr error
	var getUsersLocationCallCount int

	BeforeEach(func() {
		requestBody = &usecases.DiscoverPotentialMatchesRequestBody{}
		var err error
		requestBodyJSON, err = json.Marshal(requestBody)
		Expect(err).ToNot(HaveOccurred())

		validateJwtForUserUUID = uuid.New()
		validateJwtForUserErr = nil
		validateJwtForUserCallCount = 1

		discoverNewUsersResponse = []entities.UserDiscovery{
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

		discoverNewUsersErr = nil
		discoverNewUsersCallCount = 1

		getUsersLocationResponse = &entities.Location{
			Latitude:  gofakeit.Address().Latitude,
			Longitude: gofakeit.Address().Longitude,
		}
		getUsersLocationErr = nil
		getUsersLocationCallCount = 1
	})

	JustBeforeEach(func() {
		w = httptest.NewRecorder()

		jwtProcessor.EXPECT().ValidateJwtForUser(mockJWT).Return(validateJwtForUserUUID, validateJwtForUserErr).Times(validateJwtForUserCallCount)
		userDiscoverer.EXPECT().DiscoverNewUsers(validateJwtForUserUUID, entities.PageInfo{}).Return(discoverNewUsersResponse, discoverNewUsersErr).Times(discoverNewUsersCallCount)
		userDiscoverer.EXPECT().GetUsersLocation(validateJwtForUserUUID).Return(getUsersLocationResponse, getUsersLocationErr).Times(getUsersLocationCallCount)

		req, err := http.NewRequest("GET", "http://localhost:8080/dating-api/v1/user/discover", bytes.NewReader(requestBodyJSON))
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", mockJWT))
		Expect(err).ToNot(HaveOccurred())
		r.ServeHTTP(w, req)
	})

	It("should return an array of other users", func() {
		Expect(w.Code).To(Equal(http.StatusOK))
		var resp usecases.DiscoverPotentialMatchesResponseBody
		err := json.NewDecoder(w.Body).Decode(&resp)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Users).To(HaveLen(2))
	})

	When("the request fails to validate", func() {
		BeforeEach(func() {
			requestBodyJSON = []byte("{")
			discoverNewUsersCallCount = 0
			getUsersLocationCallCount = 0
		})

		It("should return a 400 Bad Request", func() {
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	When("getting new users returns an error", func() {
		BeforeEach(func() {
			discoverNewUsersResponse = nil
			discoverNewUsersErr = errors.New("an error occurred")
			discoverNewUsersCallCount = 1

			getUsersLocationCallCount = 0
		})

		It("should return a 500 Internal Server Error", func() {
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	When("getting the requesting users location returns an error", func() {
		BeforeEach(func() {
			getUsersLocationResponse = nil
			getUsersLocationErr = errors.New("an error occurred")
			getUsersLocationCallCount = 1
		})

		It("should return a 500 Internal Server Error", func() {
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
