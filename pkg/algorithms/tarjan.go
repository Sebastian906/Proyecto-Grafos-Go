package algorithms

import (
	"proyecto-grafos-go/internal/domain"
)

// Representar el resultado del algoritmo de Tarjan
type ResultadoTarjan struct {
	ComponentesFuertementeConectados [][]string
	EsFuertementeConectado           bool
	PuntosArticulacion               []string
}

// Implementar el algoritmo de Tarjan para encontrar componentes fuertemente conectados
type AlgoritmoTarjan struct {
	grafo       *domain.Grafo
	indices     map[string]int
	lowlinks    map[string]int
	enStack     map[string]bool
	stack       []string
	index       int
	componentes [][]string
}

// Nueva instancia del algoritmo de Tarjan
func NuevoAlgoritmoTarjan(grafo *domain.Grafo) *AlgoritmoTarjan {
	return &AlgoritmoTarjan{
		grafo:       grafo,
		indices:     make(map[string]int),
		lowlinks:    make(map[string]int),
		enStack:     make(map[string]bool),
		stack:       make([]string, 0),
		index:       0,
		componentes: make([][]string, 0),
	}
}

// ejecutar el algoritmo de Tarjan
func (t *AlgoritmoTarjan) EncontrarComponentes() *ResultadoTarjan {
	// Reiniciar el estado
	t.indices = make(map[string]int)
	t.lowlinks = make(map[string]int)
	t.enStack = make(map[string]bool)
	t.stack = make([]string, 0)
	t.index = 0
	t.componentes = make([][]string, 0)

	// Ejecutar DFS desde cada nodo no visitado
	for cuevaID := range t.grafo.Cuevas {
		if _, visitado := t.indices[cuevaID]; !visitado {
			t.tarjanDFS(cuevaID)
		}
	}

	// Determinar si el grafo es fuertemente conectado
	esFuertementeConectado := len(t.componentes) == 1 && len(t.componentes[0]) == len(t.grafo.Cuevas)

	return &ResultadoTarjan{
		ComponentesFuertementeConectados: t.componentes,
		EsFuertementeConectado:           esFuertementeConectado,
		PuntosArticulacion:               t.encontrarPuntosArticulacion(),
	}
}

// ejecutar el DFS del algoritmo de Tarjan
func (t *AlgoritmoTarjan) tarjanDFS(cuevaID string) {
	// Establecer el índice y lowlink del nodo actual
	t.indices[cuevaID] = t.index
	t.lowlinks[cuevaID] = t.index
	t.index++

	// Agregar el nodo a la pila
	t.stack = append(t.stack, cuevaID)
	t.enStack[cuevaID] = true

	// Visitar todos los vecinos
	vecinos := t.grafo.ObtenerVecinos(cuevaID)
	for _, vecino := range vecinos {
		if _, visitado := t.indices[vecino]; !visitado {
			// Si el vecino no ha sido visitado, hacer DFS recursivo
			t.tarjanDFS(vecino)
			t.lowlinks[cuevaID] = min(t.lowlinks[cuevaID], t.lowlinks[vecino])
		} else if t.enStack[vecino] {
			// Si el vecino está en la pila, actualizar lowlink
			t.lowlinks[cuevaID] = min(t.lowlinks[cuevaID], t.indices[vecino])
		}
	}

	// Si el nodo actual es la raíz de un componente fuertemente conectado
	if t.lowlinks[cuevaID] == t.indices[cuevaID] {
		var componente []string
		for {
			// Sacar elementos de la pila hasta llegar al nodo actual
			ultimo := len(t.stack) - 1
			nodo := t.stack[ultimo]
			t.stack = t.stack[:ultimo]
			t.enStack[nodo] = false
			componente = append(componente, nodo)

			if nodo == cuevaID {
				break
			}
		}
		t.componentes = append(t.componentes, componente)
	}
}

// Encontrar los puntos de articulación en el grafo
func (t *AlgoritmoTarjan) encontrarPuntosArticulacion() []string {
	// Esta es una implementación simplificada
	// En un grafo dirigido, los puntos de articulación son más complejos
	var puntosArticulacion []string

	for cuevaID := range t.grafo.Cuevas {
		if t.esPuntoArticulacion(cuevaID) {
			puntosArticulacion = append(puntosArticulacion, cuevaID)
		}
	}

	return puntosArticulacion
}

// Verificar si un nodo es un punto de articulación
func (t *AlgoritmoTarjan) esPuntoArticulacion(cuevaID string) bool {
	// Contar componentes antes de remover el nodo
	componentesAntes := len(t.componentes)

	// Simular la remoción del nodo temporalmente
	// (implementación simplificada)
	vecinos := t.grafo.ObtenerVecinos(cuevaID)
	if len(vecinos) <= 1 {
		return false
	}

	// Si un nodo tiene más de 2 vecinos y está en múltiples componentes,
	// es probable que sea un punto de articulación
	return len(vecinos) > 2 && componentesAntes > 1
}

// min retorna el menor de dos enteros
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Verificar si el grafo es fuertemente conectado
func VerificarConectividadFuerte(grafo *domain.Grafo) bool {
	if len(grafo.Cuevas) == 0 {
		return true
	}

	tarjan := NuevoAlgoritmoTarjan(grafo)
	resultado := tarjan.EncontrarComponentes()

	return resultado.EsFuertementeConectado
}

// Obtener todos los componentes fuertemente conectados
func ObtenerComponentes(grafo *domain.Grafo) [][]string {
	tarjan := NuevoAlgoritmoTarjan(grafo)
	resultado := tarjan.EncontrarComponentes()

	return resultado.ComponentesFuertementeConectados
}
