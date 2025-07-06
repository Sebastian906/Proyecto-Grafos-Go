package handler

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/service"
	"strconv"
)

// SimulationHandler maneja las operaciones de simulación de camiones
type SimulationHandler struct {
	truckService     *service.TruckService
	traversalService *service.TraversalService
	graphService     *service.ServicioGrafo
}

// NuevoSimulationHandler crea una nueva instancia del controlador de simulación
func NuevoSimulationHandler(truckService *service.TruckService, traversalService *service.TraversalService, graphService *service.ServicioGrafo) *SimulationHandler {
	return &SimulationHandler{
		truckService:     truckService,
		traversalService: traversalService,
		graphService:     graphService,
	}
}

// CrearCamion maneja la creación de un nuevo camión
func (sh *SimulationHandler) CrearCamion(id string, tipoCamion string, cuevaOrigen string) (*service.Camion, error) {
	// Validar tipo de camión
	var tipo service.TipoCamion
	switch tipoCamion {
	case "PEQUEÑO", "PEQUENO", "SMALL":
		tipo = service.CamionPequeno
	case "MEDIANO", "MEDIUM":
		tipo = service.CamionMediano
	case "GRANDE", "LARGE":
		tipo = service.CamionGrande
	default:
		return nil, fmt.Errorf("tipo de camión no válido. Use: PEQUEÑO, MEDIANO, GRANDE")
	}

	return sh.truckService.CrearCamion(id, tipo, cuevaOrigen)
}

// CargarInsumosEnCamion maneja la carga de insumos en un camión
func (sh *SimulationHandler) CargarInsumosEnCamion(camionID string, insumos map[string]string) error {
	// Convertir map[string]string a map[string]int
	insumosInt := make(map[string]int)
	for recurso, cantidadStr := range insumos {
		cantidad, err := strconv.Atoi(cantidadStr)
		if err != nil {
			return fmt.Errorf("cantidad inválida para recurso '%s': %s", recurso, cantidadStr)
		}
		if cantidad <= 0 {
			return fmt.Errorf("cantidad debe ser positiva para recurso '%s'", recurso)
		}
		insumosInt[recurso] = cantidad
	}

	return sh.truckService.CargarInsumos(camionID, insumosInt)
}

// EjecutarSimulacionDFS ejecuta una simulación de entrega usando DFS
func (sh *SimulationHandler) EjecutarSimulacionDFS(grafo *domain.Grafo, camionID string, cuevaOrigen string) (*service.SimulacionResultado, error) {
	if grafo == nil {
		return nil, fmt.Errorf("el grafo no puede ser nulo")
	}

	// Verificar que la cueva origen existe
	if _, existe := grafo.ObtenerCueva(cuevaOrigen); !existe {
		return nil, fmt.Errorf("cueva origen '%s' no existe en el grafo", cuevaOrigen)
	}

	return sh.truckService.SimularEntregaDFS(grafo, camionID, cuevaOrigen)
}

// EjecutarSimulacionBFS ejecuta una simulación de entrega usando BFS
func (sh *SimulationHandler) EjecutarSimulacionBFS(grafo *domain.Grafo, camionID string, cuevaOrigen string) (*service.SimulacionResultado, error) {
	if grafo == nil {
		return nil, fmt.Errorf("el grafo no puede ser nulo")
	}

	// Verificar que la cueva origen existe
	if _, existe := grafo.ObtenerCueva(cuevaOrigen); !existe {
		return nil, fmt.Errorf("cueva origen '%s' no existe en el grafo", cuevaOrigen)
	}

	return sh.truckService.SimularEntregaBFS(grafo, camionID, cuevaOrigen)
}

