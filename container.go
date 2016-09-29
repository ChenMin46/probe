package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Container struct {
	*State  `json:"State"` // Needed for remote api version <= 1.11
	ID      string
	Created time.Time
	Path    string
	Args    []string
	ImageID string `json:"Image"`
	Name    string
}

type State struct {
	Running           bool
	Paused            bool
	Restarting        bool
	OOMKilled         bool
	Dead              bool
	removalInProgress bool
	Pid               int
	ExitCode          int
	Error             string // contains last known error when starting the container
	StartedAt         time.Time
	FinishedAt        time.Time
}

func ContainerFromDisk(id string, root string) (*Container, error) {
	c := Container{
		ID:    id,
		State: &State{},
	}
	jsonfile := "config.json"

	jsonSource, err := os.Open(filepath.Join(root, id, jsonfile))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if os.IsNotExist(err) {
		jsonfile = "config.v2.json"
		jsonSource, err = os.Open(filepath.Join(root, id, jsonfile))
		if err != nil {
			return nil, err
		}
	}
	defer jsonSource.Close()

	dec := json.NewDecoder(jsonSource)
	if err := dec.Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

// String returns a human-readable description of the state
func (s *State) String() string {
	if s.Running {
		if s.Paused {
			return fmt.Sprintf("Up %s (Paused)", HumanDuration(time.Now().UTC().Sub(s.StartedAt)))
		}
		if s.Restarting {
			return fmt.Sprintf("Restarting (%d) %s ago", s.ExitCode, HumanDuration(time.Now().UTC().Sub(s.FinishedAt)))
		}

		return fmt.Sprintf("Up %s", HumanDuration(time.Now().UTC().Sub(s.StartedAt)))
	}

	if s.removalInProgress {
		return "Removal In Progress"
	}

	if s.Dead {
		return "Dead"
	}

	if s.StartedAt.IsZero() {
		return "Created"
	}

	if s.FinishedAt.IsZero() {
		return ""
	}

	return fmt.Sprintf("Exited (%d) %s ago", s.ExitCode, HumanDuration(time.Now().UTC().Sub(s.FinishedAt)))
}

// HumanDuration returns a human-readable approximation of a duration
// (eg. "About a minute", "4 hours ago", etc.).
func HumanDuration(d time.Duration) string {
	if seconds := int(d.Seconds()); seconds < 1 {
		return "Less than a second"
	} else if seconds < 60 {
		return fmt.Sprintf("%d seconds", seconds)
	} else if minutes := int(d.Minutes()); minutes == 1 {
		return "About a minute"
	} else if minutes < 60 {
		return fmt.Sprintf("%d minutes", minutes)
	} else if hours := int(d.Hours()); hours == 1 {
		return "About an hour"
	} else if hours < 48 {
		return fmt.Sprintf("%d hours", hours)
	} else if hours < 24*7*2 {
		return fmt.Sprintf("%d days", hours/24)
	} else if hours < 24*30*3 {
		return fmt.Sprintf("%d weeks", hours/24/7)
	} else if hours < 24*365*2 {
		return fmt.Sprintf("%d months", hours/24/30)
	}
	return fmt.Sprintf("%d years", int(d.Hours())/24/365)
}

// IsRunning returns whether the running flag is set. Used by Container to check whether a container is running.
func (s *State) IsRunning() bool {
	res := s.Running
	return res
}
