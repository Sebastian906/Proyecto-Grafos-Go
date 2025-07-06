package service

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
)

// TipoRecorrido define los tipos de recorrido disponibles
type TipoRecorrido string

const (
	DFS TipoRecorrido = "DFS" // Búsqueda en Profundidad
	BFS TipoRecorrido = "BFS" // Búsqueda en Anchura
)

// RecorridoResultado representa el resultado de un recorrido
type RecorridoResultado struct {
	TipoRecorrido  TipoRecorrido `json:"tipo_recorrido"`
	CuevasVisitas  []string      `json:"cuevas_visitadas"`
	OrdenVisita    []int         `json:"orden_visita"`
	DistanciaTotal float64       `json:"distancia_total"`
	CuevaOrigen    string        `json:"cueva_origen"`
	Completado     bool          `json:"completado"`
}

// TraversalService proporciona algoritmos de recorrido de grafos
type TraversalService struct {
	graphService *ServicioGrafo
}

// NuevoTraversalService crea una nueva instancia del servicio de recorrido
func NuevoTraversalService(graphService *ServicioGrafo) *TraversalService {
	return &TraversalService{
		graphService: graphService,
	}
}

// RealizarRecorridoDFS implementa búsqueda en profundidad
func (ts *TraversalService) RealizarRecorridoDFS(grafo *domain.Grafo, cuevaOrigen string) (*RecorridoResultado, error) {
	if grafo == nil {
		return nil, fmt.Errorf("el grafo no puede ser nulo")
	}

	if _, existe := grafo.ObtenerCueva(cuevaOrigen); !existe {
		return nil, fmt.Errorf("cueva origen '%s' no existe en el grafo", cuevaOrigen)
	}

	visitadas := make(map[string]bool)
	cuevasVisitadas := make([]string, 0)
	ordenVisita := make([]int, 0)
	distanciaTotal := 0.0

	// Realizar DFS recursivo
	ts.dfsRecursivo(grafo, cuevaOrigen, visitadas, &cuevasVisitadas, &ordenVisita, &distanciaTotal)

	// Verificar si se visitaron todas las cuevas
	todasVisitadas := len(cuevasVisitadas) == len(grafo.Cuevas)

	return &RecorridoResultado{
		TipoRecorrido:  DFS,
		CuevasVisitas:  cuevasVisitadas,
		OrdenVisita:    ordenVisita,
		DistanciaTotal: distanciaTotal,
		CuevaOrigen:    cuevaOrigen,
		Completado:     todasVisitadas,
	}, nil
}

// dfsRecursivo realiza DFS de forma recursiva
func (ts *TraversalService) dfsRecursivo(grafo *domain.Grafo, cuevaActual string, visitadas map[string]bool,
	cuevasVisitadas *[]string, ordenVisita *[]int, distanciaTotal *float64) {

	visitadas[cuevaActual] = true
	*cuevasVisitadas = append(*cuevasVisitadas, cuevaActual)
	*ordenVisita = append(*ordenVisita, len(*cuevasVisitadas))

	// Obtener vecinos de la cueva actual
	vecinos := ts.obtenerVecinos(grafo, cuevaActual)

	for _, vecino := range vecinos {
		if !visitadas[vecino.CuevaDestino] && !vecino.EsObstruido {
			*distanciaTotal += vecino.Distancia
			ts.dfsRecursivo(grafo, vecino.CuevaDestino, visitadas, cuevasVisitadas, ordenVisita, distanciaTotal)
		}
	}
}

