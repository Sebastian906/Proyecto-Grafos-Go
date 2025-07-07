package handler

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/service"
)

// AnalysisHandler maneja las operaciones de análisis del grafo
type AnalysisHandler struct {
	mstService *service.MSTService
}

// NuevoAnalysisHandler crea una nueva instancia del handler de análisis
func NuevoAnalysisHandler(mstService *service.MSTService) *AnalysisHandler {
	return &AnalysisHandler{
		mstService: mstService,
	}
}

// CalcularMSTGeneral maneja el requisito 3a: calcular MST general
func (ah *AnalysisHandler) CalcularMSTGeneral(grafo *domain.Grafo) (string, error) {
	if grafo == nil {
		return "", fmt.Errorf("no hay grafo cargado en el sistema")
	}

	// Validar prerrequisitos
	if err := ah.mstService.ValidarPrerequisitos(grafo); err != nil {
		return fmt.Sprintf("Error: %v\n\nSugerencias:\n- Verifique que el grafo tenga al menos 2 cuevas\n- Asegúrese de que existan conexiones válidas\n- Revise que no todas las conexiones estén obstruidas", err), nil
	}

	// Obtener estadísticas de la red antes del cálculo
	stats := ah.mstService.ObtenerEstadisticasRed(grafo)

	// Calcular MST
	resultado, err := ah.mstService.ObtenerMSTGeneral(grafo)
	if err != nil {
		return "", fmt.Errorf("error al calcular el árbol de expansión mínimo: %v", err)
	}

	// Formatear resultado para visualización
	output := ah.mstService.FormatearResultadoParaVisualizacion(resultado)

	// Agregar estadísticas adicionales
	output += "\n=== ESTADÍSTICAS DE LA RED ===\n"
	output += fmt.Sprintf("Total de cuevas: %v\n", stats["total_cuevas"])
	output += fmt.Sprintf("Total de conexiones: %v\n", stats["total_aristas"])
	output += fmt.Sprintf("Conexiones válidas: %v\n", stats["aristas_validas"])
	output += fmt.Sprintf("Conexiones obstruidas: %v\n", stats["aristas_obstruidas"])
	output += fmt.Sprintf("Peso total de la red completa: %.2f\n", stats["peso_total_red"])

	if stats["peso_promedio_arista"] != nil {
		output += fmt.Sprintf("Peso promedio por conexión: %.2f\n", stats["peso_promedio_arista"])
	}

	output += fmt.Sprintf("Tipo de red: ")
	if stats["es_dirigido"].(bool) {
		output += "Dirigida\n"
	} else {
		output += "No dirigida\n"
	}

	// Si el MST fue calculado exitosamente, mostrar comparación
	if resultado.MST != nil {
		pesoOriginal := stats["peso_total_red"].(float64)
		ahorro := pesoOriginal - resultado.MST.PesoTotal
		porcentajeAhorro := (ahorro / pesoOriginal) * 100

		output += "\n=== ANÁLISIS DE OPTIMIZACIÓN ===\n"
		output += fmt.Sprintf("Peso total sin optimización: %.2f\n", pesoOriginal)
		output += fmt.Sprintf("Peso total con MST: %.2f\n", resultado.MST.PesoTotal)
		output += fmt.Sprintf("Ahorro total: %.2f\n", ahorro)
		output += fmt.Sprintf("Porcentaje de optimización: %.2f%%\n", porcentajeAhorro)
		output += fmt.Sprintf("Conexiones eliminadas: %d\n", stats["aristas_validas"].(int)-resultado.MST.NumAristas)
	}

	return output, nil
}

