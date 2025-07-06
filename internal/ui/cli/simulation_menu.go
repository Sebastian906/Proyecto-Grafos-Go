package cli

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/handler"
	"proyecto-grafos-go/internal/service"
	"strconv"
	"strings"
)

// SimulationMenu maneja el menú de simulación de camiones
type SimulationMenu struct {
	simulationHandler *handler.SimulationHandler
	traversalHandler  *handler.TraversalHandler
	grafo             *domain.Grafo
}

// NuevoSimulationMenu crea una nueva instancia del menú de simulación
func NuevoSimulationMenu(simulationHandler *handler.SimulationHandler, traversalHandler *handler.TraversalHandler, grafo *domain.Grafo) *SimulationMenu {
	return &SimulationMenu{
		simulationHandler: simulationHandler,
		traversalHandler:  traversalHandler,
		grafo:             grafo,
	}
}

// MostrarMenu muestra el menú principal de simulación
func (sm *SimulationMenu) MostrarMenu() {
	for {
		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Println("MENU DE SIMULACION DE CAMIONES")
		fmt.Println(strings.Repeat("=", 60))
		fmt.Println("1. Crear nuevo camión")
		fmt.Println("2. Cargar insumos en camión")
		fmt.Println("3. Simular entrega con DFS")
		fmt.Println("4. Simular entrega con BFS")
		fmt.Println("5. Comparar algoritmos DFS vs BFS")
		fmt.Println("6. Analizar recorridos (solo navegación)")
		fmt.Println("7. Ver estado de camiones")
		fmt.Println("8. Gestionar camiones")
		fmt.Println("9. Análisis de conectividad")
		fmt.Println("0. Volver al menú principal")
		fmt.Println(strings.Repeat("-", 60))

		opcion := LeerEntrada("Seleccione una opción: ")

		switch opcion {
		case "1":
			sm.crearCamion()
		case "2":
			sm.cargarInsumos()
		case "3":
			sm.simularEntregaDFS()
		case "4":
			sm.simularEntregaBFS()
		case "5":
			sm.compararAlgoritmos()
		case "6":
			sm.analizarRecorridos()
		case "7":
			sm.verEstadoCamiones()
		case "8":
			sm.gestionarCamiones()
		case "9":
			sm.analizarConectividad()
		case "0":
			fmt.Println("Volviendo al menú principal...")
			return
		default:
			fmt.Println("ERROR: Opción no válida. Intente nuevamente.")
		}
	}
}

// crearCamion maneja la creación de un nuevo camión
func (sm *SimulationMenu) crearCamion() {
	fmt.Println("\nCREAR NUEVO CAMION")
	fmt.Println(strings.Repeat("-", 40))

	id := LeerEntrada("ID del camión: ")
	if id == "" {
		fmt.Println("ERROR: El ID del camión no puede estar vacío")
		return
	}

	fmt.Println("Tipos de camión disponibles:")
	fmt.Println("1. PEQUEÑO (100 unidades, 60 km/h)")
	fmt.Println("2. MEDIANO (200 unidades, 50 km/h)")
	fmt.Println("3. GRANDE (400 unidades, 40 km/h)")

	tipoOpcion := LeerEntrada("Seleccione tipo (1-3): ")
	var tipoCamion string

	switch tipoOpcion {
	case "1":
		tipoCamion = "PEQUEÑO"
	case "2":
		tipoCamion = "MEDIANO"
	case "3":
		tipoCamion = "GRANDE"
	default:
		fmt.Println("ERROR: Tipo de camión no válido")
		return
	}

	// Mostrar cuevas disponibles
	fmt.Println("\nCuevas disponibles:")
	contador := 1
	cuevasDisponibles := make([]string, 0)
	for cuevaID := range sm.grafo.Cuevas {
		fmt.Printf("%d. %s\n", contador, cuevaID)
		cuevasDisponibles = append(cuevasDisponibles, cuevaID)
		contador++
	}

	cuevaOrigen := LeerEntrada("Cueva origen del camión: ")
	if cuevaOrigen == "" {
		fmt.Println("ERROR: La cueva origen no puede estar vacía")
		return
	}

	// Crear camión
	camion, err := sm.simulationHandler.CrearCamion(id, tipoCamion, cuevaOrigen)
	if err != nil {
		fmt.Printf("ERROR: Error al crear camión: %s\n", err.Error())
		return
	}

	fmt.Printf("EXITO: Camión '%s' creado exitosamente\n", camion.ID)
	fmt.Printf("   Tipo: %s\n", camion.Tipo)
	fmt.Printf("   Capacidad: %d unidades\n", camion.CapacidadMaxima)
	fmt.Printf("   Velocidad: %.1f km/h\n", camion.VelocidadPromedio)
	fmt.Printf("   Ubicación inicial: %s\n", camion.CuevaActual)
}

