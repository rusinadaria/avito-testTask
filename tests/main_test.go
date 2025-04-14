package tests


import (
	"github.com/stretchr/testify/suite"
	"avito-testTask/internal/handlers"
	"avito-testTask/internal/services"
	"avito-testTask/internal/repository"
	"testing"
	"database/sql"
	"os"
	"github.com/golang/mock/gomock"
	// "avito-testTask/internal/services/mocks"
	"log/slog"
)

type APITestSuite struct {
	suite.Suite

	db *sql.DB
	handler *handlers.Handler
	services *services.Service
	repos *repository.Repository

	mockCtrl    *gomock.Controller
	// mockService *mocks.MockService
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	connStr := "user=postgres password=root dbname=test_db sslmode=disable"
	// connStr := os.Getenv("user=postgres password=root dbname=test_shop sslmode=disable")
	db, err := sql.Open("postgres", connStr)
	s.db = db
	if err != nil {
		s.FailNow("Failed connect database", err)
	}

	s.initDeps()
}

func (s *APITestSuite) TearDownSuite() {
	if s.db != nil {
		err := s.db.Close()
		if err != nil {
			s.FailNow("Failed to close database connection", err)
		}
	}
}

func (s *APITestSuite) initDeps() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	s.repos = repository.NewRepository(s.db)
	s.services = services.NewService(s.repos)
	s.handler = handlers.NewHandler(s.services, logger)

	s.mockCtrl = gomock.NewController(s.T())
	// s.mockService = mocks.NewMockService(s.mockCtrl)
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}