// RealizarRecorridoBFS implementa búsqueda en anchura
func (ts *TraversalService) RealizarRecorridoBFS(grafo *domain.Grafo, cuevaOrigen string) (*RecorridoResultado, error) {
	if grafo == nil {
		return nil, fmt.Errorf("el grafo no puede ser nulo")
	}

	if _, existe := grafo.ObtenerCueva(cuevaOrigen); !existe {
		return nil, fmt.Errorf("cueva origen '%s' no existe en el grafo", cuevaOrigen)
	}

	visitadas := make(map[string]bool)
	cuevasVisitadas := make([]string, 0)
	ordenVisita := make([]int, 0)
	distanciaTotal := 0.0
	distancias := make(map[string]float64) // Para rastrear distancias acumuladas

	// Cola para BFS
	cola := []string{cuevaOrigen}
	visitadas[cuevaOrigen] = true
	distancias[cuevaOrigen] = 0.0

	orden := 1
	for len(cola) > 0 {
		cuevaActual := cola[0]
		cola = cola[1:] // Remover primer elemento

		cuevasVisitadas = append(cuevasVisitadas, cuevaActual)
		ordenVisita = append(ordenVisita, orden)
		orden++

		// Obtener vecinos de la cueva actual
		vecinos := ts.obtenerVecinos(grafo, cuevaActual)

		for _, vecino := range vecinos {
			if !visitadas[vecino.CuevaDestino] && !vecino.EsObstruido {
				visitadas[vecino.CuevaDestino] = true
				distancias[vecino.CuevaDestino] = distancias[cuevaActual] + vecino.Distancia
				cola = append(cola, vecino.CuevaDestino)
			}
		}
	}

	// Calcular distancia total basada en el recorrido realizado
	for i := 1; i < len(cuevasVisitadas); i++ {
		distanciaArista := ts.obtenerDistanciaEntreAristas(grafo, cuevasVisitadas[i-1], cuevasVisitadas[i])
		distanciaTotal += distanciaArista
	}

	// Verificar si se visitaron todas las cuevas
	todasVisitadas := len(cuevasVisitadas) == len(grafo.Cuevas)

	return &RecorridoResultado{
		TipoRecorrido:  BFS,
		CuevasVisitas:  cuevasVisitadas,
		OrdenVisita:    ordenVisita,
		DistanciaTotal: distanciaTotal,
		CuevaOrigen:    cuevaOrigen,
		Completado:     todasVisitadas,
	}, nil
}

// Vecino representa una cueva vecina con su distancia
type Vecino struct {
	CuevaDestino string
	Distancia    float64
	EsObstruido  bool
}

// obtenerVecinos obtiene todos los vecinos de una cueva
func (ts *TraversalService) obtenerVecinos(grafo *domain.Grafo, cuevaID string) []Vecino {
	vecinos := make([]Vecino, 0)

	for _, arista := range grafo.Aristas {
		if arista.Desde == cuevaID && !arista.EsObstruido {
			vecinos = append(vecinos, Vecino{
				CuevaDestino: arista.Hasta,
				Distancia:    arista.Distancia,
				EsObstruido:  arista.EsObstruido,
			})
		}
		// Si el grafo no es dirigido, también considerar la dirección inversa
		if !grafo.EsDirigido && arista.Hasta == cuevaID && !arista.EsObstruido {
			vecinos = append(vecinos, Vecino{
				CuevaDestino: arista.Desde,
				Distancia:    arista.Distancia,
				EsObstruido:  arista.EsObstruido,
			})
		}
	}

	return vecinos
}

// obtenerDistanciaEntreAristas obtiene la distancia entre dos cuevas conectadas
func (ts *TraversalService) obtenerDistanciaEntreAristas(grafo *domain.Grafo, desde, hasta string) float64 {
	for _, arista := range grafo.Aristas {
		if arista.Desde == desde && arista.Hasta == hasta {
			return arista.Distancia
		}
		// Si el grafo no es dirigido, también buscar en dirección inversa
		if !grafo.EsDirigido && arista.Desde == hasta && arista.Hasta == desde {
			return arista.Distancia
		}
	}
	return 0.0
}

// ObtenerCuevasAccesibles obtiene todas las cuevas accesibles desde una cueva origen
func (ts *TraversalService) ObtenerCuevasAccesibles(grafo *domain.Grafo, cuevaOrigen string) ([]string, error) {
	resultado, err := ts.RealizarRecorridoDFS(grafo, cuevaOrigen)
	if err != nil {
		return nil, err
	}
	return resultado.CuevasVisitas, nil
}

// VerificarConectividad verifica si todas las cuevas son accesibles desde una cueva origen
func (ts *TraversalService) VerificarConectividad(grafo *domain.Grafo, cuevaOrigen string) (bool, []string, error) {
	cuevasAccesibles, err := ts.ObtenerCuevasAccesibles(grafo, cuevaOrigen)
	if err != nil {
		return false, nil, err
	}

	// Obtener cuevas no accesibles
	cuevasNoAccesibles := make([]string, 0)
	for cuevaID := range grafo.Cuevas {
		accesible := false
		for _, cuevaAccesible := range cuevasAccesibles {
			if cuevaID == cuevaAccesible {
				accesible = true
				break
			}
		}
		if !accesible {
			cuevasNoAccesibles = append(cuevasNoAccesibles, cuevaID)
		}
	}

	esConectado := len(cuevasNoAccesibles) == 0
	return esConectado, cuevasNoAccesibles, nil
}
