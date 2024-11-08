package flx

import (
	"bufio"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/mojocn/felix/model"
	"golang.org/x/crypto/ssh"
)

type FSInfo struct {
	MountPoint string
	Used       uint64
	Free       uint64
}

type NetIntfInfo struct {
	IPv4 string
	IPv6 string
	Rx   uint64
	Tx   uint64
}

type cpuRaw struct {
	User    uint64 // time spent in user mode
	Nice    uint64 // time spent in user mode with low priority (nice)
	System  uint64 // time spent in system mode
	Idle    uint64 // time spent in the idle task
	Iowait  uint64 // time spent waiting for I/O to complete (since Linux 2.5.41)
	Irq     uint64 // time spent servicing  interrupts  (since  2.6.0-test4)
	SoftIrq uint64 // time spent servicing softirqs (since 2.6.0-test4)
	Steal   uint64 // time spent in other OSes when running in a virtualized environment
	Guest   uint64 // time spent running a virtual CPU for guest operating systems under the control of the Linux kernel.
	Total   uint64 // total of all time fields
}

type CPUInfo struct {
	User    float32
	Nice    float32
	System  float32
	Idle    float32
	Iowait  float32
	Irq     float32
	SoftIrq float32
	Steal   float32
	Guest   float32
}

type Stats struct {
	Uptime       time.Duration
	Hostname     string
	Load1        string
	Load5        string
	Load10       string
	RunningProcs string
	TotalProcs   string
	MemTotal     uint64
	MemFree      uint64
	MemBuffers   uint64
	MemCached    uint64
	SwapTotal    uint64
	SwapFree     uint64
	FSInfos      []FSInfo
	NetIntf      map[string]NetIntfInfo
	CPU          CPUInfo // or []CPUInfo to get all the cpu-core's stats?
	OsInfo       string
}

func (stats *Stats) getAllStats(client *ssh.Client) {
	stats.getUptime(client)
	stats.getHostname(client)
	stats.getLoad(client)
	stats.getMemInfo(client)
	stats.getFSInfo(client)
	stats.getInterfaces(client)
	stats.getInterfaceInfo(client)
	stats.getCPU(client)
	stats.getOsInfo(client)
}

func (stats *Stats) getUptime(client *ssh.Client) (err error) {
	uptime, err := runCommand(client, "/bin/cat /proc/uptime")
	if err != nil {
		return
	}
	parts := strings.Fields(uptime)
	if len(parts) == 2 {
		var upsecs float64
		upsecs, err = strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return
		}
		stats.Uptime = time.Duration(upsecs * 1e9)
	}
	return
}

func (stats *Stats) getHostname(client *ssh.Client) (err error) {
	hostname, err := runCommand(client, "/bin/hostname -f")
	if err != nil {
		return
	}

	stats.Hostname = strings.TrimSpace(hostname)
	return
}

func (stats *Stats) getLoad(client *ssh.Client) (err error) {
	line, err := runCommand(client, "/bin/cat /proc/loadavg")
	if err != nil {
		return
	}

	parts := strings.Fields(line)
	if len(parts) == 5 {
		stats.Load1 = parts[0]
		stats.Load5 = parts[1]
		stats.Load10 = parts[2]
		if i := strings.Index(parts[3], "/"); i != -1 {
			stats.RunningProcs = parts[3][0:i]
			if i+1 < len(parts[3]) {
				stats.TotalProcs = parts[3][i+1:]
			}
		}
	}

	return
}

func (stats *Stats) getMemInfo(client *ssh.Client) (err error) {
	lines, err := runCommand(client, "/bin/cat /proc/meminfo")
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 3 {
			val, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				continue
			}
			val *= 1024
			switch parts[0] {
			case "MemTotal:":
				stats.MemTotal = val
			case "MemFree:":
				stats.MemFree = val
			case "Buffers:":
				stats.MemBuffers = val
			case "Cached:":
				stats.MemCached = val
			case "SwapTotal:":
				stats.SwapTotal = val
			case "SwapFree:":
				stats.SwapFree = val
			}
		}
	}

	return
}

