package service

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
)

type ServicioCueva struct {
	grafo *domain.Grafo
}

// Crea una nuevo servicio de cuevas
func ServicioNuevaCueva(grafo *domain.Grafo) *ServicioCueva {
	return &ServicioCueva{
		grafo: grafo,
	}
}

// Solicitud para crear una nueva cueva
type SolicitudCueva struct {
	ID     string  `json:"id"`
	Nombre string  `json:"nombre"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}

// Solicitud para conectar dos cuevas
type SolicitudConectarCuevas struct {
	DesdeCuevaId    string  `json:"desde_cueva_id"`
	HastaCuevaId    string  `json:"hasta_cueva_id"`
	Distancia       float64 `json:"distancia"`
	EsDirigido      bool    `json:"es_dirigido"`
	EsBidireccional bool    `json:"es_bidireccional"`
}

// crea una nueva cueva en el grafo
func (sc *ServicioCueva) CrearCueva(solicitud SolicitudCueva) error {
	if solicitud.ID == "" || solicitud.Nombre == "" {
		return fmt.Errorf("ID y nombre son requeridos")
	}

	if _, existe := sc.grafo.ObtenerCueva(solicitud.ID); existe {
		return fmt.Errorf("la cueva ya existe")
	}

	cueva := domain.NuevaCueva(solicitud.ID, solicitud.Nombre)
	cueva.X, cueva.Y = solicitud.X, solicitud.Y
	return sc.grafo.AgregarCueva(cueva)
}

// ConnectCaves conecta dos cuevas existentes
func (sc *ServicioCueva) Conectar(desdeID, hastaID string, distancia float64, esDirigido, esBidireccional bool) error {
	// Validaciones básicas
	if desdeID == hastaID {
		return fmt.Errorf("no se puede conectar una cueva consigo misma")
	}
	if distancia <= 0 {
		return fmt.Errorf("la distancia debe ser positiva")
	}

	// Crear conexión principal
	if err := sc.grafo.AgregarArista(
		domain.NuevaArista(desdeID, hastaID, distancia, esDirigido),
	); err != nil {
		return err
	}

	// Conexión inversa si es bidireccional
	if esBidireccional {
		return sc.grafo.AgregarArista(
			domain.NuevaArista(hastaID, desdeID, distancia, esDirigido),
		)
	}
	return nil
}

// ObtenerCueva obtiene una cueva por su ID
func (sc *ServicioCueva) ObtenerCueva(id string) (*domain.Cueva, error) {
	if cueva, existe := sc.grafo.ObtenerCueva(id); existe {
		return cueva, nil
	}
	return nil, fmt.Errorf("cueva no encontrada")
}

func (sc *ServicioCueva) ListarCuevas() []string {
	ids := make([]string, 0, len(sc.grafo.Cuevas))
	for id := range sc.grafo.Cuevas {
		ids = append(ids, id)
	}
	return ids
}