package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"io/ioutil"
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
	BackgroudColor        string
	LighterBackgroudColor string
	Sections              map[string][]Bookmark
}

type P struct {
	A string
}

const dataFileName = "data.yaml"

func main() {
	yamlData, err := ioutil.ReadFile(dataFileName)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %s file (%d bytes)", dataFileName, len(yamlData))
	var bookmarks Bookmarks
	err = yaml.Unmarshal(yamlData, &bookmarks)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Parsed %d bookmarks", len(bookmarks.Bookmarks))
	page := Page{
		Sections: make(map[string][]Bookmark),
	}
	colorIndex := 0
	for _, b := range bookmarks.Bookmarks {
		//log.Print(b.Section)
		b.Color = MetroColors[colorIndex%len(MetroColors)]
		colorIndex++
		page.Sections[b.Section] = append(page.Sections[b.Section], b)
	}
	page.BackgroudColor = BackgroundColor()
	page.LighterBackgroudColor = lighter(page.BackgroudColor)
	//log.Println(page)
	//	templData, err := ioutil.ReadFile("index.template")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	templ, err := template.New("HomePage").Parse(templData)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Template OK")
	err = templ.Execute(os.Stdout, &page)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done")
}
