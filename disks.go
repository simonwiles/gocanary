package main

import (
	"bufio"
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

type fileSystem struct {
	Name        string
	FsType      string
	Size        uint64
	Used        uint64
	UsedPercent uint
	Available   uint64
	MountPoint  string
}

type disk struct {
	Name        string     `json:"name"`
	FsType      string     `json:"fsType"`
	Size        uint64     `json:"size"`
	Used        uint64     `json:"used"`
	UsedPercent uint       `json:"usedPercent"`
	Available   uint64     `json:"available"`
	MountPoints []string   `json:"mountPoints"`
	Btrfs       *btrfsData `json:"btrfs,omitempty"`
}

func getDisksMap() (map[string]*disk, error) {
	dfOutput, err := getDfOutput()
	if err != nil {
		return nil, err
	}

	fileSystems, err := parseDfOutputToFileSystem(dfOutput)
	if err != nil {
		return nil, err
	}

	disksMap := make(map[string]*disk)

	for _, fs := range fileSystems {
		if disksMap[fs.Name] == nil {
			disksMap[fs.Name] = &disk{
				Name:        fs.Name,
				FsType:      fs.FsType,
				Size:        fs.Size,
				Used:        fs.Used,
				UsedPercent: fs.UsedPercent,
				Available:   fs.Available,
				MountPoints: []string{fs.MountPoint},
			}

			if fs.FsType == "btrfs" {
				disksMap[fs.Name].Btrfs, err = getBtrfsData(fs.MountPoint)
				if err != nil {
					return nil, err
				}
			}
		} else {
			disksMap[fs.Name].MountPoints = append(disksMap[fs.Name].MountPoints, fs.MountPoint)
		}
	}

	return disksMap, nil
}

func getDfOutput() ([]byte, error) {
	return exec.Command(
		"df",
		"--exclude-type=tmpfs",
		"--exclude-type=devtmpfs",
		"--exclude-type=overlay",
		"--exclude-type=efivarfs",
		"--print-type",
		"--block-size=MiB").Output()
}

func parseDfOutputToFileSystem(dfOutput []byte) ([]*fileSystem, error) {
	fileSystems := []*fileSystem{}
	scanner := bufio.NewScanner(bytes.NewReader(dfOutput))
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		if lineCount == 1 {
			continue
		}

		line := scanner.Text()
		fileSystem, err := parseDfLineToFileSystem(line)
		if err != nil {
			return nil, err
		}

		fileSystems = append(fileSystems, fileSystem)
	}

	return fileSystems, nil
}

func parseMiB(value string) (uint64, error) {
	return strconv.ParseUint(strings.TrimRight(value, "MiB"), 10, 64)
}

func parseDfLineToFileSystem(line string) (*fileSystem, error) {
	fields := strings.Fields(line)

	size, err := parseMiB(fields[2])
	if err != nil {
		return nil, err
	}

	used, err := parseMiB(fields[3])
	if err != nil {
		return nil, err
	}

	available, err := parseMiB(fields[4])
	if err != nil {
		return nil, err
	}

	usedPercent, err := strconv.ParseUint(strings.TrimRight(fields[5], "%"), 10, 32)
	if err != nil {
		return nil, err
	}

	fileSystem := fileSystem{
		Name:        fields[0],
		FsType:      fields[1],
		Size:        size,
		Used:        used,
		Available:   available,
		UsedPercent: uint(usedPercent),
		MountPoint:  fields[6],
	}

	return &fileSystem, nil
}
