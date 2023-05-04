/*
HomePage (c) 2022 by Michael Kondrashin

HomePage - generate bookmarks

main.go - main source file
*/

package main

import (
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"io"
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

func (b Bookmark) String() string {
	return fmt.Sprintf("%s: %s (%s)", b.Section, b.Label, b.Url)
}

type Bookmarks struct {
	Bookmarks []Bookmark
}

type Page struct {
	BackgroudColor string
	LighterColor   string
	Sections       map[string][]Bookmark
}

var ErrDuplicate = errors.New("duplicate")

func CheckForDuplicates(b *Bookmarks) error {
	urls := make(map[string]struct{})
	labels := make(map[string]struct{})
	for _, each := range b.Bookmarks {
		if _, ok := urls[each.Url]; ok {
			return fmt.Errorf("%v: URL %w", each, ErrDuplicate)
		}
		urls[each.Url] = struct{}{}
		if _, ok := labels[each.Label]; ok {
			return fmt.Errorf("%v: Label %w", each, ErrDuplicate)
		}
		labels[each.Label] = struct{}{}
	}
	return nil
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
	yamlData, err := io.ReadAll(inFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %s file (%d bytes)", dataFileName, len(yamlData))

	var bookmarks Bookmarks
	if err := yaml.Unmarshal(yamlData, &bookmarks); err != nil {
		log.Fatal(err)
	}
	log.Printf("Parsed %d bookmarks", len(bookmarks.Bookmarks))
	if err := CheckForDuplicates(&bookmarks); err != nil {
		log.Fatal(err)
	}

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
	if err := templ.Execute(f, &page); err != nil {
		log.Fatal(err)
	}
	log.Println("Done")
}
