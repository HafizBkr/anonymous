package commons

import (
	"github.com/jmoiron/sqlx"
)

type TxProvider interface {
	Provide() (*sqlx.Tx, error)
}

type Logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
	Info(msg string, args ...any)
}

type JWTProvider interface {
	Encode(claims map[string]interface{}) (string, error)
	Decode(token string) (map[string]interface{}, error)
}
