package main

import (
	"os"
	"strconv"
	"strings"
)

type memoryStats struct {
	Used        uint64 `json:"Used"`
	UsedPercent uint   `json:"UsedPercent"`
	Available   uint64 `json:"Available"`
	Total       uint64 `json:"Total"`
}

func getMemoryMap() (map[string]*memoryStats, error) {
	_memoryStats, err := getMemoryStats()
	if err != nil {
		return nil, err
	}

	memoryMap := make(map[string]*memoryStats)

	memoryUsed := _memoryStats["MemTotal"] - (_memoryStats["MemFree"] + _memoryStats["Buffers"] + _memoryStats["Cached"])
	memoryMap["Physical"] = &memoryStats{
		Used:        kibiToMega(memoryUsed),
		UsedPercent: uint(memoryUsed * 100 / _memoryStats["MemTotal"]),
		Available:   kibiToMega(_memoryStats["MemAvailable"]),
		Total:       kibiToMega(_memoryStats["MemTotal"]),
	}

	swapUsed := _memoryStats["SwapTotal"] - (_memoryStats["SwapFree"] + _memoryStats["SwapCached"])
	memoryMap["Swap"] = &memoryStats{
		Used:        kibiToMega(swapUsed),
		UsedPercent: uint(swapUsed * 100 / _memoryStats["SwapTotal"]),
		Available:   kibiToMega(_memoryStats["SwapFree"] + _memoryStats["SwapCached"]),
		Total:       kibiToMega(_memoryStats["SwapTotal"]),
	}

	return memoryMap, nil
}

func getMemoryStats() (map[string]uint64, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	statMap := make(map[string]uint64)

	for _, line := range lines {
		field, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		valFields := strings.Fields(value)
		val, _ := strconv.ParseUint(valFields[0], 10, 64)
		statMap[field] = val
	}

	return statMap, nil
}

func kibiToMega(kibiVal uint64) uint64 {
	return uint64(float64(kibiVal) / 1024)
}
