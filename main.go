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
	SwapPercent   float64
}

var serverStatsHistory []ServerStats
const maxHistorySize = 60 // Mantendo os últimos 60 pontos de dados

func getServerStats() *ServerStats {
	stats := &ServerStats{}
	stats.Time = time.Now()

	cpuPercent, err := cpu.Percent(time.Millisecond*200, false)
	if err == nil {
		stats.CPUPercent = cpuPercent[0]
	}

	memory, err := mem.VirtualMemory()
	if err == nil {
		stats.MemoryPercent = memory.UsedPercent
	}

	diskUsage, err := disk.Usage("/")
	if err == nil {
		stats.DiskPercent = diskUsage.UsedPercent
	}

	network, err := net.IOCounters(false)
	if err == nil {
		stats.BytesSent = network[0].BytesSent
		stats.BytesRecv = network[0].BytesRecv
	}

	swap, err := mem.SwapMemory()
    if err == nil {
        stats.SwapPercent = swap.UsedPercent
    }

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
	swapGauge := widgets.NewGauge()
	networkBar := widgets.NewBarChart()
	
	cpuGauge.Title = "CPU Usage"
	memoryGauge.Title = "Memory Usage"
	diskGauge.Title = "Disk Usage"
    swapGauge.Title = "Swap Usage"
	
	cpuGauge.SetRect(0, 0, 50, 5)
	memoryGauge.SetRect(0, 5, 50, 10)
	diskGauge.SetRect(0, 10, 50, 15)
    swapGauge.SetRect(0, 15, 50, 20)
	networkBar.SetRect(0, 15, 50, 20)
	
	networkBar.Title = "Network Traffic"
	networkDataSent := make([]float64, maxHistorySize)
	networkDataRecv := make([]float64, maxHistorySize)

	

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
			stats := getServerStats() // Coleta as estatísticas
			storeStats(*stats)        // Armazena as estatísticas
			cpuGauge.Percent = int(stats.CPUPercent)
			memoryGauge.Percent = int(stats.MemoryPercent)
			diskGauge.Percent = int(stats.DiskPercent)
			swapGauge.Percent = int(stats.SwapPercent)
			networkDataSent = append(networkDataSent[1:], float64(stats.BytesSent))
			networkDataRecv = append(networkDataRecv[1:], float64(stats.BytesRecv))
			networkBar.Data = append(networkDataSent, networkDataRecv...)
			networkBar.Labels = []string{"Sent", "Received"}
			termui.Render(cpuGauge, memoryGauge, diskGauge, swapGauge, networkBar)

		}
	}
}

func main() {
	drawDashboard()
}