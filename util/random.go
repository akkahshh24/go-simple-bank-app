package util

import (
	"math/big"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
)

// RandomInt generates a random integer between min and max
func RandomInt() int64 {
	return gofakeit.Int64()
}

// Randomusername generates a random username.
func RandomUsername() string {
	return gofakeit.Username()
}

// RandomFullName generates a random first and last name.
func RandomFullName() string {
	return gofakeit.Name()
}

// RandomEmail generates a random email address.
func RandomEmail() string {
	return gofakeit.Email()
}

// RandomAmount generates a random integer between min and max of type pgtype.Numeric.
func RandomAmount(min, max int) pgtype.Numeric {
	return pgtype.Numeric{Int: big.NewInt(int64(gofakeit.IntRange(min, max))), Exp: -2, Valid: true}
}

// RandomCurrency generates a random currency code from a predefined list.
func RandomCurrency() string {
	currencies := []string{USD, EUR, CAD, GBP, JPY, SGD}
	return currencies[gofakeit.Number(0, len(currencies)-1)]
}

// RandomString generates a random string of length n.
func RandomString(n uint) string {
	return gofakeit.LetterN(n)
}

// RandomPassword generates a random password with specified criteria.
func RandomPassword(n int) string {
	return gofakeit.Password(true, true, true, true, false, n)
}
