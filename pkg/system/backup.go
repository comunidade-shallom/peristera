package system

import (
	"bufio"
	"context"
	"io/ioutil"
	"os"

	"github.com/comunidade-shallom/peristera/pkg/database"
	"github.com/rs/zerolog"
)

type DestroyBackup func()

func noop() {}

func Backup(ctx context.Context, db database.Database) (*os.File, DestroyBackup, error) {
	logger := zerolog.Ctx(ctx)

	file, err := ioutil.TempFile(os.TempDir(), "peristera.*.bak")
	if err != nil {
		return file, noop, err
	}

	destroy := func() {
		name := file.Name()
		_ = file.Close()
		_ = os.Remove(name)
	}

	logger.Info().Msgf("Temp backup file created: %s", file.Name())

	bw := bufio.NewWriterSize(file, 64<<20) //nolint:gomnd

	// run backup
	if err = db.Backup(bw); err != nil {
		return file, destroy, err
	}

	logger.Debug().Msg("Backup generated")

	if err = bw.Flush(); err != nil {
		return file, destroy, err
	}

	if err = file.Sync(); err != nil {
		return file, destroy, err
	}

	return file, destroy, nil
}