func (stats *Stats) getFSInfo(client *ssh.Client) (err error) {
	lines, err := runCommand(client, "/bin/df -B1")
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	flag := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		n := len(parts)
		dev := n > 0 && strings.Index(parts[0], "/dev/") == 0
		if n == 1 && dev {
			flag = 1
		} else if (n == 5 && flag == 1) || (n == 6 && dev) {
			i := flag
			flag = 0
			used, err := strconv.ParseUint(parts[2-i], 10, 64)
			if err != nil {
				continue
			}
			free, err := strconv.ParseUint(parts[3-i], 10, 64)
			if err != nil {
				continue
			}
			stats.FSInfos = append(stats.FSInfos, FSInfo{
				parts[5-i], used, free,
			})
		}
	}

	return
}

func (stats *Stats) getInterfaces(client *ssh.Client) (err error) {
	var lines string
	lines, err = runCommand(client, "/bin/ip -o addr")
	if err != nil {
		// try /sbin/ip
		lines, err = runCommand(client, "/sbin/ip -o addr")
		if err != nil {
			return
		}
	}

	if stats.NetIntf == nil {
		stats.NetIntf = make(map[string]NetIntfInfo)
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 4 && (parts[2] == "inet" || parts[2] == "inet6") {
			ipv4 := parts[2] == "inet"
			intfname := parts[1]
			if info, ok := stats.NetIntf[intfname]; ok {
				if ipv4 {
					info.IPv4 = parts[3]
				} else {
					info.IPv6 = parts[3]
				}
				stats.NetIntf[intfname] = info
			} else {
				info := NetIntfInfo{}
				if ipv4 {
					info.IPv4 = parts[3]
				} else {
					info.IPv6 = parts[3]
				}
				stats.NetIntf[intfname] = info
			}
		}
	}

	return
}

func (stats *Stats) getInterfaceInfo(client *ssh.Client) (err error) {
	lines, err := runCommand(client, "/bin/cat /proc/net/dev")
	if err != nil {
		return
	}

	if stats.NetIntf == nil {
		return
	} // should have been here already

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 17 {
			intf := strings.TrimSpace(parts[0])
			intf = strings.TrimSuffix(intf, ":")
			if info, ok := stats.NetIntf[intf]; ok {
				rx, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					continue
				}
				tx, err := strconv.ParseUint(parts[9], 10, 64)
				if err != nil {
					continue
				}
				info.Rx = rx
				info.Tx = tx
				stats.NetIntf[intf] = info
			}
		}
	}

	return
}

func (stats *Stats) parseCPUFields(fields []string, stat *cpuRaw) {
	numFields := len(fields)
	for i := 1; i < numFields; i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			continue
		}

		stat.Total += val
		switch i {
		case 1:
			stat.User = val
		case 2:
			stat.Nice = val
		case 3:
			stat.System = val
		case 4:
			stat.Idle = val
		case 5:
			stat.Iowait = val
		case 6:
			stat.Irq = val
		case 7:
			stat.SoftIrq = val
		case 8:
			stat.Steal = val
		case 9:
			stat.Guest = val
		}
	}
}

func (stats *Stats) getOsInfo(client *ssh.Client) {
	lines, err := runCommand(client, "uname -a")
	if err != nil {
		return
	}
	stats.OsInfo = strings.TrimSpace(lines)
}

// the CPU stats that were fetched last time round
var preCPU cpuRaw

func (stats *Stats) getCPU(client *ssh.Client) (err error) {
	lines, err := runCommand(client, "/bin/cat /proc/stat")
	if err != nil {
		return
	}

	var (
		nowCPU cpuRaw
		total  float32
	)

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == "cpu" { // changing here if want to get every cpu-core's stats
			stats.parseCPUFields(fields, &nowCPU)
			break
		}
	}
	if preCPU.Total == 0 { // having no pre raw cpu data
		goto END
	}

	total = float32(nowCPU.Total - preCPU.Total)
	stats.CPU.User = float32(nowCPU.User-preCPU.User) / total * 100
	stats.CPU.Nice = float32(nowCPU.Nice-preCPU.Nice) / total * 100
	stats.CPU.System = float32(nowCPU.System-preCPU.System) / total * 100
	stats.CPU.Idle = float32(nowCPU.Idle-preCPU.Idle) / total * 100
	stats.CPU.Iowait = float32(nowCPU.Iowait-preCPU.Iowait) / total * 100
	stats.CPU.Irq = float32(nowCPU.Irq-preCPU.Irq) / total * 100
	stats.CPU.SoftIrq = float32(nowCPU.SoftIrq-preCPU.SoftIrq) / total * 100
	stats.CPU.Guest = float32(nowCPU.Guest-preCPU.Guest) / total * 100
