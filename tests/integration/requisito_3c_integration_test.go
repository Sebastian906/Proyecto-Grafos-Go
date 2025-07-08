package integration

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/handler"
	"proyecto-grafos-go/internal/repository"
	"proyecto-grafos-go/internal/service"
	"strings"
	"testing"
	"time"
)

// TestRequisito3c_IntegracionCompleta prueba todo el flujo del requisito 3c
func TestRequisito3c_IntegracionCompleta(t *testing.T) {
	// Crear grafo de ejemplo
	grafo := crearGrafoCompleto()

	// Crear todos los servicios necesarios
	repositorio := repository.NuevoRepositorio("../../data")
	servicioGrafo := service.NuevoServicioGrafo(grafo, repositorio)
	mstService := service.NuevoMSTService(servicioGrafo)
	analysisHandler := handler.NuevoAnalysisHandler(mstService)

	// Ejecutar el requisito 3c completo
	resultado, err := analysisHandler.CalcularMSTEnOrdenCreacion(grafo)

	// Verificaciones de integración
	if err != nil {
		t.Fatalf("Error en integración completa: %v", err)
	}

	// Verificar que el resultado contiene todos los elementos esperados
	if len(resultado) == 0 {
		t.Errorf("El resultado no debería estar vacío")
	}

	// Verificar que contiene las secciones principales
	seccionesEsperadas := []string{
		"RUTAS DE ACCESO MÍNIMAS EN ORDEN DE CREACIÓN",
		"Estado:",
		"Red conexa:",
		"Total de cuevas:",
		"ESTADÍSTICAS DE LA RED",
		"ANÁLISIS DE RUTAS DE ACCESO",
		"ANÁLISIS DE OPTIMIZACIÓN",
		"RECOMENDACIONES",
	}

	for _, seccion := range seccionesEsperadas {
		if !strings.Contains(resultado, seccion) {
			t.Errorf("El resultado debería contener la sección: %s", seccion)
		}
	}

	// Verificar que menciona las cuevas específicas
	cuevasEsperadas := []string{"C1", "C2", "C3", "C4", "C5"}
	for _, cueva := range cuevasEsperadas {
		if !strings.Contains(resultado, cueva) {
			t.Errorf("El resultado debería mencionar la cueva: %s", cueva)
		}
	}

	// Verificar que contiene información de rutas
	if !strings.Contains(resultado, "->") {
		t.Errorf("El resultado debería contener rutas con el símbolo '->'")
	}

	// Verificar que contiene información de distancias
	if !strings.Contains(resultado, "distancia:") {
		t.Errorf("El resultado debería contener información de distancias")
	}
}

// TestRequisito3c_ConGrafoDesconectado prueba con grafo desconectado
func TestRequisito3c_ConGrafoDesconectado(t *testing.T) {
	// Crear grafo desconectado
	grafo := crearGrafoDesconectadoIntegracion()

	// Crear servicios
	repositorio := repository.NuevoRepositorio("../../data")
	servicioGrafo := service.NuevoServicioGrafo(grafo, repositorio)
	mstService := service.NuevoMSTService(servicioGrafo)
	analysisHandler := handler.NuevoAnalysisHandler(mstService)

	// Ejecutar requisito 3c
	resultado, err := analysisHandler.CalcularMSTEnOrdenCreacion(grafo)

	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	// Print actual output for debugging
	t.Logf("Actual output for disconnected graph:\n%s", resultado)

	// Verificar que indica problema de conectividad
	if !strings.Contains(resultado, "no está completamente conectada") {
		t.Errorf("El resultado debería indicar un error de conectividad")
	}

	// Verificar que contiene información sobre el problema
	if !strings.Contains(resultado, "No es posible mostrar rutas") {
		t.Errorf("El resultado debería contener información sobre el problema")
	}
}

