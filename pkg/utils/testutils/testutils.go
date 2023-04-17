package testutils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func RandomTestDatabaseName() string {
	raw := make([]byte, 6)

	_, err := rand.Read(raw)
	if err != nil {
		err = fmt.Errorf("read from rand failed: %w", err)
		panic(err)
	}

	enc := hex.EncodeToString(raw)

	return fmt.Sprintf("test_%s", enc)
}
