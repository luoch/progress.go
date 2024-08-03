package main

import (
	"log"
	"bytes"
	"math"
	"strconv"
	"github.com/mattn/go-runewidth"
	"github.com/gin-gonic/gin"
	"github.com/flysnow-org/soha"
	"text/template"
)

type Progress struct {
	Title		string 	`form:"title,default="`
    TitleColor	string	`form:"color,default=428bca"`
	Scale		int		`form:"scale,default=100"`
	Progress	int		`uri:"progress" binding:"required"`
	Suffix		string	`form:"suffix,default=%"`
	Prefix		string	`form:"prefix,default="`
	Width		int		`form:"width,default=0"`
	Height		int		`form:"height,default=20"`
	FontSize	int		`form:"fontsize,default=11"`
	Size		int		`form:"size,default=17"`
}

type ProgressBar struct {
	Title			string
    TitleColor		string
	TitleWidth		int
	ProgressWidth	int
	BarWidth		int
	TotalWidth		int
	TotalHeight		int
	BarColor		string
	Suffix			string
	Progress		int
	Prefix			string
	TextPosX		int
	FontSize		int
}

type ProgressPie struct {
	TotalWidth		int
	TotalHeight		int
	PieColor		string
	CircleRadius	float32
	R				float32
	TargetX			float64
	TargetY			float64
	Suffix			string
	Prefix			string
	TextPosX		float32
	TextPosY		float32
	FontSize		int
	LargeArcFlag	int
	Progress		int
	Ratio			float32
}

func getColor(ratio float32) string {
	var color string = "#5cb85c"
	if ratio < 0.7 {
        color = "#f0ad4e"
	}
	if ratio < 0.3 {
		color = "#d9534f"
	}
    return color
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
	bar.Title = progress.Title
	bar.TitleColor = progress.TitleColor
	bar.FontSize = progress.FontSize
	bar.TotalHeight = progress.Height
	if progress.Title !="" {
		bar.TitleWidth = 10+int(math.Ceil(float64(bar.FontSize)*0.75))*runewidth.StringWidth(progress.Title)
		bar.ProgressWidth = 60
	}else{
		bar.TitleWidth = 0
		bar.ProgressWidth = 90
	}
	if progress.Width>0 {
		bar.ProgressWidth = progress.Width - bar.TitleWidth
	}

	bar.TotalWidth = bar.TitleWidth+bar.ProgressWidth
	var ratio = float32(progress.Progress)/float32(progress.Scale)
	bar.BarWidth = bar.ProgressWidth
	if ratio < 1 {
		bar.BarWidth = int(float32(bar.ProgressWidth)*ratio)
	}

	bar.BarColor = getColor(ratio)
	bar.Suffix = progress.Suffix
	bar.Prefix = progress.Prefix
	bar.Progress = progress.Progress
	bar.TextPosX = int(float32(bar.ProgressWidth)*0.5)+bar.TitleWidth
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
	var fulltext = progress.Prefix + strconv.Itoa(progress.Progress) + progress.Suffix
	var textwidth = runewidth.StringWidth(fulltext)*int(math.Ceil(float64(progress.FontSize)*0.75))
	pie.TotalWidth = textwidth+progress.Size+10
	pie.TotalHeight = progress.Size
	pie.CircleRadius = float32(progress.Size)/2.0-1.0
	pie.R = float32(progress.Size)/4.0
	pie.FontSize = progress.FontSize

	pie.Ratio = float32(progress.Progress)/float32(progress.Scale)
	pie.Progress = progress.Progress
	pie.PieColor = getColor(pie.Ratio)

	if pie.Ratio>0 && pie.Ratio<1 {
		var alpha = float64(2.0 * math.Pi * pie.Ratio)
		pie.TargetX = math.Sin(alpha)*float64(pie.R)
		pie.TargetY = math.Cos(alpha-math.Pi)*float64(pie.R)
	}else{
		pie.TargetX = 0
		pie.TargetY = 0
	}
	pie.LargeArcFlag = 0
	if pie.Ratio>0.5 {
		pie.LargeArcFlag = 1
	}

	pie.Suffix = progress.Suffix
	pie.Prefix = progress.Prefix
	
	pie.TextPosX = pie.CircleRadius+8
	pie.TextPosY = float32(progress.FontSize)/2.0-2.0

	var buf bytes.Buffer
	if execErr := tmpl.Execute(&buf, pie); execErr != nil {
		log.Println(execErr)
		c.JSON(500, gin.H{"msg": execErr})
		return
	}
	c.Header("Content-Type", "image/svg+xml")
	c.String(200, buf.String())
}

func TemplateMiddleware() gin.HandlerFunc {
	sohaFuncMap := soha.CreateFuncMap()
	progressbarTmpl, err := template.New("progressbar.svg").Funcs(sohaFuncMap).ParseFiles("template/progressbar.svg")
	progresspieTmpl, err := template.New("progresspie.svg").Funcs(sohaFuncMap).ParseFiles("template/progresspie.svg")
	if err != nil {
		log.Fatal(err)
	}

    return func(c *gin.Context) {
        c.Set("progressbar_template", progressbarTmpl)
		c.Set("progresspie_template", progresspieTmpl)
        c.Next()
    }
}

func main() {
	r := gin.Default()
	r.Use(TemplateMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/bar/:progress", getProgressbar)
	r.GET("/pie/:progress", getProgresspie)

	r.Run(":8000") // listen and serve on 0.0.0.0:8000
}
