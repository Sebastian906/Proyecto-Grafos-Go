package domain

import (
	"fmt"
	"sync"
)

// Representación de grafo en el sistema
type Grafo struct {
	Cuevas     map[string]*Cueva `json:"cuevas" xml:"cuevas"`
	Aristas    []*Arista         `json:"aristas" xml:"aristas"`
	EsDirigido bool              `json:"es_dirigido" xml:"es_dirigido"`
	mu         sync.RWMutex      // Operaciones concurrentes
}

// Función para crear un nuevo grafo
func NuevoGrafo(esDirigido bool) *Grafo {
	return &Grafo{
		Cuevas:     make(map[string]*Cueva),
		Aristas:    make([]*Arista, 0),
		EsDirigido: esDirigido,
	}
}

// Función para agregar una nueva cueva al grafo
func (g *Grafo) AgregarCueva(cueva *Cueva) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, existe := g.Cuevas[cueva.ID]; existe {
		return fmt.Errorf("cueva con ID %s ya existe", cueva.ID)
	}

	g.Cuevas[cueva.ID] = cueva
	return nil
}

// Obtener cueva por identificador
func (g *Grafo) ObtenerCueva(id string) (*Cueva, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	cueva, existe := g.Cuevas[id]
	return cueva, existe
}

// Agregar una nueva arista al grafo
func (g *Grafo) AgregarArista(arista *Arista) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	// Verificar que las cuevas existan
	if _, existe := g.Cuevas[arista.Desde]; !existe {
		return fmt.Errorf("cueva %s no existe", arista.Desde)
	}
	if _, existe := g.Cuevas[arista.Hasta]; !existe {
		return fmt.Errorf("cueva %s no existe", arista.Hasta)
	}
	// Verificar que la arista aún no exista
	for _, a := range g.Aristas {
		if a.Desde == arista.Desde && a.Hasta == arista.Hasta {
			return fmt.Errorf("arista desde %s hasta %s ya existe", arista.Desde, arista.Hasta)
		}
	}
	g.Aristas = append(g.Aristas, arista)
	// Si el grafo no es dirigido, agregar la arista inversa
	if !g.EsDirigido && !arista.EsDirigido {
		aristaInversa := arista.Reversa()
		g.Aristas = append(g.Aristas, aristaInversa)
	}
	return nil
}

// Función para obtener todas las aristas del grafo
func (g *Grafo) ObtenerAristas() []*Arista {
	g.mu.RLock()
	defer g.mu.RUnlock()
	aristas := make([]*Arista, len(g.Aristas))
	copy(aristas, g.Aristas)
	return aristas
}

// Función para obtener los vecinos de una cueva
func (g *Grafo) ObtenerVecinos(caveID string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var vecinos []string
	for _, arista := range g.Aristas {
		if arista.Desde == caveID && !arista.EsObstruido {
			vecinos = append(vecinos, arista.Hasta)
		}
	}
	return vecinos
}

// Obtener aristas entrantes hacia una cueva
func (g *Grafo) ProximasAristas(caveID string) []*Arista {
	g.mu.RLock()
    defer g.mu.RUnlock()

	var proximo []*Arista
	for _, arista := range g.Aristas {
		if arista.Hasta == caveID && !arista.EsObstruido {
			proximo = append(proximo, arista)
		}
	}
	return proximo
}

// Obtener aristas salientes de una cueva
func (g *Grafo) AristasSalientes(caveID string) []*Arista {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var saliente []*Arista
	for _, arista := range g.Aristas {
		if arista.Desde == caveID && !arista.EsObstruido {
			saliente = append(saliente, arista)
		}
	}
	return saliente
}

// Función para obtener número de cuevas
func (g *Grafo) NumeroCuevas() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return len(g.Cuevas)
}

// Función para obtener número de aristas
func (g *Grafo) NumeroAristas() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return len(g.Aristas)
}

// Función para formatear los datos del grafo
func (g *Grafo) String() string {
	return fmt.Sprintf("Grafo{Cuevas: %d, Aristas: %d, EsDirigido: %t}",
		len(g.Cuevas), len(g.Aristas), g.EsDirigido)
}