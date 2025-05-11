package util

import (
	"math/big"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
)

// RandomHolderName generates a random name.
func RandomHolderName() string {
	return gofakeit.Name()
}

// RandomBalance generates a random balance between min and max of type pgtype.Numeric.
func RandomBalance(min, max int) pgtype.Numeric {
	return pgtype.Numeric{Int: big.NewInt(int64(gofakeit.IntRange(min, max))), Exp: -2, Valid: true}
}

// RandomCurrency generates a random currency code.
// The currency code is a three-letter ISO 4217 code, such as "USD" for US dollars or "EUR" for euros.
func RandomCurrency() string {
	return gofakeit.CurrencyShort()
}
