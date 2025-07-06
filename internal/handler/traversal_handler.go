package handler

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/service"
)

// TraversalHandler maneja las operaciones de recorrido de grafos
type TraversalHandler struct {
	traversalService *service.TraversalService
	graphService     *service.ServicioGrafo
}

// NuevoTraversalHandler crea una nueva instancia del controlador de recorrido
func NuevoTraversalHandler(traversalService *service.TraversalService, graphService *service.ServicioGrafo) *TraversalHandler {
	return &TraversalHandler{
		traversalService: traversalService,
		graphService:     graphService,
	}
}

// EjecutarDFS ejecuta un recorrido DFS desde una cueva origen
func (th *TraversalHandler) EjecutarDFS(grafo *domain.Grafo, cuevaOrigen string) (*service.RecorridoResultado, error) {
	if err := th.validarParametros(grafo, cuevaOrigen); err != nil {
		return nil, err
	}

	return th.traversalService.RealizarRecorridoDFS(grafo, cuevaOrigen)
}

// EjecutarBFS ejecuta un recorrido BFS desde una cueva origen
func (th *TraversalHandler) EjecutarBFS(grafo *domain.Grafo, cuevaOrigen string) (*service.RecorridoResultado, error) {
	if err := th.validarParametros(grafo, cuevaOrigen); err != nil {
		return nil, err
	}

	return th.traversalService.RealizarRecorridoBFS(grafo, cuevaOrigen)
}

// CompararRecorridos compara DFS vs BFS desde la misma cueva origen
func (th *TraversalHandler) CompararRecorridos(grafo *domain.Grafo, cuevaOrigen string) (map[string]*service.RecorridoResultado, error) {
	if err := th.validarParametros(grafo, cuevaOrigen); err != nil {
		return nil, err
	}

	resultados := make(map[string]*service.RecorridoResultado)

	// Ejecutar DFS
	resultadoDFS, err := th.traversalService.RealizarRecorridoDFS(grafo, cuevaOrigen)
	if err != nil {
		return nil, fmt.Errorf("error en DFS: %s", err.Error())
	}
	resultados["DFS"] = resultadoDFS

	// Ejecutar BFS
	resultadoBFS, err := th.traversalService.RealizarRecorridoBFS(grafo, cuevaOrigen)
	if err != nil {
		return nil, fmt.Errorf("error en BFS: %s", err.Error())
	}
	resultados["BFS"] = resultadoBFS

	return resultados, nil
}

// AnalziarConectividad analiza la conectividad del grafo desde una cueva origen
func (th *TraversalHandler) AnalziarConectividad(grafo *domain.Grafo, cuevaOrigen string) (map[string]interface{}, error) {
	if err := th.validarParametros(grafo, cuevaOrigen); err != nil {
		return nil, err
	}

	resultado := make(map[string]interface{})

	// Verificar conectividad
	esConectado, cuevasNoAccesibles, err := th.traversalService.VerificarConectividad(grafo, cuevaOrigen)
	if err != nil {
		return nil, err
	}

	// Obtener cuevas accesibles
	cuevasAccesibles, err := th.traversalService.ObtenerCuevasAccesibles(grafo, cuevaOrigen)
	if err != nil {
		return nil, err
	}

	// Calcular estadísticas
	totalCuevas := len(grafo.Cuevas)
	cuevasAccesiblesCount := len(cuevasAccesibles)
	porcentajeConectividad := (float64(cuevasAccesiblesCount) / float64(totalCuevas)) * 100

	resultado["cueva_origen"] = cuevaOrigen
	resultado["es_conectado"] = esConectado
	resultado["total_cuevas"] = totalCuevas
	resultado["cuevas_accesibles"] = cuevasAccesibles
	resultado["cuevas_accesibles_count"] = cuevasAccesiblesCount
	resultado["cuevas_no_accesibles"] = cuevasNoAccesibles
	resultado["cuevas_no_accesibles_count"] = len(cuevasNoAccesibles)
	resultado["porcentaje_conectividad"] = porcentajeConectividad

	return resultado, nil
}

// ObtenerCuevasAccesibles obtiene todas las cuevas accesibles desde una cueva origen
func (th *TraversalHandler) ObtenerCuevasAccesibles(grafo *domain.Grafo, cuevaOrigen string) ([]string, error) {
	if err := th.validarParametros(grafo, cuevaOrigen); err != nil {
		return nil, err
	}

	return th.traversalService.ObtenerCuevasAccesibles(grafo, cuevaOrigen)
}