// cargarInsumos maneja la carga de insumos en un camión
func (sm *SimulationMenu) cargarInsumos() {
	fmt.Println("\nCARGAR INSUMOS EN CAMION")
	fmt.Println(strings.Repeat("-", 40))

	// Mostrar camiones disponibles
	camiones := sm.simulationHandler.ListarTodosLosCamiones()
	if len(camiones) == 0 {
		fmt.Println("ERROR: No hay camiones disponibles. Cree uno primero.")
		return
	}

	fmt.Println("Camiones disponibles:")
	for id, camion := range camiones {
		fmt.Printf("- %s (%s, Estado: %s)\n", id, camion.Tipo, camion.Estado)
	}

	camionID := LeerEntrada("ID del camión: ")
	if camionID == "" {
		fmt.Println("ERROR: El ID del camión no puede estar vacío")
		return
	}

	insumos := make(map[string]string)
	fmt.Println("\nIngrese los insumos (presione Enter sin escribir para terminar):")

	for {
		recurso := LeerEntrada("Tipo de recurso: ")
		if recurso == "" {
			break
		}

		cantidadStr := LeerEntrada(fmt.Sprintf("Cantidad de %s: ", recurso))
		if cantidadStr == "" {
			fmt.Println("ERROR: La cantidad no puede estar vacía")
			continue
		}

		// Validar que sea un número
		if _, err := strconv.Atoi(cantidadStr); err != nil {
			fmt.Println("ERROR: La cantidad debe ser un número válido")
			continue
		}

		insumos[recurso] = cantidadStr
	}

	if len(insumos) == 0 {
		fmt.Println("ERROR: No se especificaron insumos")
		return
	}

	err := sm.simulationHandler.CargarInsumosEnCamion(camionID, insumos)
	if err != nil {
		fmt.Printf("ERROR: Error al cargar insumos: %s\n", err.Error())
		return
	}

	fmt.Printf("EXITO: Insumos cargados exitosamente en el camión '%s'\n", camionID)
	for recurso, cantidad := range insumos {
		fmt.Printf("   - %s: %s unidades\n", recurso, cantidad)
	}
}

// simularEntregaDFS ejecuta una simulación con DFS
func (sm *SimulationMenu) simularEntregaDFS() {
	fmt.Println("\nSIMULACION DE ENTREGA CON DFS")
	fmt.Println(strings.Repeat("-", 40))

	camionID, cuevaOrigen := sm.obtenerParametrosSimulacion()
	if camionID == "" || cuevaOrigen == "" {
		return
	}

	resultado, err := sm.simulationHandler.EjecutarSimulacionDFS(sm.grafo, camionID, cuevaOrigen)
	if err != nil {
		fmt.Printf("ERROR: Error en simulación: %s\n", err.Error())
		return
	}

	sm.mostrarResultadoSimulacion(resultado)
}

