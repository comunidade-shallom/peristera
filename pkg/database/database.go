package database

import (
	"context"
	"errors"
	"io"
	"time"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/rs/zerolog"
)

type Database struct {
	db *badger.DB
}

func Open(path string) (Database, error) {
	if path == "" {
		path = "/tmp/peristera"
	}

	db, err := badger.Open(
		badger.DefaultOptions(path).
			WithIndexCacheSize(5 << 20). //nolint:gomnd // 5mb
			WithValueLogMaxEntries(50),  //nolint:gomnd
	)
	if err != nil {
		return Database{}, err
	}

	return Database{
		db: db,
	}, nil
}

func (d Database) DB() *badger.DB {
	return d.db
}

func (d Database) Close() error {
	return d.DB().Close()
}

func (d Database) Backup(w io.Writer) error {
	_, err := d.DB().Backup(w, 0)

	return err
}

// MissingKeys return missing keys in database.
func (d Database) MissingKeys(ctx context.Context, keys [][]byte) ([][]byte, error) {
	out := make([][]byte, 0)

	err := d.db.View(func(txn *badger.Txn) error {
		for _, key := range keys {
			_, err := txn.Get(key)

			if err == nil {
				// key found, ignore
				continue
			}

			if errors.Is(err, badger.ErrKeyNotFound) {
				out = append(out, key)

				continue
			}

			return err
		}

		return nil
	})

	return out, err
}

func (d Database) Worker(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute) //nolint:gomnd
	logger := zerolog.Ctx(ctx).
		With().
		Str("context", "database:worker").
		Logger()

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("Stopping worker")

			_ = d.DB().RunValueLogGC(1)

			return
		case <-ticker.C:
			logger.Info().Msg("Running CG")

			err := d.DB().RunValueLogGC(0.7) //nolint:gomnd
			if err != nil {
				logger.Error().Err(err).Msg("Fail to run GC")
			}
		}
	}
}
