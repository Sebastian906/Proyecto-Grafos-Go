package integration

import (
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/handler"
	"proyecto-grafos-go/internal/service"
	"testing"
)

func TestIntegracionDeteccionCuevasInaccesibles(t *testing.T) {
	// Setup: Crear un grafo completo con servicios y handlers
	grafo := domain.NuevoGrafo(true)

	// Crear cuevas del ejemplo caves_directed_example.json
	cuevas := []*domain.Cueva{
		domain.NuevaCueva("entrada", "Entrada Principal"),
		domain.NuevaCueva("deposito1", "Deposito Norte"),
		domain.NuevaCueva("deposito2", "Deposito Sur"),
		domain.NuevaCueva("mina1", "Mina Este"),
		domain.NuevaCueva("mina2", "Mina Oeste"),
		domain.NuevaCueva("salida", "Salida de Emergencia"),
	}

	// Agregar recursos a las cuevas
	cuevas[1].AgregarRecurso("hierro", 25)
	cuevas[1].AgregarRecurso("carbon", 15)
	cuevas[2].AgregarRecurso("oro", 30)
	cuevas[2].AgregarRecurso("plata", 20)
	cuevas[3].AgregarRecurso("diamantes", 40)
	cuevas[4].AgregarRecurso("esmeraldas", 35)

	// Agregar cuevas al grafo
	for _, cueva := range cuevas {
		err := grafo.AgregarCueva(cueva)
		if err != nil {
			t.Fatalf("Error al agregar cueva %s: %v", cueva.ID, err)
		}
	}

	// Crear conexiones dirigidas
	conexiones := []*domain.Arista{
		domain.NuevaArista("entrada", "deposito1", 100, true),
		domain.NuevaArista("entrada", "deposito2", 120, true),
		domain.NuevaArista("deposito1", "mina1", 80, true),
		domain.NuevaArista("deposito2", "mina2", 90, true),
		domain.NuevaArista("mina1", "salida", 60, true),
		domain.NuevaArista("mina2", "salida", 70, true),
	}

	// Agregar conexiones al grafo
	for _, arista := range conexiones {
		err := grafo.AgregarArista(arista)
		if err != nil {
			t.Fatalf("Error al agregar arista %s -> %s: %v", arista.Desde, arista.Hasta, err)
		}
	}

	// Crear servicios
	validacionSvc := service.NuevoServicioValidacion(grafo)
	conexionSvc := service.NuevoServicioConexion(grafo)

	// Crear handler
	controladorConexion := handler.NuevoControladorConexion(conexionSvc, validacionSvc)

	// Test 1: Verificar accesibilidad inicial desde entrada
	t.Run("Accesibilidad inicial desde entrada", func(t *testing.T) {
		resultado := validacionSvc.AnalizarAccesibilidad("entrada")

		if len(resultado.CuevasInaccesibles) != 0 {
			t.Errorf("Inicialmente todas las cuevas deberían ser accesibles desde entrada, pero hay %d inaccesibles: %v",
				len(resultado.CuevasInaccesibles), resultado.CuevasInaccesibles)
		}

		if resultado.CuevasAccesibles != 6 {
			t.Errorf("Se esperaban 6 cuevas accesibles, se obtuvieron %d", resultado.CuevasAccesibles)
		}
	})

	// Test 2: Obstruir conexión crítica y verificar impacto
	t.Run("Obstrucción de conexión crítica", func(t *testing.T) {
		// Obstruir entrada -> deposito1
		solicitud := &service.ObstruirConexion{
			DesdeCuevaID: "entrada",
			HastaCuevaID: "deposito1",
			EsObstruido:  true,
		}

		err := conexionSvc.ObstruirConexion(solicitud)
		if err != nil {
			t.Fatalf("Error al obstruir conexión: %v", err)
		}

		// Verificar impacto
		resultado := validacionSvc.AnalizarAccesibilidad("entrada")

		expectedInaccesibles := 2 // deposito1 y mina1 deberían ser inaccesibles
		if len(resultado.CuevasInaccesibles) != expectedInaccesibles {
			t.Errorf("Se esperaban %d cuevas inaccesibles, se obtuvieron %d: %v",
				expectedInaccesibles, len(resultado.CuevasInaccesibles), resultado.CuevasInaccesibles)
		}

		// Verificar que las soluciones incluyen desobstrucción
		solucionesTexto := ""
		for _, solucion := range resultado.Soluciones {
			solucionesTexto += solucion + " "
		}

		if !containsSubstring(solucionesTexto, "obstruccion") {
			t.Error("Las soluciones deberían mencionar eliminar obstrucciones")
		}

		// Restaurar para siguientes tests
		solicitud.EsObstruido = false
		conexionSvc.ObstruirConexion(solicitud)
	})

	// Test 3: Cambiar dirección y analizar impacto usando handler
	t.Run("Cambio de dirección con análisis automático", func(t *testing.T) {
		// Usar el handler para cambiar dirección con análisis
		solicitudJSON := `{
			"desde_cueva_id": "entrada",
			"hasta_cueva_id": "deposito1", 
			"nueva_direccion": false
		}`

		resultado, err := controladorConexion.ManejarCambiarDireccionConConAnalisis([]byte(solicitudJSON))
		if err != nil {
			t.Fatalf("Error en handler: %v", err)
		}

		// Verificar que el resultado incluye análisis de accesibilidad
		if !containsSubstring(resultado, "Conexion desde entrada hasta deposito1 cambiada") {
			t.Error("El resultado debería confirmar el cambio de conexión")
		}

		// El cambio a no dirigido no debería causar problemas de accesibilidad
		if containsSubstring(resultado, "cuevas inaccesibles") {
			t.Log("Resultado completo:", resultado) // Para debug si hay problemas
		}
	})

	// Test 4: Análisis desde cueva sin conexiones salientes
	t.Run("Análisis desde cueva terminal", func(t *testing.T) {
		resultado := validacionSvc.AnalizarAccesibilidad("salida")

		// Desde salida (cueva terminal), no debería poder acceder a otras en grafo dirigido
		if len(resultado.CuevasInaccesibles) == 0 {
			t.Error("Desde una cueva terminal en grafo dirigido, debería haber cuevas inaccesibles")
		}

		// Verificar que se proponen soluciones
		if len(resultado.Soluciones) == 0 {
			t.Error("Se deberían proponer soluciones para mejorar accesibilidad")
		}
	})

	// Test 5: Convertir a grafo no dirigido y verificar mejora
	t.Run("Conversión a grafo no dirigido", func(t *testing.T) {
		// Cambiar tipo de grafo
		solicitudTipo := &service.CambiarTipoGrafo{
			EsDirigido: false,
		}

		err := conexionSvc.CambiarTipoGrafo(solicitudTipo)
		if err != nil {
			t.Fatalf("Error al cambiar tipo de grafo: %v", err)
		}

		// Ahora desde cualquier cueva deberían ser accesibles las demás
		resultado := validacionSvc.AnalizarAccesibilidad("salida")

		if len(resultado.CuevasInaccesibles) != 0 {
			t.Errorf("En grafo no dirigido conectado, no debería haber cuevas inaccesibles desde 'salida', pero se encontraron: %v",
				resultado.CuevasInaccesibles)
		}
	})

	// Test 6: Verificar handler de detección automática
	t.Run("Handler detección automática", func(t *testing.T) {
		reporte, err := controladorConexion.ManejarDetectarCuevasInaccesibles()
		if err != nil {
			t.Fatalf("Error en detección automática: %v", err)
		}

		if !containsSubstring(reporte, "ANALISIS DE ACCESIBILIDAD") {
			t.Error("El reporte debería incluir encabezado de análisis")
		}

		if !containsSubstring(reporte, "Total de cuevas:") {
			t.Error("El reporte debería incluir estadísticas de cuevas")
		}
	})
}

