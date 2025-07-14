package service

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
)

// ServicioCueva maneja las operaciones relacionadas con cuevas
type ServicioCueva struct {
	grafo *domain.Grafo
}

// ServicioNuevaCueva crea una nuevo servicio de cuevas
func ServicioNuevaCueva(grafo *domain.Grafo) *ServicioCueva {
	return &ServicioCueva{
		grafo: grafo,
	}
}

// SolicitudCueva representa una solicitud para crear una nueva cueva
type SolicitudCueva struct {
	ID     string  `json:"id"`
	Nombre string  `json:"nombre"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}

// SolicitudConectarCuevas representa una solicitud para conectar dos cuevas
type SolicitudConectarCuevas struct {
	DesdeCuevaID    string  `json:"desde_cueva_id"`
	HastaCuevaID    string  `json:"hasta_cueva_id"`
	Distancia       float64 `json:"distancia"`
	EsDirigido      bool    `json:"es_dirigido"`
	EsBidireccional bool    `json:"es_bidireccional"`
}

// CrearCueva crea una nueva cueva en el grafo
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

// Conectar conecta dos cuevas existentes
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
func (sc *ServicioCueva) ObtenerCueva(id string) (*domain.Cueva, bool) {
	return sc.grafo.ObtenerCueva(id)
}

// ListarCuevas retorna la lista de IDs de todas las cuevas
func (sc *ServicioCueva) ListarCuevas() []string {
	ids := make([]string, 0, len(sc.grafo.Cuevas))
	for id := range sc.grafo.Cuevas {
		ids = append(ids, id)
	}
	return ids
}

// DetalleCueva contiene información detallada de una cueva
type DetalleCueva struct {
	ID            string         `json:"id"`
	Nombre        string         `json:"nombre"`
	X             float64        `json:"x"`
	Y             float64        `json:"y"`
	Recursos      map[string]int `json:"recursos"`
	NumConexiones int            `json:"num_conexiones"`
	Vecinos       []string       `json:"vecinos"`
}

// ConectarCuevas conecta dos cuevas (método simplificado)
func (sc *ServicioCueva) ConectarCuevas(desde, hasta string, distancia float64, esDirigido bool) error {
	return sc.Conectar(desde, hasta, distancia, esDirigido, false)
}

// EliminarCueva elimina una cueva del grafo
func (sc *ServicioCueva) EliminarCueva(id string) error {
	if _, existe := sc.grafo.ObtenerCueva(id); !existe {
		return fmt.Errorf("cueva no encontrada")
	}

	return sc.grafo.EliminarCueva(id)
}

// ModificarCueva modifica los datos de una cueva existente
func (sc *ServicioCueva) ModificarCueva(id string, solicitud SolicitudCueva) error {
	cueva, existe := sc.grafo.ObtenerCueva(id)
	if !existe {
		return fmt.Errorf("cueva no encontrada")
	}

	// Actualizar campos
	if solicitud.Nombre != "" {
		cueva.Nombre = solicitud.Nombre
	}
	cueva.X = solicitud.X
	cueva.Y = solicitud.Y

	return nil
}

// ObtenerDetalleCueva obtiene información detallada de una cueva
func (sc *ServicioCueva) ObtenerDetalleCueva(id string) (*DetalleCueva, error) {
	cueva, existe := sc.grafo.ObtenerCueva(id)
	if !existe {
		return nil, fmt.Errorf("cueva no encontrada")
	}

	detalle := &DetalleCueva{
		ID:       cueva.ID,
		Nombre:   cueva.Nombre,
		X:        cueva.X,
		Y:        cueva.Y,
		Recursos: make(map[string]int),
		Vecinos:  sc.ObtenerVecinos(id),
	}

	// Copiar recursos
	for recurso, cantidad := range cueva.Recursos {
		detalle.Recursos[recurso] = cantidad
	}

	detalle.NumConexiones = len(detalle.Vecinos)

	return detalle, nil
}

// AgregarRecurso agrega recursos a una cueva
func (sc *ServicioCueva) AgregarRecurso(idCueva, recurso string, cantidad int) error {
	cueva, existe := sc.grafo.ObtenerCueva(idCueva)
	if !existe {
		return fmt.Errorf("cueva no encontrada")
	}

	if cueva.Recursos == nil {
		cueva.Recursos = make(map[string]int)
	}

	cueva.Recursos[recurso] += cantidad
	return nil
}

// RemoverRecurso remueve recursos de una cueva
func (sc *ServicioCueva) RemoverRecurso(idCueva, recurso string, cantidad int) error {
	cueva, existe := sc.grafo.ObtenerCueva(idCueva)
	if !existe {
		return fmt.Errorf("cueva no encontrada")
	}

	if cueva.Recursos == nil {
		return fmt.Errorf("la cueva no tiene recursos")
	}

	cantidadActual, existe := cueva.Recursos[recurso]
	if !existe {
		return fmt.Errorf("recurso no encontrado en la cueva")
	}

	if cantidadActual < cantidad {
		return fmt.Errorf("cantidad insuficiente de recurso")
	}

	cueva.Recursos[recurso] -= cantidad
	if cueva.Recursos[recurso] == 0 {
		delete(cueva.Recursos, recurso)
	}

	return nil
}

// ExisteConexion verifica si existe conexión entre dos cuevas
func (sc *ServicioCueva) ExisteConexion(desde, hasta string) (bool, error) {
	_, existe1 := sc.grafo.ObtenerCueva(desde)
	_, existe2 := sc.grafo.ObtenerCueva(hasta)

	if !existe1 || !existe2 {
		return false, fmt.Errorf("una o ambas cuevas no existen")
	}

	for _, arista := range sc.grafo.Aristas {
		if arista.Desde == desde && arista.Hasta == hasta {
			return true, nil
		}
	}

	return false, nil
}

// ObtenerVecinos obtiene los vecinos de una cueva
func (sc *ServicioCueva) ObtenerVecinos(id string) []string {
	vecinos := make(map[string]bool)

	for _, arista := range sc.grafo.Aristas {
		if arista.EsObstruido {
			continue
		}

		if arista.Desde == id {
			vecinos[arista.Hasta] = true
		}

		if !sc.grafo.EsDirigido && arista.Hasta == id {
			vecinos[arista.Desde] = true
		}
	}

	result := make([]string, 0, len(vecinos))
	for vecino := range vecinos {
		result = append(result, vecino)
	}

	return result
}

// CalcularDistanciaDirecta calcula la distancia directa entre dos cuevas
func (sc *ServicioCueva) CalcularDistanciaDirecta(desde, hasta string) (float64, error) {
	for _, arista := range sc.grafo.Aristas {
		if arista.Desde == desde && arista.Hasta == hasta {
			return arista.Distancia, nil
		}
	}

	return 0, fmt.Errorf("no hay conexión directa entre las cuevas")
}

// EstablecerUbicacion establece la ubicación de una cueva
func (sc *ServicioCueva) EstablecerUbicacion(id string, x, y float64) error {
	cueva, existe := sc.grafo.ObtenerCueva(id)
	if !existe {
		return fmt.Errorf("cueva no encontrada")
	}

	cueva.X = x
	cueva.Y = y
	return nil
}

// ObtenerEstadisticas obtiene estadísticas de una cueva
func (sc *ServicioCueva) ObtenerEstadisticas(id string) (map[string]interface{}, error) {
	detalle, err := sc.ObtenerDetalleCueva(id)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"id":             detalle.ID,
		"nombre":         detalle.Nombre,
		"ubicacion":      fmt.Sprintf("(%.2f, %.2f)", detalle.X, detalle.Y),
		"num_conexiones": detalle.NumConexiones,
		"num_recursos":   len(detalle.Recursos),
		"vecinos":        detalle.Vecinos,
		"recursos":       detalle.Recursos,
	}

	return stats, nil
}

// ActualizarCueva actualiza una cueva existente
func (sc *ServicioCueva) ActualizarCueva(id string, solicitud SolicitudCueva) error {
	cueva, existe := sc.grafo.ObtenerCueva(id)
	if !existe {
		return fmt.Errorf("cueva no encontrada")
	}

	// Actualizar campos
	cueva.Nombre = solicitud.Nombre
	cueva.X = solicitud.X
	cueva.Y = solicitud.Y

	return nil
}
