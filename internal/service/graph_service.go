package service

import (
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/repository"
	"strings"
)

type ServicioGrafo struct {
	grafo       *domain.Grafo
	repositorio *repository.RepositorioArchivo
}

func NuevoServicioGrafo(grafo *domain.Grafo, repositorio *repository.RepositorioArchivo) *ServicioGrafo {
	return &ServicioGrafo{grafo: grafo, repositorio: repositorio}
}

// 1a: Cargar grafo desde archivo
func (sg *ServicioGrafo) CargarGrafo(archivo string) error {
	// Detectar tipo de archivo por extensión
	if strings.HasSuffix(strings.ToLower(archivo), ".xml") {
		grafo, err := sg.repositorio.CargarXML(archivo)
		if err != nil {
			return err
		}
		*sg.grafo = *grafo
		return nil
	} else if strings.HasSuffix(strings.ToLower(archivo), ".txt") {
		grafo, err := sg.repositorio.CargarTXT(archivo)
		if err != nil {
			return err
		}
		*sg.grafo = *grafo
		return nil
	}

	// Por defecto, cargar como JSON
	grafo, err := sg.repositorio.CargarJSON(archivo)
	if err != nil {
		return err
	}
	*sg.grafo = *grafo // Actualizar el grafo existente
	return nil
}

// 1c: Cambiar tipo de grafo (dirigido/no dirigido)
func (sg *ServicioGrafo) CambiarTipoGrafo(esDirigido bool) {
	sg.grafo.EsDirigido = esDirigido
	// Reconectar aristas si es necesario
}

// 1f: Grados de vértices
func (sg *ServicioGrafo) ObtenerGradosVertices() map[string]map[string]int {
	grados := make(map[string]map[string]int)
	for id := range sg.grafo.Cuevas {
		grados[id] = map[string]int{
			"entrante": len(sg.grafo.ProximasAristas(id)),
			"saliente": len(sg.grafo.AristasSalientes(id)),
			"total":    len(sg.grafo.ProximasAristas(id)) + len(sg.grafo.AristasSalientes(id)),
		}
	}
	return grados
}

// ObtenerGrafo devuelve el grafo actual
func (sg *ServicioGrafo) ObtenerGrafo() *domain.Grafo {
	return sg.grafo
}

// GuardarGrafo guarda el grafo actual en un archivo
func (sg *ServicioGrafo) GuardarGrafo(archivo string) error {
	// Detectar tipo de archivo por extensión
	if strings.HasSuffix(strings.ToLower(archivo), ".xml") {
		return sg.repositorio.GuardarXML(sg.grafo, archivo)
	} else if strings.HasSuffix(strings.ToLower(archivo), ".txt") {
		return sg.repositorio.GuardarTXT(sg.grafo, archivo)
	}

	// Por defecto, guardar como JSON
	return sg.repositorio.GuardarJSON(sg.grafo, archivo)
}

// ActualizarGrafo actualiza el grafo actual
func (sg *ServicioGrafo) ActualizarGrafo(nuevoGrafo *domain.Grafo) {
	if nuevoGrafo != nil {
		sg.grafo = nuevoGrafo
	}
}

// CrearGrafoVacio crea un nuevo grafo vacío
func (sg *ServicioGrafo) CrearGrafoVacio(esDirigido bool) {
	sg.grafo = domain.NuevoGrafo(esDirigido)
}

// ValidarIntegridad valida la integridad del grafo
func (sg *ServicioGrafo) ValidarIntegridad() []string {
	var errores []string

	// Validar que todas las aristas tengan cuevas válidas
	for _, arista := range sg.grafo.Aristas {
		if _, existe := sg.grafo.Cuevas[arista.Desde]; !existe {
			errores = append(errores, "Cueva origen '"+arista.Desde+"' no existe")
		}
		if _, existe := sg.grafo.Cuevas[arista.Hasta]; !existe {
			errores = append(errores, "Cueva destino '"+arista.Hasta+"' no existe")
		}
		if arista.Distancia < 0 {
			errores = append(errores, "Distancia negativa en arista "+arista.Desde+" -> "+arista.Hasta)
		}
	}

	return errores
}

