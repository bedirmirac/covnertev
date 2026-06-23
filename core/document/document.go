package document

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var customLibreOfficePath string

func init() {
	flag.StringVar(&customLibreOfficePath, "libreoffice-path", "", "LibreOffice executable (soffice) path")
}

func getLibreOfficePath() (string, error) {
	if customLibreOfficePath != "" {
		return customLibreOfficePath, nil
	}

	envPath := os.Getenv("LIBREOFFICE_PATH")
	if envPath != "" {
		return envPath, nil
	}

	var defaultPaths []string

	switch runtime.GOOS {
	case "windows":
		defaultPaths = []string{
			filepath.Join("C:", "Program Files", "LibreOffice", "program", "soffice.exe"),
			filepath.Join("C:", "Program Files (x86)", "LibreOffice", "program", "soffice.exe"),
		}
	case "darwin":
		defaultPaths = []string{
			"/Applications/LibreOffice.app/Contents/MacOS/soffice",
			"/opt/homebrew/bin/soffice", // Homebrew ile kurulanlar için
		}
	default:
		defaultPaths = []string{
			"libreoffice",
			"soffice",
		}
	}

	for _, p := range defaultPaths {
		if runtime.GOOS == "linux" || runtime.GOOS == "freebsd" {
			if path, err := exec.LookPath(p); err == nil {
				return path, nil
			}
		} else {
			if _, err := os.Stat(p); err == nil {
				return p, nil
			}
		}
	}

	return "", fmt.Errorf("libreoffice couldn't be found. Please install Libreoffice or specify its location using --libreoffice-path")
}

func PdfFromOffice(inputPath string) error {
	libreofficeExe, err := getLibreOfficePath()
	if err != nil {
		log.Fatal(err)
	}
	args := []string{
		"--headless",
		"--convert-to", "pdf",
		inputPath,
	}

	cmd := exec.Command(libreofficeExe, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("libre Office error: %v\n IN DETAIL: %s", err, string(output))
	}
	return nil
}

func DocFromPdf(inputPath string) error {
	libreofficeExe, err := getLibreOfficePath()
	if err != nil {
		log.Fatal(err)
	}
	args := []string{
		"--headless",
		"--infilter=writer_pdf_import",
		"--convert-to", "doc",
		inputPath,
	}

	cmd := exec.Command(libreofficeExe, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("libre Office Error: %v\n IN DETAIL:  %s", err, string(output))
	}
	return nil
}
