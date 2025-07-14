package unit

import (
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/service"
	"testing"
)

// TestTraversalServiceDFS prueba el algoritmo DFS
func TestTraversalServiceDFS(t *testing.T) {
	// Crear un grafo simple de prueba
	grafo := domain.NuevoGrafo(false) // No dirigido

	// Agregar cuevas
	cueva1 := domain.NuevaCueva("1", "Cueva A")
	cueva2 := domain.NuevaCueva("2", "Cueva B")
	cueva3 := domain.NuevaCueva("3", "Cueva C")

	grafo.AgregarCueva(cueva1)
	grafo.AgregarCueva(cueva2)
	grafo.AgregarCueva(cueva3)

	// Agregar aristas
	arista1 := domain.NuevaArista("1", "2", 10.0, false)
	arista2 := domain.NuevaArista("2", "3", 15.0, false)

	grafo.AgregarArista(arista1)
	grafo.AgregarArista(arista2)

	// Crear servicios
	graphService := service.NuevoServicioGrafo(grafo, nil)
	traversalService := service.NuevoTraversalService(graphService)

	// Ejecutar DFS
	resultado, err := traversalService.RealizarRecorridoDFS(grafo, "1")
	if err != nil {
		t.Fatalf("Error en DFS: %s", err.Error())
	}

	// Verificar resultado
	if resultado == nil {
		t.Fatal("Resultado no debe ser nulo")
	}

	if resultado.TipoRecorrido != service.DFS {
		t.Errorf("Tipo de recorrido esperado: %s, obtenido: %s", service.DFS, resultado.TipoRecorrido)
	}

	if resultado.CuevaOrigen != "1" {
		t.Errorf("Cueva origen esperada: 1, obtenida: %s", resultado.CuevaOrigen)
	}

	if len(resultado.CuevasVisitas) == 0 {
		t.Error("Debería haber visitado al menos una cueva")
	}

	t.Logf("DFS completado exitosamente. Cuevas visitadas: %v", resultado.CuevasVisitas)
}

// TestTraversalServiceBFS prueba el algoritmo BFS
func TestTraversalServiceBFS(t *testing.T) {
	// Crear un grafo simple de prueba
	grafo := domain.NuevoGrafo(false) // No dirigido

	// Agregar cuevas
	cueva1 := domain.NuevaCueva("1", "Cueva A")
	cueva2 := domain.NuevaCueva("2", "Cueva B")
	cueva3 := domain.NuevaCueva("3", "Cueva C")

	grafo.AgregarCueva(cueva1)
	grafo.AgregarCueva(cueva2)
	grafo.AgregarCueva(cueva3)

	// Agregar aristas
	arista1 := domain.NuevaArista("1", "2", 10.0, false)
	arista2 := domain.NuevaArista("2", "3", 15.0, false)

	grafo.AgregarArista(arista1)
	grafo.AgregarArista(arista2)

	// Crear servicios
	graphService := service.NuevoServicioGrafo(grafo, nil)
	traversalService := service.NuevoTraversalService(graphService)

	// Ejecutar BFS
	resultado, err := traversalService.RealizarRecorridoBFS(grafo, "1")
	if err != nil {
		t.Fatalf("Error en BFS: %s", err.Error())
	}

	// Verificar resultado
	if resultado == nil {
		t.Fatal("Resultado no debe ser nulo")
	}

	if resultado.TipoRecorrido != service.BFS {
		t.Errorf("Tipo de recorrido esperado: %s, obtenido: %s", service.BFS, resultado.TipoRecorrido)
	}

	if resultado.CuevaOrigen != "1" {
		t.Errorf("Cueva origen esperada: 1, obtenida: %s", resultado.CuevaOrigen)
	}

	if len(resultado.CuevasVisitas) == 0 {
		t.Error("Debería haber visitado al menos una cueva")
	}

	t.Logf("BFS completado exitosamente. Cuevas visitadas: %v", resultado.CuevasVisitas)
}

