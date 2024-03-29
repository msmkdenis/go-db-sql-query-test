package main

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	_ "modernc.org/sqlite"
)

type SQLiteSuite struct {
	suite.Suite
	db         *sql.DB
	driverName string
	dbName     string
}

func (s *SQLiteSuite) SetupSuite() {
	s.driverName = "sqlite"
	s.dbName = "demo.db"
	db, err := sql.Open(s.driverName, s.dbName)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.db = db
}

func (s *SQLiteSuite) TearDownSuite() {
	s.db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(SQLiteSuite))
}

func (s *SQLiteSuite) Test_SelectClient_WhenOk() {
	clientID := 1
	s.T().Run("Successful select", func(t *testing.T) {
		// Получить объект клиента функцией selectClient(), если функция вернула ошибку — завершить тест.
		client, err := selectClient(s.db, clientID)
		require.NoError(t, err)

		assert.Equal(t, clientID, client.ID)
		assert.NotEmpty(t, client.FIO)
		assert.NotEmpty(t, client.Login)
		assert.NotEmpty(t, client.Birthday)
		assert.NotEmpty(t, client.Email)
	})
}

func (s *SQLiteSuite) Test_SelectClient_WhenNoClient() {
	clientID := -1
	s.T().Run("Fail when no client", func(t *testing.T) {
		// Проверить, что функция вернула ошибку и ошибка равна sql.ErrNoRows. Иначе завершить тест.
		client, err := selectClient(s.db, clientID)
		require.Equal(t, err, sql.ErrNoRows)

		assert.Empty(t, client.ID)
		assert.Empty(t, client.FIO)
		assert.Empty(t, client.Login)
		assert.Empty(t, client.Birthday)
		assert.Empty(t, client.Email)
	})
}

func (s *SQLiteSuite) Test_InsertClient_ThenSelectAndCheck() {
	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}
	s.T().Run("Successful insert and then select", func(t *testing.T) {
		// Проверить, что функция вернула не пустой идентификатор и пустую ошибку. Иначе завершить тест.
		id, err := insertClient(s.db, cl)
		cl.ID = id
		require.NotEmpty(t, id)
		require.NoError(t, err)

		// Функцией selectClient() получить объект Client по идентификатору.
		// Проверить, что функция вернула пустую ошибку. Иначе завершить тест.
		client, err := selectClient(s.db, id)
		require.NoError(t, err)

		assert.Equal(t, cl.ID, client.ID)
		assert.Equal(t, cl.FIO, client.FIO)
		assert.Equal(t, cl.Login, client.Login)
		assert.Equal(t, cl.Birthday, client.Birthday)
		assert.Equal(t, cl.Email, client.Email)
	})
}

func (s *SQLiteSuite) Test_InsertClient_DeleteClient_ThenCheck() {
	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}
	s.T().Run("Successful insert, delete and then check deleted", func(t *testing.T) {
		// Проверить, что функция вернула не пустой идентификатор и пустую ошибку. Иначе завершить тест.
		id, err := insertClient(s.db, cl)
		require.NotEmpty(t, id)
		require.NoError(t, err)

		// Получить объект клиента функцией selectClient(). Если функция вернула ошибку, завершить тест.
		client, err := selectClient(s.db, id)
		require.NoError(t, err)

		// Удалить запись функцией deleteClient(). Если функция вернула ошибку, завершить тест.
		err = deleteClient(s.db, client.ID)
		require.NoError(t, err)

		// Получить объект клиента функцией selectClient().
		// Проверить, что функция вернула ошибку и ошибка равна sql.ErrNoRows. Иначе завершить тест.
		_, err = selectClient(s.db, client.ID)
		require.Equal(t, err, sql.ErrNoRows)
	})
}
