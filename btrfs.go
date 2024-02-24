package main

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

type scrubStatus struct {
	ScrubStarted string `json:"scrubStarted"`
	Status       string `json:"status"`
	Duration     string `json:"duration"`
	ErrorSummary string `json:"errorSummary"`
}

type btrfsData struct {
	ScrubStatus *scrubStatus `json:"scrubStatus"`
}

func GetBtrfsData(mountPoint string) (*btrfsData, error) {
	scrubStatusOutput, err := getScrubStatusOutput(mountPoint)
	if err != nil {
		return nil, err
	}

	scrubStatus, err := parseScrubStatusOutput(scrubStatusOutput)
	if err != nil {
		return nil, err
	}

	btrfsData := btrfsData{
		ScrubStatus: scrubStatus,
	}

	return &btrfsData, nil
}

func getScrubStatusOutput(mountPoint string) ([]byte, error) {
	return exec.Command(
		"btrfs",
		"scrub",
		"status",
		mountPoint,
	).Output()
}

func parseScrubStatusOutput(scrubStatusOutput []byte) (*scrubStatus, error) {
	scrubStatus := scrubStatus{}
	scanner := bufio.NewScanner(bytes.NewReader(scrubStatusOutput))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.SplitN(line, ":", 2)

		switch field := fields[0]; field {
		case "Scrub started":
			scrubStatus.ScrubStarted = strings.TrimSpace(fields[1])
		case "Status":
			scrubStatus.Status = strings.TrimSpace(fields[1])
		case "Duration":
			scrubStatus.Duration = strings.TrimSpace(fields[1])
		case "Error summary":
			scrubStatus.ErrorSummary = strings.TrimSpace(fields[1])
		default:
		}
	}

	return &scrubStatus, nil
}
