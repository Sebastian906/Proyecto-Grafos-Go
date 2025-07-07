package cli

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/handler"
	"proyecto-grafos-go/internal/service"
	"strings"
)

type MenuAnalisis struct {
	validacionSvc   *service.ServicioValidacion
	grafoSvc        *service.ServicioGrafo
	conexionSvc     *service.ServicioConexion
	analysisHandler *handler.AnalysisHandler
}

func NuevoMenuAnalisis(
	validacionSvc *service.ServicioValidacion,
	grafoSvc *service.ServicioGrafo,
	conexionSvc *service.ServicioConexion,
	analysisHandler *handler.AnalysisHandler,
) *MenuAnalisis {
	return &MenuAnalisis{
		validacionSvc:   validacionSvc,
		grafoSvc:        grafoSvc,
		conexionSvc:     conexionSvc,
		analysisHandler: analysisHandler,
	}
}

func (m *MenuAnalisis) Mostrar() {
	for {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("MENÚ DE ANÁLISIS DE LA RED DE CUEVAS")
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("=== ANÁLISIS DE CONECTIVIDAD ===")
		fmt.Println("1. Verificar conectividad fuerte")
		fmt.Println("2. Detectar pozos")
		fmt.Println("3. Mostrar grados de vertices")
		fmt.Println("4. Detectar cuevas inaccesibles")
		fmt.Println("5. Analizar accesibilidad desde cueva específica")
		fmt.Println("")
		fmt.Println("=== ANÁLISIS DE OPTIMIZACIÓN (MST) ===")
		fmt.Println("6. Ver estadísticas de la red")
		fmt.Println("7. Validar conectividad para MST")
		fmt.Println("8. Calcular Árbol de Expansión Mínimo General (Req. 3a)")
		fmt.Println("9. MST desde cueva específica (Req. 3b)")
		fmt.Println("10. Listar cuevas disponibles para MST")
		fmt.Println("11. Exportar MST como nuevo grafo")
		fmt.Println("")
		fmt.Println("12. Salir")
		fmt.Println(strings.Repeat("=", 50))

		opcion := ObtenerInputInt("Seleccione una opción: ")

		switch opcion {
		case 1:
			m.mostrarConectividad()
		case 2:
			m.mostrarPozos()
		case 3:
			m.mostrarGrados()
		case 4:
			m.detectarCuevasInaccesibles()
		case 5:
			m.analizarAccesibilidadEspecifica()
		case 6:
			m.mostrarEstadisticasRed()
		case 7:
			m.validarConectividadMST()
		case 8:
			m.calcularMSTGeneral()
		case 9:
			m.calcularMSTDesdeCueva()
		case 10:
			m.listarCuevasDisponibles()
		case 11:
			m.exportarMST()
		case 12:
			return
		default:
			fmt.Println("Opción inválida")
		}
	}
}

// Métodos existentes de conectividad
func (m *MenuAnalisis) mostrarConectividad() {
	if m.validacionSvc.EsFuertementeConectado() {
		fmt.Println("El grafo es fuertemente conectado")
	} else {
		fmt.Println("El grafo NO es fuertemente conectado")
	}
}

func (m *MenuAnalisis) mostrarPozos() {
	pozos := m.validacionSvc.DetectarPozos()
	if len(pozos) == 0 {
		fmt.Println("No hay pozos en el grafo")
	} else {
		fmt.Println("Pozos encontrados:")
		for _, p := range pozos {
			fmt.Println("  -", p)
		}
	}
}

func (m *MenuAnalisis) mostrarGrados() {
	grados := m.grafoSvc.ObtenerGradosVertices()
	fmt.Println("\n GRADOS DE LOS VÉRTICES:")
	fmt.Println(strings.Repeat("-", 40))
	for id, g := range grados {
		fmt.Printf(" %s: Entrantes=%d, Salientes=%d, Total=%d\n",
			id, g["entrante"], g["saliente"], g["total"])
	}
}

