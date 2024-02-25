package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type scrubStatus struct {
	LastScrubFinished  *time.Time `json:"lastScrubFinished,omitempty"`
	Duration           string     `json:"duration,omitempty"`
	Status             string     `json:"status,omitempty"`
	ErrorSummary       string     `json:"errorSummary"`
	DaysSinceLastScrub int        `json:"daysSinceLastScrub,omitempty"`
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

func daysSince(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

func parseDuration(durations string) (time.Duration, error) {
	fields := strings.Split(durations, ":")
	return time.ParseDuration(fmt.Sprintf("%sh%sm%ss", fields[0], fields[1], fields[2]))
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
	var scrubStarted time.Time
	var duration time.Duration

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.SplitN(line, ":", 2)

		switch field := fields[0]; field {

		case "Scrub started":
			var err error
			scrubStarted, err = time.Parse(time.ANSIC, strings.TrimSpace(fields[1]))
			if err != nil {
				return nil, err
			}

		case "Status":
			scrubStatus.Status = strings.TrimSpace(fields[1])

		case "Duration":
			var err error
			scrubStatus.Duration = strings.TrimSpace(fields[1])
			duration, err = parseDuration(scrubStatus.Duration)
			if err != nil {
				return nil, err
			}

		case "Error summary":
			scrubStatus.ErrorSummary = strings.TrimSpace(fields[1])
		}
	}

	if (scrubStarted != time.Time{} && duration != 0) {
		lastScrubFinished := scrubStarted.Add(duration)
		scrubStatus.LastScrubFinished = &lastScrubFinished
		scrubStatus.DaysSinceLastScrub = daysSince(lastScrubFinished)
	}

	return &scrubStatus, nil
}
