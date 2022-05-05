package main

import (
	"encoding/csv"
	"fmt"
	goquery "github.com/PuerkitoBio/goquery"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func check(err error) {
	if err != nil {
		fmt.Print(err)
	}
}

type LRA struct {
	lenguaje  string
	rating    float32
	aparicion int
}

func writeFile(data, filename string) {
	file, errors := os.Create(filename)
	defer file.Close()
	check(errors)
	file.WriteString(data)
}

func calcularRating(x, Max, Min float32) float32 {
	var R float32
	R = (((x - Min) / (Max - Min)) * 100.0)
	R = float32(math.Round(float64(R*100)) / 100)
	return R
}

func grafico(datos [20]LRA) {
	lenguajes := make([]string, 0)
	items := make([]opts.BarData, 0)
	lenguajes = append(lenguajes, "")
	items = append(items, opts.BarData{Value: 0})

	for i := 0; i < 10; i++ {
		items = append(items, opts.BarData{Value: datos[i].aparicion})
		lenguajes = append(lenguajes, datos[i].lenguaje)
	}
	for j := range lenguajes {
		println(lenguajes[j])
	}
	//repeticiones := []opts.BarData{{Value: 1}, {Value: 38}, {Value: 29}, {Value: 22}, {Value: 13}, {Value: 11}}
	//lenguajes := []string{"USA", "China", "UK", "Russia", "South Korea", "Germany"}
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{Title: "Lenguaje vs. Apariciones"}),
		charts.WithXAxisOpts(opts.XAxis{
			Type:      "category",
			Name:      "Lenguajes",
			AxisLabel: &opts.AxisLabel{Show: true, FontSize: "15", ShowMaxLabel: true, Interval: "0", Rotate: 45.0},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			AxisLabel: &opts.AxisLabel{Show: true, Formatter: "{value}"},
			Name:      "Apariciones",
		}))

	bar.SetXAxis(lenguajes)
	bar.AddSeries("Apariciones", items)
	f, _ := os.Create("bar.html")
	bar.Render(f)
}

func main() {
	file, errors := os.Create("Resultados.csv")
	check(errors)
	writer := csv.NewWriter(file)
	tiobe_github := [20]string{"python", "c", "java", "cpp", "csharp", "visual-basic", "javascript", "assembly", "sql", "php", "r", "delphi", "go", "swift", "ruby", "visual-basic-6", "objective-c", "perl", "lua", "matlab"}
	tiobe_lenguajes := [20]string{"python", "C", "Java", "C++", "C#", "Visual Basic", "JavaScript", "Assembly language", "SQL", "PHP", "R", "Delphi/Object Pascal", "Go", "Swift", "Ruby", "Classic Visual Basic", "Objective-C", "Perl", "Lua", "MATLAB"}

	datos := [20]LRA{}
	mayor := 0
	menor := math.MaxInt
	for j := range tiobe_github {
		direccion := "https://github.com/topics/"
		url := direccion + tiobe_github[j]
		response, errors := http.Get(url)
		defer response.Body.Close()
		check(errors)
		if response.StatusCode > 400 {
			fmt.Print("Status code:", response.StatusCode)
		}
		doc, errors := goquery.NewDocumentFromReader(response.Body)
		check(errors)

		var cadena string
		doc.Find("div.application-main ").Each(func(index int, item *goquery.Selection) {
			h2 := item.Find("h2")
			texto := strings.TrimSpace(h2.Text())
			for i := range texto {
				x := rune(texto[i])
				if unicode.IsNumber(x) {
					a := string(x)
					cadena = cadena + a
				}
			}
		})
		posts := []string{tiobe_lenguajes[j], cadena}
		writer.Write(posts)

		apa, errors := strconv.Atoi(cadena)
		check(errors)
		if apa > mayor {
			mayor = apa
		}
		if apa < menor {
			menor = apa
		}
		datos[j] = LRA{
			lenguaje:  tiobe_lenguajes[j],
			rating:    0,
			aparicion: apa,
		}
		writer.Flush()
		//fmt.Print(tiobe_lenguajes[j], "\t\t", cadena, "\n")
	}
	for i := range datos {
		datos[i].rating = calcularRating(float32(datos[i].aparicion), float32(mayor), float32(menor))
	}
	var auxiliar LRA
	for n := 1; n < len(datos)-1; n++ {
		for j := 0; j < len(datos)-n; j++ {
			if datos[j].rating < datos[j+1].rating {
				auxiliar = datos[j]
				datos[j] = datos[j+1]
				datos[j+1] = auxiliar
			}
		}
	}
	for i := range datos {
		fmt.Printf("%s\t\t%.2f\t\t%d\n", datos[i].lenguaje, datos[i].rating, datos[i].aparicion)
	}
	grafico(datos)
}
