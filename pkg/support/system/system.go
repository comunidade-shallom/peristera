package system

import (
	"net"
	"strconv"
	"strings"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/support"
	"github.com/matishsiao/goInfo"
	"github.com/pbnjay/memory"
	"github.com/pterm/pterm"
)

type Memory struct {
	Free  uint64
	Total uint64
}

type Data struct {
	IPs    []net.Addr
	Memory Memory
	Info   goInfo.GoInfoObject
}

func New() (Data, error) {
	ips, err := net.InterfaceAddrs()
	if err != nil {
		return Data{}, err
	}

	info, err := goInfo.GetInfo()
	if err != nil {
		return Data{}, err
	}

	free := memory.FreeMemory()
	total := memory.TotalMemory()

	return Data{
		IPs:  ips,
		Info: info,
		Memory: Memory{
			Free:  free,
			Total: total,
		},
	}, nil
}

func (d Data) MarkdownV2(head string) string {
	var builder strings.Builder

	if head != "" {
		builder.WriteString(support.AddSlashes(head))
		builder.WriteString("\n")
	}

	info := d.Info

	builder.WriteString("\n*System Details*\n")
	builder.WriteString("\n*Hostname:* " + support.AddSlashes(info.Hostname))
	builder.WriteString("\n*Platform:* " + support.AddSlashes(info.Platform))
	builder.WriteString("\n*CPUs:* " + strconv.Itoa(info.CPUs))
	builder.WriteString("\n*GoOS:* " + support.AddSlashes(info.GoOS))
	builder.WriteString("\n*Core:* " + support.AddSlashes(info.Core))
	builder.WriteString("\n*OS:* " + support.AddSlashes(info.OS))
	builder.WriteString("\n*Kernel:* " + support.AddSlashes(info.Kernel))
	builder.WriteString("\n*Memory:* " + support.AddSlashes(bytefmt.ByteSize(d.Memory.Total)))
	builder.WriteString("\n*Memory Free:* " + support.AddSlashes(bytefmt.ByteSize(d.Memory.Free)))

	builder.WriteString("\n\n*System IPs*\n")

	for _, ip := range d.IPs {
		builder.WriteString("\n`" + support.AddSlashes(ip.String()) + "`")
	}

	builder.WriteString("\n\n*System Time:* \n" + support.AddSlashes(time.Now().Format(time.RFC3339)) + "\n")
	builder.WriteString("\n*Version:*\n" + support.AddSlashes(config.Version()))

	return builder.String()
}

func (d Data) Println() error {
	info := d.Info

	infos, err := pterm.DefaultTable.
		WithBoxed().
		WithData(pterm.TableData{
			{"Hostname", info.Hostname},
			{"Platform", info.Platform},
			{"CPUs", strconv.Itoa(info.CPUs)},
			{"GoOS", info.GoOS},
			{"Core", info.Core},
			{"OS", info.OS},
			{"Kernel", info.Kernel},
			{"Memory", bytefmt.ByteSize(d.Memory.Total)},
			{"Memory Free", bytefmt.ByteSize(d.Memory.Free)},
		}).Srender()
	if err != nil {
		return err
	}

	lines := pterm.TableData{
		{"network", "IP"},
	}

	for _, v := range d.IPs {
		lines = append(lines, []string{v.Network(), v.String()})
	}

	ips, err := pterm.DefaultTable.
		WithBoxed().
		WithHasHeader().
		WithData(lines).Srender()
	if err != nil {
		return err
	}

	panels, err := pterm.DefaultPanel.WithPanels(pterm.Panels{
		{{Data: infos}},
		{{Data: ips}},
	}).Srender()
	if err != nil {
		return err
	}

	pterm.DefaultBox.
		WithTitle("System Info").
		WithTitleBottomRight().
		WithRightPadding(0).
		WithBottomPadding(0).
		Println(panels)

	return nil
}
