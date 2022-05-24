package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

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

func main() {
	yamlData, err := ioutil.ReadFile("data.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var bookmarks Bookmarks
	err = yaml.Unmarshal(yamlData, &bookmarks)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %d bookmarks", len(bookmarks.Bookmarks))
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
	templData, err := ioutil.ReadFile("index.template")
	if err != nil {
		log.Fatal(err)
	}
	templ, err := template.New("HomePage").Parse(string(templData))
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

/*
func main() {
	// "section,url,label,tooltip"
	csvF, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(csvF)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(record)
	}
}

*/
