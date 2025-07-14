package algorithms

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
)

// DFS realiza una búsqueda en profundidad desde un nodo específico
func DFS(grafo *domain.Grafo, nodoInicio string) ([]string, error) {
	if grafo == nil {
		return nil, fmt.Errorf("grafo no puede ser nil")
	}

	if _, existe := grafo.Cuevas[nodoInicio]; !existe {
		return nil, fmt.Errorf("nodo de inicio '%s' no existe en el grafo", nodoInicio)
	}

	visitados := make(map[string]bool)
	var resultado []string

	dfsRecursivo(grafo, nodoInicio, visitados, &resultado)
	return resultado, nil
}

// dfsRecursivo es la función recursiva que realiza el DFS
func dfsRecursivo(grafo *domain.Grafo, nodoActual string, visitados map[string]bool, resultado *[]string) {
	visitados[nodoActual] = true
	*resultado = append(*resultado, nodoActual)

	// Obtener vecinos del nodo actual
	for _, arista := range grafo.Aristas {
		if arista.Desde == nodoActual && !arista.EsObstruido {
			if !visitados[arista.Hasta] {
				dfsRecursivo(grafo, arista.Hasta, visitados, resultado)
			}
		}
		// Para grafos no dirigidos, también considerar la dirección opuesta
		if !grafo.EsDirigido && arista.Hasta == nodoActual && !arista.EsObstruido {
			if !visitados[arista.Desde] {
				dfsRecursivo(grafo, arista.Desde, visitados, resultado)
			}
		}
	}
}

// DFSConDistancias realiza DFS retornando también las distancias
func DFSConDistancias(grafo *domain.Grafo, nodoInicio string) ([]string, map[string]float64, error) {
	if grafo == nil {
		return nil, nil, fmt.Errorf("grafo no puede ser nil")
	}

	if _, existe := grafo.Cuevas[nodoInicio]; !existe {
		return nil, nil, fmt.Errorf("nodo de inicio '%s' no existe en el grafo", nodoInicio)
	}

	visitados := make(map[string]bool)
	distancias := make(map[string]float64)
	var resultado []string

	distancias[nodoInicio] = 0
	dfsConDistanciasRecursivo(grafo, nodoInicio, visitados, distancias, &resultado)

	return resultado, distancias, nil
}

// dfsConDistanciasRecursivo es la función recursiva que realiza DFS con distancias
func dfsConDistanciasRecursivo(grafo *domain.Grafo, nodoActual string, visitados map[string]bool, distancias map[string]float64, resultado *[]string) {
	visitados[nodoActual] = true
	*resultado = append(*resultado, nodoActual)

	// Obtener vecinos del nodo actual
	for _, arista := range grafo.Aristas {
		if arista.Desde == nodoActual && !arista.EsObstruido {
			if !visitados[arista.Hasta] {
				distancias[arista.Hasta] = distancias[nodoActual] + arista.Distancia
				dfsConDistanciasRecursivo(grafo, arista.Hasta, visitados, distancias, resultado)
			}
		}
		// Para grafos no dirigidos, también considerar la dirección opuesta
		if !grafo.EsDirigido && arista.Hasta == nodoActual && !arista.EsObstruido {
			if !visitados[arista.Desde] {
				distancias[arista.Desde] = distancias[nodoActual] + arista.Distancia
				dfsConDistanciasRecursivo(grafo, arista.Desde, visitados, distancias, resultado)
			}
		}
	}
}
