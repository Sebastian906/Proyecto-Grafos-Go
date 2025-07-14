package algorithms

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
)

// BFS realiza una búsqueda en anchura desde un nodo específico
func BFS(grafo *domain.Grafo, nodoInicio string) ([]string, error) {
	if grafo == nil {
		return nil, fmt.Errorf("grafo no puede ser nil")
	}

	if _, existe := grafo.Cuevas[nodoInicio]; !existe {
		return nil, fmt.Errorf("nodo de inicio '%s' no existe en el grafo", nodoInicio)
	}

	visitados := make(map[string]bool)
	var resultado []string
	cola := []string{nodoInicio}

	visitados[nodoInicio] = true

	for len(cola) > 0 {
		nodoActual := cola[0]
		cola = cola[1:]
		resultado = append(resultado, nodoActual)

		// Obtener vecinos del nodo actual
		for _, arista := range grafo.Aristas {
			if arista.Desde == nodoActual && !arista.EsObstruido {
				if !visitados[arista.Hasta] {
					visitados[arista.Hasta] = true
					cola = append(cola, arista.Hasta)
				}
			}
			// Para grafos no dirigidos, también considerar la dirección opuesta
			if !grafo.EsDirigido && arista.Hasta == nodoActual && !arista.EsObstruido {
				if !visitados[arista.Desde] {
					visitados[arista.Desde] = true
					cola = append(cola, arista.Desde)
				}
			}
		}
	}

	return resultado, nil
}

// BFSConDistancias realiza BFS retornando también las distancias
func BFSConDistancias(grafo *domain.Grafo, nodoInicio string) ([]string, map[string]float64, error) {
	if grafo == nil {
		return nil, nil, fmt.Errorf("grafo no puede ser nil")
	}

	if _, existe := grafo.Cuevas[nodoInicio]; !existe {
		return nil, nil, fmt.Errorf("nodo de inicio '%s' no existe en el grafo", nodoInicio)
	}

	visitados := make(map[string]bool)
	distancias := make(map[string]float64)
	var resultado []string
	cola := []string{nodoInicio}

	visitados[nodoInicio] = true
	distancias[nodoInicio] = 0

	for len(cola) > 0 {
		nodoActual := cola[0]
		cola = cola[1:]
		resultado = append(resultado, nodoActual)

		// Obtener vecinos del nodo actual
		for _, arista := range grafo.Aristas {
			if arista.Desde == nodoActual && !arista.EsObstruido {
				if !visitados[arista.Hasta] {
					visitados[arista.Hasta] = true
					distancias[arista.Hasta] = distancias[nodoActual] + arista.Distancia
					cola = append(cola, arista.Hasta)
				}
			}
			// Para grafos no dirigidos, también considerar la dirección opuesta
			if !grafo.EsDirigido && arista.Hasta == nodoActual && !arista.EsObstruido {
				if !visitados[arista.Desde] {
					visitados[arista.Desde] = true
					distancias[arista.Desde] = distancias[nodoActual] + arista.Distancia
					cola = append(cola, arista.Desde)
				}
			}
		}
	}

	return resultado, distancias, nil
}

// BFSNiveles realiza BFS organizando el resultado por niveles
func BFSNiveles(grafo *domain.Grafo, nodoInicio string) ([][]string, error) {
	if grafo == nil {
		return nil, fmt.Errorf("grafo no puede ser nil")
	}

	if _, existe := grafo.Cuevas[nodoInicio]; !existe {
		return nil, fmt.Errorf("nodo de inicio '%s' no existe en el grafo", nodoInicio)
	}

	visitados := make(map[string]bool)
	var resultado [][]string
	nivelActual := []string{nodoInicio}
	visitados[nodoInicio] = true

	for len(nivelActual) > 0 {
		resultado = append(resultado, nivelActual)
		var siguienteNivel []string

		for _, nodo := range nivelActual {
			// Obtener vecinos del nodo actual
			for _, arista := range grafo.Aristas {
				if arista.Desde == nodo && !arista.EsObstruido {
					if !visitados[arista.Hasta] {
						visitados[arista.Hasta] = true
						siguienteNivel = append(siguienteNivel, arista.Hasta)
					}
				}
				// Para grafos no dirigidos, también considerar la dirección opuesta
				if !grafo.EsDirigido && arista.Hasta == nodo && !arista.EsObstruido {
					if !visitados[arista.Desde] {
						visitados[arista.Desde] = true
						siguienteNivel = append(siguienteNivel, arista.Desde)
					}
				}
			}
		}

		nivelActual = siguienteNivel
	}

	return resultado, nil
}
