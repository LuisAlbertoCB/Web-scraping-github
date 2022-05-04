package main

import (
	"encoding/csv"
	"fmt"
	goquery "github.com/PuerkitoBio/goquery"
	"net/http"
	"os"
	"strings"
	"unicode"
)

func check(err error) {
	if err != nil {
		fmt.Print(err)
	}
}
func writeFile(data, filename string) {
	file, errors := os.Create(filename)
	defer file.Close()
	check(errors)
	file.WriteString(data)
}
func main() {
	file, errors := os.Create("Resultados.csv")
	check(errors)
	writer := csv.NewWriter(file)
	tiobe_github := [20]string{"python", "c", "java", "cpp", "csharp", "visual-basic", "javascript", "assembly", "sql", "php", "r", "delphi", "go", "swift", "ruby", "visual-basic-6", "objective-c", "perl", "lua", "matlab"}
	tiobe_lenguajes := [20]string{"python", "C", "Java", "C++", "C#", "Visual Basic", "JavaScript", "Assembly language", "SQL", "PHP", "R", "Delphi/Object Pascal", "Go", "Swift", "Ruby", "Classic Visual Basic", "Objective-C", "Perl", "Lua", "MATLAB"}
	type LRA struct {
		lenguaje  string
		rating    int
		aparicion int
	}
	//datos := [20]LRA{}

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
		//fmt.Print(tiobe_lenguajes[j], "\t\t", cadena, "\n")
		//intVar, errors := strconv.Atoi(cadena)
		check(errors)
		//fmt.Print(cadena, "\n")
		posts := []string{tiobe_lenguajes[j], cadena}
		writer.Write(posts)

		writer.Flush()
	}
}