// simularEntregaBFS ejecuta una simulación con BFS
func (sm *SimulationMenu) simularEntregaBFS() {
	fmt.Println("\nSIMULACION DE ENTREGA CON BFS")
	fmt.Println(strings.Repeat("-", 40))

	camionID, cuevaOrigen := sm.obtenerParametrosSimulacion()
	if camionID == "" || cuevaOrigen == "" {
		return
	}

	resultado, err := sm.simulationHandler.EjecutarSimulacionBFS(sm.grafo, camionID, cuevaOrigen)
	if err != nil {
		fmt.Printf("ERROR: Error en simulación: %s\n", err.Error())
		return
	}

	sm.mostrarResultadoSimulacion(resultado)
}

// compararAlgoritmos compara DFS vs BFS para la misma simulación
func (sm *SimulationMenu) compararAlgoritmos() {
	fmt.Println("\nCOMPARACION DFS vs BFS")
	fmt.Println(strings.Repeat("-", 40))

	camionID, cuevaOrigen := sm.obtenerParametrosSimulacion()
	if camionID == "" || cuevaOrigen == "" {
		return
	}

	fmt.Println("Ejecutando simulaciones...")

	resultados, err := sm.simulationHandler.CompararAlgoritmos(sm.grafo, camionID, cuevaOrigen)
	if err != nil {
		fmt.Printf("ERROR: Error en comparación: %s\n", err.Error())
		return
	}

	reporte := sm.simulationHandler.GenerarReporteComparativo(resultados)
	fmt.Println(reporte)
}

// analizarRecorridos analiza recorridos sin simulación de camiones
func (sm *SimulationMenu) analizarRecorridos() {
	fmt.Println("\nANALISIS DE RECORRIDOS")
	fmt.Println(strings.Repeat("-", 40))

	cuevaOrigen := sm.seleccionarCuevaOrigen()
	if cuevaOrigen == "" {
		return
	}

	fmt.Println("Analizando recorridos...")

	resultados, err := sm.traversalHandler.CompararRecorridos(sm.grafo, cuevaOrigen)
	if err != nil {
		fmt.Printf("ERROR: Error en análisis: %s\n", err.Error())
		return
	}

	reporte := sm.traversalHandler.GenerarReporteComparativoRecorridos(resultados)
	fmt.Println(reporte)
}

// verEstadoCamiones muestra el estado de todos los camiones
func (sm *SimulationMenu) verEstadoCamiones() {
	fmt.Println("\nESTADO DE CAMIONES")
	fmt.Println(strings.Repeat("-", 40))

	camiones := sm.simulationHandler.ListarTodosLosCamiones()
	if len(camiones) == 0 {
		fmt.Println("ERROR: No hay camiones registrados")
		return
	}

	for id, camion := range camiones {
		fmt.Printf("\nCamión: %s\n", id)
		fmt.Printf("   Tipo: %s\n", camion.Tipo)
		fmt.Printf("   Estado: %s\n", camion.Estado)
		fmt.Printf("   Capacidad: %d unidades\n", camion.CapacidadMaxima)
		fmt.Printf("   Velocidad: %.1f km/h\n", camion.VelocidadPromedio)
		fmt.Printf("   Ubicación actual: %s\n", camion.CuevaActual)
		fmt.Printf("   Distancia recorrida: %.2f km\n", camion.DistanciaRecorrida)

		if len(camion.CargaActual) > 0 {
			fmt.Printf("   Carga actual:\n")
			for recurso, cantidad := range camion.CargaActual {
				fmt.Printf("     - %s: %d unidades\n", recurso, cantidad)
			}
		} else {
			fmt.Printf("   Sin carga\n")
		}

		if camion.RutaAsignada != nil {
			fmt.Printf("   Ruta asignada: %s (%.2f km)\n",
				camion.RutaAsignada.ID, camion.RutaAsignada.DistanciaTotal)
		}
	}
}

// gestionarCamiones submenu para gestionar camiones
func (sm *SimulationMenu) gestionarCamiones() {
	for {
		fmt.Println("\nGESTION DE CAMIONES")
		fmt.Println(strings.Repeat("-", 40))
		fmt.Println("1. Reiniciar camión")
		fmt.Println("2. Eliminar camión")
		fmt.Println("3. Ver detalles de camión")
		fmt.Println("0. Volver")

		opcion := LeerEntrada("Seleccione una opción: ")

		switch opcion {
		case "1":
			sm.reiniciarCamion()
		case "2":
			sm.eliminarCamion()
		case "3":
			sm.verDetallesCamion()
		case "0":
			return
		default:
			fmt.Println("ERROR: Opción no válida")
		}
	}
}

