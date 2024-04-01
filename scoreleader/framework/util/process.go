package util

import (
	"encoding/json"
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	process2 "github.com/shirou/gopsutil/process"
)

type processInfo struct {
	Pid           int32                    `json:"pid"`
	Name          string                   `json:"name"`
	Cwd           string                   `json:"cwd"`
	CPUPercent    float64                  `json:"cpu_percent"`
	Cmdline       string                   `json:"cmd_line"`
	CreateTime    string                   `json:"create_time"`
	MemoryPercent float32                  `json:"memory_percent"`
	MemoryInfo    *process2.MemoryInfoStat `json:"memory_info"`
}

func (m processInfo) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

// ResourceInfo 시스템 정보 구조체
type ResourceInfo struct {
	ProcessInfo     *processInfo           `json:"process_info"`
	HostInfo        *host.InfoStat         `json:"host_info"`
	CPUCount        int                    `json:"cpu_count"`
	CPUCountLogical int                    `json:"cpu_count_logical"`
	CPUPercent      float64                `json:"cpu_percent"`
	PerCPUPercent   []float64              `json:"per_cpu_percent"`
	LoadAVG         *load.AvgStat          `json:"load_avg"`
	VirtualMemory   *mem.VirtualMemoryStat `json:"virtual_memory"`
	MemoryPercent   float64                `json:"memory_percent"`
	DiskPercent     float64                `json:"disk_percent"`
}

func (m ResourceInfo) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

// GetProcessInfo 시스템 정보 반환
func GetProcessInfo() *ResourceInfo {
	ret := new(ResourceInfo)

	// for memory
	virtualMemory, _ := mem.VirtualMemory()
	ret.VirtualMemory = virtualMemory
	ret.MemoryPercent = virtualMemory.UsedPercent

	// for disk
	diskStat, err := disk.Usage("/")
	if err != nil {
		diskStat, _ = disk.Usage("\\")
	}
	if diskStat != nil {
		ret.DiskPercent = diskStat.UsedPercent
	}

	// for cpu
	cpuCount, _ := cpu.Counts(false)
	ret.CPUCount = cpuCount

	cpuCountLogical, _ := cpu.Counts(true)
	ret.CPUCountLogical = cpuCountLogical

	cpuPercent, _ := cpu.Percent(0, false)
	ret.CPUPercent = cpuPercent[0]

	perCPUPercent, _ := cpu.Percent(0, true)
	ret.PerCPUPercent = perCPUPercent

	hostInfo, _ := host.Info()
	ret.HostInfo = hostInfo

	loadAvg, _ := load.Avg()
	ret.LoadAVG = loadAvg

	prs, _ := process2.Processes()
	ret.ProcessInfo = new(processInfo)
	for _, v := range prs {
		if v.Pid != int32(os.Getpid()) {
			continue
		}
		cpuPercent, _ := v.CPUPercent()
		cmdline, _ := v.Cmdline()
		msec, _ := v.CreateTime()
		memoryInfoStat, _ := v.MemoryInfo()
		memoryPercent, _ := v.MemoryPercent()
		name, _ := v.Name()
		cwd, _ := v.Cwd()

		ret.ProcessInfo.Pid = v.Pid
		ret.ProcessInfo.CPUPercent = cpuPercent
		ret.ProcessInfo.Cmdline = cmdline
		ret.ProcessInfo.CreateTime = time.Unix(0, msec*int64(time.Millisecond)).Format(time.RFC3339)
		ret.ProcessInfo.MemoryInfo = memoryInfoStat
		ret.ProcessInfo.MemoryPercent = memoryPercent
		ret.ProcessInfo.Name = name
		ret.ProcessInfo.Cwd = cwd
	}
	return ret
}
