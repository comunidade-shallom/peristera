package system

import (
	"github.com/comunidade-shallom/peristera/pkg/support/errors"
	"github.com/urfave/cli/v2"
)

var (
	NoAdminsDefined = errors.Business("No admins defined", "SY:001")
	NoRootsDefined  = errors.Business("No roots defined", "SY:002")
	BackupIsEmpty   = errors.Business("🪣 Backup is empty", "SY:003")
)

var ErrOnlyNotifyTrue = errors.Business("only notify true is supported", "SY:003")

var SystemCmd = &cli.Command{
	Name:        "system",
	Usage:       "load system info",
	Subcommands: []*cli.Command{InfoCmd, BackupCmd},
}