// TestRequisito3c_RendimientoConGrafoGrande prueba el rendimiento con un grafo más grande
func TestRequisito3c_RendimientoConGrafoGrande(t *testing.T) {
	// Crear grafo grande (20 cuevas)
	grafo := crearGrafoGrande(20)

	// Crear servicios
	repositorio := repository.NuevoRepositorio("../../data")
	servicioGrafo := service.NuevoServicioGrafo(grafo, repositorio)
	mstService := service.NuevoMSTService(servicioGrafo)
	analysisHandler := handler.NuevoAnalysisHandler(mstService)

	// Medir tiempo de ejecución
	inicio := time.Now()
	resultado, err := analysisHandler.CalcularMSTEnOrdenCreacion(grafo)
	duracion := time.Since(inicio)

	if err != nil {
		t.Fatalf("Error con grafo grande: %v", err)
	}

	// Verificar que el resultado es correcto
	if len(resultado) == 0 {
		t.Errorf("El resultado no debería estar vacío para grafo grande")
	}

	// Verificar que no toma demasiado tiempo (debería ser < 1 segundo)
	if duracion > time.Second {
		t.Errorf("El cálculo tomó demasiado tiempo: %v", duracion)
	}

	// Verificar que contiene todas las cuevas
	for i := 1; i <= 20; i++ {
		cuevaID := fmt.Sprintf("C%02d", i)
		if !strings.Contains(resultado, cuevaID) {
			t.Errorf("El resultado debería contener la cueva: %s", cuevaID)
		}
	}

	t.Logf("Tiempo de ejecución para grafo de 20 cuevas: %v", duracion)
}

// TestRequisito3c_ValidacionDatos prueba la validación de datos
func TestRequisito3c_ValidacionDatos(t *testing.T) {
	// Crear grafo válido
	grafo := crearGrafoCompleto()

	// Crear servicios
	repositorio := repository.NuevoRepositorio("../../data")
	servicioGrafo := service.NuevoServicioGrafo(grafo, repositorio)
	mstService := service.NuevoMSTService(servicioGrafo)

	// Obtener resultado del servicio directamente
	resultado, err := mstService.ObtenerMSTEnOrdenCreacion(grafo)

	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	// Validar estructura de datos del resultado
	if resultado.MST == nil {
		t.Errorf("El MST no debería ser nil")
	}

	if len(resultado.OrdenCreacion) != 5 {
		t.Errorf("El orden de creación debería tener 5 elementos, tiene %d", len(resultado.OrdenCreacion))
	}

	if len(resultado.RutasAccesoMinimas) != 5 {
		t.Errorf("Debería haber 5 rutas de acceso, hay %d", len(resultado.RutasAccesoMinimas))
	}

	if !resultado.EsConexo {
		t.Errorf("El grafo debería ser conexo")
	}

	if resultado.TotalCuevasConectadas != 5 {
		t.Errorf("Deberían haber 5 cuevas conectadas, hay %d", resultado.TotalCuevasConectadas)
	}

	// Validar que las rutas de acceso son consistentes
	for i, ruta := range resultado.RutasAccesoMinimas {
		if ruta.OrdenCreacion != i+1 {
			t.Errorf("Orden de creación inconsistente en ruta %d: esperado %d, obtenido %d",
				i, i+1, ruta.OrdenCreacion)
		}

		if i == 0 {
			// Primera ruta debería ser punto de partida
			if ruta.DistanciaTotal != 0 {
				t.Errorf("La primera ruta debería tener distancia 0, tiene %.2f", ruta.DistanciaTotal)
			}
			if !ruta.EsAccesible {
				t.Errorf("La primera ruta debería ser accesible")
			}
		} else {
			// Otras rutas deberían tener distancia > 0
			if ruta.DistanciaTotal <= 0 {
				t.Errorf("La ruta %d debería tener distancia > 0, tiene %.2f", i, ruta.DistanciaTotal)
			}
			if !ruta.EsAccesible {
				t.Errorf("La ruta %d debería ser accesible", i)
			}
		}

		// Verificar que la ruta no está vacía
		if len(ruta.Ruta) == 0 {
			t.Errorf("La ruta %d no debería estar vacía", i)
		}

		// Verificar que la ruta comienza con C1 (punto de partida)
		if ruta.Ruta[0] != "C1" {
			t.Errorf("La ruta %d debería comenzar con C1, comienza con %s", i, ruta.Ruta[0])
		}

		// Verificar que la ruta termina con el destino correcto
		ultimoNodo := ruta.Ruta[len(ruta.Ruta)-1]
		if ultimoNodo != ruta.CuevaDestino {
			t.Errorf("La ruta %d debería terminar en %s, termina en %s",
				i, ruta.CuevaDestino, ultimoNodo)
		}
	}
}

