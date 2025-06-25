package util

import (
	"math/big"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
)

// RandomName generates a random name.
func RandomName() string {
	return gofakeit.Name()
}

// RandomAmount generates a random integer between min and max of type pgtype.Numeric.
func RandomAmount(min, max int) pgtype.Numeric {
	return pgtype.Numeric{Int: big.NewInt(int64(gofakeit.IntRange(min, max))), Exp: -2, Valid: true}
}

// RandomCurrency generates a random currency code from a predefined list.
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD", "GBP", "JPY", "SGD"}
	return currencies[gofakeit.Number(0, len(currencies)-1)]
}