// CompararAlgoritmos compara el rendimiento de DFS vs BFS para la misma simulación
func (sh *SimulationHandler) CompararAlgoritmos(grafo *domain.Grafo, camionID string, cuevaOrigen string) (map[string]*service.SimulacionResultado, error) {
	if grafo == nil {
		return nil, fmt.Errorf("el grafo no puede ser nulo")
	}

	resultados := make(map[string]*service.SimulacionResultado)

	// Obtener estado inicial del camión para poder restaurarlo
	camion, err := sh.truckService.ObtenerCamion(camionID)
	if err != nil {
		return nil, err
	}

	// Hacer una copia de la carga inicial
	cargaOriginal := make(map[string]int)
	for recurso, cantidad := range camion.CargaActual {
		cargaOriginal[recurso] = cantidad
	}

	// Ejecutar simulación DFS
	resultadoDFS, err := sh.EjecutarSimulacionDFS(grafo, camionID, cuevaOrigen)
	if err != nil {
		return nil, fmt.Errorf("error en simulación DFS: %s", err.Error())
	}
	resultados["DFS"] = resultadoDFS

	// Restaurar carga del camión para la segunda simulación
	err = sh.truckService.ReiniciarCamion(camionID, cuevaOrigen)
	if err != nil {
		return nil, fmt.Errorf("error al reiniciar camión: %s", err.Error())
	}

	err = sh.truckService.CargarInsumos(camionID, cargaOriginal)
	if err != nil {
		return nil, fmt.Errorf("error al recargar insumos: %s", err.Error())
	}

	// Ejecutar simulación BFS
	resultadoBFS, err := sh.EjecutarSimulacionBFS(grafo, camionID, cuevaOrigen)
	if err != nil {
		return nil, fmt.Errorf("error en simulación BFS: %s", err.Error())
	}
	resultados["BFS"] = resultadoBFS

	return resultados, nil
}

// ObtenerEstadoCamion obtiene el estado actual de un camión
func (sh *SimulationHandler) ObtenerEstadoCamion(camionID string) (*service.Camion, error) {
	return sh.truckService.ObtenerCamion(camionID)
}

// ListarTodosLosCamiones obtiene todos los camiones registrados
func (sh *SimulationHandler) ListarTodosLosCamiones() map[string]*service.Camion {
	return sh.truckService.ListarCamiones()
}

// ReiniciarCamion reinicia un camión al estado inicial
func (sh *SimulationHandler) ReiniciarCamion(camionID string, cuevaOrigen string) error {
	return sh.truckService.ReiniciarCamion(camionID, cuevaOrigen)
}

// EliminarCamion elimina un camión del sistema
func (sh *SimulationHandler) EliminarCamion(camionID string) error {
	return sh.truckService.EliminarCamion(camionID)
}

