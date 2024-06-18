package usecases_test

import (
	"github.com/AlecSmith96/dating-api/internal/drivers"
	mock_usecases "github.com/AlecSmith96/dating-api/mocks"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHandleUsers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Users Test Suite")
}

var (
	r                 *gin.Engine
	userCreator       *mock_usecases.MockUserCreator
	userDiscoverer    *mock_usecases.MockUserDiscoverer
	jwtProcessor      *mock_usecases.MockJwtProcessor
	userAuthenticator *mock_usecases.MockUserAuthenticator
	swipeRegister     *mock_usecases.MockSwipeRegister
)

var _ = BeforeSuite(func() {
	// Put gin in test mode
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(GinkgoT())
	userCreator = mock_usecases.NewMockUserCreator(ctrl)
	userDiscoverer = mock_usecases.NewMockUserDiscoverer(ctrl)
	jwtProcessor = mock_usecases.NewMockJwtProcessor(ctrl)
	userAuthenticator = mock_usecases.NewMockUserAuthenticator(ctrl)
	swipeRegister = mock_usecases.NewMockSwipeRegister(ctrl)

	r = drivers.NewRouter(
		userCreator,
		userAuthenticator,
		jwtProcessor,
		userDiscoverer,
		swipeRegister,
	)

	go func() {
		defer GinkgoRecover()
		err := http.ListenAndServe(":8080", r)
		if err != nil {
			Expect(err).ToNot(HaveOccurred())
		}
	}()
})
