package integration

import (
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/repository"
	"proyecto-grafos-go/internal/service"
	"testing"
)

// TestAccesibilidadCon9Cuevas prueba la detección de cuevas inaccesibles con datos reales
func TestAccesibilidadCon9Cuevas(t *testing.T) {
	// Inicializar sistema
	grafo := domain.NuevoGrafo(false)
	repo := repository.NuevoRepositorio("../../data/")
	grafoSvc := service.NuevoServicioGrafo(grafo, repo)
	validacionSvc := service.NuevoServicioValidacion(grafo)

	// Cargar datos de 9 cuevas
	t.Log("Cargando datos de ejemplo (9 cuevas)...")
	if err := grafoSvc.CargarGrafo("caves_directed_example.json"); err != nil {
		t.Fatalf("Error cargando datos: %v", err)
	}

	if len(grafo.Cuevas) != 9 {
		t.Errorf("Se esperaban 9 cuevas, se cargaron %d", len(grafo.Cuevas))
	}

	t.Logf("✓ Cargadas %d cuevas y %d conexiones", len(grafo.Cuevas), len(grafo.Aristas))

	// Casos de prueba específicos
	casos := []struct {
		cueva                 string
		accesiblesEsperadas   int
		inaccesiblesEsperadas int
		descripcion           string
	}{
		{"Silvestre", 9, 0, "Debe alcanzar todas las cuevas"},
		{"Tazmania", 7, 2, "Algunas cuevas inaccesibles"},
		{"Correcaminos", 1, 8, "Solo se alcanza a sí misma"},
		{"Popeye", 2, 7, "Pocas cuevas accesibles"},
	}

	for _, caso := range casos {
		t.Run("Accesibilidad desde "+caso.cueva, func(t *testing.T) {
			resultado := validacionSvc.AnalizarAccesibilidad(caso.cueva)

			// Verificar total de cuevas
			if resultado.TotalCuevas != 9 {
				t.Errorf("Total cuevas esperado: 9, obtenido: %d", resultado.TotalCuevas)
			}

			// Verificar cuevas accesibles
			if resultado.CuevasAccesibles != caso.accesiblesEsperadas {
				t.Errorf("Cuevas accesibles esperadas: %d, obtenidas: %d",
					caso.accesiblesEsperadas, resultado.CuevasAccesibles)
			}

			// Verificar cuevas inaccesibles
			if len(resultado.CuevasInaccesibles) != caso.inaccesiblesEsperadas {
				t.Errorf("Cuevas inaccesibles esperadas: %d, obtenidas: %d",
					caso.inaccesiblesEsperadas, len(resultado.CuevasInaccesibles))
			}

			// Verificar que se generaron soluciones
			if len(resultado.CuevasInaccesibles) > 0 && len(resultado.Soluciones) == 0 {
				t.Error("Se esperaban soluciones para cuevas inaccesibles")
			}

			t.Logf("✓ Desde %s: %d accesibles, %d inaccesibles",
				caso.cueva, resultado.CuevasAccesibles, len(resultado.CuevasInaccesibles))
		})
	}
}

// TestObstruccionConexiones prueba el impacto de obstruir conexiones
func TestObstruccionConexiones(t *testing.T) {
	// Inicializar sistema
	grafo := domain.NuevoGrafo(false)
	repo := repository.NuevoRepositorio("../../data/")
	grafoSvc := service.NuevoServicioGrafo(grafo, repo)
	validacionSvc := service.NuevoServicioValidacion(grafo)
	conexionSvc := service.NuevoServicioConexion(grafo)

	// Cargar datos
	if err := grafoSvc.CargarGrafo("caves_directed_example.json"); err != nil {
		t.Fatalf("Error cargando datos: %v", err)
	}

	// Análisis inicial desde Silvestre
	resultadoInicial := validacionSvc.AnalizarAccesibilidad("Silvestre")
	t.Logf("Estado inicial desde Silvestre: %d accesibles", resultadoInicial.CuevasAccesibles)

	// Obstruir conexión crítica
	solicitud := &service.ObstruirConexion{
		DesdeCuevaID: "Silvestre",
		HastaCuevaID: "Tazmania",
		EsObstruido:  true,
	}

	if err := conexionSvc.ObstruirConexion(solicitud); err != nil {
		t.Fatalf("Error al obstruir conexión: %v", err)
	}
	t.Log("✓ Conexión Silvestre → Tazmania obstruida")

	// Analizar impacto
	resultadoDespues := validacionSvc.AnalizarAccesibilidad("Silvestre")

	// Verificar que hay impacto
	if resultadoDespues.CuevasAccesibles >= resultadoInicial.CuevasAccesibles {
		t.Error("La obstrucción debería reducir el número de cuevas accesibles")
	}

	// Verificar que hay cuevas inaccesibles
	if len(resultadoDespues.CuevasInaccesibles) == 0 {
		t.Error("Debería haber cuevas inaccesibles después de la obstrucción")
	}

	t.Logf("✓ Después de obstrucción: %d accesibles, %d inaccesibles",
		resultadoDespues.CuevasAccesibles, len(resultadoDespues.CuevasInaccesibles))
}

// TestDeteccionCuevasInaccesibles prueba la funcionalidad general de detección
func TestDeteccionCuevasInaccesibles(t *testing.T) {
	// Inicializar sistema
	grafo := domain.NuevoGrafo(false)
	repo := repository.NuevoRepositorio("../../data/")
	grafoSvc := service.NuevoServicioGrafo(grafo, repo)
	validacionSvc := service.NuevoServicioValidacion(grafo)

	// Cargar datos
	if err := grafoSvc.CargarGrafo("caves_directed_example.json"); err != nil {
		t.Fatalf("Error cargando datos: %v", err)
	}

	// Ejecutar detección general
	resultado := validacionSvc.DetectarCuevasInaccesiblesTrasChanged()

	// Verificaciones básicas
	if resultado.TotalCuevas != 9 {
		t.Errorf("Total cuevas esperado: 9, obtenido: %d", resultado.TotalCuevas)
	}

	if resultado.CuevasAccesibles == 0 {
		t.Error("Debería haber al menos algunas cuevas accesibles")
	}

	if len(resultado.CuevasInaccesibles) == 0 {
		t.Log("No hay cuevas inaccesibles en el estado actual")
	} else {
		t.Logf("Cuevas inaccesibles detectadas: %v", resultado.CuevasInaccesibles)
	}

	// Verificar que se generaron soluciones si hay problemas
	if len(resultado.CuevasInaccesibles) > 0 && len(resultado.Soluciones) == 0 {
		t.Error("Se esperaban soluciones para cuevas inaccesibles")
	}

	t.Logf("✓ Detección completada: %d accesibles, %d inaccesibles",
		resultado.CuevasAccesibles, len(resultado.CuevasInaccesibles))
}
