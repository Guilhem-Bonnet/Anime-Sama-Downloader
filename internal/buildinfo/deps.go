// Package buildinfo provides dependency checking utilities.
package buildinfo

import (
	"os/exec"
	"strings"
)

// DependencyStatus represents the status of a system dependency.
type DependencyStatus struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
	Version   string `json:"version,omitempty"`
	Required  bool   `json:"required"`
	Message   string `json:"message,omitempty"`
}

// CheckFFmpeg checks if ffmpeg is available in the system PATH.
// Returns version string if found, empty string if not.
func CheckFFmpeg() DependencyStatus {
	status := DependencyStatus{
		Name:     "ffmpeg",
		Required: true,
	}

	cmd := exec.Command("ffmpeg", "-version")
	output, err := cmd.Output()
	if err != nil {
		status.Available = false
		status.Message = "ffmpeg not found - HLS/M3U8 downloads will not work"
		return status
	}

	// Parse version from first line: "ffmpeg version X.X.X ..."
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 3 {
			status.Version = parts[2]
		}
	}

	status.Available = true
	status.Message = "ffmpeg available"
	return status
}

// CheckAllDependencies checks all required system dependencies.
func CheckAllDependencies() []DependencyStatus {
	return []DependencyStatus{
		CheckFFmpeg(),
	}
}

// HasMissingDependencies returns true if any required dependency is missing.
func HasMissingDependencies() bool {
	for _, dep := range CheckAllDependencies() {
		if dep.Required && !dep.Available {
			return true
		}
	}
	return false
}
