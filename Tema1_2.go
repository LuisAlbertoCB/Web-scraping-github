package main

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/gocolly/colly"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"
)

//Estructuras y variables globales
type LRA struct { //esta estructura almacena el nombre del lenguaje, su rating y su numero de aparicion
	lenguaje  string
	rating    float32
	aparicion int
}

//rutinas
func check(err error) { //esta rutina verifica errores
	if err != nil {
		fmt.Print("\nHay un error!!: ", err)
	}
}
func calcularRating(x, Max, Min float32) float32 { //calcula el rating, aplica la formula
	var R float32
	R = (((x - Min) / (Max - Min)) * 100.0)
	R = float32(math.Round(float64(R*100)) / 100)
	return R
}
func ordenar(datos []LRA) []LRA { //rutina que sirver para ordenar de mayor a menor
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
	return datos //retorna la lista ya ordenada
}
func grafico(datos []LRA, nombreGrafico string) { //rutina que se encarga de recibir una lista de datos de tipo LRA y un nombre archivo
	lenguajes := make([]string, 0)   //lista de string donde van a estar los nombres de los lenguajes
	items := make([]opts.BarData, 0) //lista de datos donde van a estar las apariciones de los lenguajes
	lenguajes = append(lenguajes, "")
	items = append(items, opts.BarData{Value: 0})

	for i := 0; i < len(datos); i++ { //se pasa a la lista de tipo LRA a la lista de la libreria del graficador
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
	f, _ := os.Create(nombreGrafico) //se crea el archivo html
	bar.Render(f)                    //se renderiza grafico en el html
	println("\nAbriendo grafico...")
	//se llama al comando del sistema para que habra el archivo html en el navegador predeterminado
	if salida, err := exec.Command("powershell", "C:\\Users\\Luis\\go\\src\\TP-Estructura.\\"+nombreGrafico).CombinedOutput(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", salida)
	}
}
func timeSub(t1, t2 time.Time) int { //rutina que sirve para hallar la diferencia entre dos fechas
	t1 = t1.UTC().Truncate(24 * time.Hour)
	t2 = t2.UTC().Truncate(24 * time.Hour)
	return int(t1.Sub(t2).Hours() / 24)
}
func Tema1() { //tema 1
	println("Iniciando Tema 1....")
	file, errors := os.Create("Resultados.csv") //se crea un tipo de dato archivo
	check(errors)                               //se verifica errores
	writer := csv.NewWriter(file)               //el tipo archivo pasa a ser de tipo archivo csv
	//listas de nombres y nombres en la url
	tiobe_github := [20]string{"python", "c", "java", "cpp", "csharp", "visual-basic", "javascript", "assembly", "sql", "php", "r", "delphi", "go", "swift", "ruby", "visual-basic-6", "objective-c", "perl", "lua", "matlab"}
	tiobe_lenguajes := [20]string{"python", "C", "Java", "C++", "C#", "Visual Basic", "JavaScript", "Assembly language", "SQL", "PHP", "R", "Delphi/Object Pascal", "Go", "Swift", "Ruby", "Classic Visual Basic", "Objective-C", "Perl", "Lua", "MATLAB"}
	datos := []LRA{}              //se crea una lista de tipo LRA
	mayor := 0                    //variable que sirver para hallar el mayor y se inicializa
	menor := math.MaxInt          //variable que sirve para hallar el menor y se inicializa
	var cadena string             //cadena que va tener string de lo que extrae del bloque del datos que nos interesa
	for j := range tiobe_github { //se recorre las url
		time.Sleep(2 * time.Second)               //deley de 2 segundo para evitar bloqueos
		direccion := "https://github.com/topics/" //url principal
		url := direccion + tiobe_github[j]        //se concatena los nombres de los lenguaje con la url principal
		println("Visiting " + url)
		response, errors := http.Get(url) //se obtiene el resultado de la peticion
		defer response.Body.Close()
		check(errors)                  //se verifica errores
		if response.StatusCode > 400 { //se verifica si la pagina no tiro un codigo de error
			fmt.Print("\nStatus code:", response.StatusCode)
		}
		doc, errors := goquery.NewDocumentFromReader(response.Body)                       //se pasa a un tipo de dato go query document
		check(errors)                                                                     //se verifica error
		doc.Find("div.application-main ").Each(func(index int, item *goquery.Selection) { //se selecciona u div mas superior y se intera anidadamente buscando lo que se necesita
			h2 := item.Find("h2")                 //se busca el h2 del html porque es ahi donde esta el dato que buscamos
			texto := strings.TrimSpace(h2.Text()) //se pasa a formato texto el dato encontrado
			for i := range texto {                //se intera en el texto
				x := rune(texto[i])      //se castea a numero
				if unicode.IsNumber(x) { //se pregunta si es una representacion de un numero si los se pasa a string y se va concatenando
					a := string(x)
					cadena = cadena + a
				}
			}

		})
		posts := []string{tiobe_lenguajes[j], cadena} //posts es lo que se va a escribir en el arhivo, (lenguaje, aparicion)
		writer.Write(posts)                           //se prepara para escribir el dato
		apa, errors := strconv.Atoi(cadena)           //casteamos a tipo de dato numero para poder hallar el mayor y el menor
		cadena = ""                                   //se setea a vacio la cadena para el siguiente lenguaje
		check(errors)                                 //se chequea error
		if apa > mayor {                              //se busca la mayor aparicion
			mayor = apa
		}
		if apa < menor { //se busca la menor aparicion
			menor = apa
		}
		datos = append(datos, LRA{lenguaje: tiobe_lenguajes[j], rating: 0, aparicion: apa}) //se guarda el lenguaje y su aparcion en la lista de tipo LRA
		writer.Flush()                                                                      //se escribe el archivo
		//fmt.Print(tiobe_lenguajes[j], "\t\t", cadena, "\n")
	} //fin de la extraccion
	for i := 0; i < len(datos); i++ { //se calucla rating
		datos[i].rating = calcularRating(float32(datos[i].aparicion), float32(mayor), float32(menor))
	}
	var auxiliar LRA                    //auxiliar para ordenar
	for n := 1; n < len(datos)-1; n++ { //se ordena por rating
		for j := 0; j < len(datos)-n; j++ {
			if datos[j].rating < datos[j+1].rating {
				auxiliar = datos[j]
				datos[j] = datos[j+1]
				datos[j+1] = auxiliar
			}
		}
	}
	//se muestra los datos
	fmt.Println("\n\tLenguajes\t\tRating\t\tNºAparicion")
	for i := range datos {
		fmt.Printf("%20s\t\t%.2f\t\t%d\n", datos[i].lenguaje, datos[i].rating, datos[i].aparicion)
	}
	grafico(datos[:10], "Grafico1.html") //se llama a la rutina grafico y se le pasa 10 datos mayores
}

func Tema2() {
	println("Iniciando Tema 2....")
	file, errors := os.Create("TopicDeInteres.csv")               //se tipo de dato archivo
	check(errors)                                                 //se verifica errores
	writer := csv.NewWriter(file)                                 //se crea un tipo de datos archivo csv
	topicSeleccionado := "https://github.com/topics/golang?page=" //url del topic seleccionado
	var i = 0                                                     //se declara variable i sirve para contar hasta cuantas paginas se debe ir
	var TopicApariciones = make(map[string]int)                   //se inicializa lista de tipo map
	c := colly.NewCollector()                                     //se crea un tipo de dato colly colector

	c.OnHTML("article.border.rounded.color-shadow-small.color-bg-subtle.my-4", func(e *colly.HTMLElement) { //se selecciona el articulo del html y se va interando
		horaActualizacion := strings.TrimSpace(e.ChildAttr("relative-time.no-wrap", "datetime")) //se extrae la fecha del atributo datetime del html
		Topic := strings.TrimSpace(e.ChildText("a.topic-tag.topic-tag-link.f6.mb-2"))            //se extrae los topic asociado
		listaTopic := strings.Fields(Topic)                                                      //se hace un split de los topic, asi tambien se elimina los espacion y salto de lineas
		FechaExtraido := horaActualizacion[:10]                                                  //fecha extraida
		fecha, error := time.Parse("2006-01-02", FechaExtraido)                                  //parseo de fecha time.time
		check(error)
		diaDiferencia := timeSub(time.Now(), fecha) //se realiza la resta de la fecha
		if listaTopic != nil {                      //como hay elementos que no contienen la lista de topic asociado entonces si verifica que no sea nulo
			if diaDiferencia <= 30 { //se verifica que sea menor a 30 dias
				for _, s := range listaTopic { //se recorre los topics extraidos
					ap := TopicApariciones[strings.TrimSpace(s)]      //como TopicApariciones es un map, se verifica si ese topic ya existe y tenga apariciones y se extrae
					TopicApariciones[strings.TrimSpace(s)] = (ap + 1) //si topic ya existe se le suma +1 y si no esta entonces ap tiene cero y si le suma 1 porque sera su aparicion
				}
			}
		}
		i++ //se incrementa la cantidad de veces que ya se paso a la pagina siguiente

		if i < 35 { //el limite es ir hasta la pagina 34
			time.Sleep(2 * time.Second)                  //delay de 2 segundos
			c.Visit(topicSeleccionado + strconv.Itoa(i)) //se visita la siguiente pagina, y como una recursion
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Recopilando datos de: ", r.URL.String()) //cada ves que se llama al metodo visit se imprime el mensaje y la url
	})
	c.Visit(topicSeleccionado + strconv.Itoa(0)) //se llama la primera ves a la pagina

	TopicsOrdenados := []LRA{}           //lista de topic para ordenar
	for c, v := range TopicApariciones { //se intera el la lista map de forma de clave, valor
		TopicsOrdenados = append(TopicsOrdenados, LRA{lenguaje: c, rating: 0, aparicion: v}) //se cargan a la lista de tipo LRA para ordenar despues
		posts := []string{c, strconv.Itoa(v)}
		writer.Write(posts) //se va imprimiendo en tipo de dato archivo
	}
	writer.Flush()
	TopicsOrdenados = ordenar(TopicsOrdenados) //se ordena de mayor a menor
	println("\tTopics\t\tNºApariciones")       //se imprime
	for i := 0; i < len(TopicsOrdenados); i++ {
		fmt.Printf("\n%20s\t%d", TopicsOrdenados[i].lenguaje, TopicsOrdenados[i].aparicion)
	}
	grafico(TopicsOrdenados[:20], "Grafico2.html") //se grafica

}

func main() {

	Tema1()
	Tema2()

}
