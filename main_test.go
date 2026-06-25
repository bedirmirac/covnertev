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
		errContains string // substring that must appear in the error message
	}{
		{
			name: "empty mode",
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
			name: "empty input",
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
			name: "doc mode with output flag set",
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
			name: "img mode missing output",
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
			name: "media mode missing output",
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
		// valid cases
		{
			name: "doc mode valid",
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
			name: "img mode valid",
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
			name: "media mode with ffmpeg and ffprobe already set, no auto-detection",
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
					t.Errorf("Validate() expected error but got nil")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Validate() error = %q, expected to contain: %q", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}