// TestRequisito3c_ComparacionAlgoritmos compara resultados con otros requisitos
func TestRequisito3c_ComparacionAlgoritmos(t *testing.T) {
	// Crear grafo
	grafo := crearGrafoCompleto()

	// Crear servicios
	repositorio := repository.NuevoRepositorio("../../data")
	servicioGrafo := service.NuevoServicioGrafo(grafo, repositorio)
	mstService := service.NuevoMSTService(servicioGrafo)

	// Obtener resultados de todos los requisitos
	resultado3a, err := mstService.ObtenerMSTGeneral(grafo)
	if err != nil {
		t.Fatalf("Error en requisito 3a: %v", err)
	}

	resultado3b, err := mstService.ObtenerMSTDesdeCueva(grafo, "C1")
	if err != nil {
		t.Fatalf("Error en requisito 3b: %v", err)
	}

	resultado3c, err := mstService.ObtenerMSTEnOrdenCreacion(grafo)
	if err != nil {
		t.Fatalf("Error en requisito 3c: %v", err)
	}

	// Comparar pesos totales de MST (deberían ser iguales)
	if resultado3a.MST.PesoTotal != resultado3c.MST.PesoTotal {
		t.Errorf("Los pesos de MST deberían ser iguales: 3a=%.2f, 3c=%.2f",
			resultado3a.MST.PesoTotal, resultado3c.MST.PesoTotal)
	}

	if resultado3b.MST.MST.PesoTotal != resultado3c.MST.PesoTotal {
		t.Errorf("Los pesos de MST deberían ser iguales: 3b=%.2f, 3c=%.2f",
			resultado3b.MST.MST.PesoTotal, resultado3c.MST.PesoTotal)
	}

	// Comparar número de aristas (deberían ser iguales: n-1)
	if resultado3a.MST.NumAristas != resultado3c.MST.NumAristas {
		t.Errorf("El número de aristas debería ser igual: 3a=%d, 3c=%d",
			resultado3a.MST.NumAristas, resultado3c.MST.NumAristas)
	}

	t.Logf("Comparación exitosa - Peso MST: %.2f, Aristas: %d",
		resultado3c.MST.PesoTotal, resultado3c.MST.NumAristas)
}

// Funciones auxiliares para las pruebas de integración

func crearGrafoCompleto() *domain.Grafo {
	grafo := domain.NuevoGrafo(false)

	// Agregar cuevas
	cuevas := []struct {
		id     string
		nombre string
	}{
		{"C1", "Entrada Principal"},
		{"C2", "Sala de Recursos"},
		{"C3", "Depósito Norte"},
		{"C4", "Sala de Herramientas"},
		{"C5", "Depósito Sur"},
	}

	for _, c := range cuevas {
		cueva := domain.NuevaCueva(c.id, c.nombre)
		grafo.AgregarCueva(cueva)
	}

	// Agregar conexiones densas
	conexiones := []struct {
		desde, hasta string
		distancia    float64
	}{
		{"C1", "C2", 10.0},
		{"C1", "C3", 15.0},
		{"C2", "C3", 8.0},
		{"C2", "C4", 12.0},
		{"C3", "C4", 6.0},
		{"C3", "C5", 20.0},
		{"C4", "C5", 9.0},
		{"C1", "C4", 25.0}, // Conexión adicional
		{"C2", "C5", 22.0}, // Conexión adicional
	}

	for _, conn := range conexiones {
		arista := &domain.Arista{
			Desde:       conn.desde,
			Hasta:       conn.hasta,
			Distancia:   conn.distancia,
			EsDirigido:  false,
			EsObstruido: false,
		}
		grafo.AgregarArista(arista)
	}

	return grafo
}