func (m *MenuAnalisis) detectarCuevasInaccesibles() {
	resultado := m.validacionSvc.DetectarCuevasInaccesiblesTrasChanged()

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println(" ANÁLISIS DE ACCESIBILIDAD")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf(" Total de cuevas: %d\n", resultado.TotalCuevas)
	fmt.Printf(" Cuevas accesibles: %d\n", resultado.CuevasAccesibles)
	fmt.Printf(" Cuevas inaccesibles: %d\n", len(resultado.CuevasInaccesibles))

	if len(resultado.CuevasInaccesibles) > 0 {
		fmt.Println("\n CUEVAS INACCESIBLES:")
		for i, cueva := range resultado.CuevasInaccesibles {
			fmt.Printf("   %d. %s\n", i+1, cueva)
		}
	}

	fmt.Println("\n SOLUCIONES PROPUESTAS:")
	for _, solucion := range resultado.Soluciones {
		fmt.Println("  •", solucion)
	}

	fmt.Println("\nPresione Enter para continuar...")
	ObtenerInputString("")
}

func (m *MenuAnalisis) analizarAccesibilidadEspecifica() {
	cuevaInicio := ObtenerInputString("Ingrese el ID de la cueva de inicio: ")
	if cuevaInicio == "" {
		fmt.Println(" ID de cueva no puede estar vacío")
		return
	}

	resultado := m.validacionSvc.AnalizarAccesibilidad(cuevaInicio)

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf(" ANÁLISIS DE ACCESIBILIDAD DESDE '%s'\n", cuevaInicio)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf(" Total de cuevas: %d\n", resultado.TotalCuevas)
	fmt.Printf(" Cuevas accesibles: %d\n", resultado.CuevasAccesibles)
	fmt.Printf(" Cuevas inaccesibles: %d\n", len(resultado.CuevasInaccesibles))

	if len(resultado.CuevasInaccesibles) > 0 {
		fmt.Println("\n CUEVAS INACCESIBLES:")
		for i, cueva := range resultado.CuevasInaccesibles {
			fmt.Printf("   %d. %s\n", i+1, cueva)
		}
	}

	fmt.Println("\n SOLUCIONES PROPUESTAS:")
	for _, solucion := range resultado.Soluciones {
		fmt.Println("  •", solucion)
	}

	fmt.Println("\nPresione Enter para continuar...")
	ObtenerInputString("")
}

// Nuevos métodos para MST (Requisito 3a)
func (m *MenuAnalisis) mostrarEstadisticasRed() {
	grafo := m.grafoSvc.ObtenerGrafo()
	if grafo == nil {
		fmt.Println(" No hay grafo cargado en el sistema")
		return
	}

	fmt.Println("\n Obteniendo estadísticas de la red...")

	resultado, err := m.analysisHandler.ObtenerEstadisticasRed(grafo)
	if err != nil {
		fmt.Printf(" Error: %v\n", err)
		return
	}

	fmt.Println(resultado)
	fmt.Println("\nPresione Enter para continuar...")
	ObtenerInputString("")
}

func (m *MenuAnalisis) validarConectividadMST() {
	grafo := m.grafoSvc.ObtenerGrafo()
	if grafo == nil {
		fmt.Println(" No hay grafo cargado en el sistema")
		return
	}

	fmt.Println("\n Validando conectividad para MST...")

	resultado, err := m.analysisHandler.ValidarConectividad(grafo)
	if err != nil {
		fmt.Printf(" Error: %v\n", err)
		return
	}

	fmt.Println(resultado)
	fmt.Println("\nPresione Enter para continuar...")
	ObtenerInputString("")
}

func (m *MenuAnalisis) calcularMSTGeneral() {
	grafo := m.grafoSvc.ObtenerGrafo()
	if grafo == nil {
		fmt.Println(" No hay grafo cargado en el sistema")
		return
	}

	fmt.Println("\n🌲 Calculando Árbol de Expansión Mínimo...")
	fmt.Println("Este proceso encuentra las conexiones mínimas necesarias")
	fmt.Println("para mantener toda la red conectada con el menor costo total.")

	resultado, err := m.analysisHandler.CalcularMSTGeneral(grafo)
	if err != nil {
		fmt.Printf(" Error: %v\n", err)
		return
	}

	fmt.Println(resultado)

	// Preguntar si desea ver detalles adicionales
	if SolicitarConfirmacion("¿Desea ver una explicación detallada del algoritmo utilizado?") {
		m.mostrarExplicacionAlgoritmo()
	}

	fmt.Println("\nPresione Enter para continuar...")
	ObtenerInputString("")
}

