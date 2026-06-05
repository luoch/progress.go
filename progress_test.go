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

func TestGetTextPosition(t *testing.T) {
	tests := []struct {
		align      string
		wantX      int
		wantShadow int
		wantAnchor string
	}{
		{align: "left", wantX: 105, wantShadow: 106, wantAnchor: "start"},
		{align: "center", wantX: 149, wantShadow: 150, wantAnchor: "middle"},
		{align: "right", wantX: 195, wantShadow: 196, wantAnchor: "end"},
		{align: "bad-value", wantX: 149, wantShadow: 150, wantAnchor: "middle"},
	}

	for _, tt := range tests {
		t.Run(tt.align, func(t *testing.T) {
			gotX, gotShadow, gotAnchor := getTextPosition(100, 100, 11, tt.align)
			if gotX != tt.wantX || gotShadow != tt.wantShadow || gotAnchor != tt.wantAnchor {
				t.Fatalf("getTextPosition align %q = (%d, %d, %q), want (%d, %d, %q)", tt.align, gotX, gotShadow, gotAnchor, tt.wantX, tt.wantShadow, tt.wantAnchor)
			}
		})
	}
}

func TestGetOuterPath(t *testing.T) {
	if got := getOuterPath(0, 0, 40, 20, 4, true, false); got != `M4 0H40V20H4Q0 20 0 16V4Q0 0 4 0z` {
		t.Fatalf("left-rounded path = %q", got)
	}
	if got := getOuterPath(40, 0, 80, 20, 4, false, true); got != `M40 0H116Q120 0 120 4V16Q120 20 116 20H40z` {
		t.Fatalf("right-rounded path = %q", got)
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
		`id="preview-background"`,
		`id="generated-url"`,
		`id="generated-markdown"`,
		`data-copy="generated-url"`,
		`data-copy="generated-markdown"`,
		`data-theme="neon"`,
		`data-skin="pill"`,
		`data-align="right"`,
		`id="titleheight"`,
		`addParam(params, "theme", state.theme, "classic")`,
		`addParam(params, "skin", state.skin, "badge")`,
		`addParam(params, "titleheight", positiveInteger(fields.titleHeight.value, 0), 0)`,
		`addParam(params, "align", state.align, "center")`,
		`previewShell.style.backgroundColor = fields.previewBackground.value || "#ffffff"`,
		"`/${state.type}/${progress}${query ? `?${query}` : \"\"}`",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("GET / body does not contain %q", want)
		}
	}

	if strings.Contains(body, ".preview-shell[data-theme=") {
		t.Fatal("GET / body should not bind preview background to the selected theme")
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
		`d="M6 0H37V20H6Q0 20 0 14V6Q0 0 6 0z"`,
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
		`d="M4 0H120V28H4Q0 28 0 24V4Q0 0 4 0z"`,
		`d="M120 0H236Q240 0 240 4V24Q240 28 236 28H120z"`,
		`d="M120 0h4v28h-4z"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("title-sized bar body does not contain %q", want)
		}
	}
}

func TestProgressBarTextAlignAndTitleHeight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := newRouter()

	req := httptest.NewRequest(http.MethodGet, "/bar/50?title=progress&titlewidth=120&titleheight=16&height=28&width=240&align=right", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /bar text align status = %d, want %d", w.Code, http.StatusOK)
	}
	body := w.Body.String()
	for _, want := range []string{
		`<svg width="240" height="28"`,
		`d="M4 6H120V22H4Q0 22 0 18V10Q0 6 4 6z"`,
		`text-anchor="end"`,
		`<text x="235" y="18">`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("aligned bar body does not contain %q", want)
		}
	}
}
