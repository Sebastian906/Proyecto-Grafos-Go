package service

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/pkg/algorithms"
)

// ServicioValidacion maneja la validación y análisis de conectividad del grafo
type ServicioValidacion struct {
	grafo *domain.Grafo
}

// ResultadoAccesibilidad contiene el resultado del análisis de accesibilidad
type ResultadoAccesibilidad struct {
	CuevasInaccesibles []string
	Soluciones         []string
	TotalCuevas        int
	CuevasAccesibles   int
}

// SolucionConectividad representa una solución propuesta para mejorar la conectividad
type SolucionConectividad struct {
	Tipo        string // "agregar_arista", "cambiar_direccion", "eliminar_obstruccion"
	Descripcion string
	Desde       string
	Hasta       string
}

// NuevoServicioValidacion crea una nueva instancia del servicio de validación
func NuevoServicioValidacion(grafo *domain.Grafo) *ServicioValidacion {
	return &ServicioValidacion{grafo: grafo}
}

// EsFuertementeConectado verifica si el grafo es fuertemente conectado
func (sv *ServicioValidacion) EsFuertementeConectado() bool {
	return algorithms.VerificarConectividadFuerte(sv.grafo)
}

// DetectarPozos encuentra las cuevas que son pozos (sin conexiones salientes)
func (sv *ServicioValidacion) DetectarPozos() []string {
	var pozos []string
	for id := range sv.grafo.Cuevas {
		if len(sv.grafo.AristasSalientes(id)) == 0 &&
			len(sv.grafo.ProximasAristas(id)) == 0 {
			pozos = append(pozos, id)
		}
	}
	return pozos
}

// AnalizarAccesibilidad analiza qué cuevas son accesibles desde un punto de inicio
func (sv *ServicioValidacion) AnalizarAccesibilidad(cuevaInicio string) *ResultadoAccesibilidad {
	if _, existe := sv.grafo.ObtenerCueva(cuevaInicio); !existe {
		return &ResultadoAccesibilidad{
			CuevasInaccesibles: []string{},
			Soluciones:         []string{fmt.Sprintf("La cueva de inicio '%s' no existe", cuevaInicio)},
			TotalCuevas:        len(sv.grafo.Cuevas),
			CuevasAccesibles:   0,
		}
	}

	cuevasAccesibles := sv.obtenerCuevasAccesibles(cuevaInicio)
	cuevasInaccesibles := sv.encontrarCuevasInaccesibles(cuevasAccesibles)
	soluciones := sv.generarSoluciones(cuevasInaccesibles, cuevasAccesibles)

	return &ResultadoAccesibilidad{
		CuevasInaccesibles: cuevasInaccesibles,
		Soluciones:         soluciones,
		TotalCuevas:        len(sv.grafo.Cuevas),
		CuevasAccesibles:   len(cuevasAccesibles),
	}
}

// DetectarCuevasInaccesiblesTrasChanged detecta cuevas inaccesibles después de cambios en las conexiones
func (sv *ServicioValidacion) DetectarCuevasInaccesiblesTrasChanged() *ResultadoAccesibilidad {
	// Encontrar una cueva de referencia (idealmente la primera o una entrada principal)
	var cuevaReferencia string
	for id := range sv.grafo.Cuevas {
		cuevaReferencia = id
		break
	}

	if cuevaReferencia == "" {
		return &ResultadoAccesibilidad{
			CuevasInaccesibles: []string{},
			Soluciones:         []string{"No hay cuevas en el grafo"},
			TotalCuevas:        0,
			CuevasAccesibles:   0,
		}
	}

	return sv.AnalizarAccesibilidad(cuevaReferencia)
}

// obtenerCuevasAccesibles utiliza DFS para encontrar todas las cuevas accesibles desde un punto
func (sv *ServicioValidacion) obtenerCuevasAccesibles(inicio string) map[string]bool {
	visitadas := make(map[string]bool)
	sv.dfs(inicio, visitadas)
	return visitadas
}

// dfs implementa búsqueda en profundidad
func (sv *ServicioValidacion) dfs(cuevaID string, visitadas map[string]bool) {
	visitadas[cuevaID] = true

	// Recorrer todas las aristas no obstruidas
	for _, arista := range sv.grafo.Aristas {
		if arista.EsObstruido {
			continue
		}

		// Arista saliente
		if arista.Desde == cuevaID {
			if !visitadas[arista.Hasta] {
				sv.dfs(arista.Hasta, visitadas)
			}
		}

		// Arista entrante (si la arista es no dirigida O si el grafo es no dirigido)
		if (!arista.EsDirigido || !sv.grafo.EsDirigido) && arista.Hasta == cuevaID {
			if !visitadas[arista.Desde] {
				sv.dfs(arista.Desde, visitadas)
			}
		}
	}
}

// encontrarCuevasInaccesibles identifica las cuevas que no fueron visitadas
func (sv *ServicioValidacion) encontrarCuevasInaccesibles(cuevasAccesibles map[string]bool) []string {
	var inaccesibles []string
	for id := range sv.grafo.Cuevas {
		if !cuevasAccesibles[id] {
			inaccesibles = append(inaccesibles, id)
		}
	}
	return inaccesibles
}

