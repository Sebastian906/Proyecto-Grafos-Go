package service

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/pkg/algorithms"
	"strings"
)

// MSTService maneja las operaciones de Minimum Spanning Tree
type MSTService struct {
	servicioGrafo *ServicioGrafo
}

// NuevoMSTService crea una nueva instancia del servicio MST
func NuevoMSTService(servicioGrafo *ServicioGrafo) *MSTService {
	return &MSTService{
		servicioGrafo: servicioGrafo,
	}
}

// MSTResult contiene el resultado del cálculo del MST con información adicional
type MSTResult struct {
	MST                *algorithms.MST `json:"mst"`
	EsConexo           bool            `json:"es_conexo"`
	ComponentesConexos int             `json:"componentes_conexos"`
	Mensaje            string          `json:"mensaje"`
	Detalles           []string        `json:"detalles"`
}

// MSTDesdeCuevaResult contiene el resultado del cálculo del MST desde una cueva específica
type MSTDesdeCuevaResult struct {
	MST               *algorithms.MSTDesdeCueva `json:"mst_desde_cueva"`
	CuevaOrigen       string                    `json:"cueva_origen"`
	TotalAlcanzables  int                       `json:"total_alcanzables"`
	TotalNoAlcanzable int                       `json:"total_no_alcanzable"`
	EsCompleto        bool                      `json:"es_completo"`
	Mensaje           string                    `json:"mensaje"`
	Detalles          []string                  `json:"detalles"`
	RutasMinimas      map[string][]string       `json:"rutas_minimas"`
}

// ObtenerMSTGeneral implementa el requisito 3a:
// Obtiene el árbol de expansión mínimo general de toda la red
func (ms *MSTService) ObtenerMSTGeneral(grafo *domain.Grafo) (*MSTResult, error) {
	if grafo == nil {
		return nil, fmt.Errorf("el grafo no puede ser nil")
	}

	// Verificar que el grafo tenga al menos una cueva
	if len(grafo.Cuevas) == 0 {
		return &MSTResult{
			MST:      nil,
			EsConexo: false,
			Mensaje:  "El grafo está vacío",
			Detalles: []string{"No hay cuevas en la red"},
		}, nil
	}

	// Verificar conectividad
	esConexo := algorithms.EsConexo(grafo)
	if !esConexo {
		componentes := ms.contarComponentesConexos(grafo)
		return &MSTResult{
			MST:                nil,
			EsConexo:           false,
			ComponentesConexos: componentes,
			Mensaje:            "La red de cuevas no está completamente conectada",
			Detalles: []string{
				fmt.Sprintf("Se encontraron %d componentes conexos separados", componentes),
				"No es posible crear un árbol de expansión mínimo para una red desconectada",
				"Sugerencia: Agregar conexiones entre componentes aislados",
			},
		}, nil
	}

	// Calcular MST usando algoritmo de Kruskal
	mst, err := algorithms.Kruskal(grafo)
	if err != nil {
		return nil, fmt.Errorf("error al calcular MST: %v", err)
	}

	// Generar detalles del resultado
	detalles := ms.generarDetallesMST(mst, grafo)

	return &MSTResult{
		MST:                mst,
		EsConexo:           true,
		ComponentesConexos: 1,
		Mensaje:            "Árbol de expansión mínimo calculado exitosamente",
		Detalles:           detalles,
	}, nil
}

