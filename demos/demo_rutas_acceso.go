package main

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/handler"
	"proyecto-grafos-go/internal/repository"
	"proyecto-grafos-go/internal/service"
	"strings"
)

// main permite ejecutar este demo individualmente
func main() {
	fmt.Println("=== DEMO: RUTAS DE ACCESO MÍNIMAS EN ORDEN DE CREACIÓN ===")

	// Crear un grafo de ejemplo
	grafo := domain.NuevoGrafo(false)

	// Agregar cuevas en orden de creación
	cuevas := []struct {
		id     string
		nombre string
	}{
		{"C1", "Entrada Principal"},
		{"C2", "Sala de Recursos"},
		{"C3", "Depósito Norte"},
		{"C4", "Sala de Herramientas"},
		{"C5", "Depósito Sur"},
	}

	for _, c := range cuevas {
		cueva := domain.NuevaCueva(c.id, c.nombre)
		grafo.AgregarCueva(cueva)
	}

	// Agregar conexiones con diferentes distancias
	conexiones := []struct {
		desde, hasta string
		distancia    float64
	}{
		{"C1", "C2", 10.0},
		{"C1", "C3", 15.0},
		{"C2", "C3", 8.0},
		{"C2", "C4", 12.0},
		{"C3", "C4", 6.0},
		{"C3", "C5", 20.0},
		{"C4", "C5", 9.0},
	}

	for _, conn := range conexiones {
		arista := &domain.Arista{
			Desde:       conn.desde,
			Hasta:       conn.hasta,
			Distancia:   conn.distancia,
			EsDirigido:  false,
			EsObstruido: false,
		}
		grafo.AgregarArista(arista)
	} // Crear servicios
	repositorio := repository.NuevoRepositorio("data")
	servicioGrafo := service.NuevoServicioGrafo(grafo, repositorio)

	mstService := service.NuevoMSTService(servicioGrafo)
	analysisHandler := handler.NuevoAnalysisHandler(mstService)

	// Ejecutar la funcionalidad del requisito 3c
	fmt.Println("\n1. Calculando rutas de acceso mínimas en orden de creación...")

	resultado, err := analysisHandler.CalcularMSTEnOrdenCreacion(grafo)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(resultado)

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("DEMO COMPLETADO EXITOSAMENTE")
	fmt.Println(strings.Repeat("=", 70))
}
