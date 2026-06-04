package main

import (
	"bytes"
	"github.com/flysnow-org/soha"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-runewidth"
	"log"
	"math"
	"strconv"
	"strings"
	"text/template"
)

type Progress struct {
	Title      string `form:"title,default="`
	TitleColor string `form:"color,default=428bca"`
	Theme      string `form:"theme,default=classic"`
	Skin       string `form:"skin,default=badge"`
	Scale      int    `form:"scale,default=100"`
	Progress   int    `uri:"progress" binding:"required"`
	Suffix     string `form:"suffix,default=%"`
	Prefix     string `form:"prefix,default="`
	Width      int    `form:"width,default=0"`
	Height     int    `form:"height,default=20"`
	FontSize   int    `form:"fontsize,default=11"`
	Size       int    `form:"size,default=17"`
}

type ProgressBar struct {
	Title            string
	TitleColor       string
	TitleWidth       int
	ProgressWidth    int
	BarWidth         int
	TotalWidth       int
	TotalHeight      int
	BarColor         string
	TrackColor       string
	TextColor        string
	ShadowColor      string
	Radius           int
	GradientID       string
	GradientStart    string
	GradientEnd      string
	UseGradient      bool
	HighlightOpacity string
	ShadowOpacity    string
	Suffix           string
	Progress         int
	Prefix           string
	TextPosX         int
	FontSize         int
}

type ProgressPie struct {
	TotalWidth       int
	TotalHeight      int
	PieColor         string
	TextColor        string
	TrackColor       string
	GradientID       string
	GradientStart    string
	GradientEnd      string
	UseGradient      bool
	TrackStrokeWidth string
	LineCap          string
	CircleRadius     float32
	R                float32
	TargetX          float64
	TargetY          float64
	Suffix           string
	Prefix           string
	TextPosX         float32
	TextPosY         float32
	FontSize         int
	LargeArcFlag     int
	Progress         int
	Ratio            float32
}

type ProgressTheme struct {
	Name          string
	TitleColor    string
	TrackColor    string
	LowColor      string
	MidColor      string
	HighColor     string
	TextColor     string
	ShadowColor   string
	GradientStart string
	GradientEnd   string
	UseGradient   bool
	Radius        int
}

type ProgressSkin struct {
	Name             string
	HighlightOpacity string
	ShadowOpacity    string
	Radius           int
	PieTrackWidth    string
	PieLineCap       string
}

var progressThemes = map[string]ProgressTheme{
	"classic": {
		Name:        "classic",
		TitleColor:  "428bca",
		TrackColor:  "#555",
		LowColor:    "#d9534f",
		MidColor:    "#f0ad4e",
		HighColor:   "#5cb85c",
		TextColor:   "#fff",
		ShadowColor: "#010101",
		Radius:      4,
	},
	"slate": {
		Name:        "slate",
		TitleColor:  "111827",
		TrackColor:  "#e5e7eb",
		LowColor:    "#e11d48",
		MidColor:    "#d97706",
		HighColor:   "#2563eb",
		TextColor:   "#fff",
		ShadowColor: "#020617",
		Radius:      6,
	},
	"mint": {
		Name:        "mint",
		TitleColor:  "065f46",
		TrackColor:  "#d1fae5",
		LowColor:    "#f43f5e",
		MidColor:    "#14b8a6",
		HighColor:   "#059669",
		TextColor:   "#fff",
		ShadowColor: "#064e3b",
		Radius:      8,
	},
	"amber": {
		Name:        "amber",
		TitleColor:  "78350f",
		TrackColor:  "#fde68a",
		LowColor:    "#dc2626",
		MidColor:    "#f59e0b",
		HighColor:   "#84cc16",
		TextColor:   "#fff",
		ShadowColor: "#451a03",
		Radius:      5,
	},
	"neon": {
		Name:          "neon",
		TitleColor:    "020617",
		TrackColor:    "#1e293b",
		LowColor:      "#fb7185",
		MidColor:      "#22d3ee",
		HighColor:     "#a3e635",
		TextColor:     "#f8fafc",
		ShadowColor:   "#020617",
		GradientStart: "#22d3ee",
		GradientEnd:   "#a3e635",
		UseGradient:   true,
		Radius:        6,
	},
	"mono": {
		Name:        "mono",
		TitleColor:  "18181b",
		TrackColor:  "#e4e4e7",
		LowColor:    "#52525b",
		MidColor:    "#3f3f46",
		HighColor:   "#18181b",
		TextColor:   "#fff",
		ShadowColor: "#000",
		Radius:      0,
	},
}

