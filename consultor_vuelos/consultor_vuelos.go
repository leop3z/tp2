package consulta_vuelos

import (
	TDALista "tdas/lista"
	"time"
)

type ConsultorVuelos interface {
	AgregarArchivo(vuelos TDALista.Lista[[]string])

	VerTablero(cantidad int, modo string, desde time.Time, hasta time.Time) TDALista.Lista[[]string]

	InfoVuelo(codigo string) []string

	PrioridadVuelos(cantidad int) TDALista.Lista[[]string]

	SiguienteVuelo(origen string, destino string, fecha string) []string

	Borrar(desde string, hasta string)
}
