package interfaz_consultor

import (
	"bufio"
	"fmt"
	"os"
	TDALista "tdas/lista"
	"time"
	Consultor "tp2/consultor_vuelos"
	Texto "tp2/manipulacion_texto"
)

const TIME_LAYOUT = "2006-01-02T15:04:05"

type interfaz struct {
	consultor Consultor.ConsultorVuelos
}

func CrearInterfaz() Interfaz {
	return &interfaz{consultor: Consultor.CrearConsultor()}
}

func (interfaz *interfaz) Iniciar() {
	comando := bufio.NewScanner(os.Stdin)
	for comando.Scan() {
		if comando.Text() == "" {
			break
		}
		argumentos := Texto.Split(comando.Text(), ' ')
		switch argumentos[0] {
		case "agregar_archivo":
			interfaz.agregar_archivo(argumentos)
		case "ver_tablero":
			interfaz.ver_tablero(argumentos)
		case "info_vuelo":
			interfaz.info_vuelo(argumentos)
		case "siguiente_vuelo":
			interfaz.siguiente_vuelo(argumentos)
		case "prioridad_vuelos":
			interfaz.prioridad_vuelos(argumentos)
		case "borrar":
			interfaz.borrar(argumentos)
		default:
			interfaz.error(argumentos[0])
		}
		fmt.Scan(&comando)
	}
}

func (interfaz *interfaz) agregar_archivo(argumentos []string) {
	ruta := argumentos[1]
	archivo, err := os.Open(ruta)
	if err != nil {
		interfaz.error(argumentos[0])
		return
	}
	defer archivo.Close()

	vuelos := TDALista.CrearListaEnlazada[[]string]()
	scan := bufio.NewScanner(archivo)
	for scan.Scan() {
		vuelos.InsertarUltimo(Texto.Split(scan.Text(), ','))
	}
	interfaz.consultor.AgregarArchivo(vuelos)
	interfaz.ok()
}

func (interfaz *interfaz) ver_tablero(argumentos []string) {
	if len(argumentos) != 5 {
		interfaz.error(argumentos[0])
		return
	}
	cantidad := Texto.StringToInt(argumentos[1])
	modo := argumentos[2]
	desde, _ := time.Parse(TIME_LAYOUT, argumentos[3])
	hasta, _ := time.Parse(TIME_LAYOUT, argumentos[4])
	if cantidad < 1 || (modo != "asc" && modo != "desc") || hasta.Before(desde) {
		interfaz.error(argumentos[0])
		return
	}
	lista := interfaz.consultor.VerTablero(cantidad, modo, desde, hasta)
	iter := lista.Iterador()
	for iter.HaySiguiente() {
		vuelo := iter.VerActual()
		fmt.Printf("%s - %s\n", vuelo[6], vuelo[0])
		iter.Siguiente()
	}
	interfaz.ok()
}

func (interfaz *interfaz) info_vuelo(argumentos []string) {
	if len(argumentos) != 2 {
		interfaz.error(argumentos[0])
		return
	}
	codigo := argumentos[1]
	vuelo := interfaz.consultor.InfoVuelo(codigo)
	Texto.PrintSlice(vuelo, ' ')
	interfaz.ok()
}

func (interfaz *interfaz) siguiente_vuelo(argumentos []string) {
}

func (interfaz *interfaz) prioridad_vuelos(argumentos []string) {
	cantidad := Texto.StringToInt(argumentos[1])
	if cantidad <= 0 {
		interfaz.error(argumentos[0])
		return
	}
	lista := interfaz.consultor.PrioridadVuelos(cantidad)
	iter := lista.Iterador()
	for iter.HaySiguiente() {
		vuelos := iter.VerActual()
		fmt.Printf("%s - %s\n", vuelos[5], vuelos[0])
		iter.Siguiente()
	}
	interfaz.ok()
}

func (interfaz *interfaz) borrar(argumentos []string) {
	if len(argumentos) != 3 {
		interfaz.error(argumentos[0])
		return
	}
	desde := argumentos[1]
	hasta := argumentos[2]
	vuelos_borrados := interfaz.consultor.Borrar(desde, hasta)
	for iter := vuelos_borrados.Iterador(); iter.HaySiguiente(); iter.Siguiente() {
		Texto.PrintSlice(iter.VerActual(), ' ')
	}
	interfaz.ok()
}

func (interfaz *interfaz) ok() {
	fmt.Println("OK")
}

func (interfaz *interfaz) error(comando string) {
	fmt.Fprintf(os.Stderr, "Error en comando %s\n", comando)
}
