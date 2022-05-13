package main

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/gocolly/colly"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

//Estructuras y variables globales
type LRA struct {
	lenguaje  string
	rating    float32
	aparicion int
}

//rutinas
func check(err error) {
	if err != nil {
		fmt.Print("\nHay un error!!: ", err)
	}
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
func ordenar(datos []LRA) []LRA {
	var auxiliar LRA
	for n := 1; n < len(datos)-1; n++ {
		for j := 0; j < len(datos)-n; j++ {
			if datos[j].aparicion < datos[j+1].aparicion {
				auxiliar = datos[j]
				datos[j] = datos[j+1]
				datos[j+1] = auxiliar
			}
		}
	}
	return datos
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
func timeSub(t1, t2 time.Time) int {
	t1 = t1.UTC().Truncate(24 * time.Hour)
	t2 = t2.UTC().Truncate(24 * time.Hour)
	return int(t1.Sub(t2).Hours() / 24)
}
func Tema1() {
	println("Iniciando Tema 1....")
	file, errors := os.Create("Resultados.csv")
	check(errors)
	writer := csv.NewWriter(file)
	tiobe_github := [20]string{"python", "c", "java", "cpp", "csharp", "visual-basic", "javascript", "assembly", "sql", "php", "r", "delphi", "go", "swift", "ruby", "visual-basic-6", "objective-c", "perl", "lua", "matlab"}
	tiobe_lenguajes := [20]string{"python", "C", "Java", "C++", "C#", "Visual Basic", "JavaScript", "Assembly language", "SQL", "PHP", "R", "Delphi/Object Pascal", "Go", "Swift", "Ruby", "Classic Visual Basic", "Objective-C", "Perl", "Lua", "MATLAB"}
	datos := [20]LRA{}
	mayor := 0
	menor := math.MaxInt
	var cadena string
	for j := range tiobe_github {
		time.Sleep(3 * time.Second)
		direccion := "https://github.com/topics/"
		url := direccion + tiobe_github[j]
		println("Visiting " + url)
		response, errors := http.Get(url)
		defer response.Body.Close()
		check(errors)
		if response.StatusCode > 400 {
			fmt.Print("\nStatus code:", response.StatusCode)
		}
		doc, errors := goquery.NewDocumentFromReader(response.Body)
		check(errors)
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
		cadena = ""
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
	fmt.Println("\n\tLenguajes\t\tRating\t\tNºAparicion")
	for i := range datos {
		fmt.Printf("%20s\t\t%.2f\t\t%d\n", datos[i].lenguaje, datos[i].rating, datos[i].aparicion)
	}
	grafico(datos)
}

func Tema2() {
	println("Iniciando Tema 2....")
	file, errors := os.Create("TopicDeInteres.csv")
	check(errors)
	writer := csv.NewWriter(file)
	topicSeleccionado := "https://github.com/topics/go?page="
	var i = 0
	var TopicApariciones = make(map[string]int)
	c := colly.NewCollector(
	//colly.AllowedDomains("https://github.com"),
	)
	c.OnHTML("article.border.rounded.color-shadow-small.color-bg-subtle.my-4", func(e *colly.HTMLElement) {
		horaActualizacion := strings.TrimSpace(e.ChildAttr("relative-time.no-wrap", "datetime"))
		Topic := strings.TrimSpace(e.ChildText("a.topic-tag.topic-tag-link.f6.mb-2"))
		//Topic = strings.ReplaceAll(Topic, "\n\n", " ")
		listaTopic := strings.Fields(Topic)
		FechaExtraido := horaActualizacion[:10] //fecha extraida
		//println(FechaExtraido)
		fecha, error := time.Parse("2006-01-02", FechaExtraido) //parseo de fecha time.time
		check(error)
		diaDiferencia := timeSub(time.Now(), fecha) //se realiza la resta de la fecha
		if listaTopic != nil {
			if diaDiferencia <= 30 {
				for _, s := range listaTopic {
					ap := TopicApariciones[strings.TrimSpace(s)]
					TopicApariciones[strings.TrimSpace(s)] = (ap + 1)
				}
			}
		}
		i++

		if i < 35 {
			time.Sleep(3 * time.Second)
			c.Visit(topicSeleccionado + strconv.Itoa(i))
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Recopilando datos de: ", r.URL.String())
	})
	//c.OnResponse(func(r *colly.Response) {
	//	log.Println("response received", r.StatusCode)
	//	re = r.StatusCode
	//})
	c.Visit(topicSeleccionado + strconv.Itoa(0))

	TopicsOrdenados := []LRA{}
	for c, v := range TopicApariciones {
		TopicsOrdenados = append(TopicsOrdenados, LRA{lenguaje: c, rating: 0, aparicion: v})
		posts := []string{c, strconv.Itoa(v)}
		writer.Write(posts)
	}
	writer.Flush()
	TopicsOrdenados = ordenar(TopicsOrdenados)
	println("\tTopics\t\tNºApariciones")
	for i := 0; i < 20; i++ {
		fmt.Printf("\n%20s\t%d", TopicsOrdenados[i].lenguaje, TopicsOrdenados[i].aparicion)
	}
}
func main() {

	//Tema1()
	Tema2()

}
