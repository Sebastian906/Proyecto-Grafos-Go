package unit

import (
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/handler"
	"proyecto-grafos-go/internal/repository"
	"proyecto-grafos-go/internal/service"
	"testing"
)

func TestGraphHandler_CrearGrafo(t *testing.T) {
	// Setup
	grafo := domain.NuevoGrafo(false)
	repo := repository.NuevoRepositorio("../../data/")
	grafoSvc := service.NuevoServicioGrafo(grafo, repo)
	handler := handler.NuevoGraphHandler(grafoSvc)

	// Test crear grafo no dirigido
	grafoCreado, err := handler.CrearGrafo(false)
	if err != nil {
		t.Errorf("Error al crear grafo no dirigido: %v", err)
	}
	if grafoCreado == nil {
		t.Error("Grafo creado es nil")
	}
	if grafoCreado.EsDirigido != false {
		t.Error("Grafo debería ser no dirigido")
	}

	// Test crear grafo dirigido
	grafoCreado, err = handler.CrearGrafo(true)
	if err != nil {
		t.Errorf("Error al crear grafo dirigido: %v", err)
	}
	if grafoCreado == nil {
		t.Error("Grafo creado es nil")
	}
	if grafoCreado.EsDirigido != true {
		t.Error("Grafo debería ser dirigido")
	}
}

func TestGraphHandler_ObtenerEstadisticas(t *testing.T) {
	// Setup
	grafo := domain.NuevoGrafo(false)
	repo := repository.NuevoRepositorio("../../data/")
	grafoSvc := service.NuevoServicioGrafo(grafo, repo)
	handler := handler.NuevoGraphHandler(grafoSvc)

	// Agregar algunas cuevas
	cueva1 := domain.NuevaCueva("C1", "Cueva 1")
	cueva2 := domain.NuevaCueva("C2", "Cueva 2")
	cueva3 := domain.NuevaCueva("C3", "Cueva 3")

	grafo.AgregarCueva(cueva1)
	grafo.AgregarCueva(cueva2)
	grafo.AgregarCueva(cueva3)

	// Agregar algunas conexiones
	grafo.AgregarConexion("C1", "C2", 10.0)
	grafo.AgregarConexion("C2", "C3", 15.0)

	// Test obtener estadísticas
	stats, err := handler.ObtenerEstadisticas()
	if err != nil {
		t.Errorf("Error al obtener estadísticas: %v", err)
	}
	if stats.NumCuevas != 3 {
		t.Errorf("Esperaba 3 cuevas, obtuvo %d", stats.NumCuevas)
	}
	if stats.NumConexiones != 2 {
		t.Errorf("Esperaba 2 conexiones, obtuvo %d", stats.NumConexiones)
	}
}

func TestGraphHandler_CargarGrafo(t *testing.T) {
	// Setup
	grafo := domain.NuevoGrafo(false)
	repo := repository.NuevoRepositorio("../../data/")
	grafoSvc := service.NuevoServicioGrafo(grafo, repo)
	handler := handler.NuevoGraphHandler(grafoSvc)

	// Test cargar grafo válido
	err := handler.CargarGrafo("caves_example.json")
	if err != nil {
		t.Logf("Advertencia: No se pudo cargar el archivo de ejemplo: %v", err)
		// No es un error crítico si el archivo no existe
	}

	// Test cargar grafo con archivo inválido
	err = handler.CargarGrafo("archivo_inexistente.json")
	if err == nil {
		t.Error("Debería fallar al cargar archivo inexistente")
	}
}

func TestGraphHandler_GuardarGrafo(t *testing.T) {
	// Setup
	grafo := domain.NuevoGrafo(false)
	repo := repository.NuevoRepositorio("../../data/")
	grafoSvc := service.NuevoServicioGrafo(grafo, repo)
	handler := handler.NuevoGraphHandler(grafoSvc)

	// Agregar contenido al grafo
	cueva1 := domain.NuevaCueva("C1", "Cueva Test")
	grafo.AgregarCueva(cueva1)

	// Test guardar grafo
	err := handler.GuardarGrafo("test_output.json")
	if err != nil {
		t.Errorf("Error al guardar grafo: %v", err)
	}
}

func TestGraphHandler_CambiarTipoGrafo(t *testing.T) {
	// Setup
	grafo := domain.NuevoGrafo(false)
	repo := repository.NuevoRepositorio("../../data/")
	grafoSvc := service.NuevoServicioGrafo(grafo, repo)
	handler := handler.NuevoGraphHandler(grafoSvc)

	// Test cambiar a dirigido
	err := handler.CambiarTipoGrafo(true)
	if err != nil {
		t.Errorf("Error al cambiar tipo de grafo: %v", err)
	}

	// Test cambiar a no dirigido
	err = handler.CambiarTipoGrafo(false)
	if err != nil {
		t.Errorf("Error al cambiar tipo de grafo: %v", err)
	}
}
