package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func getServerStats() {
	// CPU
	cpuPercent, _ := cpu.Percent(time.Second, false)
	fmt.Printf("CPU usage: %.2f%%\n", cpuPercent[0])

	// Memory
	memory, _ := mem.VirtualMemory()
	fmt.Printf("Memory usage: %.2f%%\n", memory.UsedPercent)

	// Disk
	diskUsage, _ := disk.Usage("/")
	fmt.Printf("Disk usage: %.2f%%\n", diskUsage.UsedPercent)

	// Network
	network, _ := net.IOCounters(false)
	fmt.Printf("Network sent: %v bytes\n", network[0].BytesSent)
	fmt.Printf("Network received: %v bytes\n", network[0].BytesRecv)
}

func main() {
	for {
		getServerStats()
		time.Sleep(time.Second)
	}
}
