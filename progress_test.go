package main

import "testing"

func TestGetTitleWidth(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		fontSize int
		want     int
	}{
		{
			name:     "empty title",
			title:    "",
			fontSize: 11,
			want:     0,
		},
		{
			name:     "ascii title uses proportional width",
			title:    "progress",
			fontSize: 11,
			want:     63,
		},
		{
			name:     "wide letters use extra room",
			title:    "www",
			fontSize: 11,
			want:     39,
		},
		{
			name:     "narrow letters stay compact",
			title:    "ill",
			fontSize: 11,
			want:     22,
		},
		{
			name:     "cjk title keeps full-width room",
			title:    "进度",
			fontSize: 11,
			want:     32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getTitleWidth(tt.title, tt.fontSize)
			if got != tt.want {
				t.Fatalf("getTitleWidth(%q, %d) = %d, want %d", tt.title, tt.fontSize, got, tt.want)
			}
		})
	}
}
