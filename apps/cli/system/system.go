package system

import (
	"github.com/comunidade-shallom/peristera/pkg/support/errors"
	"github.com/urfave/cli/v2"
)

var NoAdminsDefined = errors.Business("No admins defined", "SY:001")

var SystemCmd = &cli.Command{
	Name:        "system",
	Usage:       "load system info",
	Subcommands: []*cli.Command{InfoCmd},
}