// generarSoluciones propone soluciones para conectar las cuevas inaccesibles
func (sv *ServicioValidacion) generarSoluciones(cuevasInaccesibles []string, cuevasAccesibles map[string]bool) []string {
	var soluciones []string

	if len(cuevasInaccesibles) == 0 {
		soluciones = append(soluciones, "Todas las cuevas son accesibles. No se requieren cambios.")
		return soluciones
	}

	soluciones = append(soluciones, fmt.Sprintf("Se detectaron %d cuevas inaccesibles:", len(cuevasInaccesibles)))

	// Analizar cada cueva inaccesible
	for _, cuevaInaccesible := range cuevasInaccesibles {
		solucionesCueva := sv.analizarSolucionesParaCueva(cuevaInaccesible, cuevasAccesibles)
		soluciones = append(soluciones, solucionesCueva...)
	}

	// Soluciones generales
	soluciones = append(soluciones, "")
	soluciones = append(soluciones, "SOLUCIONES GENERALES:")
	soluciones = append(soluciones, "1. Revisar conexiones obstruidas y eliminar obstrucciones innecesarias")
	soluciones = append(soluciones, "2. Cambiar direcciones de aristas para permitir acceso bidireccional")
	soluciones = append(soluciones, "3. Agregar nuevas conexiones entre cuevas accesibles e inaccesibles")
	soluciones = append(soluciones, "4. Verificar si el cambio a grafo no dirigido resuelve la conectividad")

	return soluciones
}

// analizarSolucionesParaCueva propone soluciones específicas para una cueva inaccesible
func (sv *ServicioValidacion) analizarSolucionesParaCueva(cuevaInaccesible string, cuevasAccesibles map[string]bool) []string {
	var soluciones []string

	// Verificar si tiene conexiones obstruidas
	conexionesObstruidas := sv.verificarConexionesObstruidas(cuevaInaccesible)
	if len(conexionesObstruidas) > 0 {
		soluciones = append(soluciones, fmt.Sprintf("- Cueva '%s': Eliminar obstrucciones en: %v", cuevaInaccesible, conexionesObstruidas))
	}

	// Verificar si tiene aristas con direcciones problemáticas
	aristasProblematicas := sv.verificarAristasDireccionales(cuevaInaccesible, cuevasAccesibles)
	if len(aristasProblematicas) > 0 {
		soluciones = append(soluciones, fmt.Sprintf("- Cueva '%s': Cambiar dirección de: %v", cuevaInaccesible, aristasProblematicas))
	}

	// Sugerir conexiones nuevas
	conexionesSugeridas := sv.sugerirNuevasConexiones(cuevaInaccesible, cuevasAccesibles)
	if len(conexionesSugeridas) > 0 {
		soluciones = append(soluciones, fmt.Sprintf("- Cueva '%s': Agregar conexiones hacia: %v", cuevaInaccesible, conexionesSugeridas))
	}

	if len(soluciones) == 0 {
		soluciones = append(soluciones, fmt.Sprintf("- Cueva '%s': Requiere análisis manual para determinar la mejor conexión", cuevaInaccesible))
	}

	return soluciones
}

// verificarConexionesObstruidas encuentra conexiones obstruidas que podrían resolver la accesibilidad
func (sv *ServicioValidacion) verificarConexionesObstruidas(cuevaID string) []string {
	var conexionesObstruidas []string

	for _, arista := range sv.grafo.Aristas {
		if arista.EsObstruido {
			// Conexión saliente obstruida
			if arista.Desde == cuevaID {
				conexionesObstruidas = append(conexionesObstruidas, fmt.Sprintf("%s -> %s", arista.Desde, arista.Hasta))
			}
			// Conexión entrante obstruida (importante para accesibilidad)
			if arista.Hasta == cuevaID {
				conexionesObstruidas = append(conexionesObstruidas, fmt.Sprintf("%s -> %s", arista.Desde, arista.Hasta))
			}
		}
	}

	return conexionesObstruidas
}

// verificarAristasDireccionales encuentra aristas que podrían beneficiarse de un cambio de dirección
func (sv *ServicioValidacion) verificarAristasDireccionales(cuevaID string, cuevasAccesibles map[string]bool) []string {
	var aristasProblematicas []string

	for _, arista := range sv.grafo.Aristas {
		// Solo considerar aristas dirigidas
		if !arista.EsDirigido || arista.EsObstruido {
			continue
		}

		// Si hay una arista que va desde una cueva accesible hacia la inaccesible, pero en dirección contraria
		if arista.Desde == cuevaID && cuevasAccesibles[arista.Hasta] {
			aristasProblematicas = append(aristasProblematicas, fmt.Sprintf("%s -> %s (invertir)", arista.Desde, arista.Hasta))
		}
	}

	return aristasProblematicas
}

// sugerirNuevasConexiones sugiere conexiones nuevas con cuevas accesibles cercanas
func (sv *ServicioValidacion) sugerirNuevasConexiones(cuevaID string, cuevasAccesibles map[string]bool) []string {
	var sugerencias []string
	contador := 0
	maxSugerencias := 3

	// Sugerir conexiones con las primeras cuevas accesibles encontradas
	for cuevaAccesible := range cuevasAccesibles {
		if contador >= maxSugerencias {
			break
		}

		// Verificar que no existe ya una conexión
		if !sv.existeConexion(cuevaAccesible, cuevaID) {
			sugerencias = append(sugerencias, cuevaAccesible)
			contador++
		}
	}

	return sugerencias
}

// existeConexion verifica si existe una conexión entre dos cuevas
func (sv *ServicioValidacion) existeConexion(desde, hasta string) bool {
	for _, arista := range sv.grafo.Aristas {
		if (arista.Desde == desde && arista.Hasta == hasta) ||
			(!sv.grafo.EsDirigido && arista.Desde == hasta && arista.Hasta == desde) {
			return true
		}
	}
	return false
}
