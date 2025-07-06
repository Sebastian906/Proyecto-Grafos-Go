package main

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/handler"
	"proyecto-grafos-go/internal/repository"
	"proyecto-grafos-go/internal/service"
	"proyecto-grafos-go/internal/ui/cli"
	"strings"
)

func main() {
	fmt.Println("Sistema de Gestión de Cuevas y Simulación de Camiones")
	fmt.Println("=======================================================")

	// Inicialización
	grafo := domain.NuevoGrafo(false)
	repo := repository.NuevoRepositorio("data/")

	// Servicios básicos
	grafoSvc := service.NuevoServicioGrafo(grafo, repo)
	cuevaSvc := service.ServicioNuevaCueva(grafo)
	validacionSvc := service.NuevoServicioValidacion(grafo)
	conexionSvc := service.NuevoServicioConexion(grafo)

	// Nuevos servicios para simulación
	traversalSvc := service.NuevoTraversalService(grafoSvc)
	truckSvc := service.NuevoTruckService(traversalSvc, grafoSvc)

	// Controladores (handlers)
	simulationHandler := handler.NuevoSimulationHandler(truckSvc, traversalSvc, grafoSvc)
	traversalHandler := handler.NuevoTraversalHandler(traversalSvc, grafoSvc)

	// Cargar datos de ejemplo si existe
	if err := grafoSvc.CargarGrafo("caves_example.json"); err != nil {
		fmt.Printf("INFO: No se pudo cargar archivo de ejemplo: %s\n", err.Error())
		fmt.Println("NOTA: Puede crear cuevas manualmente desde el menú")
	} else {
		fmt.Println("EXITO: Datos de ejemplo cargados exitosamente")
	}

	// Crear menús actualizados
	mainMenu := cli.NuevoMainMenu(grafoSvc, cuevaSvc, validacionSvc, conexionSvc)

	// Agregar funcionalidad de simulación
	fmt.Println("\nIniciando interfaz de usuario...")

	// Mostrar menú principal mejorado con opciones de simulación
	mostrarMenuPrincipalMejorado(mainMenu, simulationHandler, traversalHandler, grafo)
}

// mostrarMenuPrincipalMejorado extiende el menú principal con opciones de simulación
func mostrarMenuPrincipalMejorado(mainMenu *cli.MainMenu, simulationHandler *handler.SimulationHandler, traversalHandler *handler.TraversalHandler, grafo *domain.Grafo) {
	for {
		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Println("MENU PRINCIPAL - SISTEMA DE CUEVAS")
		fmt.Println(strings.Repeat("=", 60))
		fmt.Println("1. Gestión de Grafos y Cuevas")
		fmt.Println("2. Simulación de Camiones (NUEVO)")
		fmt.Println("3. Análisis de Recorridos (NUEVO)")
		fmt.Println("4. Información del Sistema")
		fmt.Println("0. Salir")
		fmt.Println(strings.Repeat("=", 60))

		opcion := cli.LeerEntrada("Seleccione una opción: ")

		switch opcion {
		case "1":
			// Usar el menú original
			mainMenu.Mostrar()
		case "2":
			// Nuevo menú de simulación
			simulationMenu := cli.NuevoSimulationMenu(simulationHandler, traversalHandler, grafo)
			simulationMenu.MostrarMenu()
		case "3":
			// Menú de análisis de recorridos
			mostrarMenuAnalisisRecorridos(traversalHandler, grafo)
		case "4":
			mostrarInformacionSistema(grafo)
		case "0":
			fmt.Println("Gracias por usar el sistema! Hasta la vista.")
			return
		default:
			fmt.Println("ERROR: Opción no válida. Intente nuevamente.")
		}
	}
}

// mostrarMenuAnalisisRecorridos muestra opciones para análisis de recorridos
func mostrarMenuAnalisisRecorridos(traversalHandler *handler.TraversalHandler, grafo *domain.Grafo) {
	for {
		fmt.Println("\nANALISIS DE RECORRIDOS")
		fmt.Println(strings.Repeat("=", 40))
		fmt.Println("1. Ejecutar DFS")
		fmt.Println("2. Ejecutar BFS")
		fmt.Println("3. Comparar DFS vs BFS")
		fmt.Println("4. Análisis de conectividad")
		fmt.Println("0. Volver")

		opcion := cli.LeerEntrada("Seleccione una opción: ")

		switch opcion {
		case "1":
			ejecutarDFSSolo(traversalHandler, grafo)
		case "2":
			ejecutarBFSSolo(traversalHandler, grafo)
		case "3":
			compararRecorridos(traversalHandler, grafo)
		case "4":
			analizarConectividad(traversalHandler, grafo)
		case "0":
			return
		default:
			fmt.Println("ERROR: Opción no válida")
		}
	}
}

