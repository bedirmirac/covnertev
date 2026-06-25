package main

import (
	"strings"
	"testing"
)

// strPtr is a helper to create *string values inline.
func strPtr(s string) *string {
	return &s
}

func TestArgsValidate(t *testing.T) {
	tests := []struct {
		name        string
		args        Args
		wantErr     bool
		errContains string // Hata mesajında bulunması gereken substring
	}{
		{
			name: "mode boş",
			args: Args{
				Mode:        strPtr(""),
				Input:       strPtr("input.txt"),
				Output:      strPtr(""),
				FFmpegPath:  strPtr(""),
				FFprobePath: strPtr(""),
			},
			wantErr:     true,
			errContains: "usage",
		},
		{
			name: "input boş",
			args: Args{
				Mode:        strPtr("doc"),
				Input:       strPtr(""),
				Output:      strPtr(""),
				FFmpegPath:  strPtr(""),
				FFprobePath: strPtr(""),
			},
			wantErr:     true,
			errContains: "usage",
		},
		{
			name: "doc modunda output verilmiş — yasak",
			args: Args{
				Mode:        strPtr("doc"),
				Input:       strPtr("input.docx"),
				Output:      strPtr("output.pdf"),
				FFmpegPath:  strPtr(""),
				FFprobePath: strPtr(""),
			},
			wantErr:     true,
			errContains: "Do NOT use '-o'",
		},
		{
			name: "img modunda output eksik",
			args: Args{
				Mode:        strPtr("img"),
				Input:       strPtr("input.png"),
				Output:      strPtr(""),
				FFmpegPath:  strPtr(""),
				FFprobePath: strPtr(""),
			},
			wantErr:     true,
			errContains: "Output is needed",
		},
		{
			name: "media modunda output eksik",
			args: Args{
				Mode:        strPtr("media"),
				Input:       strPtr("input.mp4"),
				Output:      strPtr(""),
				FFmpegPath:  strPtr("ffmpeg"),
				FFprobePath: strPtr("ffprobe"),
			},
			wantErr:     true,
			errContains: "Output is needed",
		},
		// --- Geçerli senaryolar ---
		{
			name: "doc modu geçerli",
			args: Args{
				Mode:        strPtr("doc"),
				Input:       strPtr("input.docx"),
				Output:      strPtr(""),
				FFmpegPath:  strPtr(""),
				FFprobePath: strPtr(""),
			},
			wantErr: false,
		},
		{
			name: "img modu geçerli",
			args: Args{
				Mode:        strPtr("img"),
				Input:       strPtr("input.png"),
				Output:      strPtr("output.jpg"),
				FFmpegPath:  strPtr(""),
				FFprobePath: strPtr(""),
			},
			wantErr: false,
		},
		{
			name: "media modu — ffmpeg ve ffprobe önceden verilmiş, oto-arama yapılmamalı",
			args: Args{
				Mode:        strPtr("media"),
				Input:       strPtr("input.mp4"),
				Output:      strPtr("output.mp4"),
				FFmpegPath:  strPtr("/usr/bin/ffmpeg"),
				FFprobePath: strPtr("/usr/bin/ffprobe"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() hata bekleniyor ama nil döndü")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Validate() hata = %q, içermesi beklenen: %q", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() beklenmedik hata: %v", err)
				}
			}
		})
	}
}