// analizarConectividad analiza la conectividad del grafo
func (sm *SimulationMenu) analizarConectividad() {
	fmt.Println("\nANALISIS DE CONECTIVIDAD")
	fmt.Println(strings.Repeat("-", 40))

	cuevaOrigen := sm.seleccionarCuevaOrigen()
	if cuevaOrigen == "" {
		return
	}

	fmt.Println("Analizando conectividad...")

	analisis, err := sm.traversalHandler.AnalziarConectividad(sm.grafo, cuevaOrigen)
	if err != nil {
		fmt.Printf("ERROR: Error en análisis: %s\n", err.Error())
		return
	}

	reporte := sm.traversalHandler.GenerarReporteConectividad(analisis)
	fmt.Println(reporte)
}

// Métodos auxiliares

func (sm *SimulationMenu) obtenerParametrosSimulacion() (string, string) {
	// Verificar que hay camiones
	camiones := sm.simulationHandler.ListarTodosLosCamiones()
	if len(camiones) == 0 {
		fmt.Println("ERROR: No hay camiones disponibles. Cree uno primero.")
		return "", ""
	}

	// Mostrar camiones disponibles
	fmt.Println("Camiones disponibles:")
	for id, camion := range camiones {
		fmt.Printf("- %s (%s, Estado: %s)\n", id, camion.Tipo, camion.Estado)
	}

	camionID := LeerEntrada("ID del camión: ")
	if camionID == "" {
		fmt.Println("ERROR: El ID del camión no puede estar vacío")
		return "", ""
	}

	cuevaOrigen := sm.seleccionarCuevaOrigen()
	return camionID, cuevaOrigen
}

func (sm *SimulationMenu) seleccionarCuevaOrigen() string {
	if len(sm.grafo.Cuevas) == 0 {
		fmt.Println("ERROR: No hay cuevas en el grafo")
		return ""
	}

	fmt.Println("Cuevas disponibles:")
	contador := 1
	for cuevaID := range sm.grafo.Cuevas {
		fmt.Printf("%d. %s\n", contador, cuevaID)
		contador++
	}

	cuevaOrigen := LeerEntrada("Cueva origen: ")
	if cuevaOrigen == "" {
		fmt.Println("ERROR: La cueva origen no puede estar vacía")
		return ""
	}

	if _, existe := sm.grafo.ObtenerCueva(cuevaOrigen); !existe {
		fmt.Printf("ERROR: La cueva '%s' no existe\n", cuevaOrigen)
		return ""
	}

	return cuevaOrigen
}

func (sm *SimulationMenu) mostrarResultadoSimulacion(resultado *service.SimulacionResultado) {
	reporte := sm.simulationHandler.GenerarReporteSimulacion(resultado)
	fmt.Println(reporte)

	if resultado.Exitoso {
		fmt.Println("EXITO: Simulación completada exitosamente")
	} else {
		fmt.Println("ERROR: La simulación tuvo problemas")
	}
}

func (sm *SimulationMenu) reiniciarCamion() {
	camiones := sm.simulationHandler.ListarTodosLosCamiones()
	if len(camiones) == 0 {
		fmt.Println("ERROR: No hay camiones disponibles")
		return
	}

	fmt.Println("Camiones disponibles:")
	for id := range camiones {
		fmt.Printf("- %s\n", id)
	}

	camionID := LeerEntrada("ID del camión a reiniciar: ")
	if camionID == "" {
		fmt.Println("ERROR: El ID del camión no puede estar vacío")
		return
	}

	cuevaOrigen := sm.seleccionarCuevaOrigen()
	if cuevaOrigen == "" {
		return
	}

	err := sm.simulationHandler.ReiniciarCamion(camionID, cuevaOrigen)
	if err != nil {
		fmt.Printf("ERROR: Error al reiniciar camión: %s\n", err.Error())
		return
	}

	fmt.Printf("EXITO: Camión '%s' reiniciado exitosamente\n", camionID)
}

