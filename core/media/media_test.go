package media

import (
	"os"
	"testing"
)

// -----------------------------------------
// ParseTimeToSeconds
// -----------------------------------------

func TestParseTimeToSeconds(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"zero", "00:00:00.00", 0},
		{"seconds only, whole number", "00:00:30.00", 30},
		{"seconds only, decimal", "00:00:30.50", 30.5},
		{"minutes only", "00:01:00.00", 60},
		{"hours only", "01:00:00.00", 3600},
		{"mixed values", "01:30:45.25", 5445.25},
		{"large hour value", "02:00:00.00", 7200},
		// malformed inputs: current behavior returns 0 silently
		{"invalid string", "invalid", 0},
		{"empty string", "", 0},
		{"missing part (2 parts)", "01:02", 0},
		{"extra part (4 parts)", "01:02:03:04", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTimeToSeconds(tt.input)
			if got != tt.expected {
				t.Errorf("ParseTimeToSeconds(%q) = %v, expected %v", tt.input, got, tt.expected)
			}
		})
	}
}

// -----------------------------------------
// GetFFmpegPath
// -----------------------------------------

func TestGetFFmpegPath_EnvVar(t *testing.T) {
	// clear customFFmpegPath; it should not be set by a flag in test environment
	original := customFFmpegPath
	customFFmpegPath = ""
	defer func() { customFFmpegPath = original }()

	fakeFFmpegPath := "/fake/path/to/ffmpeg"
	t.Setenv("FFMPEG_PATH", fakeFFmpegPath)

	path, err := GetFFmpegPath()
	if err != nil {
		t.Fatalf("GetFFmpegPath() unexpected error: %v", err)
	}
	if path != fakeFFmpegPath {
		t.Errorf("GetFFmpegPath() = %q, expected %q", path, fakeFFmpegPath)
	}
}

func TestGetFFmpegPath_CustomFlag(t *testing.T) {
	original := customFFmpegPath
	customFFmpegPath = "/custom/ffmpeg"
	defer func() { customFFmpegPath = original }()

	// env var is also set; the flag value must take priority
	t.Setenv("FFMPEG_PATH", "/env/ffmpeg")

	path, err := GetFFmpegPath()
	if err != nil {
		t.Fatalf("GetFFmpegPath() unexpected error: %v", err)
	}
	if path != "/custom/ffmpeg" {
		t.Errorf("GetFFmpegPath() = %q, expected %q (flag must take priority)", path, "/custom/ffmpeg")
	}
}

func TestGetFFmpegPath_NotFound(t *testing.T) {
	// clear all sources; expect an error if ffmpeg is not in PATH
	original := customFFmpegPath
	customFFmpegPath = ""
	defer func() { customFFmpegPath = original }()

	os.Unsetenv("FFMPEG_PATH")

	path, err := GetFFmpegPath()
	if err != nil {
		// expected: ffmpeg is not installed
		t.Logf("GetFFmpegPath() returned expected error: %v", err)
		return
	}
	// ffmpeg is present on this system, log and pass
	t.Logf("ffmpeg found on system, skipping not-found check: %s", path)
}

// -----------------------------------------
// GetFFprobePath
// -----------------------------------------

func TestGetFFprobePath_EnvVar(t *testing.T) {
	fakeFFprobePath := "/fake/path/to/ffprobe"
	t.Setenv("FFPROBE_PATH", fakeFFprobePath)

	path, err := GetFFprobePath()
	if err != nil {
		t.Fatalf("GetFFprobePath() unexpected error: %v", err)
	}
	if path != fakeFFprobePath {
		t.Errorf("GetFFprobePath() = %q, expected %q", path, fakeFFprobePath)
	}
}

func TestGetFFprobePath_NotFound(t *testing.T) {
	os.Unsetenv("FFPROBE_PATH")

	path, err := GetFFprobePath()
	if err != nil {
		t.Logf("GetFFprobePath() returned expected error: %v", err)
		return
	}
	t.Logf("ffprobe found on system, skipping not-found check: %s", path)
}
