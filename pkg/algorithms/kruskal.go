package algorithms

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"sort"
)

// UnionFind estructura para el algoritmo de Kruskal
type UnionFind struct {
	parent map[string]string
	rank   map[string]int
}

// NuevoUnionFind crea una nueva estructura Union-Find
func NuevoUnionFind() *UnionFind {
	return &UnionFind{
		parent: make(map[string]string),
		rank:   make(map[string]int),
	}
}

// MakeSet crea un conjunto para un nodo
func (uf *UnionFind) MakeSet(x string) {
	uf.parent[x] = x
	uf.rank[x] = 0
}

// Find encuentra el representante del conjunto que contiene x
func (uf *UnionFind) Find(x string) string {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x]) // Compresión de caminos
	}
	return uf.parent[x]
}

// Union une dos conjuntos
func (uf *UnionFind) Union(x, y string) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)

	if rootX == rootY {
		return false // Ya están en el mismo conjunto
	}

	// Union por rango
	if uf.rank[rootX] < uf.rank[rootY] {
		uf.parent[rootX] = rootY
	} else if uf.rank[rootX] > uf.rank[rootY] {
		uf.parent[rootY] = rootX
	} else {
		uf.parent[rootY] = rootX
		uf.rank[rootX]++
	}

	return true
}

// MST contiene el resultado del árbol de expansión mínimo
type MST struct {
	Aristas    []*domain.Arista `json:"aristas"`
	PesoTotal  float64          `json:"peso_total"`
	NumNodos   int              `json:"num_nodos"`
	NumAristas int              `json:"num_aristas"`
}

// Kruskal implementa el algoritmo de Kruskal para encontrar el MST
func Kruskal(grafo *domain.Grafo) (*MST, error) {
	if grafo == nil {
		return nil, fmt.Errorf("el grafo no puede ser nil")
	}

	// Obtener todas las aristas válidas (no obstruidas)
	aristas := make([]*domain.Arista, 0)
	for _, arista := range grafo.Aristas {
		if !arista.EsObstruido {
			aristas = append(aristas, arista)
		}
	}

	// Si es un grafo dirigido, necesitamos convertirlo a no dirigido para MST
	// Agregamos las aristas inversas si no existen
	if grafo.EsDirigido {
		aristasExtendidas := make([]*domain.Arista, 0)
		aristasMap := make(map[string]bool)

		// Primero agregamos todas las aristas existentes
		for _, arista := range aristas {
			key := arista.Desde + "->" + arista.Hasta
			if !aristasMap[key] {
				aristasExtendidas = append(aristasExtendidas, arista)
				aristasMap[key] = true
			}
		}

		// Luego agregamos las aristas inversas si no existen
		for _, arista := range aristas {
			keyInversa := arista.Hasta + "->" + arista.Desde
			if !aristasMap[keyInversa] {
				aristaInversa := &domain.Arista{
					Desde:       arista.Hasta,
					Hasta:       arista.Desde,
					Distancia:   arista.Distancia,
					EsDirigido:  false, // Para MST tratamos como no dirigido
					EsObstruido: false,
				}
				aristasExtendidas = append(aristasExtendidas, aristaInversa)
				aristasMap[keyInversa] = true
			}
		}
		aristas = aristasExtendidas
	}

	// Ordenar aristas por peso (distancia)
	sort.Slice(aristas, func(i, j int) bool {
		return aristas[i].Distancia < aristas[j].Distancia
	})

	// Inicializar Union-Find
	uf := NuevoUnionFind()
	for id := range grafo.Cuevas {
		uf.MakeSet(id)
	}

	// Aplicar algoritmo de Kruskal
	mstAristas := make([]*domain.Arista, 0)
	pesoTotal := 0.0
	numNodos := len(grafo.Cuevas)

	for _, arista := range aristas {
		if uf.Union(arista.Desde, arista.Hasta) {
			mstAristas = append(mstAristas, arista)
			pesoTotal += arista.Distancia

			// Si ya tenemos n-1 aristas, hemos terminado
			if len(mstAristas) == numNodos-1 {
				break
			}
		}
	}

	return &MST{
		Aristas:    mstAristas,
		PesoTotal:  pesoTotal,
		NumNodos:   numNodos,
		NumAristas: len(mstAristas),
	}, nil
}

// EsConexo verifica si el grafo es conexo usando Union-Find
func EsConexo(grafo *domain.Grafo) bool {
	if len(grafo.Cuevas) <= 1 {
		return true
	}

	uf := NuevoUnionFind()
	for id := range grafo.Cuevas {
		uf.MakeSet(id)
	}

	// Unir nodos conectados por aristas no obstruidas
	for _, arista := range grafo.Aristas {
		if !arista.EsObstruido {
			uf.Union(arista.Desde, arista.Hasta)
		}
	}

	// Verificar si todos los nodos están en el mismo componente
	var primerNodo string
	for id := range grafo.Cuevas {
		primerNodo = id
		break
	}

	raizPrincipal := uf.Find(primerNodo)
	for id := range grafo.Cuevas {
		if uf.Find(id) != raizPrincipal {
			return false
		}
	}

	return true
}
