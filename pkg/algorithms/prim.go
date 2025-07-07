package algorithms

import (
	"container/heap"
	"fmt"
	"proyecto-grafos-go/internal/domain"
)

// AristaPrim representa una arista para el algoritmo de Prim con prioridad
type AristaPrim struct {
	Desde     string
	Hasta     string
	Distancia float64
	Index     int // Para el heap
}

// PriorityQueue implementa heap.Interface para aristas ordenadas por distancia
type PriorityQueue []*AristaPrim

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Distancia < pq[j].Distancia
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

// Push agrega un elemento al heap
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*AristaPrim)
	item.Index = n
	*pq = append(*pq, item)
}

// Pop remueve y retorna el elemento de menor prioridad del heap
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

// MSTDesdeCueva contiene el resultado del MST desde una cueva específica
type MSTDesdeCueva struct {
	MST          *MST     `json:"mst"`
	CuevaOrigen  string   `json:"cueva_origen"`
	Alcanzables  []string `json:"alcanzables"`
	NoAlcanzable []string `json:"no_alcanzable"`
	EsCompleto   bool     `json:"es_completo"`
}

// Prim implementa el algoritmo de Prim para encontrar MST desde una cueva específica
func Prim(grafo *domain.Grafo, cuevaInicio string) (*MSTDesdeCueva, error) {
	if grafo == nil {
		return nil, fmt.Errorf("el grafo no puede ser nil")
	}

	if _, existe := grafo.Cuevas[cuevaInicio]; !existe {
		return nil, fmt.Errorf("la cueva de inicio '%s' no existe en el grafo", cuevaInicio)
	}

	// Inicializar estructuras
	visitado := make(map[string]bool)
	mstAristas := make([]*domain.Arista, 0)
	pesoTotal := 0.0

	// Priority queue para aristas candidatas
	pq := &PriorityQueue{}
	heap.Init(pq)

	// Marcar cueva inicial como visitada
	visitado[cuevaInicio] = true
	numVisitados := 1

	// Agregar todas las aristas desde la cueva inicial
	agregarAristasDesdeCueva(grafo, cuevaInicio, visitado, pq)

	// Algoritmo de Prim
	for pq.Len() > 0 && numVisitados < len(grafo.Cuevas) {
		// Obtener la arista de menor peso
		aristaMin := heap.Pop(pq).(*AristaPrim)

		// Si el nodo destino ya fue visitado, continuar
		if visitado[aristaMin.Hasta] {
			continue
		}

		// Agregar la arista al MST
		arista := &domain.Arista{
			Desde:       aristaMin.Desde,
			Hasta:       aristaMin.Hasta,
			Distancia:   aristaMin.Distancia,
			EsDirigido:  false, // MST siempre es no dirigido
			EsObstruido: false,
		}
		mstAristas = append(mstAristas, arista)
		pesoTotal += aristaMin.Distancia

		// Marcar el nuevo nodo como visitado
		visitado[aristaMin.Hasta] = true
		numVisitados++

		// Agregar aristas desde el nuevo nodo
		agregarAristasDesdeCueva(grafo, aristaMin.Hasta, visitado, pq)
	}

	// Crear MST
	mst := &MST{
		Aristas:    mstAristas,
		PesoTotal:  pesoTotal,
		NumNodos:   numVisitados,
		NumAristas: len(mstAristas),
	}

	// Determinar nodos alcanzables y no alcanzables
	alcanzables := make([]string, 0, numVisitados)
	noAlcanzable := make([]string, 0)

	for id := range grafo.Cuevas {
		if visitado[id] {
			alcanzables = append(alcanzables, id)
		} else {
			noAlcanzable = append(noAlcanzable, id)
		}
	}

	esCompleto := numVisitados == len(grafo.Cuevas)

	return &MSTDesdeCueva{
		MST:          mst,
		CuevaOrigen:  cuevaInicio,
		Alcanzables:  alcanzables,
		NoAlcanzable: noAlcanzable,
		EsCompleto:   esCompleto,
	}, nil
}

// agregarAristasDesdeCueva agrega todas las aristas válidas desde una cueva al priority queue
func agregarAristasDesdeCueva(grafo *domain.Grafo, cueva string, visitado map[string]bool, pq *PriorityQueue) {
	for _, arista := range grafo.Aristas {
		if arista.EsObstruido {
			continue
		}

		var destino string
		valida := false

		// Verificar si la arista parte desde la cueva actual
		if arista.Desde == cueva && !visitado[arista.Hasta] {
			destino = arista.Hasta
			valida = true
		} else if arista.Hasta == cueva && !arista.EsDirigido && !visitado[arista.Desde] {
			// Para grafos no dirigidos, también considerar la dirección inversa
			destino = arista.Desde
			valida = true
		}

		if valida {
			aristaPrim := &AristaPrim{
				Desde:     cueva,
				Hasta:     destino,
				Distancia: arista.Distancia,
			}
			heap.Push(pq, aristaPrim)
		}
	}
}

// ObtenerRutasDesdeOrigen obtiene todas las rutas desde el origen hacia cada nodo alcanzable
func (mst *MSTDesdeCueva) ObtenerRutasDesdeOrigen() map[string][]string {
	rutas := make(map[string][]string)

	// Construir grafo de adyacencia del MST
	adyacencia := make(map[string][]string)
	for _, arista := range mst.MST.Aristas {
		adyacencia[arista.Desde] = append(adyacencia[arista.Desde], arista.Hasta)
		adyacencia[arista.Hasta] = append(adyacencia[arista.Hasta], arista.Desde)
	}

	// Para cada nodo alcanzable, encontrar la ruta desde el origen
	for _, destino := range mst.Alcanzables {
		if destino == mst.CuevaOrigen {
			rutas[destino] = []string{destino}
			continue
		}

		ruta := encontrarRutaEnMST(adyacencia, mst.CuevaOrigen, destino)
		if len(ruta) > 0 {
			rutas[destino] = ruta
		}
	}

	return rutas
}

// encontrarRutaEnMST encuentra la ruta única entre dos nodos en un MST usando DFS
func encontrarRutaEnMST(adyacencia map[string][]string, origen, destino string) []string {
	visitado := make(map[string]bool)
	ruta := make([]string, 0)

	if dfsRuta(adyacencia, origen, destino, visitado, &ruta) {
		return ruta
	}

	return nil
}

// dfsRuta realiza DFS para encontrar una ruta específica
func dfsRuta(adyacencia map[string][]string, actual, destino string, visitado map[string]bool, ruta *[]string) bool {
	*ruta = append(*ruta, actual)
	visitado[actual] = true

	if actual == destino {
		return true
	}

	for _, vecino := range adyacencia[actual] {
		if !visitado[vecino] {
			if dfsRuta(adyacencia, vecino, destino, visitado, ruta) {
				return true
			}
		}
	}

	// Backtrack
	*ruta = (*ruta)[:len(*ruta)-1]
	return false
}