// GenerarReporteRecorrido genera un reporte detallado de un recorrido
func (th *TraversalHandler) GenerarReporteRecorrido(resultado *service.RecorridoResultado) string {
	if resultado == nil {
		return "Error: Resultado de recorrido nulo"
	}

	reporte := fmt.Sprintf("=== REPORTE DE RECORRIDO %s ===\n", resultado.TipoRecorrido)
	reporte += fmt.Sprintf("Cueva Origen: %s\n", resultado.CuevaOrigen)
	reporte += fmt.Sprintf("Completado: %t\n", resultado.Completado)
	reporte += fmt.Sprintf("Distancia Total: %.2f km\n", resultado.DistanciaTotal)
	reporte += fmt.Sprintf("Cuevas Visitadas: %d\n", len(resultado.CuevasVisitas))

	reporte += "\n--- SECUENCIA DE VISITAS ---\n"
	for i, cueva := range resultado.CuevasVisitas {
		orden := resultado.OrdenVisita[i]
		reporte += fmt.Sprintf("%d. %s (orden: %d)\n", i+1, cueva, orden)
	}

	reporte += "\n--- ESTADÍSTICAS ---\n"
	if len(resultado.CuevasVisitas) > 0 {
		distanciaPromedio := resultado.DistanciaTotal / float64(len(resultado.CuevasVisitas))
		reporte += fmt.Sprintf("Distancia promedio por cueva: %.2f km\n", distanciaPromedio)
	}

	return reporte
}

// GenerarReporteComparativoRecorridos genera un reporte comparativo entre DFS y BFS
func (th *TraversalHandler) GenerarReporteComparativoRecorridos(resultados map[string]*service.RecorridoResultado) string {
	if len(resultados) != 2 {
		return "Error: Se requieren exactamente 2 resultados para comparar"
	}

	dfs, existeDFS := resultados["DFS"]
	bfs, existeBFS := resultados["BFS"]

	if !existeDFS || !existeBFS {
		return "Error: Se requieren resultados de DFS y BFS para comparar"
	}

	reporte := "=== COMPARACIÓN DE ALGORITMOS DE RECORRIDO ===\n\n"

	// Información general
	reporte += "--- INFORMACIÓN GENERAL ---\n"
	reporte += fmt.Sprintf("Cueva Origen: %s\n", dfs.CuevaOrigen)
	reporte += fmt.Sprintf("DFS - Completado: %t, Distancia: %.2f km, Cuevas: %d\n",
		dfs.Completado, dfs.DistanciaTotal, len(dfs.CuevasVisitas))
	reporte += fmt.Sprintf("BFS - Completado: %t, Distancia: %.2f km, Cuevas: %d\n",
		bfs.Completado, bfs.DistanciaTotal, len(bfs.CuevasVisitas))

	// Comparación de eficiencia
	reporte += "\n--- ANÁLISIS COMPARATIVO ---\n"

	// Completitud
	if dfs.Completado && bfs.Completado {
		reporte += "RECOMENDACION: Ambos algoritmos completaron el recorrido\n"
	} else if dfs.Completado {
		reporte += "RECOMENDACION: Solo DFS completó el recorrido\n"
	} else if bfs.Completado {
		reporte += "RECOMENDACION: Solo BFS completó el recorrido\n"
	} else {
		reporte += "ERROR: Ningún algoritmo completó el recorrido\n"
	}

	// Distancia
	if dfs.DistanciaTotal < bfs.DistanciaTotal {
		reporte += "DFS recorrió menor distancia\n"
	} else if bfs.DistanciaTotal < dfs.DistanciaTotal {
		reporte += "BFS recorrió menor distancia\n"
	} else {
		reporte += "Ambos algoritmos recorrieron la misma distancia\n"
	}

	// Número de cuevas visitadas
	if len(dfs.CuevasVisitas) > len(bfs.CuevasVisitas) {
		reporte += "DFS visitó más cuevas\n"
	} else if len(bfs.CuevasVisitas) > len(dfs.CuevasVisitas) {
		reporte += "BFS visitó más cuevas\n"
	} else {
		reporte += "Ambos algoritmos visitaron la misma cantidad de cuevas\n"
	}

	// Secuencias de visita
	reporte += "\n--- SECUENCIAS DE VISITA ---\n"
	reporte += "DFS: " + fmt.Sprintf("%v", dfs.CuevasVisitas) + "\n"
	reporte += "BFS: " + fmt.Sprintf("%v", bfs.CuevasVisitas) + "\n"

	// Análisis de orden
	reporte += "\n--- ANÁLISIS DE ORDEN ---\n"
	if th.sonSecuenciasIguales(dfs.CuevasVisitas, bfs.CuevasVisitas) {
		reporte += "Ambos algoritmos siguieron la misma secuencia\n"
	} else {
		reporte += "Los algoritmos siguieron secuencias diferentes\n"

		// Cuevas comunes al inicio
		cuevasComunes := th.contarCuevasComunesAlInicio(dfs.CuevasVisitas, bfs.CuevasVisitas)
		if cuevasComunes > 0 {
			reporte += fmt.Sprintf("Las primeras %d cuevas fueron las mismas\n", cuevasComunes)
		}
	}

	// Recomendación
	reporte += "\n--- RECOMENDACIÓN ---\n"
	if dfs.Completado && bfs.Completado {
		if dfs.DistanciaTotal <= bfs.DistanciaTotal {
			reporte += "RECOMENDACION: Para este grafo, DFS parece más eficiente en distancia\n"
		} else {
			reporte += "RECOMENDACION: Para este grafo, BFS parece más eficiente en distancia\n"
		}
	} else if dfs.Completado {
		reporte += "RECOMENDACION: Use DFS para este tipo de grafo\n"
	} else if bfs.Completado {
		reporte += "RECOMENDACION: Use BFS para este tipo de grafo\n"
	} else {
		reporte += "ADVERTENCIA: Revisar la conectividad del grafo\n"
	}

	return reporte
}

