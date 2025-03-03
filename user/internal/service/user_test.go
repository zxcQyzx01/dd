package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "dd/pkg/user"
)

func setupTest(t *testing.T) (*UserService, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	service := &UserService{db: db}

	cleanup := func() {
		db.Close()
	}

	return service, mock, cleanup
}

func TestUserService_CreateUser(t *testing.T) {
	service, mock, cleanup := setupTest(t)
	defer cleanup()

	t.Run("successful creation", func(t *testing.T) {
		// Подготавливаем мок для БД
		mock.ExpectQuery("INSERT INTO users").
			WithArgs("test@example.com", sqlmock.AnyArg()). // Пароль будет хэширован, поэтому используем AnyArg
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at"}).
				AddRow("123", "test@example.com", "hashed_password", time.Now()))

		// Выполняем запрос
		resp, err := service.CreateUser(context.Background(), &pb.CreateUserRequest{
			Email:    "test@example.com",
			Password: "password123",
		})

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "test@example.com", resp.User.Email)
		assert.Equal(t, "123", resp.User.Id)
	})

	t.Run("duplicate email", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO users").
			WithArgs("existing@example.com", sqlmock.AnyArg()).
			WillReturnError(&pq.Error{Code: "23505"}) // Имитируем ошибку дубликата

		resp, err := service.CreateUser(context.Background(), &pb.CreateUserRequest{
			Email:    "existing@example.com",
			Password: "password123",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestUserService_GetProfile(t *testing.T) {
	service, mock, cleanup := setupTest(t)
	defer cleanup()

	t.Run("existing user", func(t *testing.T) {
		createdAt := time.Now()
		mock.ExpectQuery("SELECT (.+) FROM users WHERE").
			WithArgs("test@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "created_at"}).
				AddRow("123", "test@example.com", createdAt))

		resp, err := service.GetProfile(context.Background(), &pb.GetProfileRequest{
			Email: "test@example.com",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "test@example.com", resp.User.Email)
		assert.Equal(t, "123", resp.User.Id)
	})

	t.Run("non-existing user", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users WHERE").
			WithArgs("nonexistent@example.com").
			WillReturnError(sql.ErrNoRows)

		resp, err := service.GetProfile(context.Background(), &pb.GetProfileRequest{
			Email: "nonexistent@example.com",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.NotFound, status.Code(err))
	})
}

func TestUserService_ListUsers(t *testing.T) {
	service, mock, cleanup := setupTest(t)
	defer cleanup()

	t.Run("successful listing", func(t *testing.T) {
		createdAt := time.Now()
		// Мокаем запрос на получение списка пользователей
		mock.ExpectQuery("SELECT (.+) FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "created_at"}).
				AddRow("1", "user1@example.com", createdAt).
				AddRow("2", "user2@example.com", createdAt))

		// Мокаем запрос на подсчет общего количества
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		resp, err := service.ListUsers(context.Background(), &pb.ListUsersRequest{
			Page:    1,
			PerPage: 10,
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Users, 2)
		assert.Equal(t, int32(2), resp.Total)
		assert.Equal(t, "user1@example.com", resp.Users[0].Email)
		assert.Equal(t, "user2@example.com", resp.Users[1].Email)
	})

	t.Run("empty list", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "created_at"}))

		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		resp, err := service.ListUsers(context.Background(), &pb.ListUsersRequest{
			Page:    1,
			PerPage: 10,
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Users, 0)
		assert.Equal(t, int32(0), resp.Total)
	})
}

func TestNew(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	service := New(db)
	assert.NotNil(t, service)
	assert.Equal(t, db, service.db)
}
