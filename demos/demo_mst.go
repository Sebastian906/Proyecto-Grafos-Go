package demos

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/handler"
	"proyecto-grafos-go/internal/repository"
	"proyecto-grafos-go/internal/service"
	"strings"
)

// DemoMST ejecuta demostración del Árbol de Expansión Mínimo (MST)
func DemoMST() {
	fmt.Println(" DEMOSTRACIÓN: ÁRBOL DE EXPANSIÓN MÍNIMO (MST)")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Requisito 3a: Visualizar rutas mínimas para visitar toda la red")
	fmt.Println("")

	// Crear repositorio y cargar datos
	repositorio := repository.NuevoRepositorio("data")

	// Cargar el grafo de demostración
	fmt.Println(" Cargando archivo de demostración: caves_mst_demo.json")
	grafo, err := repositorio.CargarJSON("caves_mst_demo.json")
	if err != nil {
		fmt.Printf(" Error al cargar archivo: %v\n", err)
		return
	}

	fmt.Printf(" Archivo cargado exitosamente\n")
	fmt.Printf("   • Cuevas: %d\n", len(grafo.Cuevas))
	fmt.Printf("   • Conexiones: %d\n", len(grafo.Aristas))
	fmt.Printf("   • Tipo: ")
	if grafo.EsDirigido {
		fmt.Println("Dirigido")
	} else {
		fmt.Println("No dirigido")
	}
	fmt.Println("")

	// Mostrar estructura del grafo original
	mostrarEstructuraGrafo(grafo)

	// Crear servicios y handler
	servicioGrafo := service.NuevoServicioGrafo(grafo, repositorio)
	mstService := service.NuevoMSTService(servicioGrafo)
	analysisHandler := handler.NuevoAnalysisHandler(mstService)

	// Demostrar estadísticas de la red
	fmt.Println(" ESTADÍSTICAS DE LA RED ORIGINAL")
	fmt.Println(strings.Repeat("=", 50))
	estadisticas, err := analysisHandler.ObtenerEstadisticasRed(grafo)
	if err != nil {
		fmt.Printf(" Error: %v\n", err)
		return
	}
	fmt.Println(estadisticas)

	// Validar conectividad
	fmt.Println(" VALIDACIÓN DE CONECTIVIDAD")
	fmt.Println(strings.Repeat("=", 50))
	validacion, err := analysisHandler.ValidarConectividad(grafo)
	if err != nil {
		fmt.Printf(" Error: %v\n", err)
		return
	}
	fmt.Println(validacion)

	// Calcular MST (Requisito 3a)
	fmt.Println(" CÁLCULO DEL ÁRBOL DE EXPANSIÓN MÍNIMO")
	fmt.Println(strings.Repeat("=", 50))
	mstResult, err := analysisHandler.CalcularMSTGeneral(grafo)
	if err != nil {
		fmt.Printf(" Error: %v\n", err)
		return
	}
	fmt.Println(mstResult)

	// Exportar MST
	fmt.Println(" EXPORTACIÓN DEL MST")
	fmt.Println(strings.Repeat("=", 50))
	grafoMST, resumen, err := analysisHandler.ExportarMST(grafo)
	if err != nil {
		fmt.Printf(" Error: %v\n", err)
		return
	}

	fmt.Println(" " + resumen)
	fmt.Println("")

	// Mostrar estructura del MST
	mostrarEstructuraMST(grafoMST)

	// Comparación final
	mostrarComparacion(grafo, grafoMST)

	fmt.Println("\n DEMOSTRACIÓN COMPLETADA")
	fmt.Println("El algoritmo de Kruskal ha calculado exitosamente el MST")
	fmt.Println("que representa las conexiones mínimas necesarias para")
	fmt.Println("mantener toda la red de cuevas conectada.")
}

func mostrarEstructuraGrafo(grafo *domain.Grafo) {
	fmt.Println("  ESTRUCTURA DEL GRAFO ORIGINAL")
	fmt.Println(strings.Repeat("=", 40))

	fmt.Println(" Cuevas:")
	for id, cueva := range grafo.Cuevas {
		fmt.Printf("   • %s: %s\n", id, cueva.Nombre)
	}

	fmt.Println("\n Conexiones:")
	pesoTotal := 0.0
	for i, arista := range grafo.Aristas {
		direccion := "↔"
		if arista.EsDirigido {
			direccion = "→"
		}
		estado := ""
		if arista.EsObstruido {
			estado = " (OBSTRUIDA)"
		}
		fmt.Printf("   %d. %s %s %s (%.2f)%s\n",
			i+1, arista.Desde, direccion, arista.Hasta, arista.Distancia, estado)
		if !arista.EsObstruido {
			pesoTotal += arista.Distancia
		}
	}

	fmt.Printf("\n Peso total de la red: %.2f\n", pesoTotal)
	fmt.Println("")
}

func mostrarEstructuraMST(grafoMST *domain.Grafo) {
	fmt.Println(" ESTRUCTURA DEL MST CALCULADO")
	fmt.Println(strings.Repeat("=", 40))

	fmt.Printf(" Cuevas conectadas: %d\n", len(grafoMST.Cuevas))

	fmt.Println(" Conexiones mínimas requeridas:")
	pesoMST := 0.0
	for i, arista := range grafoMST.Aristas {
		fmt.Printf("   %d. %s ↔ %s (%.2f)\n",
			i+1, arista.Desde, arista.Hasta, arista.Distancia)
		pesoMST += arista.Distancia
	}

	fmt.Printf("\n Peso total del MST: %.2f\n", pesoMST)
	fmt.Println("")
}

func mostrarComparacion(grafoOriginal, grafoMST *domain.Grafo) {
	fmt.Println(" COMPARACIÓN: RED ORIGINAL vs MST")
	fmt.Println(strings.Repeat("=", 50))

	// Calcular pesos
	pesoOriginal := 0.0
	conexionesOriginales := 0
	for _, arista := range grafoOriginal.Aristas {
		if !arista.EsObstruido {
			pesoOriginal += arista.Distancia
			conexionesOriginales++
		}
	}

	pesoMST := 0.0
	for _, arista := range grafoMST.Aristas {
		pesoMST += arista.Distancia
	}

	ahorro := pesoOriginal - pesoMST
	porcentajeAhorro := (ahorro / pesoOriginal) * 100
	conexionesEliminadas := conexionesOriginales - len(grafoMST.Aristas)

	fmt.Printf(" Red original:\n")
	fmt.Printf("   • Conexiones: %d\n", conexionesOriginales)
	fmt.Printf("   • Peso total: %.2f\n", pesoOriginal)

	fmt.Printf("\n MST optimizado:\n")
	fmt.Printf("   • Conexiones: %d\n", len(grafoMST.Aristas))
	fmt.Printf("   • Peso total: %.2f\n", pesoMST)

	fmt.Printf("\n Optimización lograda:\n")
	fmt.Printf("   • Ahorro total: %.2f\n", ahorro)
	fmt.Printf("   • Porcentaje de optimización: %.2f%%\n", porcentajeAhorro)
	fmt.Printf("   • Conexiones eliminadas: %d\n", conexionesEliminadas)

	fmt.Printf("\n Interpretación:\n")
	fmt.Printf("   El MST mantiene la conectividad completa de la red\n")
	fmt.Printf("   eliminando %d conexiones redundantes y reduciendo\n", conexionesEliminadas)
	fmt.Printf("   el costo total en %.2f%% (%.2f unidades)\n", porcentajeAhorro, ahorro)
}