// ObtenerMSTDesdeCueva implementa el requisito 3b:
// Obtiene el árbol de expansión mínimo desde una cueva específica usando Prim
func (ms *MSTService) ObtenerMSTDesdeCueva(grafo *domain.Grafo, cuevaOrigen string) (*MSTDesdeCuevaResult, error) {
	if grafo == nil {
		return nil, fmt.Errorf("el grafo no puede ser nil")
	}

	// Validar que la cueva origen existe
	if _, existe := grafo.Cuevas[cuevaOrigen]; !existe {
		return nil, fmt.Errorf("la cueva origen '%s' no existe en el grafo", cuevaOrigen)
	}

	// Verificar prerequisitos básicos
	if err := ms.ValidarPrerequisitos(grafo); err != nil {
		return &MSTDesdeCuevaResult{
			MST:         nil,
			CuevaOrigen: cuevaOrigen,
			EsCompleto:  false,
			Mensaje:     fmt.Sprintf("Error en prerequisitos: %v", err),
			Detalles:    []string{"No se puede calcular MST debido a prerequisitos no cumplidos"},
		}, nil
	}

	// Calcular MST desde la cueva específica usando Prim
	mstResult, err := algorithms.Prim(grafo, cuevaOrigen)
	if err != nil {
		return nil, fmt.Errorf("error al calcular MST desde cueva '%s': %v", cuevaOrigen, err)
	}

	// Obtener rutas mínimas desde el origen
	rutasMinimas := mstResult.ObtenerRutasDesdeOrigen()

	// Generar detalles del resultado
	detalles := ms.generarDetallesMSTDesdeCueva(mstResult, grafo)

	// Determinar mensaje basado en si es completo o parcial
	var mensaje string
	if mstResult.EsCompleto {
		mensaje = fmt.Sprintf("MST completo desde cueva '%s' calculado exitosamente", cuevaOrigen)
	} else {
		mensaje = fmt.Sprintf("MST parcial desde cueva '%s' - Red no completamente conectada desde este punto", cuevaOrigen)
	}

	return &MSTDesdeCuevaResult{
		MST:               mstResult,
		CuevaOrigen:       cuevaOrigen,
		TotalAlcanzables:  len(mstResult.Alcanzables),
		TotalNoAlcanzable: len(mstResult.NoAlcanzable),
		EsCompleto:        mstResult.EsCompleto,
		Mensaje:           mensaje,
		Detalles:          detalles,
		RutasMinimas:      rutasMinimas,
	}, nil
}

// ValidarPrerequisitos verifica que se cumplan las condiciones para calcular MST
func (ms *MSTService) ValidarPrerequisitos(grafo *domain.Grafo) error {
	if grafo == nil {
		return fmt.Errorf("el grafo no puede ser nil")
	}

	if len(grafo.Cuevas) == 0 {
		return fmt.Errorf("el grafo debe tener al menos una cueva")
	}

	if len(grafo.Cuevas) == 1 {
		return fmt.Errorf("el grafo debe tener al menos dos cuevas para calcular MST")
	}

	// Verificar que existan aristas no obstruidas
	aristasValidas := 0
	for _, arista := range grafo.Aristas {
		if !arista.EsObstruido {
			aristasValidas++
		}
	}

	if aristasValidas == 0 {
		return fmt.Errorf("no existen conexiones válidas (no obstruidas) en el grafo")
	}

	return nil
}

// ObtenerEstadisticasRed proporciona información estadística de la red
func (ms *MSTService) ObtenerEstadisticasRed(grafo *domain.Grafo) map[string]interface{} {
	stats := make(map[string]interface{})

	stats["total_cuevas"] = len(grafo.Cuevas)
	stats["total_aristas"] = len(grafo.Aristas)

	// Contar aristas válidas (no obstruidas)
	aristasValidas := 0
	pesoTotal := 0.0
	for _, arista := range grafo.Aristas {
		if !arista.EsObstruido {
			aristasValidas++
			pesoTotal += arista.Distancia
		}
	}

	stats["aristas_validas"] = aristasValidas
	stats["aristas_obstruidas"] = len(grafo.Aristas) - aristasValidas
	stats["peso_total_red"] = pesoTotal
	stats["es_dirigido"] = grafo.EsDirigido
	stats["es_conexo"] = algorithms.EsConexo(grafo)

	if aristasValidas > 0 {
		stats["peso_promedio_arista"] = pesoTotal / float64(aristasValidas)
	}

	return stats
}

// contarComponentesConexos cuenta el número de componentes conexos en el grafo
func (ms *MSTService) contarComponentesConexos(grafo *domain.Grafo) int {
	visitados := make(map[string]bool)
	componentes := 0

	for id := range grafo.Cuevas {
		if !visitados[id] {
			componentes++
			ms.dfsComponente(grafo, id, visitados)
		}
	}

	return componentes
}

