package service_test

import (
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/repository"
	"proyecto-grafos-go/internal/service"
	"strings"
	"testing"
)

// TestObtenerMSTEnOrdenCreacion_CasoBasico prueba el caso básico con 5 cuevas conectadas
func TestObtenerMSTEnOrdenCreacion_CasoBasico(t *testing.T) {
	// Crear grafo de prueba
	grafo := crearGrafoEjemplo()

	// Crear servicio
	repositorio := repository.NuevoRepositorio("../../../data")
	servicioGrafo := service.NuevoServicioGrafo(grafo, repositorio)
	mstService := service.NuevoMSTService(servicioGrafo)

	// Ejecutar función
	resultado, err := mstService.ObtenerMSTEnOrdenCreacion(grafo)

	// Verificar que no hay error
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	// Verificar que la red es conexa
	if !resultado.EsConexo {
		t.Errorf("Se esperaba que la red fuera conexa")
	}

	// Verificar que hay 5 cuevas
	if resultado.TotalCuevasConectadas != 5 {
		t.Errorf("Se esperaban 5 cuevas, se obtuvo %d", resultado.TotalCuevasConectadas)
	}

	// Verificar que hay 5 rutas de acceso (incluyendo punto de partida)
	if len(resultado.RutasAccesoMinimas) != 5 {
		t.Errorf("Se esperaban 5 rutas de acceso, se obtuvo %d", len(resultado.RutasAccesoMinimas))
	}

	// Verificar que el orden de creación es correcto
	ordenEsperado := []string{"C1", "C2", "C3", "C4", "C5"}
	for i, cueva := range resultado.OrdenCreacion {
		if cueva != ordenEsperado[i] {
			t.Errorf("Orden incorrecto en posición %d: esperado %s, obtenido %s", i, ordenEsperado[i], cueva)
		}
	}

	// Verificar que el MST tiene 4 aristas (n-1 para n nodos)
	if resultado.MST.NumAristas != 4 {
		t.Errorf("Se esperaban 4 aristas en MST, se obtuvo %d", resultado.MST.NumAristas)
	}

	// Verificar que la primera cueva es el punto de partida
	if resultado.RutasAccesoMinimas[0].CuevaDestino != "C1" {
		t.Errorf("La primera cueva debería ser C1, se obtuvo %s", resultado.RutasAccesoMinimas[0].CuevaDestino)
	}

	if resultado.RutasAccesoMinimas[0].DistanciaTotal != 0 {
		t.Errorf("La distancia del punto de partida debería ser 0, se obtuvo %.2f", resultado.RutasAccesoMinimas[0].DistanciaTotal)
	}
}

// TestObtenerMSTEnOrdenCreacion_GrafoDesconectado prueba con un grafo desconectado
func TestObtenerMSTEnOrdenCreacion_GrafoDesconectado(t *testing.T) {
	// Crear grafo desconectado
	grafo := crearGrafoDesconectado()

	// Crear servicio
	repositorio := repository.NuevoRepositorio("../../../data")
	servicioGrafo := service.NuevoServicioGrafo(grafo, repositorio)
	mstService := service.NuevoMSTService(servicioGrafo)

	// Ejecutar función
	resultado, err := mstService.ObtenerMSTEnOrdenCreacion(grafo)

	// Verificar que no hay error
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	// Verificar que la red NO es conexa
	if resultado.EsConexo {
		t.Errorf("Se esperaba que la red NO fuera conexa")
	}

	// Verificar que el MST es nil
	if resultado.MST != nil {
		t.Errorf("Se esperaba que el MST fuera nil para grafo desconectado")
	}

	// Verificar que el mensaje indica problema de conectividad
	if resultado.Mensaje != "La red no está completamente conectada" {
		t.Errorf("Mensaje incorrecto: %s", resultado.Mensaje)
	}
}

// TestFormatearMSTOrdenCreacionParaVisualizacion prueba el formato de salida
func TestFormatearMSTOrdenCreacionParaVisualizacion(t *testing.T) {
	// Crear grafo y calcular resultado
	grafo := crearGrafoEjemplo()
	repositorio := repository.NuevoRepositorio("../../../data")
	servicioGrafo := service.NuevoServicioGrafo(grafo, repositorio)
	mstService := service.NuevoMSTService(servicioGrafo)

	resultado, err := mstService.ObtenerMSTEnOrdenCreacion(grafo)
	if err != nil {
		t.Fatalf("Error al calcular MST: %v", err)
	}

	// Formatear resultado
	output := mstService.FormatearMSTOrdenCreacionParaVisualizacion(resultado)

	// Verificar que la salida no está vacía
	if output == "" {
		t.Errorf("La salida no debería estar vacía")
	}

	// Print the actual output for debugging
	t.Logf("Actual output:\n%s", output)

	// Verificar que contiene las secciones esperadas básicas
	seccionesEsperadas := []string{
		"RUTAS DE ACCESO MÍNIMAS EN ORDEN DE CREACIÓN",
		"Estado:",
		"Red conexa:",
		"Total de cuevas:",
		"C1:",
		"C2:",
		"C3:",
		"C4:",
		"C5:",
	}

	for _, seccion := range seccionesEsperadas {
		if !strings.Contains(output, seccion) {
			t.Errorf("La salida no contiene la sección esperada: %s", seccion)
		}
	}
}

// Helper functions for creating test graphs

func crearGrafoEjemplo() *domain.Grafo {
	grafo := domain.NuevoGrafo(false)

	// Agregar cuevas
	cuevas := []struct {
		id     string
		nombre string
	}{
		{"C1", "Cueva Entrada"},
		{"C2", "Cueva Norte"},
		{"C3", "Cueva Sur"},
		{"C4", "Cueva Este"},
		{"C5", "Cueva Oeste"},
	}

	for _, c := range cuevas {
		cueva := domain.NuevaCueva(c.id, c.nombre)
		grafo.AgregarCueva(cueva)
	}

	// Agregar aristas
	aristas := []struct {
		desde, hasta string
		distancia    float64
	}{
		{"C1", "C2", 10.0},
		{"C1", "C3", 15.0},
		{"C2", "C3", 8.0},
		{"C2", "C4", 12.0},
		{"C3", "C4", 9.0},
		{"C3", "C5", 11.0},
		{"C4", "C5", 7.0},
	}

	for _, a := range aristas {
		arista := &domain.Arista{
			Desde:       a.desde,
			Hasta:       a.hasta,
			Distancia:   a.distancia,
			EsDirigido:  false,
			EsObstruido: false,
		}
		grafo.AgregarArista(arista)
	}

	return grafo
}

func crearGrafoDesconectado() *domain.Grafo {
	grafo := domain.NuevoGrafo(false)

	// Agregar cuevas
	grafo.AgregarCueva(domain.NuevaCueva("C1", "Cueva 1"))
	grafo.AgregarCueva(domain.NuevaCueva("C2", "Cueva 2"))
	grafo.AgregarCueva(domain.NuevaCueva("C3", "Cueva 3"))
	grafo.AgregarCueva(domain.NuevaCueva("C4", "Cueva 4"))

	// Agregar solo algunas aristas para hacer el grafo desconectado
	arista1 := &domain.Arista{
		Desde:       "C1",
		Hasta:       "C2",
		Distancia:   5.0,
		EsDirigido:  false,
		EsObstruido: false,
	}
	grafo.AgregarArista(arista1)

	// C3 y C4 quedan desconectados

	return grafo
}