// GenerarReporteConectividad genera un reporte de análisis de conectividad
func (th *TraversalHandler) GenerarReporteConectividad(analisis map[string]interface{}) string {
	reporte := "=== ANÁLISIS DE CONECTIVIDAD ===\n\n"

	reporte += fmt.Sprintf("Cueva Origen: %s\n", analisis["cueva_origen"])
	reporte += fmt.Sprintf("Grafo Conectado: %t\n", analisis["es_conectado"])
	reporte += fmt.Sprintf("Total de Cuevas: %d\n", analisis["total_cuevas"])
	reporte += fmt.Sprintf("Cuevas Accesibles: %d\n", analisis["cuevas_accesibles_count"])
	reporte += fmt.Sprintf("Cuevas No Accesibles: %d\n", analisis["cuevas_no_accesibles_count"])
	reporte += fmt.Sprintf("Porcentaje de Conectividad: %.2f%%\n", analisis["porcentaje_conectividad"])

	if cuevasAccesibles, ok := analisis["cuevas_accesibles"].([]string); ok {
		reporte += "\n--- CUEVAS ACCESIBLES ---\n"
		for i, cueva := range cuevasAccesibles {
			reporte += fmt.Sprintf("%d. %s\n", i+1, cueva)
		}
	}

	if cuevasNoAccesibles, ok := analisis["cuevas_no_accesibles"].([]string); ok && len(cuevasNoAccesibles) > 0 {
		reporte += "\n--- CUEVAS NO ACCESIBLES ---\n"
		for i, cueva := range cuevasNoAccesibles {
			reporte += fmt.Sprintf("%d. %s\n", i+1, cueva)
		}

		reporte += "\n--- RECOMENDACIONES ---\n"
		reporte += "ADVERTENCIA: Hay cuevas no accesibles desde el origen\n"
		reporte += "Considere agregar conexiones para mejorar la conectividad\n"
		reporte += "Verifique si hay obstáculos bloqueando las rutas\n"
	} else {
		reporte += "\nRECOMENDACION: Todas las cuevas son accesibles desde el origen\n"
	}

	return reporte
}

// validarParametros valida los parámetros comunes de entrada
func (th *TraversalHandler) validarParametros(grafo *domain.Grafo, cuevaOrigen string) error {
	if grafo == nil {
		return fmt.Errorf("el grafo no puede ser nulo")
	}

	if cuevaOrigen == "" {
		return fmt.Errorf("la cueva origen no puede estar vacía")
	}

	if _, existe := grafo.ObtenerCueva(cuevaOrigen); !existe {
		return fmt.Errorf("cueva origen '%s' no existe en el grafo", cuevaOrigen)
	}

	return nil
}

// sonSecuenciasIguales verifica si dos secuencias de cuevas son iguales
func (th *TraversalHandler) sonSecuenciasIguales(seq1, seq2 []string) bool {
	if len(seq1) != len(seq2) {
		return false
	}

	for i := range seq1 {
		if seq1[i] != seq2[i] {
			return false
		}
	}

	return true
}

// contarCuevasComunesAlInicio cuenta cuántas cuevas son iguales al inicio de ambas secuencias
func (th *TraversalHandler) contarCuevasComunesAlInicio(seq1, seq2 []string) int {
	minLen := len(seq1)
	if len(seq2) < minLen {
		minLen = len(seq2)
	}

	contador := 0
	for i := 0; i < minLen; i++ {
		if seq1[i] == seq2[i] {
			contador++
		} else {
			break
		}
	}

	return contador
}