var progressSkins = map[string]ProgressSkin{
	"badge": {
		Name:             "badge",
		HighlightOpacity: "1",
		ShadowOpacity:    ".3",
		Radius:           -1,
		PieTrackWidth:    "2",
		PieLineCap:       "butt",
	},
	"flat": {
		Name:             "flat",
		HighlightOpacity: "0",
		ShadowOpacity:    "0",
		Radius:           -1,
		PieTrackWidth:    "2",
		PieLineCap:       "butt",
	},
	"soft": {
		Name:             "soft",
		HighlightOpacity: ".55",
		ShadowOpacity:    ".18",
		Radius:           8,
		PieTrackWidth:    "2.5",
		PieLineCap:       "round",
	},
	"pill": {
		Name:             "pill",
		HighlightOpacity: ".6",
		ShadowOpacity:    ".18",
		Radius:           999,
		PieTrackWidth:    "3",
		PieLineCap:       "round",
	},
}

func getTheme(name string) ProgressTheme {
	key := strings.ToLower(strings.TrimSpace(name))
	if theme, ok := progressThemes[key]; ok {
		return theme
	}
	return progressThemes["classic"]
}

func getSkin(name string) ProgressSkin {
	key := strings.ToLower(strings.TrimSpace(name))
	if skin, ok := progressSkins[key]; ok {
		return skin
	}
	return progressSkins["badge"]
}

func getColor(ratio float32, theme ProgressTheme) string {
	var color string = theme.HighColor
	if ratio < 0.7 {
		color = theme.MidColor
	}
	if ratio < 0.3 {
		color = theme.LowColor
	}
	return color
}

func getTitleWidth(title string, fontSize int) int {
	if title == "" {
		return 0
	}
	return 10 + getSVGTextWidth(title, fontSize)
}

func getSVGTextWidth(text string, fontSize int) int {
	if text == "" || fontSize <= 0 {
		return 0
	}

	width := 0.0
	for _, r := range text {
		runeWidth := runewidth.RuneWidth(r)
		switch {
		case runeWidth == 0:
			continue
		case runeWidth > 1:
			width += float64(fontSize)
		case r == ' ':
			width += float64(fontSize) * 0.35
		case strings.ContainsRune("ilI.,:;!|`'", r):
			width += float64(fontSize) * 0.35
		case strings.ContainsRune("mwMW@#%&", r):
			width += float64(fontSize) * 0.85
		default:
			width += float64(fontSize) * 0.6
		}
	}
	return int(math.Ceil(width))
}

func getProgressbar(c *gin.Context) {
	tmpl, tmplerr := c.MustGet("progressbar_template").(*template.Template)
	if !tmplerr {
		log.Println(tmplerr)
		c.JSON(500, gin.H{"msg": tmplerr})
		return
	}

	var progress Progress
	var bar ProgressBar
	if err := c.ShouldBindUri(&progress); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	c.Bind(&progress)
	theme := getTheme(progress.Theme)
	skin := getSkin(progress.Skin)
	bar.Title = progress.Title
	bar.TitleColor = theme.TitleColor
	if titleColor := c.Query("color"); titleColor != "" {
		bar.TitleColor = strings.TrimPrefix(titleColor, "#")
	}
	bar.FontSize = progress.FontSize
	bar.TotalHeight = progress.Height
	if progress.Title != "" {
		bar.TitleWidth = getTitleWidth(progress.Title, bar.FontSize)
		bar.ProgressWidth = 60
	} else {
		bar.TitleWidth = 0
		bar.ProgressWidth = 90
	}
	if progress.Width > 0 {
		bar.ProgressWidth = progress.Width - bar.TitleWidth
	}

	bar.TotalWidth = bar.TitleWidth + bar.ProgressWidth
	var ratio = float32(progress.Progress) / float32(progress.Scale)
	bar.BarWidth = bar.ProgressWidth
	if ratio < 1 {
		bar.BarWidth = int(float32(bar.ProgressWidth) * ratio)
	}

	bar.BarColor = getColor(ratio, theme)
	bar.TrackColor = theme.TrackColor
	bar.TextColor = theme.TextColor
	bar.ShadowColor = theme.ShadowColor
	bar.Radius = theme.Radius
	if skin.Radius >= 0 {
		bar.Radius = skin.Radius
	}
	bar.GradientID = "progress-gradient-" + theme.Name
	bar.GradientStart = theme.GradientStart
	bar.GradientEnd = theme.GradientEnd
	bar.UseGradient = theme.UseGradient
	bar.HighlightOpacity = skin.HighlightOpacity
	bar.ShadowOpacity = skin.ShadowOpacity
	bar.Suffix = progress.Suffix
	bar.Prefix = progress.Prefix
	bar.Progress = progress.Progress
	bar.TextPosX = int(float32(bar.ProgressWidth)*0.5) + bar.TitleWidth
	var buf bytes.Buffer
	if execErr := tmpl.Execute(&buf, bar); execErr != nil {
		log.Println(execErr)
		c.JSON(500, gin.H{"msg": execErr})
		return
	}
	c.Header("Content-Type", "image/svg+xml")
	c.String(200, buf.String())
}

