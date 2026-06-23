package media

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var customFFmpegPath string

func init() {
	flag.StringVar(&customFFmpegPath, "ffmpeg-path", "", "path to ffmpeg executable")
}

func GetFFmpegPath() (string, error) {
	if customFFmpegPath != "" {
		return customFFmpegPath, nil
	}
	if envPath := os.Getenv("FFMPEG_PATH"); envPath != "" {
		return envPath, nil
	}
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return path, nil
	}
	return "", fmt.Errorf("ffmpeg not found. Please install FFmpeg or specify its location using --ffmpeg-path or FFMPEG_PATH")
}

func GetFFprobePath() (string, error) {
	if envPath := os.Getenv("FFPROBE_PATH"); envPath != "" {
		return envPath, nil
	}

	if path, err := exec.LookPath("ffprobe"); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("ffprobe not found. Please install FFmpeg or specify FFPROBE_PATH")
}

func GetMediaDuration(inputPath string, ffprobePath string) (float64, error) {
	args := []string{
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputPath,
	}

	cmd := exec.Command(ffprobePath, args...)

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("an error ocurred while ffprobe running: %v", err)
	}

	durationStr := strings.TrimSpace(string(output))

	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, fmt.Errorf("time duration couldn't convert to a mathmatical number: %v", err)
	}

	return duration, nil
}

func ParseTimeToSeconds(timeStr string) float64 {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0
	}

	hours, _ := strconv.ParseFloat(parts[0], 64)
	minutes, _ := strconv.ParseFloat(parts[1], 64)
	seconds, _ := strconv.ParseFloat(parts[2], 64)

	return (hours * 3600) + (minutes * 60) + seconds
}

// timeRegex extracts the time=HH:MM:SS.xx portion from FFmpeg progress output
var timeRegex = regexp.MustCompile(`time=(\d+:\d+:\d+\.\d+)`)

// MediaConverter runs FFmpeg to convert inputPath → outputPath.
// onProgress is called with a percentage (0–100) whenever FFmpeg reports progress.
// If totalDuration <= 0, progress tracking is skipped (callback is not called).
func MediaConverter(inputPath, outputPath, ffmpeg string, totalDuration float64, onProgress func(percent float64)) error {
	cmd := exec.Command(ffmpeg, "-y", "-i", inputPath, outputPath)

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe couldn't open: %v", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("ffmpeg couldn't start: %v", err)
	}

	scanner := bufio.NewScanner(stderrPipe)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if totalDuration > 0 && onProgress != nil && strings.Contains(line, "time=") {
			matches := timeRegex.FindStringSubmatch(line)
			if len(matches) >= 2 {
				currentSeconds := ParseTimeToSeconds(matches[1])
				percent := (currentSeconds / totalDuration) * 100
				if percent > 100 {
					percent = 100
				}
				onProgress(percent)
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("an error occured while converting: %v", err)
	}

	return nil
}
