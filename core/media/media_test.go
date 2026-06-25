package media

import (
	"os"
	"testing"
)

// ─────────────────────────────────────────
// ParseTimeToSeconds
// ─────────────────────────────────────────

func TestParseTimeToSeconds(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"sıfır", "00:00:00.00", 0},
		{"sadece saniye — tam", "00:00:30.00", 30},
		{"sadece saniye — ondalık", "00:00:30.50", 30.5},
		{"sadece dakika", "00:01:00.00", 60},
		{"sadece saat", "01:00:00.00", 3600},
		{"karışık değer", "01:30:45.25", 5445.25},
		{"büyük saat değeri", "02:00:00.00", 7200},
		// Hatalı formatlar — mevcut davranış: sessizce 0 döner
		{"geçersiz string", "invalid", 0},
		{"boş string", "", 0},
		{"eksik kısım (2 parça)", "01:02", 0},
		{"fazla kısım (4 parça)", "01:02:03:04", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTimeToSeconds(tt.input)
			if got != tt.expected {
				t.Errorf("ParseTimeToSeconds(%q) = %v, beklenen %v", tt.input, got, tt.expected)
			}
		})
	}
}

// ─────────────────────────────────────────
// GetFFmpegPath
// ─────────────────────────────────────────

func TestGetFFmpegPath_EnvVar(t *testing.T) {
	// customFFmpegPath'i temizle; flag ile set edilmemiş olmalı test ortamında
	original := customFFmpegPath
	customFFmpegPath = ""
	defer func() { customFFmpegPath = original }()

	fakeFFmpegPath := "/fake/path/to/ffmpeg"
	t.Setenv("FFMPEG_PATH", fakeFFmpegPath)

	path, err := GetFFmpegPath()
	if err != nil {
		t.Fatalf("GetFFmpegPath() beklenmedik hata: %v", err)
	}
	if path != fakeFFmpegPath {
		t.Errorf("GetFFmpegPath() = %q, beklenen %q", path, fakeFFmpegPath)
	}
}

func TestGetFFmpegPath_CustomFlag(t *testing.T) {
	original := customFFmpegPath
	customFFmpegPath = "/custom/ffmpeg"
	defer func() { customFFmpegPath = original }()

	// Env var da set edilmiş olsa bile flag öncelikli olmalı
	t.Setenv("FFMPEG_PATH", "/env/ffmpeg")

	path, err := GetFFmpegPath()
	if err != nil {
		t.Fatalf("GetFFmpegPath() beklenmedik hata: %v", err)
	}
	if path != "/custom/ffmpeg" {
		t.Errorf("GetFFmpegPath() = %q, beklenen %q (flag öncelikli olmalı)", path, "/custom/ffmpeg")
	}
}

func TestGetFFmpegPath_NotFound(t *testing.T) {
	// Tüm kaynakları temizle; sistem PATH'inde ffmpeg yoksa hata döner
	original := customFFmpegPath
	customFFmpegPath = ""
	defer func() { customFFmpegPath = original }()

	os.Unsetenv("FFMPEG_PATH")

	path, err := GetFFmpegPath()
	if err != nil {
		// Beklenen durum: ffmpeg kurulu değil
		t.Logf("GetFFmpegPath() hata döndü (bekleniyor): %v", err)
		return
	}
	// ffmpeg kurulu olan sistemde test geçerli, path loglanır
	t.Logf("ffmpeg sistemde mevcut, atlıyor: %s", path)
}

// ─────────────────────────────────────────
// GetFFprobePath
// ─────────────────────────────────────────

func TestGetFFprobePath_EnvVar(t *testing.T) {
	fakeFFprobePath := "/fake/path/to/ffprobe"
	t.Setenv("FFPROBE_PATH", fakeFFprobePath)

	path, err := GetFFprobePath()
	if err != nil {
		t.Fatalf("GetFFprobePath() beklenmedik hata: %v", err)
	}
	if path != fakeFFprobePath {
		t.Errorf("GetFFprobePath() = %q, beklenen %q", path, fakeFFprobePath)
	}
}

func TestGetFFprobePath_NotFound(t *testing.T) {
	os.Unsetenv("FFPROBE_PATH")

	path, err := GetFFprobePath()
	if err != nil {
		t.Logf("GetFFprobePath() hata döndü (bekleniyor): %v", err)
		return
	}
	t.Logf("ffprobe sistemde mevcut, atlıyor: %s", path)
}
