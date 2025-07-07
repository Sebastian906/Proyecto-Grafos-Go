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
