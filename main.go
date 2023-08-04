package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type ServerStats struct {
	Time          time.Time
	CPUPercent    float64
	MemoryPercent float64
	DiskPercent   float64
	BytesSent     uint64
	BytesRecv     uint64
}

var serverStatsHistory []ServerStats
const maxHistorySize = 60 // Mantendo os últimos 60 pontos de dados

func getServerStats() ServerStats {
	stats := ServerStats{}
	stats.Time = time.Now()

	cpuPercent, _ := cpu.Percent(time.Millisecond*200, false)
	stats.CPUPercent = cpuPercent[0]

	memory, _ := mem.VirtualMemory()
	stats.MemoryPercent = memory.UsedPercent

	diskUsage, _ := disk.Usage("/")
	stats.DiskPercent = diskUsage.UsedPercent

	network, _ := net.IOCounters(false)
	stats.BytesSent = network[0].BytesSent
	stats.BytesRecv = network[0].BytesRecv

	return stats
}

func storeStats(stats ServerStats) {
	serverStatsHistory = append(serverStatsHistory, stats)
	if len(serverStatsHistory) > maxHistorySize {
		serverStatsHistory = serverStatsHistory[1:]
	}
}

func drawDashboard() {
	if err := termui.Init(); err != nil {
		fmt.Printf("Failed to initialize termui: %v", err)
		return
	}
	defer termui.Close()

	cpuGauge := widgets.NewGauge()
	memoryGauge := widgets.NewGauge()
	diskGauge := widgets.NewGauge()
	cpuGauge.Title = "CPU Usage"
	memoryGauge.Title = "Memory Usage"
	diskGauge.Title = "Disk Usage"
	cpuGauge.SetRect(0, 0, 50, 5)
	memoryGauge.SetRect(0, 5, 50, 10)
	diskGauge.SetRect(0, 10, 50, 15)

	uiEvents := termui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:
			stats := getServerStats()
			storeStats(stats)
			cpuGauge.Percent = int(stats.CPUPercent)
			memoryGauge.Percent = int(stats.MemoryPercent)
			diskGauge.Percent = int(stats.DiskPercent)
			termui.Render(cpuGauge, memoryGauge, diskGauge)
		}
	}
}

func main() {
	drawDashboard()
}