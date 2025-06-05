package manipulacion_texto

import (
	"fmt"
	"strconv"
	"strings"
)

func Split(cadena string, separador rune) []string {
	res := []string{}
	separado := ""
	for _, s := range cadena {
		if s == separador {
			res = append(res, separado)
			separado = ""
			continue
		}
		separado += string(s)
	}
	res = append(res, separado)
	return res
}

func StringToInt(cadena string) int {
	numero, _ := strconv.Atoi(cadena)
	return numero
}

func PrintSlice(cadena []string, separador rune) {
	for _, i := range cadena {
		if len(i) > 1 {
			i = strings.TrimLeft(i, "0")
		}
		fmt.Printf("%s%c", i, separador)
	}
	fmt.Print("\n")
}