func (sm *SimulationMenu) eliminarCamion() {
	camiones := sm.simulationHandler.ListarTodosLosCamiones()
	if len(camiones) == 0 {
		fmt.Println("ERROR: No hay camiones disponibles")
		return
	}

	fmt.Println("Camiones disponibles:")
	for id := range camiones {
		fmt.Printf("- %s\n", id)
	}

	camionID := LeerEntrada("ID del camión a eliminar: ")
	if camionID == "" {
		fmt.Println("ERROR: El ID del camión no puede estar vacío")
		return
	}

	confirmacion := LeerEntrada(fmt.Sprintf("¿Está seguro de eliminar el camión '%s'? (s/N): ", camionID))
	if strings.ToLower(confirmacion) != "s" && strings.ToLower(confirmacion) != "si" {
		fmt.Println("Operación cancelada")
		return
	}

	err := sm.simulationHandler.EliminarCamion(camionID)
	if err != nil {
		fmt.Printf("ERROR: Error al eliminar camión: %s\n", err.Error())
		return
	}

	fmt.Printf("EXITO: Camión '%s' eliminado exitosamente\n", camionID)
}

func (sm *SimulationMenu) verDetallesCamion() {
	camiones := sm.simulationHandler.ListarTodosLosCamiones()
	if len(camiones) == 0 {
		fmt.Println("ERROR: No hay camiones disponibles")
		return
	}

	fmt.Println("Camiones disponibles:")
	for id := range camiones {
		fmt.Printf("- %s\n", id)
	}

	camionID := LeerEntrada("ID del camión: ")
	if camionID == "" {
		fmt.Println("ERROR: El ID del camión no puede estar vacío")
		return
	}

	camion, err := sm.simulationHandler.ObtenerEstadoCamion(camionID)
	if err != nil {
		fmt.Printf("ERROR: Error: %s\n", err.Error())
		return
	}

	fmt.Printf("\nDETALLES DEL CAMIÓN '%s'\n", camion.ID)
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("Tipo: %s\n", camion.Tipo)
	fmt.Printf("Estado: %s\n", camion.Estado)
	fmt.Printf("Capacidad máxima: %d unidades\n", camion.CapacidadMaxima)
	fmt.Printf("Velocidad promedio: %.1f km/h\n", camion.VelocidadPromedio)
	fmt.Printf("Ubicación actual: %s\n", camion.CuevaActual)
	fmt.Printf("Distancia recorrida: %.2f km\n", camion.DistanciaRecorrida)

	if !camion.TiempoInicio.IsZero() {
		fmt.Printf("Tiempo de inicio: %s\n", camion.TiempoInicio.Format("15:04:05"))
	}
	if !camion.TiempoFin.IsZero() {
		fmt.Printf("Tiempo de fin: %s\n", camion.TiempoFin.Format("15:04:05"))
	}

	if len(camion.CargaActual) > 0 {
		fmt.Printf("\nCarga actual:\n")
		for recurso, cantidad := range camion.CargaActual {
			fmt.Printf("- %s: %d unidades\n", recurso, cantidad)
		}
	} else {
		fmt.Printf("\nSin carga actual\n")
	}

	if camion.RutaAsignada != nil {
		fmt.Printf("\nRuta asignada:\n")
		fmt.Printf("ID: %s\n", camion.RutaAsignada.ID)
		fmt.Printf("Distancia total: %.2f km\n", camion.RutaAsignada.DistanciaTotal)
		fmt.Printf("Cuevas en ruta: %v\n", camion.RutaAsignada.CuevaIDs)
		fmt.Printf("Completado: %t\n", camion.RutaAsignada.EstaCompleto)
	}
}
