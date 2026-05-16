// ClusterGuard Agent — install on each monitored server.
// Exposes system metrics over HTTP for the desktop app.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"ClusterGuard/pkg/agentapi"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	addr := flag.String("addr", ":9100", "listen address")
	token := flag.String("token", "", "optional bearer token (env CLUSTERGUARD_TOKEN)")
	flag.Parse()

	if *token == "" {
		*token = os.Getenv("CLUSTERGUARD_TOKEN")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		if *token != "" && r.Header.Get("X-ClusterGuard-Token") != *token {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		m, err := collectMetrics()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(m)
	})

	log.Printf("ClusterGuard agent listening on %s", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal(err)
	}
}

func collectMetrics() (*agentapi.Metrics, error) {
	cpuPercents, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("cpu: %w", err)
	}
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("memory: %w", err)
	}
	du, err := disk.Usage("/")
	if err != nil {
		return nil, fmt.Errorf("disk: %w", err)
	}

	cpuPct := 0.0
	if len(cpuPercents) > 0 {
		cpuPct = cpuPercents[0]
	}

	return &agentapi.Metrics{
		Timestamp:            time.Now().Unix(),
		CPUPercent:           cpuPct,
		MemoryUsedPercent:    vm.UsedPercent,
		MemoryAvailableBytes: vm.Available,
		MemoryTotalBytes:     vm.Total,
		DiskUsedPercent:      du.UsedPercent,
		DiskFreeBytes:        du.Free,
		DiskTotalBytes:       du.Total,
	}, nil
}