// ObtenerEstadisticasRed proporciona información estadística sobre la red de cuevas
func (ah *AnalysisHandler) ObtenerEstadisticasRed(grafo *domain.Grafo) (string, error) {
	if grafo == nil {
		return "", fmt.Errorf("no hay grafo cargado en el sistema")
	}

	stats := ah.mstService.ObtenerEstadisticasRed(grafo)

	output := "=== ESTADÍSTICAS DE LA RED DE CUEVAS ===\n\n"
	output += fmt.Sprintf(" Información general:\n")
	output += fmt.Sprintf("   • Total de cuevas: %v\n", stats["total_cuevas"])
	output += fmt.Sprintf("   • Total de conexiones: %v\n", stats["total_aristas"])
	output += fmt.Sprintf("   • Conexiones activas: %v\n", stats["aristas_validas"])
	output += fmt.Sprintf("   • Conexiones obstruidas: %v\n", stats["aristas_obstruidas"])

	output += fmt.Sprintf("\n🔗 Características de la red:\n")
	if stats["es_dirigido"].(bool) {
		output += "   • Tipo: Red dirigida (conexiones unidireccionales)\n"
	} else {
		output += "   • Tipo: Red no dirigida (conexiones bidireccionales)\n"
	}

	if stats["es_conexo"].(bool) {
		output += "   • Conectividad: Todas las cuevas están conectadas\n"
	} else {
		output += "   • Conectividad:  Existen cuevas aisladas\n"
	}

	output += fmt.Sprintf("\n Métricas de distancia:\n")
	output += fmt.Sprintf("   • Distancia total de la red: %.2f\n", stats["peso_total_red"])

	if stats["peso_promedio_arista"] != nil {
		output += fmt.Sprintf("   • Distancia promedio por conexión: %.2f\n", stats["peso_promedio_arista"])
	}

	// Calcular densidad del grafo
	numCuevas := stats["total_cuevas"].(int)
	numAristas := stats["aristas_validas"].(int)

	if numCuevas > 1 {
		var maxAristas int
		if stats["es_dirigido"].(bool) {
			maxAristas = numCuevas * (numCuevas - 1)
		} else {
			maxAristas = numCuevas * (numCuevas - 1) / 2
		}

		densidad := float64(numAristas) / float64(maxAristas) * 100
		output += fmt.Sprintf("   • Densidad de la red: %.2f%%\n", densidad)

		if densidad < 30 {
			output += "     (Red dispersa - pocas conexiones)\n"
		} else if densidad > 70 {
			output += "     (Red densa - muchas conexiones)\n"
		} else {
			output += "     (Red moderadamente conectada)\n"
		}
	}

	return output, nil
}

// ValidarConectividad verifica si la red permite crear un MST
func (ah *AnalysisHandler) ValidarConectividad(grafo *domain.Grafo) (string, error) {
	if grafo == nil {
		return "", fmt.Errorf("no hay grafo cargado en el sistema")
	}

	output := "=== VALIDACIÓN DE CONECTIVIDAD ===\n\n"

	// Validar prerrequisitos básicos
	if err := ah.mstService.ValidarPrerequisitos(grafo); err != nil {
		output += fmt.Sprintf(" Error: %v\n\n", err)
		output += "🔧 Soluciones sugeridas:\n"

		if len(grafo.Cuevas) < 2 {
			output += "   • Agregar más cuevas a la red\n"
		}

		// Verificar aristas obstruidas
		todasObstruidas := true
		for _, arista := range grafo.Aristas {
			if !arista.EsObstruido {
				todasObstruidas = false
				break
			}
		}

		if todasObstruidas {
			output += "   • Desbloquear al menos algunas conexiones\n"
		}

		if len(grafo.Aristas) == 0 {
			output += "   • Crear conexiones entre las cuevas\n"
		}

		return output, nil
	}

	// Obtener resultado del MST para verificar conectividad
	resultado, err := ah.mstService.ObtenerMSTGeneral(grafo)
	if err != nil {
		return "", fmt.Errorf("error al validar conectividad: %v", err)
	}

	if resultado.EsConexo {
		output += " La red está completamente conectada\n"
		output += " Es posible calcular un árbol de expansión mínimo\n\n"
		output += fmt.Sprintf(" Información del MST potencial:\n")
		output += fmt.Sprintf("   • Conexiones requeridas: %d\n", len(grafo.Cuevas)-1)
		output += fmt.Sprintf("   • Peso total mínimo: %.2f\n", resultado.MST.PesoTotal)
	} else {
		output += " La red NO está completamente conectada\n"
		output += fmt.Sprintf(" Se encontraron %d componentes separados\n\n", resultado.ComponentesConexos)
		output += " Para calcular MST necesita:\n"
		output += "   • Conectar todos los componentes aislados\n"
		output += "   • Agregar conexiones entre grupos de cuevas separados\n"
	}

	return output, nil
}

