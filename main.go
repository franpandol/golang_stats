package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"net/http"
)

func main() {
	// Initialize the Gin router
	router := gin.Default()
	// handle cors

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return

		}
		c.Next()
	})
	
	// Define a route to retrieve system statistics
	router.GET("/stats", func(c *gin.Context) {
		// Get CPU information
		cpuInfo, err := cpu.Info()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve CPU information"})
			return
		}
		
		// Get memory information
		memInfo, err := mem.VirtualMemory()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve memory information"})
			return
		}

		// Get disk usage information
		partitions, err := disk.PartitionsWithContext(context.TODO(), false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve disk partition information"})
			return
		}

		// Prepare the response data
		data := gin.H{
			"cpu": gin.H{
				"model": cpuInfo[0].ModelName,
				"cores": cpuInfo[0].Cores,
			},
			"memory": gin.H{
				"total":     memInfo.Total,
				"available": memInfo.Available,
				"used":      memInfo.Used,
			},
			"disk": gin.H{},
		}

		// Populate disk usage information
		for _, partition := range partitions {
			usage, err := disk.UsageWithContext(context.TODO(), partition.Mountpoint)
			if err != nil {
				fmt.Println("Failed to retrieve disk usage for", partition.Mountpoint, ":", err)
				continue
			}
			data["disk"].(gin.H)[partition.Mountpoint] = gin.H{
				"total": usage.Total,
				"free":  usage.Free,
				"used":  usage.Used,
			}
		}

		// Return the system statistics in the response
		c.JSON(http.StatusOK, data)
	})

	// Run the server
	router.Run(":8087")
}