func (m *MenuAnalisis) exportarMST() {
	grafo := m.grafoSvc.ObtenerGrafo()
	if grafo == nil {
		fmt.Println(" No hay grafo cargado en el sistema")
		return
	}

	fmt.Println("\n Exportando MST como nuevo grafo...")

	grafoMST, resumen, err := m.analysisHandler.ExportarMST(grafo)
	if err != nil {
		fmt.Printf(" Error: %v\n", err)
		return
	}

	fmt.Println(" " + resumen)

	// Mostrar opciones para el grafo exportado
	if SolicitarConfirmacion("¿Desea ver la estructura del MST exportado?") {
		m.mostrarEstructuraMST(grafoMST)
	}

	fmt.Println("\nPresione Enter para continuar...")
	ObtenerInputString("")
}

func (m *MenuAnalisis) mostrarEstructuraMST(grafoMST *domain.Grafo) {
	fmt.Println("\n ESTRUCTURA DEL ÁRBOL DE EXPANSIÓN MÍNIMO")
	fmt.Println(strings.Repeat("=", 55))

	fmt.Printf(" Cuevas en el MST: %d\n", len(grafoMST.Cuevas))
	for id, cueva := range grafoMST.Cuevas {
		fmt.Printf("   • %s: %s\n", id, cueva.Nombre)
	}

	fmt.Printf("\n Conexiones del MST: %d\n", len(grafoMST.Aristas))
	pesoTotal := 0.0
	for i, arista := range grafoMST.Aristas {
		fmt.Printf("   %d. %s ↔ %s (distancia: %.2f)\n",
			i+1, arista.Desde, arista.Hasta, arista.Distancia)
		pesoTotal += arista.Distancia
	}

	fmt.Printf("\n Peso total del MST: %.2f\n", pesoTotal)
}

func (m *MenuAnalisis) mostrarExplicacionAlgoritmo() {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println(" EXPLICACIÓN: ALGORITMO DE KRUSKAL PARA MST")
	fmt.Println(strings.Repeat("=", 70))

	fmt.Println(" Objetivo:")
	fmt.Println("   El algoritmo de Kruskal encuentra el árbol de expansión mínimo")
	fmt.Println("   (MST), que conecta todos los nodos con el menor costo total.")

	fmt.Println("\n  Funcionamiento:")
	fmt.Println("   1. Ordena todas las aristas por peso (distancia) de menor a mayor")
	fmt.Println("   2. Examina cada arista en orden y la incluye si:")
	fmt.Println("      • No forma un ciclo con las aristas ya seleccionadas")
	fmt.Println("      • Conecta componentes separados")
	fmt.Println("   3. Continúa hasta tener n-1 aristas (donde n = número de nodos)")

	fmt.Println("\n Estructuras utilizadas:")
	fmt.Println("   • Union-Find: Para detectar y prevenir ciclos eficientemente")
	fmt.Println("   • Ordenamiento: Para procesar aristas por orden de peso")

	fmt.Println("\n Ventajas:")
	fmt.Println("   • Garantiza la solución óptima (menor peso total)")
	fmt.Println("   • Eficiente para grafos dispersos")
	fmt.Println("   • Complejidad: O(E log E) donde E = número de aristas")

	fmt.Println("\n Aplicación en cuevas:")
	fmt.Println("   • Encuentra las conexiones mínimas para mantener toda la red unida")
	fmt.Println("   • Minimiza la distancia total de construcción/mantenimiento")
	fmt.Println("   • Identifica conexiones redundantes que pueden eliminarse")

	fmt.Println(strings.Repeat("=", 70))
}

