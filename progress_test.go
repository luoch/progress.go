package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

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
			got := getTitleWidth(tt.title, tt.fontSize, 10)
			if got != tt.want {
				t.Fatalf("getTitleWidth(%q, %d) = %d, want %d", tt.title, tt.fontSize, got, tt.want)
			}
		})
	}
}

func TestGetTitleWidthUsesPadding(t *testing.T) {
	got := getTitleWidth("go", 11, 24)
	want := 38
	if got != want {
		t.Fatalf("getTitleWidth(%q, %d, %d) = %d, want %d", "go", 11, 24, got, want)
	}

	got = getTitleWidth("go", 11, -1)
	want = 14
	if got != want {
		t.Fatalf("getTitleWidth with negative padding = %d, want %d", got, want)
	}
}

func TestGetTheme(t *testing.T) {
	if got := getTheme("mint"); got.Name != "mint" {
		t.Fatalf("getTheme(%q).Name = %q, want mint", "mint", got.Name)
	}
	if got := getTheme("unknown"); got.Name != "classic" {
		t.Fatalf("getTheme(%q).Name = %q, want classic", "unknown", got.Name)
	}
}

func TestGetSkin(t *testing.T) {
	if got := getSkin("pill"); got.Name != "pill" {
		t.Fatalf("getSkin(%q).Name = %q, want pill", "pill", got.Name)
	}
	if got := getSkin("unknown"); got.Name != "badge" {
		t.Fatalf("getSkin(%q).Name = %q, want badge", "unknown", got.Name)
	}
}

func TestIndexPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := newRouter()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET / status = %d, want %d", w.Code, http.StatusOK)
	}
	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Fatalf("GET / Content-Type = %q, want text/html", contentType)
	}

	body := w.Body.String()
	for _, want := range []string{
		"Progress in Markdown",
		`id="preview"`,
		`id="generated-url"`,
		`id="generated-markdown"`,
		`data-copy="generated-url"`,
		`data-copy="generated-markdown"`,
		`data-theme="neon"`,
		`data-skin="pill"`,
		`addParam(params, "theme", state.theme, "classic")`,
		`addParam(params, "skin", state.skin, "badge")`,
		"`/${state.type}/${progress}${query ? `?${query}` : \"\"}`",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("GET / body does not contain %q", want)
		}
	}
}

func TestProgressThemeRendering(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := newRouter()

	req := httptest.NewRequest(http.MethodGet, "/bar/88?title=done&theme=neon&skin=flat", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /bar themed status = %d, want %d", w.Code, http.StatusOK)
	}
	body := w.Body.String()
	for _, want := range []string{
		`id="progress-gradient-neon"`,
		`fill="#020617"`,
		`url(#progress-gradient-neon)`,
		`rx="6"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("themed bar body does not contain %q", want)
		}
	}
	if strings.Contains(body, `fill="url(#a)"`) {
		t.Fatalf("flat themed bar should not contain highlight overlay")
	}

	req = httptest.NewRequest(http.MethodGet, "/pie/88?theme=mint&skin=pill", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /pie themed status = %d, want %d", w.Code, http.StatusOK)
	}
	body = w.Body.String()
	for _, want := range []string{
		`stroke: #d1fae5`,
		`fill="#065f46"`,
		`stroke-width: 3`,
		`stroke-linecap: round`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("themed pie body does not contain %q", want)
		}
	}
}

func TestProgressBarTitleSizing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := newRouter()

	req := httptest.NewRequest(http.MethodGet, "/bar/50?title=progress&titlewidth=120&titlepadding=40&height=28&width=240", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /bar title sizing status = %d, want %d", w.Code, http.StatusOK)
	}
	body := w.Body.String()
	for _, want := range []string{
		`<svg width="240" height="28"`,
		`x="120" width="120" height="28"`,
		`d="M120 0h4v28h-4z"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("title-sized bar body does not contain %q", want)
		}
	}
}
