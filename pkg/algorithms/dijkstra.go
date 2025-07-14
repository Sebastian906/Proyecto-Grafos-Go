package algorithms

import (
	"fmt"
	"math"
	"proyecto-grafos-go/internal/domain"
)

// Dijkstra implementa el algoritmo de Dijkstra para encontrar el camino más corto
func Dijkstra(grafo *domain.Grafo, nodoInicio string) (map[string]float64, map[string]string, error) {
	if grafo == nil {
		return nil, nil, fmt.Errorf("grafo no puede ser nil")
	}

	if _, existe := grafo.Cuevas[nodoInicio]; !existe {
		return nil, nil, fmt.Errorf("nodo de inicio '%s' no existe en el grafo", nodoInicio)
	}

	distancias := make(map[string]float64)
	predecesores := make(map[string]string)
	visitados := make(map[string]bool)

	// Inicializar distancias con infinito
	for id := range grafo.Cuevas {
		distancias[id] = math.Inf(1)
	}
	distancias[nodoInicio] = 0

	for {
		// Encontrar el nodo no visitado con menor distancia
		nodoActual := ""
		menorDistancia := math.Inf(1)

		for id := range grafo.Cuevas {
			if !visitados[id] && distancias[id] < menorDistancia {
				menorDistancia = distancias[id]
				nodoActual = id
			}
		}

		if nodoActual == "" {
			break // No hay más nodos alcanzables
		}

		visitados[nodoActual] = true

		// Actualizar distancias de los vecinos
		for _, arista := range grafo.Aristas {
			if arista.EsObstruido {
				continue
			}

			var vecino string
			if arista.Desde == nodoActual {
				vecino = arista.Hasta
			} else if !grafo.EsDirigido && arista.Hasta == nodoActual {
				vecino = arista.Desde
			} else {
				continue
			}

			if visitados[vecino] {
				continue
			}

			nuevaDistancia := distancias[nodoActual] + arista.Distancia
			if nuevaDistancia < distancias[vecino] {
				distancias[vecino] = nuevaDistancia
				predecesores[vecino] = nodoActual
			}
		}
	}

	return distancias, predecesores, nil
}

// DijkstraRuta encuentra la ruta más corta entre dos nodos específicos
func DijkstraRuta(grafo *domain.Grafo, nodoInicio, nodoDestino string) ([]string, float64, error) {
	distancias, predecesores, err := Dijkstra(grafo, nodoInicio)
	if err != nil {
		return nil, 0, err
	}

	if _, existe := grafo.Cuevas[nodoDestino]; !existe {
		return nil, 0, fmt.Errorf("nodo de destino '%s' no existe en el grafo", nodoDestino)
	}

	if math.IsInf(distancias[nodoDestino], 1) {
		return nil, 0, fmt.Errorf("no hay ruta desde '%s' hasta '%s'", nodoInicio, nodoDestino)
	}

	// Reconstruir la ruta
	var ruta []string
	nodoActual := nodoDestino

	for nodoActual != nodoInicio {
		ruta = append([]string{nodoActual}, ruta...)
		nodoActual = predecesores[nodoActual]
	}
	ruta = append([]string{nodoInicio}, ruta...)

	return ruta, distancias[nodoDestino], nil
}

// DijkstraTodasLasRutas encuentra todas las rutas más cortas desde un nodo
func DijkstraTodasLasRutas(grafo *domain.Grafo, nodoInicio string) (map[string][]string, map[string]float64, error) {
	distancias, predecesores, err := Dijkstra(grafo, nodoInicio)
	if err != nil {
		return nil, nil, err
	}

	rutas := make(map[string][]string)

	for nodoDestino := range grafo.Cuevas {
		if nodoDestino == nodoInicio {
			rutas[nodoDestino] = []string{nodoInicio}
			continue
		}

		if math.IsInf(distancias[nodoDestino], 1) {
			continue // No hay ruta
		}

		// Reconstruir la ruta
		var ruta []string
		nodoActual := nodoDestino

		for nodoActual != nodoInicio {
			ruta = append([]string{nodoActual}, ruta...)
			nodoActual = predecesores[nodoActual]
		}
		ruta = append([]string{nodoInicio}, ruta...)
		rutas[nodoDestino] = ruta
	}

	return rutas, distancias, nil
}