// dfsComponente realiza DFS para marcar todos los nodos de un componente conexo
func (ms *MSTService) dfsComponente(grafo *domain.Grafo, nodoActual string, visitados map[string]bool) {
	visitados[nodoActual] = true

	// Visitar todos los vecinos conectados por aristas no obstruidas
	for _, arista := range grafo.Aristas {
		if arista.EsObstruido {
			continue
		}

		var vecino string
		if arista.Desde == nodoActual {
			vecino = arista.Hasta
		} else if arista.Hasta == nodoActual && !arista.EsDirigido {
			vecino = arista.Desde
		} else {
			continue
		}

		if !visitados[vecino] {
			ms.dfsComponente(grafo, vecino, visitados)
		}
	}
}

// generarDetallesMST genera información detallada sobre el MST calculado
func (ms *MSTService) generarDetallesMST(mst *algorithms.MST, grafo *domain.Grafo) []string {
	var detalles []string

	detalles = append(detalles, fmt.Sprintf("Número total de cuevas en la red: %d", len(grafo.Cuevas)))
	detalles = append(detalles, fmt.Sprintf("Aristas en el MST: %d", mst.NumAristas))
	detalles = append(detalles, fmt.Sprintf("Peso total del MST: %.2f", mst.PesoTotal))

	if len(mst.Aristas) > 0 {
		pesoPromedio := mst.PesoTotal / float64(len(mst.Aristas))
		detalles = append(detalles, fmt.Sprintf("Peso promedio por arista: %.2f", pesoPromedio))
	}

	// Estadísticas de la red original
	stats := ms.ObtenerEstadisticasRed(grafo)
	pesoOriginal := stats["peso_total_red"].(float64)
	if pesoOriginal > 0 {
		porcentajeAhorro := ((pesoOriginal - mst.PesoTotal) / pesoOriginal) * 100
		detalles = append(detalles, fmt.Sprintf("Ahorro respecto a la red completa: %.2f%%", porcentajeAhorro))
	}

	// Listar las conexiones del MST
	if len(mst.Aristas) > 0 {
		detalles = append(detalles, "")
		detalles = append(detalles, "Conexiones mínimas requeridas:")
		for i, arista := range mst.Aristas {
			detalles = append(detalles, fmt.Sprintf("  %d. %s ↔ %s (distancia: %.2f)",
				i+1, arista.Desde, arista.Hasta, arista.Distancia))
		}
	}

	return detalles
}

// generarDetallesMSTDesdeCueva genera información detallada sobre el MST calculado desde una cueva específica
func (ms *MSTService) generarDetallesMSTDesdeCueva(mst *algorithms.MSTDesdeCueva, grafo *domain.Grafo) []string {
	var detalles []string

	detalles = append(detalles, fmt.Sprintf("Cueva de origen: %s", mst.CuevaOrigen))
	detalles = append(detalles, fmt.Sprintf("Cuevas alcanzables: %d de %d", len(mst.Alcanzables), len(grafo.Cuevas)))

	if mst.MST != nil {
		detalles = append(detalles, fmt.Sprintf("Conexiones en el MST: %d", mst.MST.NumAristas))
		detalles = append(detalles, fmt.Sprintf("Peso total del MST: %.2f", mst.MST.PesoTotal))

		if len(mst.MST.Aristas) > 0 {
			pesoPromedio := mst.MST.PesoTotal / float64(len(mst.MST.Aristas))
			detalles = append(detalles, fmt.Sprintf("Peso promedio por conexión: %.2f", pesoPromedio))
		}
	}

	// Información sobre cobertura
	if mst.EsCompleto {
		detalles = append(detalles, "Cobertura: COMPLETA - Todas las cuevas son alcanzables")
	} else {
		detalles = append(detalles, "Cobertura: PARCIAL - Algunas cuevas no son alcanzables")
		if len(mst.NoAlcanzable) > 0 {
			detalles = append(detalles, fmt.Sprintf("Cuevas no alcanzables: %v", mst.NoAlcanzable))
		}
	}

	// Listar las conexiones del MST
	if mst.MST != nil && len(mst.MST.Aristas) > 0 {
		detalles = append(detalles, "")
		detalles = append(detalles, "Conexiones mínimas requeridas:")
		for i, arista := range mst.MST.Aristas {
			detalles = append(detalles, fmt.Sprintf("  %d. %s -> %s (distancia: %.2f)",
				i+1, arista.Desde, arista.Hasta, arista.Distancia))
		}
	}

	return detalles
}

