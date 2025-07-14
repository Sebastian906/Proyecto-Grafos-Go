package utils

import (
	"proyecto-grafos-go/internal/domain"
)

// EsConexo verifica si el grafo es conexo usando DFS
func EsConexo(grafo *domain.Grafo) bool {
	if grafo == nil || len(grafo.Cuevas) == 0 {
		return false
	}

	if len(grafo.Cuevas) == 1 {
		return true
	}

	// Obtener el primer vértice para empezar DFS
	var startVertex string
	for id := range grafo.Cuevas {
		startVertex = id
		break
	}

	// Hacer DFS y contar vértices alcanzables
	visited := make(map[string]bool)
	dfs(grafo, startVertex, visited)

	// El grafo es conexo si todos los vértices son alcanzables
	return len(visited) == len(grafo.Cuevas)
}

// dfs es una función auxiliar para el recorrido en profundidad
func dfs(grafo *domain.Grafo, vertex string, visited map[string]bool) {
	visited[vertex] = true

	// Visitar todos los vecinos no visitados
	for _, arista := range grafo.Aristas {
		var vecino string
		if arista.Desde == vertex && !arista.EsObstruido {
			vecino = arista.Hasta
		} else if !grafo.EsDirigido && arista.Hasta == vertex && !arista.EsObstruido {
			vecino = arista.Desde
		}

		if vecino != "" && !visited[vecino] {
			dfs(grafo, vecino, visited)
		}
	}
}

// TieneCiclos verifica si el grafo tiene ciclos
func TieneCiclos(grafo *domain.Grafo) bool {
	if grafo == nil || len(grafo.Aristas) == 0 {
		return false
	}

	visited := make(map[string]bool)

	// Para grafos no dirigidos, usar una versión diferente del algoritmo
	if !grafo.EsDirigido {
		// En un grafo no dirigido, cada arista se almacena dos veces
		aristasUnicas := len(grafo.Aristas) / 2
		// Un árbol con N vértices tiene N-1 aristas
		if aristasUnicas < len(grafo.Cuevas)-1 {
			return false // No hay suficientes aristas para formar un ciclo
		}

		for vertex := range grafo.Cuevas {
			if !visited[vertex] {
				if tieneCiclosNoDirigido(grafo, vertex, "", visited) {
					return true
				}
			}
		}
		return false
	}

	// Para grafos dirigidos
	recStack := make(map[string]bool)
	for vertex := range grafo.Cuevas {
		if !visited[vertex] {
			if tieneCiclosDFS(grafo, vertex, visited, recStack) {
				return true
			}
		}
	}

	return false
}

// tieneCiclosNoDirigido detecta ciclos en grafos no dirigidos
func tieneCiclosNoDirigido(grafo *domain.Grafo, vertex, parent string, visited map[string]bool) bool {
	visited[vertex] = true

	// Buscar todos los vértices adyacentes
	for _, arista := range grafo.Aristas {
		if arista.EsObstruido {
			continue
		}

		var vecino string
		// Solo considerar aristas que salen del vértice actual
		if arista.Desde == vertex {
			vecino = arista.Hasta
		}

		if vecino != "" && vecino != parent {
			if !visited[vecino] {
				if tieneCiclosNoDirigido(grafo, vecino, vertex, visited) {
					return true
				}
			} else {
				// Si el vecino fue visitado y no es el padre, hay un ciclo
				return true
			}
		}
	}

	return false
}

// tieneCiclosDFS es una función auxiliar para detectar ciclos
func tieneCiclosDFS(grafo *domain.Grafo, vertex string, visited, recStack map[string]bool) bool {
	visited[vertex] = true
	recStack[vertex] = true

	// Visitar todos los adyacentes
	for _, arista := range grafo.Aristas {
		if arista.Desde == vertex && !arista.EsObstruido {
			vecino := arista.Hasta

			// Si el vecino no ha sido visitado, recursivamente verificar por ciclos
			if !visited[vecino] {
				if tieneCiclosDFS(grafo, vecino, visited, recStack) {
					return true
				}
			} else if recStack[vecino] {
				// Si el vecino está en la pila de recursión, hay un ciclo
				return true
			}
		}

		// Para grafos no dirigidos, verificar en ambas direcciones
		if !grafo.EsDirigido && arista.Hasta == vertex && !arista.EsObstruido {
			vecino := arista.Desde

			if !visited[vecino] {
				if tieneCiclosDFS(grafo, vecino, visited, recStack) {
					return true
				}
			} else if recStack[vecino] {
				return true
			}
		}
	}

	recStack[vertex] = false
	return false
}

// ObtenerGrado obtiene el grado de un vértice
func ObtenerGrado(grafo *domain.Grafo, vertice string) int {
	if grafo == nil {
		return 0
	}

	grado := 0
	for _, arista := range grafo.Aristas {
		if arista.EsObstruido {
			continue
		}

		if grafo.EsDirigido {
			// Para grafos dirigidos, contar solo las aristas salientes
			if arista.Desde == vertice {
				grado++
			}
		} else {
			// Para grafos no dirigidos, contar solo las aristas que salen del vértice
			// (ya que cada conexión se almacena en ambas direcciones)
			if arista.Desde == vertice {
				grado++
			}
		}
	}
	return grado
}

// EsAciclico verifica si el grafo es acíclico
func EsAciclico(grafo *domain.Grafo) bool {
	return !TieneCiclos(grafo)
}

// ObtenerComponentesConexas obtiene las componentes conexas del grafo
func ObtenerComponentesConexas(grafo *domain.Grafo) [][]string {
	if grafo == nil || len(grafo.Cuevas) == 0 {
		return [][]string{}
	}

	visited := make(map[string]bool)
	var componentes [][]string

	// Para cada vértice no visitado, encontrar su componente
	for vertex := range grafo.Cuevas {
		if !visited[vertex] {
			var componente []string
			dfsComponente(grafo, vertex, visited, &componente)
			if len(componente) > 0 {
				componentes = append(componentes, componente)
			}
		}
	}

	return componentes
}

// dfsComponente es una función auxiliar para encontrar componentes conexas
func dfsComponente(grafo *domain.Grafo, vertex string, visited map[string]bool, componente *[]string) {
	visited[vertex] = true
	*componente = append(*componente, vertex)

	// Visitar todos los vecinos no visitados
	for _, arista := range grafo.Aristas {
		if arista.EsObstruido {
			continue
		}

		var vecino string
		if arista.Desde == vertex {
			vecino = arista.Hasta
		} else if !grafo.EsDirigido && arista.Hasta == vertex {
			vecino = arista.Desde
		}

		if vecino != "" && !visited[vecino] {
			dfsComponente(grafo, vecino, visited, componente)
		}
	}
}

// EsCompleto verifica si el grafo es completo
func EsCompleto(grafo *domain.Grafo) bool {
	if grafo == nil || len(grafo.Cuevas) < 2 {
		return false
	}

	n := len(grafo.Cuevas)
	aristasNecesarias := n * (n - 1) / 2
	if grafo.EsDirigido {
		aristasNecesarias = n * (n - 1)
	}

	return len(grafo.Aristas) >= aristasNecesarias
}