func getProgresspie(c *gin.Context) {
	tmpl, tmplerr := c.MustGet("progresspie_template").(*template.Template)
	if !tmplerr {
		c.JSON(500, gin.H{"msg": tmplerr})
		return
	}

	var progress Progress
	var pie ProgressPie
	if err := c.ShouldBindUri(&progress); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"msg": err})
		return
	}
	c.Bind(&progress)
	theme := getTheme(progress.Theme)
	skin := getSkin(progress.Skin)
	var fulltext = progress.Prefix + strconv.Itoa(progress.Progress) + progress.Suffix
	var textwidth = runewidth.StringWidth(fulltext) * int(math.Ceil(float64(progress.FontSize)*0.75))
	pie.TotalWidth = textwidth + progress.Size + 10
	pie.TotalHeight = progress.Size
	pie.CircleRadius = float32(progress.Size)/2.0 - 1.0
	pie.R = float32(progress.Size) / 4.0
	pie.FontSize = progress.FontSize

	pie.Ratio = float32(progress.Progress) / float32(progress.Scale)
	pie.Progress = progress.Progress
	pie.PieColor = getColor(pie.Ratio, theme)
	pie.TrackColor = theme.TrackColor
	pie.TextColor = "#" + theme.TitleColor
	if theme.Name == "neon" {
		pie.TextColor = theme.TextColor
	}
	pie.GradientID = "progress-gradient-" + theme.Name
	pie.GradientStart = theme.GradientStart
	pie.GradientEnd = theme.GradientEnd
	pie.UseGradient = theme.UseGradient
	pie.TrackStrokeWidth = skin.PieTrackWidth
	pie.LineCap = skin.PieLineCap

	if pie.Ratio > 0 && pie.Ratio < 1 {
		var alpha = float64(2.0 * math.Pi * pie.Ratio)
		pie.TargetX = math.Sin(alpha) * float64(pie.R)
		pie.TargetY = math.Cos(alpha-math.Pi) * float64(pie.R)
	} else {
		pie.TargetX = 0
		pie.TargetY = 0
	}
	pie.LargeArcFlag = 0
	if pie.Ratio > 0.5 {
		pie.LargeArcFlag = 1
	}

	pie.Suffix = progress.Suffix
	pie.Prefix = progress.Prefix

	pie.TextPosX = pie.CircleRadius + 8
	pie.TextPosY = float32(progress.FontSize)/2.0 - 2.0

	var buf bytes.Buffer
	if execErr := tmpl.Execute(&buf, pie); execErr != nil {
		log.Println(execErr)
		c.JSON(500, gin.H{"msg": execErr})
		return
	}
	c.Header("Content-Type", "image/svg+xml")
	c.String(200, buf.String())
}

func getIndex(c *gin.Context) {
	tmpl, tmplerr := c.MustGet("index_template").(*template.Template)
	if !tmplerr {
		c.JSON(500, gin.H{"msg": tmplerr})
		return
	}

	var buf bytes.Buffer
	if execErr := tmpl.Execute(&buf, nil); execErr != nil {
		log.Println(execErr)
		c.JSON(500, gin.H{"msg": execErr})
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, buf.String())
}

func TemplateMiddleware() gin.HandlerFunc {
	sohaFuncMap := soha.CreateFuncMap()
	progressbarTmpl, err := template.New("progressbar.svg").Funcs(sohaFuncMap).ParseFiles("template/progressbar.svg")
	if err != nil {
		log.Fatal(err)
	}
	progresspieTmpl, err := template.New("progresspie.svg").Funcs(sohaFuncMap).ParseFiles("template/progresspie.svg")
	if err != nil {
		log.Fatal(err)
	}
	indexTmpl, err := template.New("index.html").ParseFiles("template/index.html")
	if err != nil {
		log.Fatal(err)
	}

	return func(c *gin.Context) {
		c.Set("progressbar_template", progressbarTmpl)
		c.Set("progresspie_template", progresspieTmpl)
		c.Set("index_template", indexTmpl)
		c.Next()
	}
}

func newRouter() *gin.Engine {
	r := gin.Default()
	r.Use(TemplateMiddleware())

	r.GET("/", getIndex)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/bar/:progress", getProgressbar)
	r.GET("/pie/:progress", getProgresspie)

	return r
}

func main() {
	r := newRouter()
	r.Run(":8000") // listen and serve on 0.0.0.0:8000
}
