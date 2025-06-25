package util

import "github.com/jackc/pgx/v5/pgtype"

func Int64ToPgTypeInt8(value int64) pgtype.Int8 {
	int8Value := pgtype.Int8{}
	int8Value.Scan(value)
	return int8Value
}
