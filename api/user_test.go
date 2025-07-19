package api

import (
	"testing"

	db "github.com/akkahshh24/go-simple-bank-app/db/sqlc"
	"github.com/akkahshh24/go-simple-bank-app/util"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomPassword(12)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomUsername(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomFullName(),
		Email:          util.RandomEmail(),
	}
	return
}
