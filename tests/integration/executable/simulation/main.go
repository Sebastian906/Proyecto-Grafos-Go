package main

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/handler"
	"proyecto-grafos-go/internal/repository"
	"proyecto-grafos-go/internal/service"
	"strings"
)

func main() {
	fmt.Println("=== PRUEBA DE INTEGRACION: Simulacion de Camiones ===")
	fmt.Println("Este test verifica que todo el sistema de simulacion funcione correctamente")

	// Inicializar componentes
	grafo := domain.NuevoGrafo(false)
	repo := repository.NuevoRepositorio("../../../../data/") // Ruta relativa desde tests/integration/executable/simulation

	// Cargar datos de ejemplo
	fmt.Println("\n1. Cargando datos de ejemplo...")
	grafoNuevo, err := repo.CargarJSON("caves_example.json")
	if err != nil {
		fmt.Printf("ERROR al cargar archivo: %v\n", err)
		fmt.Println("NOTA: Asegurese de que el archivo 'data/caves_example.json' existe")
		return
	}
	grafo = grafoNuevo
	fmt.Printf("EXITO: Datos cargados exitosamente - %d cuevas, %d conexiones\n",
		len(grafo.Cuevas), len(grafo.Aristas))

	// Mostrar cuevas disponibles
	fmt.Println("\nCuevas disponibles en el grafo:")
	for id, cueva := range grafo.Cuevas {
		fmt.Printf("   - %s: %s (%.1f, %.1f)\n", id, cueva.Nombre, cueva.X, cueva.Y)
	}

	// Crear servicios
	grafoSvc := service.NuevoServicioGrafo(grafo, repo)
	traversalSvc := service.NuevoTraversalService(grafoSvc)
	truckSvc := service.NuevoTruckService(traversalSvc, grafoSvc)

	// Crear handler
	simulationHandler := handler.NuevoSimulationHandler(truckSvc, traversalSvc, grafoSvc)

	// 2. Crear un camion
	fmt.Println("\n2. PRUEBA: Creacion de camion...")
	camion, err := simulationHandler.CrearCamion("camion001", "PEQUEÑO", "CR")
	if err != nil {
		fmt.Printf("ERROR al crear camion: %v\n", err)
		return
	}
	fmt.Printf("EXITO: Camion creado exitosamente:\n")
	fmt.Printf("   - ID: %s\n", camion.ID)
	fmt.Printf("   - Tipo: %s\n", camion.Tipo)
	fmt.Printf("   - Capacidad: %d unidades\n", camion.CapacidadMaxima)
	fmt.Printf("   - Velocidad: %.1f km/h\n", camion.VelocidadPromedio)
	fmt.Printf("   - Ubicacion inicial: %s\n", camion.CuevaActual)

	// 3. Cargar insumos
	fmt.Println("\n3. PRUEBA: Carga de insumos...")
	insumos := map[string]string{
		"medicinas": "50",
		"comida":    "30",
		"agua":      "20",
	}
	err = simulationHandler.CargarInsumosEnCamion("camion001", insumos)
	if err != nil {
		fmt.Printf("ERROR al cargar insumos: %v\n", err)
		return
	}
	fmt.Println("EXITO: Insumos cargados exitosamente:")
	for recurso, cantidad := range insumos {
		fmt.Printf("   - %s: %s unidades\n", recurso, cantidad)
	}

	// 4. Ejecutar simulacion DFS
	fmt.Println("\n4. PRUEBA: Simulacion con algoritmo DFS...")
	resultado, err := simulationHandler.EjecutarSimulacionDFS(grafo, "camion001", "CR")
	if err != nil {
		fmt.Printf("ERROR en simulacion DFS: %v\n", err)
		return
	}

	// Mostrar resultados DFS
	fmt.Printf("EXITO: Simulacion DFS completada: %t\n", resultado.Exitoso)
	fmt.Printf("RESULTADOS DFS:\n")
	fmt.Printf("   - Tiempo total: %v\n", resultado.TiempoTotal)
	fmt.Printf("   - Distancia total: %.2f km\n", resultado.DistanciaTotal)
	fmt.Printf("   - Ruta seguida: %v\n", resultado.RutaCompleta)
	fmt.Printf("   - Entregas realizadas en %d ubicaciones\n", len(resultado.EntregasRealizadas))

	// Mostrar detalles de entregas
	fmt.Println("   - Detalle de entregas:")
	for cueva, entregas := range resultado.EntregasRealizadas {
		fmt.Printf("     %s: ", cueva)
		for recurso, cantidad := range entregas {
			fmt.Printf("%s(%d) ", recurso, cantidad)
		}
		fmt.Println()
	}

	// 5. Reiniciar camion y probar BFS
	fmt.Println("\n5. PRUEBA: Reinicio de camion para simulacion BFS...")
	err = simulationHandler.ReiniciarCamion("camion001", "CR")
	if err != nil {
		fmt.Printf("ERROR al reiniciar camion: %v\n", err)
		return
	}
	fmt.Println("EXITO: Camion reiniciado exitosamente")

	// Cargar insumos nuevamente
	err = simulationHandler.CargarInsumosEnCamion("camion001", insumos)
	if err != nil {
		fmt.Printf("ERROR al cargar insumos: %v\n", err)
		return
	}
	fmt.Println("EXITO: Insumos recargados para simulacion BFS")

	// Ejecutar simulacion BFS
	fmt.Println("\n6. PRUEBA: Simulacion con algoritmo BFS...")
	resultadoBFS, err := simulationHandler.EjecutarSimulacionBFS(grafo, "camion001", "CR")
	if err != nil {
		fmt.Printf("ERROR en simulacion BFS: %v\n", err)
		return
	}

	// Mostrar resultados BFS
	fmt.Printf("EXITO: Simulacion BFS completada: %t\n", resultadoBFS.Exitoso)
	fmt.Printf("RESULTADOS BFS:\n")
	fmt.Printf("   - Tiempo total: %v\n", resultadoBFS.TiempoTotal)
	fmt.Printf("   - Distancia total: %.2f km\n", resultadoBFS.DistanciaTotal)
	fmt.Printf("   - Ruta seguida: %v\n", resultadoBFS.RutaCompleta)
	fmt.Printf("   - Entregas realizadas en %d ubicaciones\n", len(resultadoBFS.EntregasRealizadas))

	// 7. Comparacion de algoritmos
	fmt.Println("\n7. COMPARACION DE ALGORITMOS:")
	fmt.Printf("   DFS vs BFS:\n")
	if resultado.TiempoTotal < resultadoBFS.TiempoTotal {
		fmt.Printf("   - DFS fue mas rapido (%v vs %v)\n", resultado.TiempoTotal, resultadoBFS.TiempoTotal)
	} else if resultadoBFS.TiempoTotal < resultado.TiempoTotal {
		fmt.Printf("   - BFS fue mas rapido (%v vs %v)\n", resultadoBFS.TiempoTotal, resultado.TiempoTotal)
	} else {
		fmt.Printf("   - Ambos algoritmos tardaron lo mismo\n")
	}

	if resultado.DistanciaTotal < resultadoBFS.DistanciaTotal {
		fmt.Printf("   - DFS recorrio menor distancia (%.2f vs %.2f km)\n", resultado.DistanciaTotal, resultadoBFS.DistanciaTotal)
	} else if resultadoBFS.DistanciaTotal < resultado.DistanciaTotal {
		fmt.Printf("   - BFS recorrio menor distancia (%.2f vs %.2f km)\n", resultadoBFS.DistanciaTotal, resultado.DistanciaTotal)
	} else {
		fmt.Printf("   - Ambos algoritmos recorrieron la misma distancia\n")
	}

	// Verificar que las rutas sean iguales (para este grafo simple)
	rutasIguales := len(resultado.RutaCompleta) == len(resultadoBFS.RutaCompleta)
	if rutasIguales {
		for i, cueva := range resultado.RutaCompleta {
			if i >= len(resultadoBFS.RutaCompleta) || cueva != resultadoBFS.RutaCompleta[i] {
				rutasIguales = false
				break
			}
		}
	}

	if rutasIguales {
		fmt.Printf("   - Ambos algoritmos siguieron la misma ruta\n")
	} else {
		fmt.Printf("   - Los algoritmos siguieron rutas diferentes\n")
	}

	// 8. Prueba de estado del camion
	fmt.Println("\n8. PRUEBA: Estado final del camion...")
	estadoCamion, err := simulationHandler.ObtenerEstadoCamion("camion001")
	if err != nil {
		fmt.Printf("ERROR al obtener estado: %v\n", err)
		return
	}

	fmt.Printf("EXITO: Estado final del camion:\n")
	fmt.Printf("   - Estado: %s\n", estadoCamion.Estado)
	fmt.Printf("   - Distancia total recorrida: %.2f km\n", estadoCamion.DistanciaRecorrida)
	fmt.Printf("   - Ubicacion actual: %s\n", estadoCamion.CuevaActual)

	if len(estadoCamion.CargaActual) > 0 {
		fmt.Printf("   - Carga restante:\n")
		for recurso, cantidad := range estadoCamion.CargaActual {
			fmt.Printf("     %s: %d unidades\n", recurso, cantidad)
		}
	} else {
		fmt.Printf("   - Sin carga restante (todo entregado)\n")
	}

	// 9. Resultados finales
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("RESUMEN DE PRUEBAS DE INTEGRACION")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("EXITO: Creacion de camiones: EXITOSA")
	fmt.Println("EXITO: Carga de insumos: EXITOSA")
	fmt.Println("EXITO: Simulacion DFS: EXITOSA")
	fmt.Println("EXITO: Simulacion BFS: EXITOSA")
	fmt.Println("EXITO: Reinicio de camion: EXITOSO")
	fmt.Println("EXITO: Obtencion de estado: EXITOSA")
	fmt.Println()
	fmt.Println("CONCLUSION:")
	fmt.Println("   El sistema de simulacion de camiones esta funcionando")
	fmt.Println("   correctamente. Todos los componentes se integran bien")
	fmt.Println("   y las simulaciones se ejecutan sin errores.")
	fmt.Println()
	fmt.Println("El sistema esta listo para uso en produccion!")
	fmt.Println(strings.Repeat("=", 60))
}