// ObtenerEstadisticas obtiene estadísticas del grafo
func (sg *ServicioGrafo) ObtenerEstadisticas() (*EstadisticasGrafo, error) {
	numConexiones := sg.grafo.NumeroAristas()
	if !sg.grafo.EsDirigido {
		// En grafos no dirigidos, cada conexión se almacena dos veces
		numConexiones = numConexiones / 2
	}

	stats := &EstadisticasGrafo{
		NumCuevas:     sg.grafo.NumeroCuevas(),
		NumConexiones: numConexiones,
		EsDirigido:    sg.grafo.EsDirigido,
		EsConexo:      sg.esConexo(),
		TieneCiclos:   sg.tieneCiclos(),
		Densidad:      sg.calcularDensidad(),
	}
	return stats, nil
}

// EstadisticasGrafo contiene información estadística del grafo
type EstadisticasGrafo struct {
	NumCuevas     int     `json:"num_cuevas"`
	NumConexiones int     `json:"num_conexiones"`
	EsDirigido    bool    `json:"es_dirigido"`
	EsConexo      bool    `json:"es_conexo"`
	TieneCiclos   bool    `json:"tiene_ciclos"`
	Densidad      float64 `json:"densidad"`
}

// esConexo verifica si el grafo es conexo
func (sg *ServicioGrafo) esConexo() bool {
	if sg.grafo.NumeroCuevas() == 0 {
		return true
	}

	// Obtener primera cueva
	var primeraCueva string
	for id := range sg.grafo.Cuevas {
		primeraCueva = id
		break
	}

	// DFS para ver cuántas cuevas son alcanzables
	visitadas := make(map[string]bool)
	sg.dfs(primeraCueva, visitadas)

	return len(visitadas) == sg.grafo.NumeroCuevas()
}

// dfs realiza búsqueda en profundidad
func (sg *ServicioGrafo) dfs(cueva string, visitadas map[string]bool) {
	visitadas[cueva] = true

	for _, vecino := range sg.grafo.ObtenerVecinos(cueva) {
		if !visitadas[vecino] {
			sg.dfs(vecino, visitadas)
		}
	}
}

// tieneCiclos verifica si el grafo tiene ciclos
func (sg *ServicioGrafo) tieneCiclos() bool {
	visitadas := make(map[string]bool)
	enPila := make(map[string]bool)

	for id := range sg.grafo.Cuevas {
		if !visitadas[id] {
			if sg.dfsCiclo(id, visitadas, enPila) {
				return true
			}
		}
	}

	return false
}

// dfsCiclo detecta ciclos usando DFS
func (sg *ServicioGrafo) dfsCiclo(cueva string, visitadas, enPila map[string]bool) bool {
	visitadas[cueva] = true
	enPila[cueva] = true

	for _, vecino := range sg.grafo.ObtenerVecinos(cueva) {
		if !visitadas[vecino] {
			if sg.dfsCiclo(vecino, visitadas, enPila) {
				return true
			}
		} else if enPila[vecino] {
			return true
		}
	}

	enPila[cueva] = false
	return false
}

// calcularDensidad calcula la densidad del grafo
func (sg *ServicioGrafo) calcularDensidad() float64 {
	numCuevas := sg.grafo.NumeroCuevas()
	if numCuevas < 2 {
		return 0.0
	}

	numAristas := sg.grafo.NumeroAristas()
	var maxAristas int

	if sg.grafo.EsDirigido {
		maxAristas = numCuevas * (numCuevas - 1)
	} else {
		maxAristas = numCuevas * (numCuevas - 1) / 2
	}

	if maxAristas == 0 {
		return 0.0
	}

	return float64(numAristas) / float64(maxAristas)
}
