package domain

import "fmt"

// Representación de conexión entre cuevas (arista)
type Arista struct {
	Desde        string  `json:"desde" xml:"desde"`
	Hasta        string  `json:"hasta" xml:"hasta"`
	Distancia    float64 `json:"distancia" xml:"distancia"`
	EsDirigido   bool    `json:"es_dirigido" xml:"es_dirigido"`
	EsObstruido  bool    `json:"es_obstruido" xml:"es_obstruido"`
}

// Función para crear una nueva arista
func NuevaArista(desde, hasta string, distancia float64, esDirigido bool) *Arista {
	return &Arista{
		Desde:        desde,
		Hasta:        hasta,
		Distancia:    distancia,
		EsDirigido:   esDirigido,
		EsObstruido: false,
	}
}

// Función para devolver nueva arista inversa
func (a *Arista) Reversa() *Arista {
	return &Arista{
		Desde:        a.Hasta,
		Hasta:        a.Desde,
		Distancia:    a.Distancia,
		EsDirigido:   a.EsDirigido,
		EsObstruido:  a.EsObstruido,
	}
}

// Función para formatear los datos de la arista
func (a *Arista) String() string {
	direccion := "no dirigido"
	if a.EsDirigido {
		direccion = "dirigido"
	}
	estado := "abierto"
	if a.EsObstruido {
		estado = "obstruido"
	}
	return fmt.Sprintf("Arista{%s -> %s, distancia: %.2f, %s, %s}",
		a.Desde, a.Hasta, a.Distancia, direccion, estado)
}