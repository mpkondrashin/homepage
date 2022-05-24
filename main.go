/*
HomePage (c) 2022 by Michael Kondrashin

HomePage - generate bookmarks

main.go - main source file
*/

package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

//go:embed index.template
var templData string

var (
	MetroBackground = []string{"#004050", "#0e6d38", "#001941", "#260930", "#261300"}
	MetroColors     = []string{"#A200FF", "#FF0097", "#00ABA9", "#8CBF26", "#A05000", "#E671B8", "#F09609", "#1BA1E2", "#E51400", "#339933"}
)

func BackgroundColor() string {
	hour := time.Now().Hour()
	return MetroBackground[hour%len(MetroBackground)]
}

func lighter(color string) string {
	result := "#"
	for i := 0; i < 3; i++ {
		v, err := strconv.ParseInt(color[i*2+1:i*2+3], 16, 8)
		if err != nil {
			return color
		}
		vLighter := (0xFF + 2*v) / 3
		result += fmt.Sprintf("%02x", vLighter)
	}
	return result
}

type Bookmark struct {
	Section string
	Url     string
	Label   string
	Tooltip string
	Color   string
}

type Bookmarks struct {
	Bookmarks []Bookmark
}

type Page struct {
	BackgroudColor string
	LighterColor   string
	Sections       map[string][]Bookmark
}

type P struct {
	A string
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Homepage - generate HTML for your bookmarks using yaml file.\nUsage: %s <input yaml file> <output HTML file>\n",
			os.Args[0])
		return
	}
	log.Println("Homepage started")
	dataFileName := os.Args[1]
	htmlFileName := os.Args[2]
	inFile := os.Stdin
	if dataFileName != "-" {
		var err error
		inFile, err = os.Open(dataFileName)
		if err != nil {
			log.Fatal(err)
		}
	}
	var yamlData []byte
	size, err := inFile.Read(yamlData)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %s file (%d bytes)", dataFileName, size)
	var bookmarks Bookmarks
	err = yaml.Unmarshal(yamlData, &bookmarks)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Parsed %d bookmarks", len(bookmarks.Bookmarks))
	bc := BackgroundColor()
	page := Page{
		BackgroudColor: bc,
		LighterColor:   lighter(bc),
		Sections:       make(map[string][]Bookmark),
	}
	colorIndex := 0
	for _, b := range bookmarks.Bookmarks {
		b.Color = MetroColors[colorIndex%len(MetroColors)]
		colorIndex++
		page.Sections[b.Section] = append(page.Sections[b.Section], b)
	}
	templ, err := template.New("HomePage").Parse(templData)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Template OK")
	f := os.Stdout
	if htmlFileName != "-" {
		f, err = os.Create(htmlFileName)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = templ.Execute(f, &page)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done")
}
