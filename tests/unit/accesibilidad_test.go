package service_test

import (
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/service"
	"testing"
)

func TestDetectarCuevasInaccesibles(t *testing.T) {
	// Crear grafo dirigido para testing
	grafo := domain.NuevoGrafo(true)

	// Agregar cuevas
	cueva1 := domain.NuevaCueva("A", "Cueva A")
	cueva2 := domain.NuevaCueva("B", "Cueva B")
	cueva3 := domain.NuevaCueva("C", "Cueva C")
	cueva4 := domain.NuevaCueva("D", "Cueva D")

	// Agregar recursos
	cueva1.AgregarRecurso("oro", 10)
	cueva2.AgregarRecurso("plata", 15)
	cueva3.AgregarRecurso("hierro", 20)
	cueva4.AgregarRecurso("carbon", 25)

	grafo.AgregarCueva(cueva1)
	grafo.AgregarCueva(cueva2)
	grafo.AgregarCueva(cueva3)
	grafo.AgregarCueva(cueva4)

	// Agregar conexiones: A -> B -> C, D queda aislada
	arista1 := domain.NuevaArista("A", "B", 10, true)
	arista2 := domain.NuevaArista("B", "C", 15, true)

	grafo.AgregarArista(arista1)
	grafo.AgregarArista(arista2)

	// Crear servicio de validación
	validacionSvc := service.NuevoServicioValidacion(grafo)

	// Test: Detectar cuevas inaccesibles desde A
	resultado := validacionSvc.AnalizarAccesibilidad("A")

	if resultado == nil {
		t.Fatal("El resultado no debería ser nil")
	}

	if resultado.TotalCuevas != 4 {
		t.Errorf("Se esperaban 4 cuevas totales, se obtuvieron %d", resultado.TotalCuevas)
	}

	if resultado.CuevasAccesibles != 3 {
		t.Errorf("Se esperaban 3 cuevas accesibles, se obtuvieron %d", resultado.CuevasAccesibles)
	}

	if len(resultado.CuevasInaccesibles) != 1 {
		t.Errorf("Se esperaba 1 cueva inaccesible, se obtuvieron %d", len(resultado.CuevasInaccesibles))
	}

	if len(resultado.CuevasInaccesibles) > 0 && resultado.CuevasInaccesibles[0] != "D" {
		t.Errorf("Se esperaba que la cueva 'D' fuera inaccesible, se obtuvo '%s'", resultado.CuevasInaccesibles[0])
	}

	if len(resultado.Soluciones) == 0 {
		t.Error("Se esperaban soluciones propuestas")
	}
}

func TestDetectarCuevasInaccesiblesTrasObstruccion(t *testing.T) {
	// Crear grafo dirigido
	grafo := domain.NuevoGrafo(true)

	// Agregar cuevas
	cueva1 := domain.NuevaCueva("A", "Cueva A")
	cueva2 := domain.NuevaCueva("B", "Cueva B")
	cueva3 := domain.NuevaCueva("C", "Cueva C")

	cueva1.AgregarRecurso("oro", 10)
	cueva2.AgregarRecurso("plata", 15)
	cueva3.AgregarRecurso("hierro", 20)

	grafo.AgregarCueva(cueva1)
	grafo.AgregarCueva(cueva2)
	grafo.AgregarCueva(cueva3)

	// Crear conexiones: A -> B -> C
	arista1 := domain.NuevaArista("A", "B", 10, true)
	arista2 := domain.NuevaArista("B", "C", 15, true)

	grafo.AgregarArista(arista1)
	grafo.AgregarArista(arista2)

	// Servicios
	validacionSvc := service.NuevoServicioValidacion(grafo)
	conexionSvc := service.NuevoServicioConexion(grafo)

	// Verificar que inicialmente todas son accesibles desde A
	resultado := validacionSvc.AnalizarAccesibilidad("A")
	if len(resultado.CuevasInaccesibles) != 0 {
		t.Error("Inicialmente todas las cuevas deberían ser accesibles")
	}

	// Obstruir la conexión A -> B
	solicitudObstruir := &service.ObstruirConexion{
		DesdeCuevaID: "A",
		HastaCuevaID: "B",
		EsObstruido:  true,
	}

	err := conexionSvc.ObstruirConexion(solicitudObstruir)
	if err != nil {
		t.Fatalf("Error al obstruir conexión: %v", err)
	}

	// Ahora verificar cuevas inaccesibles
	resultado = validacionSvc.AnalizarAccesibilidad("A")

	if len(resultado.CuevasInaccesibles) != 2 {
		t.Errorf("Se esperaban 2 cuevas inaccesibles tras obstrucción, se obtuvieron %d", len(resultado.CuevasInaccesibles))
	}

	// Verificar que las soluciones incluyen eliminar la obstrucción
	solucionesTexto := ""
	for _, solucion := range resultado.Soluciones {
		solucionesTexto += solucion + " "
	}

	if !contains(solucionesTexto, "obstruccion") && !contains(solucionesTexto, "obstrucci") {
		t.Error("Las soluciones deberían mencionar eliminar obstrucciones")
	}
}

