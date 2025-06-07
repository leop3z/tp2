package consulta_vuelos

import (
	"strings"
	TDAHeap "tdas/cola_prioridad"
	TDADiccionario "tdas/diccionario"
	TDALista "tdas/lista"
	"time"
	Texto "tp2/manipulacion_texto"
)

const TIME_LAYOUT = "2006-01-02T15:04:05"

type vuelo struct {
	codigo    string
	fecha     time.Time
	datos     []string
	prioridad int
}

type consultor struct {
	vuelos TDADiccionario.Diccionario[string, *vuelo]
	orden  TDADiccionario.DiccionarioOrdenado[string, nulo]
}

type nulo = struct{}

func (consultor consultor) compararVuelos(a, b string) int {
	vueloA := consultor.vuelos.Obtener(a)
	vueloB := consultor.vuelos.Obtener(b)
	if vueloA.fecha.Before(vueloB.fecha) {
		return -1
	} else if vueloA.fecha.After(vueloB.fecha) {
		return 1
	} else {
		return strings.Compare(vueloA.codigo, vueloB.codigo)
	}
}

func (consultor consultor) compararPrioridades(a, b string) int {
	vueloA := consultor.vuelos.Obtener(a)
	vueloB := consultor.vuelos.Obtener(b)
	if vueloA.prioridad > vueloB.prioridad {
		return 1
	} else if vueloA.prioridad < vueloB.prioridad {
		return -1
	} else {
		return -strings.Compare(vueloA.codigo, vueloB.codigo)
	}
}

func crearVuelo(datos []string) *vuelo {
	vuelo := &vuelo{datos: make([]string, 10)}
	vuelo.codigo = datos[0]
	vuelo.fecha, _ = time.Parse(TIME_LAYOUT, datos[6])
	vuelo.prioridad = Texto.StringToInt(datos[5])
	copy(vuelo.datos, datos)
	return vuelo
}

func CrearConsultor() ConsultorVuelos {
	consultor := &consultor{}
	consultor.vuelos = TDADiccionario.CrearHash[string, *vuelo]()
	consultor.orden = TDADiccionario.CrearABB[string, nulo](consultor.compararVuelos)
	return consultor
}

func (consultor *consultor) AgregarArchivo(lista TDALista.Lista[[]string]) {
	iter := lista.Iterador()
	for iter.HaySiguiente() {
		vuelo := crearVuelo(iter.Borrar())
		consultor.vuelos.Guardar(vuelo.codigo, vuelo)
		consultor.orden.Guardar(vuelo.codigo, nulo{})
	}
}

func (consultor *consultor) VerTablero(cantidad int, modo string, desde time.Time, hasta time.Time) TDALista.Lista[[]string] {
	tablero := TDALista.CrearListaEnlazada[[]string]()
	contador := 0
	var temp [][]string

	iter := consultor.orden.Iterador()
	for iter.HaySiguiente() {
		codigo, _ := iter.VerActual()
		vuelo := consultor.vuelos.Obtener(codigo)
		if vuelo.fecha.After(hasta) {
			iter.Siguiente()
			continue
		}
		if vuelo.fecha.Before(desde) {
			iter.Siguiente()
			continue
		}
		temp = append(temp, vuelo.datos)
		iter.Siguiente()
	}

	if modo == "desc" {
		for i := len(temp) - 1; i >= 0 && (cantidad <= 0 || contador < cantidad); i-- {
			tablero.InsertarUltimo(temp[i])
			contador++
		}
	} else {
		for i := 0; i < len(temp) && (cantidad <= 0 || contador < cantidad); i++ {
			tablero.InsertarUltimo(temp[i])
			contador++
		}
	}
	return tablero
}

func (consultor *consultor) InfoVuelo(codigo string) []string {
	vuelo := []string{}
	if consultor.vuelos.Pertenece(codigo) {
		vuelo = consultor.vuelos.Obtener(codigo).datos
	}
	return vuelo
}

func (consultor *consultor) SiguienteVuelo(origen string, destino string, fecha string) []string {
	iter := consultor.vuelos.Iterador()
	_, res := iter.VerActual()
	return res.datos
}

func (consultor *consultor) PrioridadVuelos(cantidad int) TDALista.Lista[[]string] {
	arr := []string{}
	for iter := consultor.vuelos.Iterador(); iter.HaySiguiente(); iter.Siguiente() {
		codigo, _ := iter.VerActual()
		arr = append(arr, codigo)
	}
	prioridades := TDAHeap.CrearHeapArr(arr, consultor.compararPrioridades)
	tablero := TDALista.CrearListaEnlazada[[]string]()
	for range cantidad {
		tablero.InsertarUltimo(consultor.vuelos.Obtener(prioridades.Desencolar()).datos)
	}
	return tablero
}

func (consultor *consultor) Borrar(desde string, hasta string) TDALista.Lista[[]string] {
	iter := consultor.orden.IteradorRango(&desde, &hasta)
	claves_borrar := []string{}
	borrados := TDALista.CrearListaEnlazada[[]string]()

	for iter.HaySiguiente() {
		codigo, _ := iter.VerActual()
		borrados.InsertarUltimo(consultor.vuelos.Obtener(codigo).datos)
		claves_borrar = append(claves_borrar, codigo)
	}

	for _, codigo := range claves_borrar {
		consultor.orden.Borrar(codigo)
		consultor.vuelos.Borrar(codigo)
	}
	return borrados
}