// Nuevos métodos para MST desde cueva específica (Requisito 3b)
func (m *MenuAnalisis) calcularMSTDesdeCueva() {
	grafo := m.grafoSvc.ObtenerGrafo()
	if grafo == nil {
		fmt.Println(" No hay grafo cargado en el sistema")
		return
	}

	// Mostrar cuevas disponibles primero
	fmt.Println("\n Cuevas disponibles en la red:")
	for id, cueva := range grafo.Cuevas {
		nombre := cueva.Nombre
		if nombre == "" || nombre == id {
			fmt.Printf("   • %s\n", id)
		} else {
			fmt.Printf("   • %s (%s)\n", id, nombre)
		}
	}

	// Solicitar cueva origen
	cuevaOrigen := ObtenerInputString("\nIngrese el ID de la cueva origen para el MST: ")
	if cuevaOrigen == "" {
		fmt.Println(" ID de cueva no puede estar vacío")
		return
	}

	fmt.Printf("\n Calculando MST desde cueva '%s'...\n", cuevaOrigen)
	fmt.Println("Este proceso encuentra el árbol de expansión mínimo")
	fmt.Println("que conecta todas las cuevas alcanzables desde el punto de origen.")

	resultado, err := m.analysisHandler.CalcularMSTDesdeCueva(grafo, cuevaOrigen)
	if err != nil {
		fmt.Printf(" Error: %v\n", err)
		return
	}

	fmt.Println(resultado)

	// Preguntar si desea ver información adicional
	if SolicitarConfirmacion("¿Desea ver una explicación del algoritmo de Prim utilizado?") {
		m.mostrarExplicacionPrim()
	}

	fmt.Println("\nPresione Enter para continuar...")
	ObtenerInputString("")
}

func (m *MenuAnalisis) listarCuevasDisponibles() {
	grafo := m.grafoSvc.ObtenerGrafo()
	if grafo == nil {
		fmt.Println(" No hay grafo cargado en el sistema")
		return
	}

	fmt.Println("\n Analizando cuevas disponibles...")

	resultado, err := m.analysisHandler.ListarCuevasDisponibles(grafo)
	if err != nil {
		fmt.Printf(" Error: %v\n", err)
		return
	}

	fmt.Println(resultado)
	fmt.Println("\nPresione Enter para continuar...")
	ObtenerInputString("")
}

func (m *MenuAnalisis) mostrarExplicacionPrim() {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println(" EXPLICACIÓN: ALGORITMO DE PRIM PARA MST")
	fmt.Println(strings.Repeat("=", 70))

	fmt.Println(" Objetivo:")
	fmt.Println("   El algoritmo de Prim construye un árbol de expansión mínimo")
	fmt.Println("   creciendo desde un nodo específico, ideal para analizar")
	fmt.Println("   accesibilidad desde un punto de origen dado.")

	fmt.Println("\n Funcionamiento:")
	fmt.Println("   1. Inicia desde la cueva origen especificada")
	fmt.Println("   2. Mantiene un conjunto de nodos ya incluidos en el MST")
	fmt.Println("   3. En cada paso, selecciona la arista de menor peso que:")
	fmt.Println("      • Conecta un nodo ya incluido con uno no incluido")
	fmt.Println("      • No está obstruida")
	fmt.Println("   4. Continúa hasta que no puedan agregarse más nodos")

	fmt.Println("\n Estructuras utilizadas:")
	fmt.Println("   • Priority Queue (Min-Heap): Para seleccionar arista de menor peso")
	fmt.Println("   • Conjunto de visitados: Para rastrear nodos incluidos")
	fmt.Println("   • Mapa de adyacencia: Para encontrar conexiones válidas")

	fmt.Println("\n Ventajas para análisis por origen:")
	fmt.Println("   • Muestra exactamente qué cuevas son alcanzables desde el origen")
	fmt.Println("   • Identifica componentes desconectados desde el punto de vista del origen")
	fmt.Println("   • Proporciona las rutas mínimas desde el origen a cada destino")
	fmt.Println("   • Útil para planificar expansiones desde ubicaciones específicas")

	fmt.Println("\n Diferencias con Kruskal:")
	fmt.Println("   • Prim: Crece desde un punto específico (mejor para análisis por origen)")
	fmt.Println("   • Kruskal: Considera todas las aristas globalmente (mejor para MST total)")
	fmt.Println("   • Ambos garantizan el mismo peso total en grafos completamente conectados")

	fmt.Println("\n Aplicación en cuevas:")
	fmt.Println("   • Planificación de rutas desde una base o entrada principal")
	fmt.Println("   • Análisis de accesibilidad desde puntos estratégicos")
	fmt.Println("   • Identificación de cuevas aisladas desde ubicaciones específicas")

	fmt.Println(strings.Repeat("=", 70))
}