// TestTruckService prueba la creación de camiones y simulación básica
func TestTruckService(t *testing.T) {
	// Crear un grafo simple de prueba
	grafo := domain.NuevoGrafo(false)

	// Agregar cuevas
	cueva1 := domain.NuevaCueva("1", "Almacén")
	cueva2 := domain.NuevaCueva("2", "Destino")

	grafo.AgregarCueva(cueva1)
	grafo.AgregarCueva(cueva2)

	// Agregar arista
	arista := domain.NuevaArista("1", "2", 20.0, false)
	grafo.AgregarArista(arista)

	// Crear servicios
	graphService := service.NuevoServicioGrafo(grafo, nil)
	traversalService := service.NuevoTraversalService(graphService)
	truckService := service.NuevoTruckService(traversalService, graphService)

	// Crear camión
	camion, err := truckService.CrearCamion("camion1", service.CamionPequeno, "1")
	if err != nil {
		t.Fatalf("Error al crear camión: %s", err.Error())
	}

	// Verificar camión
	if camion.ID != "camion1" {
		t.Errorf("ID esperado: camion1, obtenido: %s", camion.ID)
	}

	if camion.Tipo != service.CamionPequeno {
		t.Errorf("Tipo esperado: %s, obtenido: %s", service.CamionPequeno, camion.Tipo)
	}

	// Cargar insumos
	insumos := map[string]int{
		"alimentos": 50,
		"medicinas": 30,
	}

	err = truckService.CargarInsumos("camion1", insumos)
	if err != nil {
		t.Fatalf("Error al cargar insumos: %s", err.Error())
	}

	// Verificar carga
	if camion.CargaActual["alimentos"] != 50 {
		t.Errorf("Alimentos esperados: 50, obtenidos: %d", camion.CargaActual["alimentos"])
	}

	// Simular entrega con DFS
	resultado, err := truckService.SimularEntregaDFS(grafo, "camion1", "1")
	if err != nil {
		t.Fatalf("Error en simulación DFS: %s", err.Error())
	}

	if resultado.CamionID != "camion1" {
		t.Errorf("Camión ID esperado: camion1, obtenido: %s", resultado.CamionID)
	}

	if resultado.TipoRecorrido != service.DFS {
		t.Errorf("Tipo de recorrido esperado: %s, obtenido: %s", service.DFS, resultado.TipoRecorrido)
	}

	t.Logf("Simulación completada. Exitoso: %t, Distancia: %.2f", resultado.Exitoso, resultado.DistanciaTotal)
}

// TestConectividad prueba el análisis de conectividad
func TestConectividad(t *testing.T) {
	// Crear un grafo con cuevas conectadas y desconectadas
	grafo := domain.NuevoGrafo(false)

	// Grupo conectado
	cueva1 := domain.NuevaCueva("1", "Cueva 1")
	cueva2 := domain.NuevaCueva("2", "Cueva 2")
	// Cueva aislada
	cueva3 := domain.NuevaCueva("3", "Cueva Aislada")

	grafo.AgregarCueva(cueva1)
	grafo.AgregarCueva(cueva2)
	grafo.AgregarCueva(cueva3)

	// Solo conectar 1 y 2
	arista := domain.NuevaArista("1", "2", 10.0, false)
	grafo.AgregarArista(arista)

	// Crear servicio
	graphService := service.NuevoServicioGrafo(grafo, nil)
	traversalService := service.NuevoTraversalService(graphService)

	// Verificar conectividad desde cueva 1
	esConectado, cuevasNoAccesibles, err := traversalService.VerificarConectividad(grafo, "1")
	if err != nil {
		t.Fatalf("Error verificando conectividad: %s", err.Error())
	}

	// No debería ser completamente conectado (cueva 3 está aislada)
	if esConectado {
		t.Error("El grafo no debería ser completamente conectado")
	}

	// Debería haber una cueva no accesible
	if len(cuevasNoAccesibles) != 1 {
		t.Errorf("Esperadas 1 cueva no accesible, obtenidas: %d", len(cuevasNoAccesibles))
	}

	if len(cuevasNoAccesibles) > 0 && cuevasNoAccesibles[0] != "3" {
		t.Errorf("Cueva no accesible esperada: 3, obtenida: %s", cuevasNoAccesibles[0])
	}

	t.Logf("Conectividad verificada correctamente. Cuevas no accesibles: %v", cuevasNoAccesibles)
}
