package pstore

import (
	"fmt"
)

const error_prefix = "pstore: "

func IsPStoreError(err error) bool {
	return err != nil && err.Error()[:len(error_prefix)] == error_prefix
}

func (ps *PersistentStorage) errorf(format string, args ...interface{}) error {
	prefix := fmt.Sprintf("%s%s: ", error_prefix, ps.name)
	return fmt.Errorf(prefix+format, args...)
}
