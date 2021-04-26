package badger

import (
	"os"
	"path"
)

var testBadgerBaseDir = os.Getenv("BADGER_DIR")

func newTestDB(testName string) (*DB, func(), error) {
	badgeDir := path.Join(testBadgerBaseDir, testName)
	db, close, err := newDB(badgeDir, "webshell")
	if err != nil {
		return nil, nil, err
	}

	return db, func() {
		close()
		_ = os.RemoveAll(badgeDir)
	}, err
}
