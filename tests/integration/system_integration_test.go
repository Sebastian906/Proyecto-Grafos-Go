package integration

import (
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/handler"
	"proyecto-grafos-go/internal/repository"
	"proyecto-grafos-go/internal/service"
	"testing"
)

func TestIntegracionCompleta(t *testing.T) {
	// Setup del sistema completo
	grafo := domain.NuevoGrafo(false)
	repo := repository.NuevoRepositorio("../../data/")

	// Servicios
	grafoSvc := service.NuevoServicioGrafo(grafo, repo)
	cuevaSvc := service.ServicioNuevaCueva(grafo)
	validacionSvc := service.NuevoServicioValidacion(grafo)
	conexionSvc := service.NuevoServicioConexion(grafo)

	// Handlers
	grafoHandler := handler.NuevoGraphHandler(grafoSvc)
	cuevaHandler := handler.NuevoCaveHandler(cuevaSvc)

	// Test 1: Creación de cuevas a través de handlers
	t.Run("CrearCuevasMedianteHandler", func(t *testing.T) {
		solicitudes := []service.SolicitudCueva{
			{ID: "C1", Nombre: "Cueva Principal", X: 0, Y: 0},
			{ID: "C2", Nombre: "Cueva Secundaria", X: 10, Y: 10},
			{ID: "C3", Nombre: "Cueva Auxiliar", X: 20, Y: 0},
		}

		for _, sol := range solicitudes {
			err := cuevaHandler.CrearCueva(sol)
			if err != nil {
				t.Errorf("Error al crear cueva %s: %v", sol.ID, err)
			}
		}

		// Verificar que las cuevas fueron creadas
		cuevas, err := cuevaHandler.ListarCuevas()
		if err != nil {
			t.Errorf("Error al listar cuevas: %v", err)
		}
		if len(cuevas) != 3 {
			t.Errorf("Esperaba 3 cuevas, obtuvo %d", len(cuevas))
		}
	})

	// Test 2: Creación de conexiones
	t.Run("CrearConexiones", func(t *testing.T) {
		conexiones := []service.SolicitudConectarCuevas{
			{DesdeCuevaID: "C1", HastaCuevaID: "C2", Distancia: 14.14, EsBidireccional: true},
			{DesdeCuevaID: "C2", HastaCuevaID: "C3", Distancia: 14.14, EsBidireccional: true},
			{DesdeCuevaID: "C1", HastaCuevaID: "C3", Distancia: 20.0, EsBidireccional: true},
		}

		for _, conn := range conexiones {
			err := cuevaSvc.ConectarCuevas(conn)
			if err != nil {
				t.Errorf("Error al conectar %s con %s: %v", conn.DesdeCuevaID, conn.HastaCuevaID, err)
			}
		}
	})

	// Test 3: Validación del sistema
	t.Run("ValidarSistema", func(t *testing.T) {
		// Verificar que el grafo es conexo
		if !validacionSvc.EsConexo() {
			t.Error("El grafo debería ser conexo")
		}

		// Verificar estadísticas del grafo
		stats, err := grafoHandler.ObtenerEstadisticas()
		if err != nil {
			t.Errorf("Error al obtener estadísticas: %v", err)
		}

		if stats.NumCuevas != 3 {
			t.Errorf("Esperaba 3 cuevas, obtuvo %d", stats.NumCuevas)
		}

		if stats.NumConexiones != 3 {
			t.Errorf("Esperaba 3 conexiones, obtuvo %d", stats.NumConexiones)
		}
	})

	// Test 4: Operaciones de modificación
	t.Run("ModificarGrafo", func(t *testing.T) {
		// Actualizar una cueva
		solicitudActualizacion := service.SolicitudCueva{
			ID:     "C1",
			Nombre: "Cueva Principal Modificada",
			X:      5,
			Y:      5,
		}

		err := cuevaHandler.ActualizarCueva("C1", solicitudActualizacion)
		if err != nil {
			t.Errorf("Error al actualizar cueva: %v", err)
		}

		// Verificar actualización
		cueva, err := cuevaHandler.ObtenerCueva("C1")
		if err != nil {
			t.Errorf("Error al obtener cueva actualizada: %v", err)
		}

		if cueva.Nombre != "Cueva Principal Modificada" {
			t.Errorf("Nombre no actualizado correctamente")
		}
	})

	// Test 5: Operaciones de eliminación
	t.Run("EliminarElementos", func(t *testing.T) {
		// Eliminar una conexión
		err := conexionSvc.EliminarConexion("C1", "C3")
		if err != nil {
			t.Errorf("Error al eliminar conexión: %v", err)
		}

		// Verificar que la conexión fue eliminada
		if grafo.ExisteConexion("C1", "C3") {
			t.Error("La conexión C1-C3 debería haber sido eliminada")
		}

		// Verificar que el grafo sigue siendo conexo a través de C2
		if !validacionSvc.EsConexo() {
			t.Error("El grafo debería seguir siendo conexo")
		}
	})

	// Test 6: Operaciones de archivo
	t.Run("GuardarYCargarGrafo", func(t *testing.T) {
		// Guardar el grafo actual
		err := grafoHandler.GuardarGrafo("test_integration.json")
		if err != nil {
			t.Errorf("Error al guardar grafo: %v", err)
		}

		// Crear un nuevo grafo y cargar desde archivo
		nuevoGrafo := domain.NuevoGrafo(false)
		nuevoRepo := repository.NuevoRepositorio("../../data/")
		nuevoGrafoSvc := service.NuevoServicioGrafo(nuevoGrafo, nuevoRepo)
		nuevoGrafoHandler := handler.NuevoGraphHandler(nuevoGrafoSvc)

		err = nuevoGrafoHandler.CargarGrafo("test_integration.json")
		if err != nil {
			t.Errorf("Error al cargar grafo: %v", err)
		}

		// Verificar que el grafo cargado tiene la misma estructura
		stats, err := nuevoGrafoHandler.ObtenerEstadisticas()
		if err != nil {
			t.Errorf("Error al obtener estadísticas del grafo cargado: %v", err)
		}

		if stats.NumCuevas != 3 {
			t.Errorf("Grafo cargado tiene %d cuevas, esperaba 3", stats.NumCuevas)
		}
	})

	// Test 7: Cambio de tipo de grafo
	t.Run("CambiarTipoGrafo", func(t *testing.T) {
		// Cambiar a grafo dirigido
		err := grafoHandler.CambiarTipoGrafo(true)
		if err != nil {
			t.Errorf("Error al cambiar tipo de grafo: %v", err)
		}

		// Verificar que el cambio se aplicó
		if !grafo.EsDirigido {
			t.Error("El grafo debería ser dirigido")
		}

		// Cambiar de vuelta a no dirigido
		err = grafoHandler.CambiarTipoGrafo(false)
		if err != nil {
			t.Errorf("Error al cambiar tipo de grafo: %v", err)
		}

		if grafo.EsDirigido {
			t.Error("El grafo debería ser no dirigido")
		}
	})
}