func TestIntegracionSolucionesEspecificas(t *testing.T) {
	// Crear escenario específico para probar generación de soluciones
	grafo := domain.NuevoGrafo(true)

	// Crear cuevas
	cuevas := []*domain.Cueva{
		domain.NuevaCueva("centro", "Centro de Control"),
		domain.NuevaCueva("norte", "Sector Norte"),
		domain.NuevaCueva("sur", "Sector Sur"),
		domain.NuevaCueva("aislada", "Cueva Aislada"),
	}

	// Agregar recursos
	cuevas[1].AgregarRecurso("hierro", 20)
	cuevas[2].AgregarRecurso("oro", 25)
	cuevas[3].AgregarRecurso("diamantes", 30)

	for _, cueva := range cuevas {
		grafo.AgregarCueva(cueva)
	}

	// Conexiones: centro -> norte, centro -> sur, aislada sin conexiones
	aristas := []*domain.Arista{
		domain.NuevaArista("centro", "norte", 50, true),
		domain.NuevaArista("centro", "sur", 60, true),
	}

	for _, arista := range aristas {
		grafo.AgregarArista(arista)
	}

	// Agregar conexión obstruida hacia cueva aislada
	aristaObstruida := domain.NuevaArista("norte", "aislada", 40, true)
	aristaObstruida.EsObstruido = true
	grafo.AgregarArista(aristaObstruida)

	validacionSvc := service.NuevoServicioValidacion(grafo)
	conexionSvc := service.NuevoServicioConexion(grafo)
	controladorConexion := handler.NuevoControladorConexion(conexionSvc, validacionSvc)

	// Test: Análisis específico de cueva aislada
	reporte, err := controladorConexion.ManejarAnalizarAccesibilidadDesde("centro")
	if err != nil {
		t.Fatalf("Error en análisis específico: %v", err)
	}

	// Verificar que identifica la cueva aislada
	if !containsSubstring(reporte, "aislada") {
		t.Error("Debería identificar la cueva aislada como inaccesible")
	}

	// Verificar que propone soluciones específicas
	if !containsSubstring(reporte, "obstruccion") {
		t.Error("Debería proponer eliminar obstrucciones como solución")
	}

	if !containsSubstring(reporte, "SOLUCIONES") {
		t.Error("Debería incluir sección de soluciones")
	}
}

// Función auxiliar para verificar subcadenas
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