// mostrarInformacionSistema muestra información del estado actual del sistema
func mostrarInformacionSistema(grafo *domain.Grafo) {
	fmt.Println("\nINFORMACION DEL SISTEMA")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("Total de cuevas: %d\n", len(grafo.Cuevas))
	fmt.Printf("Total de aristas: %d\n", len(grafo.Aristas))
	fmt.Printf("Tipo de grafo: ")
	if grafo.EsDirigido {
		fmt.Println("Dirigido")
	} else {
		fmt.Println("No dirigido")
	}

	fmt.Println("\nCuevas registradas:")
	for id, cueva := range grafo.Cuevas {
		fmt.Printf("- %s: %s\n", id, cueva.Nombre)
	}

	fmt.Println("\nConexiones:")
	for i, arista := range grafo.Aristas {
		estado := "OK"
		if arista.EsObstruido {
			estado = "OBSTRUIDO"
		}
		fmt.Printf("%d. %s -> %s (%.2f km) %s\n", i+1, arista.Desde, arista.Hasta, arista.Distancia, estado)
	}

	cli.LeerEntrada("\nPresione Enter para continuar...")
}

// Funciones auxiliares para análisis de recorridos

func ejecutarDFSSolo(traversalHandler *handler.TraversalHandler, grafo *domain.Grafo) {
	cuevaOrigen := seleccionarCuevaOrigen(grafo)
	if cuevaOrigen == "" {
		return
	}

	fmt.Println("Ejecutando DFS...")
	resultado, err := traversalHandler.EjecutarDFS(grafo, cuevaOrigen)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	reporte := traversalHandler.GenerarReporteRecorrido(resultado)
	fmt.Println(reporte)
	cli.LeerEntrada("\nPresione Enter para continuar...")
}

func ejecutarBFSSolo(traversalHandler *handler.TraversalHandler, grafo *domain.Grafo) {
	cuevaOrigen := seleccionarCuevaOrigen(grafo)
	if cuevaOrigen == "" {
		return
	}

	fmt.Println("Ejecutando BFS...")
	resultado, err := traversalHandler.EjecutarBFS(grafo, cuevaOrigen)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	reporte := traversalHandler.GenerarReporteRecorrido(resultado)
	fmt.Println(reporte)
	cli.LeerEntrada("\nPresione Enter para continuar...")
}

func compararRecorridos(traversalHandler *handler.TraversalHandler, grafo *domain.Grafo) {
	cuevaOrigen := seleccionarCuevaOrigen(grafo)
	if cuevaOrigen == "" {
		return
	}

	fmt.Println("Comparando DFS vs BFS...")
	resultados, err := traversalHandler.CompararRecorridos(grafo, cuevaOrigen)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	reporte := traversalHandler.GenerarReporteComparativoRecorridos(resultados)
	fmt.Println(reporte)
	cli.LeerEntrada("\nPresione Enter para continuar...")
}

func analizarConectividad(traversalHandler *handler.TraversalHandler, grafo *domain.Grafo) {
	cuevaOrigen := seleccionarCuevaOrigen(grafo)
	if cuevaOrigen == "" {
		return
	}

	fmt.Println("Analizando conectividad...")
	analisis, err := traversalHandler.AnalziarConectividad(grafo, cuevaOrigen)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	reporte := traversalHandler.GenerarReporteConectividad(analisis)
	fmt.Println(reporte)
	cli.LeerEntrada("\nPresione Enter para continuar...")
}

func seleccionarCuevaOrigen(grafo *domain.Grafo) string {
	if len(grafo.Cuevas) == 0 {
		fmt.Println("ERROR: No hay cuevas en el grafo")
		return ""
	}

	fmt.Println("Cuevas disponibles:")
	contador := 1
	for cuevaID := range grafo.Cuevas {
		fmt.Printf("%d. %s\n", contador, cuevaID)
		contador++
	}

	cuevaOrigen := cli.LeerEntrada("Cueva origen: ")
	if cuevaOrigen == "" {
		fmt.Println("ERROR: La cueva origen no puede estar vacía")
		return ""
	}

	if _, existe := grafo.ObtenerCueva(cuevaOrigen); !existe {
		fmt.Printf("ERROR: La cueva '%s' no existe\n", cuevaOrigen)
		return ""
	}

	return cuevaOrigen
}
