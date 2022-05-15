package db

type CACertificate struct {
	Id         int64  `db:"id"`
	Type       string `db:"type"`
	PublicKey  string `db:"public_key"`
	PrivateKey string `db:"private_key"`
}
