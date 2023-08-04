package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
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

func main() {
	for {
		stats := getServerStats()
		storeStats(stats)
		fmt.Printf("Server Stats: %+v\n", stats)
		time.Sleep(time.Second)
	}
}