// ExportarMST exporta el MST como un nuevo grafo independiente
func (ah *AnalysisHandler) ExportarMST(grafo *domain.Grafo) (*domain.Grafo, string, error) {
	if grafo == nil {
		return nil, "", fmt.Errorf("no hay grafo cargado en el sistema")
	}

	// Calcular MST
	resultado, err := ah.mstService.ObtenerMSTGeneral(grafo)
	if err != nil {
		return nil, "", fmt.Errorf("error al calcular MST: %v", err)
	}

	if !resultado.EsConexo {
		return nil, "", fmt.Errorf("no se puede exportar MST: la red no está conectada")
	}

	// Crear grafo MST
	grafoMST := ah.mstService.ExportarMSTComoGrafo(resultado.MST, grafo)

	resumen := fmt.Sprintf("MST exportado exitosamente:\n")
	resumen += fmt.Sprintf("• Cuevas: %d\n", len(grafoMST.Cuevas))
	resumen += fmt.Sprintf("• Conexiones: %d\n", len(grafoMST.Aristas))
	resumen += fmt.Sprintf("• Peso total: %.2f\n", resultado.MST.PesoTotal)

	return grafoMST, resumen, nil
}

// CalcularMSTDesdeCueva maneja el requisito 3b: calcular MST desde cueva específica
func (ah *AnalysisHandler) CalcularMSTDesdeCueva(grafo *domain.Grafo, cuevaOrigen string) (string, error) {
	if grafo == nil {
		return "", fmt.Errorf("no hay grafo cargado en el sistema")
	}

	// Validar que la cueva origen existe
	if _, existe := grafo.Cuevas[cuevaOrigen]; !existe {
		cuevasDisponibles := make([]string, 0, len(grafo.Cuevas))
		for id := range grafo.Cuevas {
			cuevasDisponibles = append(cuevasDisponibles, id)
		}
		return fmt.Sprintf("Error: La cueva '%s' no existe en el grafo\n\nCuevas disponibles: %v", cuevaOrigen, cuevasDisponibles), nil
	}

	// Obtener estadísticas de la red antes del cálculo
	stats := ah.mstService.ObtenerEstadisticasRed(grafo)

	// Calcular MST desde la cueva específica
	resultado, err := ah.mstService.ObtenerMSTDesdeCueva(grafo, cuevaOrigen)
	if err != nil {
		return "", fmt.Errorf("error al calcular MST desde cueva '%s': %v", cuevaOrigen, err)
	}

	// Formatear resultado para visualización
	output := ah.mstService.FormatearResultadoMSTDesdeCuevaParaVisualizacion(resultado)

	// Agregar estadísticas adicionales
	output += "\n=== ESTADÍSTICAS DE LA RED ===\n"
	output += fmt.Sprintf("Total de cuevas en la red: %v\n", stats["total_cuevas"])
	output += fmt.Sprintf("Total de conexiones: %v\n", stats["total_aristas"])
	output += fmt.Sprintf("Conexiones válidas: %v\n", stats["aristas_validas"])

	if resultado.MST != nil && resultado.MST.MST != nil {
		output += "\n=== ANÁLISIS DE COBERTURA ===\n"
		cobertura := float64(resultado.TotalAlcanzables) / float64(stats["total_cuevas"].(int)) * 100
		output += fmt.Sprintf("Porcentaje de cobertura: %.2f%%\n", cobertura)

		if resultado.EsCompleto {
			output += "Resultado: ÓPTIMO - Todas las cuevas son alcanzables desde el origen\n"
		} else {
			output += "Resultado: PARCIAL - Algunas cuevas no son alcanzables desde el origen\n"
			output += fmt.Sprintf("Cuevas aisladas: %d\n", resultado.TotalNoAlcanzable)
			output += "\nSugerencias:\n"
			output += "• Verificar conectividad entre componentes de la red\n"
			output += "• Considerar agregar conexiones para mejorar la cobertura\n"
			output += "• Seleccionar una cueva origen diferente en un componente más grande\n"
		}

		// Análisis de eficiencia
		if resultado.EsCompleto {
			pesoOriginal := stats["peso_total_red"].(float64)
			ahorro := pesoOriginal - resultado.MST.MST.PesoTotal
			porcentajeAhorro := (ahorro / pesoOriginal) * 100

			output += "\n=== ANÁLISIS DE EFICIENCIA ===\n"
			output += fmt.Sprintf("Peso total de la red completa: %.2f\n", pesoOriginal)
			output += fmt.Sprintf("Peso del MST desde origen: %.2f\n", resultado.MST.MST.PesoTotal)
			output += fmt.Sprintf("Ahorro logrado: %.2f (%.2f%%)\n", ahorro, porcentajeAhorro)
			output += fmt.Sprintf("Conexiones eliminadas: %d\n", stats["aristas_validas"].(int)-resultado.MST.MST.NumAristas)
		}
	}

	return output, nil
}

