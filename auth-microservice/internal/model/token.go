package model

import "time"

type BlacklistedToken struct {
	Token     string    `bson:"token"`
	ExpiredAt time.Time `bson:"expired_at"`
}
