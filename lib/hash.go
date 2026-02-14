package lib

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/google/uuid"
)

func Hash(id uuid.UUID) string {
	hash := md5.Sum(id[:])
	return hex.EncodeToString(hash[:4])
}
