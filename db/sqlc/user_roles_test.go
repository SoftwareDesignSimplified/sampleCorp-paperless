package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomUserRoles(t *testing.T) UserRole {
	user := CreateRandomUser(t)
	role := createRandomRole(t)
	userRoles, err := testQueries.CreateUserRoles(context.Background(), CreateUserRolesParams{
		UserID: user.ID,
		RoleID: role.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, userRoles)
	return userRoles
}

func TestQueries_CreateUserRoles(t *testing.T) {
	createRandomUserRoles(t)
}

func TestQueries_GetUserRole(t *testing.T) {
	userRoles := createRandomUserRoles(t)
	userRoles2, err := testQueries.GetUserRole(context.Background(), GetUserRoleParams{
		UserID: userRoles.UserID,
		RoleID: userRoles.RoleID,
	})
	fmt.Println(userRoles2)
	require.NoError(t, err)
	require.NotEmpty(t, userRoles2)
	require.Equal(t, userRoles.ID, userRoles2.ID)
	require.Equal(t, userRoles.UserID, userRoles2.UserID)
	require.Equal(t, userRoles.RoleID, userRoles2.RoleID)
	require.Equal(t, userRoles.CreatedAt, userRoles2.CreatedAt)
	require.Equal(t, userRoles.TerminatedAt, userRoles2.TerminatedAt)
}

func TestQueries_UpdateUserRole(t *testing.T) {
	userRoles := createRandomUserRoles(t)
	userRoles2, err := testQueries.UpdateUserRoles(context.Background(), UpdateUserRolesParams{
		UserID: userRoles.UserID,
		RoleID: userRoles.RoleID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, userRoles2)
	require.Equal(t, userRoles.ID, userRoles2.ID)
	require.Equal(t, userRoles.UserID, userRoles2.UserID)
	require.Equal(t, userRoles.RoleID, userRoles2.RoleID)
	require.Equal(t, userRoles.CreatedAt, userRoles2.CreatedAt)
	require.Equal(t, userRoles.TerminatedAt, userRoles2.TerminatedAt)

}