func crearGrafoDesconectadoIntegracion() *domain.Grafo {
	grafo := domain.NuevoGrafo(false)

	// Grupo 1: C1, C2, C3
	for i := 1; i <= 3; i++ {
		cueva := domain.NuevaCueva(fmt.Sprintf("C%d", i), fmt.Sprintf("Cueva %d", i))
		grafo.AgregarCueva(cueva)
	}

	// Conexiones grupo 1
	arista1 := &domain.Arista{
		Desde: "C1", Hasta: "C2", Distancia: 10.0, EsDirigido: false, EsObstruido: false,
	}
	arista2 := &domain.Arista{
		Desde: "C2", Hasta: "C3", Distancia: 15.0, EsDirigido: false, EsObstruido: false,
	}
	grafo.AgregarArista(arista1)
	grafo.AgregarArista(arista2)

	// Grupo 2: C4, C5 (aislado)
	for i := 4; i <= 5; i++ {
		cueva := domain.NuevaCueva(fmt.Sprintf("C%d", i), fmt.Sprintf("Cueva %d", i))
		grafo.AgregarCueva(cueva)
	}

	// Conexión grupo 2
	arista3 := &domain.Arista{
		Desde: "C4", Hasta: "C5", Distancia: 12.0, EsDirigido: false, EsObstruido: false,
	}
	grafo.AgregarArista(arista3)

	return grafo
}

func crearGrafoGrande(numCuevas int) *domain.Grafo {
	grafo := domain.NuevoGrafo(false)

	// Agregar cuevas
	for i := 1; i <= numCuevas; i++ {
		cueva := domain.NuevaCueva(fmt.Sprintf("C%02d", i), fmt.Sprintf("Cueva %02d", i))
		grafo.AgregarCueva(cueva)
	}

	// Crear conexiones en forma de estrella desde C01
	for i := 2; i <= numCuevas; i++ {
		arista := &domain.Arista{
			Desde:       "C01",
			Hasta:       fmt.Sprintf("C%02d", i),
			Distancia:   float64(i * 5),
			EsDirigido:  false,
			EsObstruido: false,
		}
		grafo.AgregarArista(arista)
	}

	// Agregar algunas conexiones adicionales para hacer el grafo más interesante
	for i := 2; i < numCuevas; i++ {
		if i%3 == 0 { // Cada tercera cueva se conecta con la siguiente
			arista := &domain.Arista{
				Desde:       fmt.Sprintf("C%02d", i),
				Hasta:       fmt.Sprintf("C%02d", i+1),
				Distancia:   float64(i * 2),
				EsDirigido:  false,
				EsObstruido: false,
			}
			grafo.AgregarArista(arista)
		}
	}

	return grafo
}

// BenchmarkRequisito3c_RendimientoMST benchmark para medir rendimiento
func BenchmarkRequisito3c_RendimientoMST(b *testing.B) {
	grafo := crearGrafoCompleto()
	repositorio := repository.NuevoRepositorio("../../data")
	servicioGrafo := service.NuevoServicioGrafo(grafo, repositorio)
	mstService := service.NuevoMSTService(servicioGrafo)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := mstService.ObtenerMSTEnOrdenCreacion(grafo)
		if err != nil {
			b.Fatalf("Error en benchmark: %v", err)
		}
	}
}

// BenchmarkRequisito3c_RendimientoGrafoGrande benchmark con grafo grande
func BenchmarkRequisito3c_RendimientoGrafoGrande(b *testing.B) {
	grafo := crearGrafoGrande(50)
	repositorio := repository.NuevoRepositorio("../../data")
	servicioGrafo := service.NuevoServicioGrafo(grafo, repositorio)
	mstService := service.NuevoMSTService(servicioGrafo)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := mstService.ObtenerMSTEnOrdenCreacion(grafo)
		if err != nil {
			b.Fatalf("Error en benchmark: %v", err)
		}
	}
}
