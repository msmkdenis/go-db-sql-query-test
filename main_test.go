package main

import (
	"database/sql"
	"fmt"
	"testing"

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
		client, err := selectClient(s.db, clientID)
		require.NoError(t, err)
		require.Equal(t, clientID, client.ID)
		require.NotEmpty(t, client.FIO)
		require.NotEmpty(t, client.Login)
		require.NotEmpty(t, client.Birthday)
		require.NotEmpty(t, client.Email)
	})
}

func (s *SQLiteSuite) Test_SelectClient_WhenNoClient() {
	clientID := -1
	s.T().Run("Fail when no client", func(t *testing.T) {
		client, err := selectClient(s.db, clientID)
		require.Equal(t, err, sql.ErrNoRows)
		require.Empty(t, client.ID)
		require.Empty(t, client.FIO)
		require.Empty(t, client.Login)
		require.Empty(t, client.Birthday)
		require.Empty(t, client.Email)
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
		id, err := insertClient(s.db, cl)
		cl.ID = id
		require.NotEmpty(t, id)
		require.NoError(t, err)

		client, err := selectClient(s.db, id)
		require.NoError(t, err)
		require.Equal(t, cl.ID, client.ID)
		require.Equal(t, cl.FIO, client.FIO)
		require.Equal(t, cl.Login, client.Login)
		require.Equal(t, cl.Birthday, client.Birthday)
		require.Equal(t, cl.Email, client.Email)
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
		id, err := insertClient(s.db, cl)
		require.NotEmpty(t, id)
		require.NoError(t, err)

		client, err := selectClient(s.db, id)
		require.NoError(t, err)

		err = deleteClient(s.db, client.ID)
		require.NoError(t, err)

		_, err = selectClient(s.db, client.ID)
		require.Equal(t, err, sql.ErrNoRows)
	})
}
