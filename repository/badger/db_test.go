package badger

import (
	"fmt"
	"os"
	"path"
)

var testBadgerBaseDir = os.Getenv("BADGER_DIR")

func newTestDB(testName string) (*DB, func(), error) {
	badgeDir := path.Join(testBadgerBaseDir, "test_"+testName)

	db, closeFn, err := newDB(badgeDir, testName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create %s: %v", badgeDir, err)
	}

	return db, func() {
		closeFn()
		_ = os.RemoveAll(badgeDir)
	}, err
}