END:
	preCPU = nowCPU
	return
}

func fmtUptime(stats *Stats) string {
	dur := stats.Uptime
	dur = dur - (dur % time.Second)
	var days int
	for dur.Hours() > 24.0 {
		days++
		dur -= 24 * time.Hour
	}
	s1 := dur.String()
	s2 := ""
	if days > 0 {
		s2 = fmt.Sprintf("%dd ", days)
	}
	for _, ch := range s1 {
		s2 += string(ch)
		if ch == 'h' || ch == 'm' {
			s2 += " "
		}
	}
	return s2
}

func fmtBytes(val uint64) string {
	if val < 1024 {
		return fmt.Sprintf("%d bytes", val)
	} else if val < 1024*1024 {
		return fmt.Sprintf("%6.2f KiB", float64(val)/1024.0)
	} else if val < 1024*1024*1024 {
		return fmt.Sprintf("%6.2f MiB", float64(val)/1024.0/1024.0)
	} else {
		return fmt.Sprintf("%6.2f GiB", float64(val)/1024.0/1024.0/1024.0)
	}
}

func ShowHardwareInfo(mac *model.Machine) error {
	client, err := NewSshClient(mac)
	if err != nil {
		return err
	}
	showStats(client)
	return nil
}
func FetchHardwareInfo(mc *model.Machine) (*Stats, error) {
	client, err := NewSshClient(mc)
	if err != nil {
		return nil, err
	}
	stats := Stats{}
	stats.getAllStats(client)
	return &stats, nil
}
func showStats(client *ssh.Client) {
	stats := Stats{}
	stats.getAllStats(client)
	used := stats.MemTotal - stats.MemFree - stats.MemBuffers - stats.MemCached
	color.HiCyan("Machine Hardware Info:")
	color.Cyan("HostName:\t%s\tUpTime:\t%s", stats.Hostname, fmtUptime(&stats))
	color.Yellow("%s", stats.OsInfo)
	color.HiCyan("Load:")
	color.Yellow("%s\t%s\t%s", stats.Load1, stats.Load5, stats.Load10)
	color.HiCyan("CPU:")
	color.HiYellow("User:%.2f\t SYS: %.2f\t Nice: %.2f\t Idle: %.2f\t Iowait:%.2f\t Irq:%.2f\t SoftIrq:%.2f\t Guest:%.2f", stats.CPU.User, stats.CPU.System, stats.CPU.Nice, stats.CPU.Idle, stats.CPU.Iowait, stats.CPU.Irq, stats.CPU.SoftIrq, stats.CPU.Guest)
	color.HiCyan("Processes:")
	color.Cyan("Running:\t%s\tTotal:\t%s", stats.RunningProcs, stats.TotalProcs)
	color.HiCyan("Memory:")
	color.Cyan("Total: %s\t Free: %s\t Used: %s\t Buffers: %s\t Cached: %s\t Swap: %d of %d", fmtBytes(stats.MemTotal), fmtBytes(stats.MemFree), fmtBytes(used), fmtBytes(stats.MemBuffers), fmtBytes(stats.MemCached), stats.SwapFree, stats.SwapTotal)
	if len(stats.FSInfos) > 0 {
		color.HiCyan("Filesystems:")
		for _, fs := range stats.FSInfos {
			color.Cyan("%s:\t%s free of %s", fs.MountPoint, fmtBytes(fs.Free), fmtBytes(fs.Used+fs.Free))
		}
	}
	if len(stats.NetIntf) > 0 {
		color.HiCyan("Network Interfaces:")
		keys := make([]string, 0, len(stats.NetIntf))
		for intf := range stats.NetIntf {
			keys = append(keys, intf)
		}
		sort.Strings(keys)
		for _, intf := range keys {
			info := stats.NetIntf[intf]
			color.Cyan("%s:\tIPv4: %s\tIPv6: %s", intf, info.IPv4, info.IPv6)
			color.Yellow("rx: %s\ttx: %s", fmtBytes(info.Rx), fmtBytes(info.Tx))
		}
	}
}