// ExportarMSTComoGrafo crea un nuevo grafo que contiene solo las aristas del MST
func (ms *MSTService) ExportarMSTComoGrafo(mst *algorithms.MST, grafoOriginal *domain.Grafo) *domain.Grafo {
	// Crear nuevo grafo para el MST
	grafoMST := domain.NuevoGrafo(false) // MST siempre es no dirigido

	// Agregar todas las cuevas del grafo original
	for id, cueva := range grafoOriginal.Cuevas {
		grafoMST.Cuevas[id] = &domain.Cueva{
			ID:       cueva.ID,
			Nombre:   cueva.Nombre,
			Recursos: cueva.Recursos,
		}
	}

	// Agregar solo las aristas del MST
	for _, arista := range mst.Aristas {
		grafoMST.Aristas = append(grafoMST.Aristas, &domain.Arista{
			Desde:       arista.Desde,
			Hasta:       arista.Hasta,
			Distancia:   arista.Distancia,
			EsDirigido:  false,
			EsObstruido: false,
		})
	}

	return grafoMST
}

// FormatearResultadoParaVisualizacion formatea el resultado para mostrar en CLI
func (ms *MSTService) FormatearResultadoParaVisualizacion(resultado *MSTResult) string {
	var sb strings.Builder

	sb.WriteString("=== ÁRBOL DE EXPANSIÓN MÍNIMO (MST) ===\n\n")
	sb.WriteString(fmt.Sprintf("Estado: %s\n", resultado.Mensaje))
	sb.WriteString(fmt.Sprintf("Red conexa: %v\n", resultado.EsConexo))

	if !resultado.EsConexo {
		sb.WriteString(fmt.Sprintf("Componentes conexos: %d\n", resultado.ComponentesConexos))
	}

	sb.WriteString("\n")

	if resultado.MST != nil {
		sb.WriteString(fmt.Sprintf("Peso total del MST: %.2f\n", resultado.MST.PesoTotal))
		sb.WriteString(fmt.Sprintf("Número de conexiones: %d\n", resultado.MST.NumAristas))
		sb.WriteString("\n")
	}

	// Agregar detalles
	for _, detalle := range resultado.Detalles {
		sb.WriteString(detalle + "\n")
	}

	return sb.String()
}

// FormatearResultadoMSTDesdeCuevaParaVisualizacion formatea el resultado para mostrar en CLI
func (ms *MSTService) FormatearResultadoMSTDesdeCuevaParaVisualizacion(resultado *MSTDesdeCuevaResult) string {
	var sb strings.Builder

	sb.WriteString("=== ÁRBOL DE EXPANSIÓN MÍNIMO DESDE CUEVA ESPECÍFICA ===\n\n")
	sb.WriteString(fmt.Sprintf("Cueva de origen: %s\n", resultado.CuevaOrigen))
	sb.WriteString(fmt.Sprintf("Estado: %s\n", resultado.Mensaje))
	sb.WriteString(fmt.Sprintf("Cobertura completa: %v\n", resultado.EsCompleto))
	sb.WriteString(fmt.Sprintf("Cuevas alcanzables: %d\n", resultado.TotalAlcanzables))

	if resultado.TotalNoAlcanzable > 0 {
		sb.WriteString(fmt.Sprintf("Cuevas no alcanzables: %d\n", resultado.TotalNoAlcanzable))
	}

	if resultado.MST != nil && resultado.MST.MST != nil {
		sb.WriteString(fmt.Sprintf("Peso total del MST: %.2f\n", resultado.MST.MST.PesoTotal))
		sb.WriteString(fmt.Sprintf("Número de conexiones: %d\n", resultado.MST.MST.NumAristas))
	}

	sb.WriteString("\n")

	// Agregar detalles
	for _, detalle := range resultado.Detalles {
		sb.WriteString(detalle + "\n")
	}

	// Mostrar rutas mínimas si existen
	if len(resultado.RutasMinimas) > 0 {
		sb.WriteString("\n=== RUTAS MÍNIMAS DESDE ORIGEN ===\n")
		for destino, ruta := range resultado.RutasMinimas {
			if destino != resultado.CuevaOrigen {
				sb.WriteString(fmt.Sprintf("%s: %v\n", destino, ruta))
			}
		}
	}

	return sb.String()
}
