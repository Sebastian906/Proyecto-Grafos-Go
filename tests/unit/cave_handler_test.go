package unit

import (
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/handler"
	"proyecto-grafos-go/internal/service"
	"testing"
)

func TestCaveHandler_CrearCueva(t *testing.T) {
	// Setup
	grafo := domain.NuevoGrafo(false)
	cuevaSvc := service.ServicioNuevaCueva(grafo)
	handler := handler.NuevoCaveHandler(cuevaSvc)

	// Test crear cueva válida
	solicitud := service.SolicitudCueva{
		ID:     "C1",
		Nombre: "Cueva Test",
		X:      10.0,
		Y:      20.0,
	}

	err := handler.CrearCueva(solicitud)
	if err != nil {
		t.Errorf("Error al crear cueva: %v", err)
	}

	// Verificar que la cueva fue creada
	cueva, existe := grafo.ObtenerCueva("C1")
	if !existe {
		t.Error("La cueva no fue creada correctamente")
	}
	if cueva.Nombre != "Cueva Test" {
		t.Errorf("Nombre incorrecto, esperaba 'Cueva Test', obtuvo '%s'", cueva.Nombre)
	}
}

func TestCaveHandler_CrearCuevaInvalida(t *testing.T) {
	// Setup
	grafo := domain.NuevoGrafo(false)
	cuevaSvc := service.ServicioNuevaCueva(grafo)
	handler := handler.NuevoCaveHandler(cuevaSvc)

	// Test crear cueva sin ID
	solicitud := service.SolicitudCueva{
		ID:     "",
		Nombre: "Cueva Test",
		X:      10.0,
		Y:      20.0,
	}

	err := handler.CrearCueva(solicitud)
	if err == nil {
		t.Error("Debería fallar al crear cueva sin ID")
	}

	// Test crear cueva sin nombre
	solicitud = service.SolicitudCueva{
		ID:     "C1",
		Nombre: "",
		X:      10.0,
		Y:      20.0,
	}

	err = handler.CrearCueva(solicitud)
	if err == nil {
		t.Error("Debería fallar al crear cueva sin nombre")
	}
}

func TestCaveHandler_ListarCuevas(t *testing.T) {
	// Setup
	grafo := domain.NuevoGrafo(false)
	cuevaSvc := service.ServicioNuevaCueva(grafo)
	handler := handler.NuevoCaveHandler(cuevaSvc)

	// Agregar algunas cuevas
	cueva1 := domain.NuevaCueva("C1", "Cueva 1")
	cueva2 := domain.NuevaCueva("C2", "Cueva 2")
	grafo.AgregarCueva(cueva1)
	grafo.AgregarCueva(cueva2)

	// Test listar cuevas
	cuevas, err := handler.ListarCuevas()
	if err != nil {
		t.Errorf("Error al listar cuevas: %v", err)
	}
	if len(cuevas) != 2 {
		t.Errorf("Esperaba 2 cuevas, obtuvo %d", len(cuevas))
	}
}

func TestCaveHandler_ObtenerCueva(t *testing.T) {
	// Setup
	grafo := domain.NuevoGrafo(false)
	cuevaSvc := service.ServicioNuevaCueva(grafo)
	handler := handler.NuevoCaveHandler(cuevaSvc)

	// Agregar una cueva
	cueva := domain.NuevaCueva("C1", "Cueva Test")
	grafo.AgregarCueva(cueva)

	// Test obtener cueva existente
	cuevaObtenida, err := handler.ObtenerCueva("C1")
	if err != nil {
		t.Errorf("Error al obtener cueva: %v", err)
	}
	if cuevaObtenida.ID != "C1" {
		t.Errorf("ID incorrecto, esperaba 'C1', obtuvo '%s'", cuevaObtenida.ID)
	}

	// Test obtener cueva inexistente
	_, err = handler.ObtenerCueva("C999")
	if err == nil {
		t.Error("Debería fallar al obtener cueva inexistente")
	}
}

func TestCaveHandler_EliminarCueva(t *testing.T) {
	// Setup
	grafo := domain.NuevoGrafo(false)
	cuevaSvc := service.ServicioNuevaCueva(grafo)
	handler := handler.NuevoCaveHandler(cuevaSvc)

	// Agregar una cueva
	cueva := domain.NuevaCueva("C1", "Cueva Test")
	grafo.AgregarCueva(cueva)

	// Test eliminar cueva existente
	err := handler.EliminarCueva("C1")
	if err != nil {
		t.Errorf("Error al eliminar cueva: %v", err)
	}

	// Verificar que la cueva fue eliminada
	_, existe := grafo.ObtenerCueva("C1")
	if existe {
		t.Error("La cueva no fue eliminada correctamente")
	}

	// Test eliminar cueva inexistente
	err = handler.EliminarCueva("C999")
	if err == nil {
		t.Error("Debería fallar al eliminar cueva inexistente")
	}
}

func TestCaveHandler_ActualizarCueva(t *testing.T) {
	// Setup
	grafo := domain.NuevoGrafo(false)
	cuevaSvc := service.ServicioNuevaCueva(grafo)
	handler := handler.NuevoCaveHandler(cuevaSvc)

	// Agregar una cueva
	cueva := domain.NuevaCueva("C1", "Cueva Original")
	grafo.AgregarCueva(cueva)

	// Test actualizar cueva
	solicitud := service.SolicitudCueva{
		ID:     "C1",
		Nombre: "Cueva Actualizada",
		X:      50.0,
		Y:      60.0,
	}

	err := handler.ActualizarCueva("C1", solicitud)
	if err != nil {
		t.Errorf("Error al actualizar cueva: %v", err)
	}

	// Verificar que la cueva fue actualizada
	cuevaActualizada, existe := grafo.ObtenerCueva("C1")
	if !existe {
		t.Error("La cueva no existe después de la actualización")
	}
	if cuevaActualizada.Nombre != "Cueva Actualizada" {
		t.Errorf("Nombre no actualizado, esperaba 'Cueva Actualizada', obtuvo '%s'", cuevaActualizada.Nombre)
	}
}
