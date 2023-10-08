package db

import (
	"context"
	"testing"

	"github.com/qwetu_petro/backend/utils"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	name := utils.RandomString(10)
	email := utils.RandomEmail()

	arg := CreateUserParams{
		Username:   name,
		Email:      email,
		Department: "ACCOUNTS",
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	return user
}

func TestQueries_CreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestQueries_GetUserByUserNameOrEmail(t *testing.T) {
	user := CreateRandomUser(t)

	args := GetUserByUserNameOrEmailParams{
		Username: user.Username,
		Email:    user.Email,
	}
	user2, err := testQueries.GetUserByUserNameOrEmail(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.Email, user2.Email)

}

func TestQueries_UpdateUser(t *testing.T) {
	user := CreateRandomUser(t)
	name := utils.RandomString(10)
	email := utils.RandomEmail()
	arg := UpdateUserParams{
		ID:       user.ID,
		Username: name,
		Email:    email,
	}
	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, arg.Username, user2.Username)
	require.Equal(t, arg.Email, user2.Email)

}

func TestQueries_DeleteUser(t *testing.T) {
	user := CreateRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)
	args := GetUserByUserNameOrEmailParams{
		Username: user.Username,
		Email:    user.Email,
	}
	user2, err := testQueries.GetUserByUserNameOrEmail(context.Background(), args)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, user2)
}
