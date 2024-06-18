package usecases_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlecSmith96/dating-api/internal/entities"
	"github.com/AlecSmith96/dating-api/internal/usecases"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
)

const (
	mockJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYWI4ZjFiZTAtOWJjMC00YjAxLTg1NzgtMzlkNWM0MDlmMWFlIiwiZW1haWwiOiJhZG1pbiIsIm5hbWUiOiJhZG1pbiIsImlzcyI6ImRhdGluZy1hcGkiLCJzdWIiOiJ1c2VyIGF1dGhlbnRpY2F0aW9uIiwiYXVkIjpbImRhdGluZy1hcGkgdXNlcnMiXSwiZXhwIjoxNzE4NjMyOTI5LCJpYXQiOjE3MTg2MzI2MjksImp0aSI6InVuaXF1ZS1pZC0xMjM0NSJ9.wWukfV2Pc_qFCAxyK6v6g4U8daODLZqaF4wXU4XHGOE"
)

var _ = Describe("creating a user", func() {
	var w *httptest.ResponseRecorder

	var validateJwtForUserUUID uuid.UUID
	var validateJwtForUserErr error
	var validateJwtForUserCallCount int

	var createUserResponse *entities.User
	var createUserErr error
	var createUserCallCount int

	BeforeEach(func() {
		validateJwtForUserUUID = uuid.New()
		validateJwtForUserErr = nil
		validateJwtForUserCallCount = 1

		createUserResponse = &entities.User{
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
		createUserErr = nil
		createUserCallCount = 1
	})

	JustBeforeEach(func() {
		w = httptest.NewRecorder()

		jwtProcessor.EXPECT().ValidateJwtForUser(mockJWT).Return(validateJwtForUserUUID, validateJwtForUserErr).Times(validateJwtForUserCallCount)
		userCreator.EXPECT().CreateUser(gomock.AssignableToTypeOf(&entities.User{})).Return(createUserResponse, createUserErr).Times(createUserCallCount)

		req, err := http.NewRequest("POST", "http://localhost:8080/dating-api/v1/user/create", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", mockJWT))
		Expect(err).ToNot(HaveOccurred())
		r.ServeHTTP(w, req)
	})

	It("should return the newly generated user", func() {
		Expect(w.Code).To(Equal(http.StatusOK))
		var user usecases.CreateUserResponseBody
		err := json.NewDecoder(w.Body).Decode(&user)
		Expect(err).ToNot(HaveOccurred())

		Expect(user.ID).To(Equal(createUserResponse.ID.String()))
		Expect(user.Email).To(Equal(createUserResponse.Email))
		Expect(user.Password).To(Equal(createUserResponse.Password))
		Expect(user.Name).To(Equal(createUserResponse.Name))
		Expect(user.Gender).To(Equal(createUserResponse.Gender))
		Expect(user.Age).To(Equal(createUserResponse.GetAge()))
		Expect(user.Location.Latitude).To(Equal(createUserResponse.Location.Latitude))
		Expect(user.Location.Longitude).To(Equal(createUserResponse.Location.Longitude))
	})

	When("the jwt cannot be validated", func() {
		BeforeEach(func() {
			validateJwtForUserUUID = uuid.UUID{}
			validateJwtForUserErr = errors.New("unable to validate jwt")
			createUserCallCount = 0
		})

		It("should return a 401 Unauthorized error", func() {
			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	When("the adapter returns an error", func() {
		BeforeEach(func() {
			createUserResponse = nil
			createUserErr = errors.New("an error occurred")
		})

		It("should return a 500 Internal Server Error", func() {
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