// ListarCuevasDisponibles proporciona una lista de cuevas disponibles para usar como origen
func (ah *AnalysisHandler) ListarCuevasDisponibles(grafo *domain.Grafo) (string, error) {
	if grafo == nil {
		return "", fmt.Errorf("no hay grafo cargado en el sistema")
	}

	if len(grafo.Cuevas) == 0 {
		return "No hay cuevas disponibles en la red", nil
	}

	output := "=== CUEVAS DISPONIBLES PARA MST ===\n\n"

	// Contar conexiones por cueva para sugerir mejores orígenes
	conexionesPorCueva := make(map[string]int)
	for _, arista := range grafo.Aristas {
		if !arista.EsObstruido {
			conexionesPorCueva[arista.Desde]++
			if !arista.EsDirigido {
				conexionesPorCueva[arista.Hasta]++
			}
		}
	}

	// Listar cuevas con información relevante
	for id, cueva := range grafo.Cuevas {
		conexiones := conexionesPorCueva[id]
		output += fmt.Sprintf("• %s", id)
		if cueva.Nombre != "" && cueva.Nombre != id {
			output += fmt.Sprintf(" (%s)", cueva.Nombre)
		}
		output += fmt.Sprintf(" - %d conexiones", conexiones)

		if conexiones == 0 {
			output += " [AISLADA]"
		} else if conexiones >= 3 {
			output += " [BIEN CONECTADA]"
		}
		output += "\n"
	}

	// Sugerencias
	output += "\n=== RECOMENDACIONES ===\n"

	// Encontrar la cueva con más conexiones
	maxConexiones := 0
	mejorCueva := ""
	for id, conexiones := range conexionesPorCueva {
		if conexiones > maxConexiones {
			maxConexiones = conexiones
			mejorCueva = id
		}
	}

	if mejorCueva != "" {
		output += fmt.Sprintf("• Cueva recomendada: %s (%d conexiones)\n", mejorCueva, maxConexiones)
		output += "  Razón: Mayor número de conexiones puede resultar en mejor cobertura\n"
	}

	// Verificar si hay cuevas aisladas
	cuevasAisladas := 0
	for id := range grafo.Cuevas {
		if conexionesPorCueva[id] == 0 {
			cuevasAisladas++
		}
	}

	if cuevasAisladas > 0 {
		output += fmt.Sprintf("\n⚠  Advertencia: %d cueva(s) aislada(s) sin conexiones\n", cuevasAisladas)
		output += "  MST desde cualquier origen no podrá alcanzar cuevas aisladas\n"
	}

	return output, nil
}
