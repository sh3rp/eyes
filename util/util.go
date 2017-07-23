package util

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

func GenID() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}