func TestIntegracionConObstrucciones(t *testing.T) {
	// Setup del sistema
	grafo := domain.NuevoGrafo(false)
	repo := repository.NuevoRepositorio("../../data/")

	grafoSvc := service.NuevoServicioGrafo(grafo, repo)
	cuevaSvc := service.ServicioNuevaCueva(grafo)
	conexionSvc := service.NuevoServicioConexion(grafo)
	validacionSvc := service.NuevoServicioValidacion(grafo)

	// Crear cuevas
	cuevas := []service.SolicitudCueva{
		{ID: "A", Nombre: "Cueva A", X: 0, Y: 0},
		{ID: "B", Nombre: "Cueva B", X: 10, Y: 0},
		{ID: "C", Nombre: "Cueva C", X: 20, Y: 0},
		{ID: "D", Nombre: "Cueva D", X: 30, Y: 0},
	}

	for _, cueva := range cuevas {
		err := cuevaSvc.CrearCueva(cueva)
		if err != nil {
			t.Errorf("Error al crear cueva %s: %v", cueva.ID, err)
		}
	}

	// Crear conexiones
	conexiones := []service.SolicitudConectarCuevas{
		{DesdeCuevaID: "A", HastaCuevaID: "B", Distancia: 10, EsBidireccional: true},
		{DesdeCuevaID: "B", HastaCuevaID: "C", Distancia: 10, EsBidireccional: true},
		{DesdeCuevaID: "C", HastaCuevaID: "D", Distancia: 10, EsBidireccional: true},
		{DesdeCuevaID: "A", HastaCuevaID: "D", Distancia: 30, EsBidireccional: true},
	}

	for _, conn := range conexiones {
		err := cuevaSvc.ConectarCuevas(conn)
		if err != nil {
			t.Errorf("Error al conectar %s con %s: %v", conn.DesdeCuevaID, conn.HastaCuevaID, err)
		}
	}

	// Test obstrucciones
	t.Run("ObstruccionesYDesbloqueos", func(t *testing.T) {
		// Verificar que inicialmente es conexo
		if !validacionSvc.EsConexo() {
			t.Error("El grafo debería ser conexo inicialmente")
		}

		// Obstruir conexión crítica B-C
		err := conexionSvc.ObstruirConexion("B", "C")
		if err != nil {
			t.Errorf("Error al obstruir conexión: %v", err)
		}

		// El grafo debería seguir siendo conexo (hay ruta A-D)
		if !validacionSvc.EsConexo() {
			t.Error("El grafo debería seguir siendo conexo")
		}

		// Obstruir también A-D
		err = conexionSvc.ObstruirConexion("A", "D")
		if err != nil {
			t.Errorf("Error al obstruir conexión: %v", err)
		}

		// Ahora el grafo debería ser desconexo
		if validacionSvc.EsConexo() {
			t.Error("El grafo debería ser desconexo")
		}

		// Desobstruir B-C
		err = conexionSvc.DesobstruirConexion("B", "C")
		if err != nil {
			t.Errorf("Error al desobstruir conexión: %v", err)
		}

		// El grafo debería ser conexo de nuevo
		if !validacionSvc.EsConexo() {
			t.Error("El grafo debería ser conexo después de desobstruir")
		}
	})
}

func TestIntegracionConDatos(t *testing.T) {
	// Test con datos reales del archivo
	grafo := domain.NuevoGrafo(false)
	repo := repository.NuevoRepositorio("../../data/")
	grafoSvc := service.NuevoServicioGrafo(grafo, repo)
	grafoHandler := handler.NuevoGraphHandler(grafoSvc)

	t.Run("CargarDatosReales", func(t *testing.T) {
		// Intentar cargar el archivo de ejemplo
		err := grafoHandler.CargarGrafo("caves_example.json")
		if err != nil {
			t.Logf("Advertencia: No se pudo cargar archivo de ejemplo: %v", err)
			t.Skip("Saltando test con datos reales")
		}

		// Verificar que se cargaron datos
		stats, err := grafoHandler.ObtenerEstadisticas()
		if err != nil {
			t.Errorf("Error al obtener estadísticas: %v", err)
		}

		if stats.NumCuevas == 0 {
			t.Error("No se cargaron cuevas del archivo")
		}

		t.Logf("Datos cargados: %d cuevas, %d conexiones", stats.NumCuevas, stats.NumConexiones)
	})
}