// GenerarReporteComparativo genera un reporte comparativo entre DFS y BFS
func (sh *SimulationHandler) GenerarReporteComparativo(resultados map[string]*service.SimulacionResultado) string {
	if len(resultados) != 2 {
		return "Error: Se requieren exactamente 2 resultados para comparar"
	}

	dfs, existeDFS := resultados["DFS"]
	bfs, existeBFS := resultados["BFS"]

	if !existeDFS || !existeBFS {
		return "Error: Se requieren resultados de DFS y BFS para comparar"
	}

	reporte := "=== REPORTE COMPARATIVO DFS vs BFS ===\n\n"

	// Comparación general
	reporte += "--- COMPARACIÓN GENERAL ---\n"
	reporte += fmt.Sprintf("DFS - Éxito: %t, Tiempo: %v, Distancia: %.2f km, Cuevas: %d\n",
		dfs.Exitoso, dfs.TiempoTotal, dfs.DistanciaTotal, len(dfs.RutaCompleta))
	reporte += fmt.Sprintf("BFS - Éxito: %t, Tiempo: %v, Distancia: %.2f km, Cuevas: %d\n",
		bfs.Exitoso, bfs.TiempoTotal, bfs.DistanciaTotal, len(bfs.RutaCompleta))

	// Análisis de eficiencia
	reporte += "\n--- ANÁLISIS DE EFICIENCIA ---\n"
	if dfs.TiempoTotal < bfs.TiempoTotal {
		reporte += "DFS fue más rápido\n"
	} else if bfs.TiempoTotal < dfs.TiempoTotal {
		reporte += "BFS fue más rápido\n"
	} else {
		reporte += " Ambos algoritmos tardaron lo mismo\n"
	}

	if dfs.DistanciaTotal < bfs.DistanciaTotal {
		reporte += "DFS recorrió menor distancia\n"
	} else if bfs.DistanciaTotal < dfs.DistanciaTotal {
		reporte += "BFS recorrió menor distancia\n"
	} else {
		reporte += "Ambos algoritmos recorrieron la misma distancia\n"
	}

	// Comparación de entregas
	reporte += "\n--- COMPARACIÓN DE ENTREGAS ---\n"
	entregasDFS := len(dfs.EntregasRealizadas)
	entregasBFS := len(bfs.EntregasRealizadas)

	reporte += fmt.Sprintf("DFS realizó entregas en %d cuevas\n", entregasDFS)
	reporte += fmt.Sprintf("BFS realizó entregas en %d cuevas\n", entregasBFS)

	if entregasDFS > entregasBFS {
		reporte += "DFS realizó más entregas\n"
	} else if entregasBFS > entregasDFS {
		reporte += "BFS realizó más entregas\n"
	} else {
		reporte += "Ambos algoritmos realizaron la misma cantidad de entregas\n"
	}

	// Rutas seguidas
	reporte += "\n--- RUTAS SEGUIDAS ---\n"
	reporte += "DFS: " + fmt.Sprintf("%v", dfs.RutaCompleta) + "\n"
	reporte += "BFS: " + fmt.Sprintf("%v", bfs.RutaCompleta) + "\n"

	// Recomendación
	reporte += "\n--- RECOMENDACIÓN ---\n"
	if dfs.Exitoso && bfs.Exitoso {
		if dfs.TiempoTotal < bfs.TiempoTotal && dfs.DistanciaTotal <= bfs.DistanciaTotal {
			reporte += "RECOMENDACION: Se recomienda usar DFS por mejor rendimiento en tiempo y distancia\n"
		} else if bfs.TiempoTotal < dfs.TiempoTotal && bfs.DistanciaTotal <= dfs.DistanciaTotal {
			reporte += "RECOMENDACION: Se recomienda usar BFS por mejor rendimiento en tiempo y distancia\n"
		} else {
			reporte += "Ambos algoritmos son viables, elija según sus preferencias específicas\n"
		}
	} else if dfs.Exitoso {
		reporte += "RECOMENDACION: Use DFS (BFS falló)\n"
	} else if bfs.Exitoso {
		reporte += "RECOMENDACION: Use BFS (DFS falló)\n"
	} else {
		reporte += "ERROR: Ningún algoritmo completó exitosamente la simulación\n"
	}

	return reporte
}

// GenerarReporteSimulacion genera un reporte detallado de una simulación
func (sh *SimulationHandler) GenerarReporteSimulacion(resultado *service.SimulacionResultado) string {
	return sh.truckService.GenerarReporteSimulacion(resultado)
}

// ValidarParametrosSimulacion valida los parámetros antes de ejecutar una simulación
func (sh *SimulationHandler) ValidarParametrosSimulacion(grafo *domain.Grafo, camionID string, cuevaOrigen string) error {
	if grafo == nil {
		return fmt.Errorf("el grafo no puede ser nulo")
	}

	if camionID == "" {
		return fmt.Errorf("el ID del camión no puede estar vacío")
	}

	if cuevaOrigen == "" {
		return fmt.Errorf("la cueva origen no puede estar vacía")
	}

	// Verificar que el camión existe
	_, err := sh.truckService.ObtenerCamion(camionID)
	if err != nil {
		return fmt.Errorf("camión no encontrado: %s", err.Error())
	}

	// Verificar que la cueva origen existe
	if _, existe := grafo.ObtenerCueva(cuevaOrigen); !existe {
		return fmt.Errorf("cueva origen '%s' no existe en el grafo", cuevaOrigen)
	}

	// Verificar conectividad básica
	cuevasAccesibles, err := sh.traversalService.ObtenerCuevasAccesibles(grafo, cuevaOrigen)
	if err != nil {
		return fmt.Errorf("error verificando conectividad: %s", err.Error())
	}

	if len(cuevasAccesibles) == 0 {
		return fmt.Errorf("no hay cuevas accesibles desde '%s'", cuevaOrigen)
	}

	return nil
}