func TestCambiarDireccionConImpactoEnAccesibilidad(t *testing.T) {
	// Crear grafo dirigido
	grafo := domain.NuevoGrafo(true)

	// Agregar cuevas
	cueva1 := domain.NuevaCueva("A", "Cueva A")
	cueva2 := domain.NuevaCueva("B", "Cueva B")

	cueva1.AgregarRecurso("oro", 10)
	cueva2.AgregarRecurso("plata", 15)

	grafo.AgregarCueva(cueva1)
	grafo.AgregarCueva(cueva2)

	// Crear conexión dirigida A -> B
	arista1 := domain.NuevaArista("A", "B", 10, true)
	grafo.AgregarArista(arista1)

	// Servicios
	validacionSvc := service.NuevoServicioValidacion(grafo)
	conexionSvc := service.NuevoServicioConexion(grafo)

	// Desde A, B debería ser accesible
	resultado := validacionSvc.AnalizarAccesibilidad("A")
	if len(resultado.CuevasInaccesibles) != 0 {
		t.Error("Desde A, B debería ser accesible")
	}

	// Desde B, A debería ser inaccesible (grafo dirigido)
	resultado = validacionSvc.AnalizarAccesibilidad("B")
	if len(resultado.CuevasInaccesibles) != 1 {
		t.Error("Desde B, A debería ser inaccesible en grafo dirigido")
	}

	// Cambiar conexión a no dirigida
	solicitudCambio := &service.CambiarDireccion{
		DesdeCuevaID:   "A",
		HastaCuevaID:   "B",
		NuevaDireccion: false, // no dirigida
	}

	err := conexionSvc.CambiarDireccionConexion(solicitudCambio)
	if err != nil {
		t.Fatalf("Error al cambiar dirección: %v", err)
	}

	// Ahora desde B, A debería ser accesible
	resultado = validacionSvc.AnalizarAccesibilidad("B")
	if len(resultado.CuevasInaccesibles) != 0 {
		t.Error("Después del cambio a no dirigida, desde B, A debería ser accesible")
	}
}

func TestGenerarSolucionesEspecificas(t *testing.T) {
	// Crear grafo con cueva aislada
	grafo := domain.NuevoGrafo(true)

	cueva1 := domain.NuevaCueva("A", "Cueva A")
	cueva2 := domain.NuevaCueva("B", "Cueva B")
	cueva3 := domain.NuevaCueva("C", "Cueva C")

	cueva1.AgregarRecurso("oro", 10)
	cueva2.AgregarRecurso("plata", 15)
	cueva3.AgregarRecurso("hierro", 20)

	grafo.AgregarCueva(cueva1)
	grafo.AgregarCueva(cueva2)
	grafo.AgregarCueva(cueva3)

	// Solo conectar A -> B, C queda aislada
	arista1 := domain.NuevaArista("A", "B", 10, true)
	grafo.AgregarArista(arista1)

	// Crear conexión obstruida hacia C
	aristaObstruida := domain.NuevaArista("B", "C", 15, true)
	aristaObstruida.EsObstruido = true
	grafo.AgregarArista(aristaObstruida)

	validacionSvc := service.NuevoServicioValidacion(grafo)

	resultado := validacionSvc.AnalizarAccesibilidad("A")

	// Verificar que se proponen soluciones específicas
	solucionesTexto := ""
	for _, solucion := range resultado.Soluciones {
		solucionesTexto += solucion + " "
	}

	// Debería sugerir eliminar obstrucciones
	if !contains(solucionesTexto, "obstruccion") && !contains(solucionesTexto, "obstrucci") {
		t.Error("Debería sugerir eliminar obstrucciones")
	}

	// Debería sugerir agregar conexiones
	if !contains(solucionesTexto, "Agregar") && !contains(solucionesTexto, "agregar") {
		t.Error("Debería sugerir agregar nuevas conexiones")
	}
}

// Función auxiliar para verificar si una cadena contiene una subcadena
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsAt(s, substr)))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
