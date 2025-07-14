package unit

import (
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/service"
	"testing"
)

func TestCambiarSentidoRuta(t *testing.T) {
	// Crear un grafo dirigido de prueba
	grafo := domain.NuevoGrafo(true)

	// Agregar cuevas
	cueva1 := domain.NuevaCueva("C1", "Cueva 1")
	cueva1.X, cueva1.Y = 0, 0
	cueva2 := domain.NuevaCueva("C2", "Cueva 2")
	cueva2.X, cueva2.Y = 1, 1
	grafo.AgregarCueva(cueva1)
	grafo.AgregarCueva(cueva2)

	// Agregar arista dirigida de C1 a C2
	arista := domain.NuevaArista("C1", "C2", 10.0, true)
	grafo.AgregarArista(arista)

	// Crear servicio de conexión
	servicioConexion := service.NuevoServicioConexion(grafo)

	// Verificar arista original
	aristasOriginales := grafo.Aristas
	if len(aristasOriginales) != 1 {
		t.Fatalf("Se esperaba 1 arista, pero se encontraron %d", len(aristasOriginales))
	}

	aristaOriginal := aristasOriginales[0]
	if aristaOriginal.Desde != "C1" || aristaOriginal.Hasta != "C2" {
		t.Fatalf("Arista original incorrecta: %s -> %s", aristaOriginal.Desde, aristaOriginal.Hasta)
	}

	// Cambiar sentido de la ruta
	solicitud := &service.CambiarSentidoRuta{
		DesdeCuevaID: "C1",
		HastaCuevaID: "C2",
	}

	err := servicioConexion.CambiarSentidoRuta(solicitud)
	if err != nil {
		t.Fatalf("Error al cambiar sentido de ruta: %v", err)
	}

	// Verificar que la arista se invirtió
	aristasInvertidas := grafo.Aristas
	if len(aristasInvertidas) != 1 {
		t.Fatalf("Se esperaba 1 arista después de invertir, pero se encontraron %d", len(aristasInvertidas))
	}

	aristaInvertida := aristasInvertidas[0]
	if aristaInvertida.Desde != "C2" || aristaInvertida.Hasta != "C1" {
		t.Fatalf("Arista no se invirtió correctamente: %s -> %s", aristaInvertida.Desde, aristaInvertida.Hasta)
	}

	// Verificar que mantiene las mismas propiedades
	if aristaInvertida.Distancia != 10.0 {
		t.Fatalf("Distancia no se mantuvo: esperada 10.0, obtenida %f", aristaInvertida.Distancia)
	}

	if !aristaInvertida.EsDirigido {
		t.Fatal("La arista debería seguir siendo dirigida")
	}
}

func TestCambiarSentidoRutaNoDirigida(t *testing.T) {
	// Crear un grafo dirigido de prueba
	grafo := domain.NuevoGrafo(true)

	// Agregar cuevas
	cueva1 := domain.NuevaCueva("C1", "Cueva 1")
	cueva1.X, cueva1.Y = 0, 0
	cueva2 := domain.NuevaCueva("C2", "Cueva 2")
	cueva2.X, cueva2.Y = 1, 1
	grafo.AgregarCueva(cueva1)
	grafo.AgregarCueva(cueva2)

	// Agregar arista NO dirigida
	arista := domain.NuevaArista("C1", "C2", 10.0, false)
	grafo.AgregarArista(arista)

	// Crear servicio de conexión
	servicioConexion := service.NuevoServicioConexion(grafo)

	// Intentar cambiar sentido de la ruta no dirigida (debería fallar)
	solicitud := &service.CambiarSentidoRuta{
		DesdeCuevaID: "C1",
		HastaCuevaID: "C2",
	}

	err := servicioConexion.CambiarSentidoRuta(solicitud)
	if err == nil {
		t.Fatal("Se esperaba un error al intentar cambiar el sentido de una ruta no dirigida")
	}

	expectedError := "no se puede cambiar el sentido de una conexión no dirigida"
	if err.Error() != expectedError {
		t.Fatalf("Error inesperado: esperado '%s', obtenido '%s'", expectedError, err.Error())
	}
}

func TestCambiarSentidoRutaInexistente(t *testing.T) {
	// Crear un grafo dirigido de prueba
	grafo := domain.NuevoGrafo(true)

	// Agregar cuevas
	cueva1 := domain.NuevaCueva("C1", "Cueva 1")
	cueva1.X, cueva1.Y = 0, 0
	cueva2 := domain.NuevaCueva("C2", "Cueva 2")
	cueva2.X, cueva2.Y = 1, 1
	grafo.AgregarCueva(cueva1)
	grafo.AgregarCueva(cueva2)

	// Crear servicio de conexión (sin agregar aristas)
	servicioConexion := service.NuevoServicioConexion(grafo)

	// Intentar cambiar sentido de una ruta que no existe
	solicitud := &service.CambiarSentidoRuta{
		DesdeCuevaID: "C1",
		HastaCuevaID: "C2",
	}

	err := servicioConexion.CambiarSentidoRuta(solicitud)
	if err == nil {
		t.Fatal("Se esperaba un error al intentar cambiar el sentido de una ruta inexistente")
	}

	expectedError := "conexión desde C1 hasta C2 no existe"
	if err.Error() != expectedError {
		t.Fatalf("Error inesperado: esperado '%s', obtenido '%s'", expectedError, err.Error())
	}
}

func TestCambiarSentidoMultiplesRutas(t *testing.T) {
	// Crear un grafo dirigido de prueba
	grafo := domain.NuevoGrafo(true)

	// Agregar cuevas
	cueva1 := domain.NuevaCueva("C1", "Cueva 1")
	cueva1.X, cueva1.Y = 0, 0
	cueva2 := domain.NuevaCueva("C2", "Cueva 2")
	cueva2.X, cueva2.Y = 1, 1
	cueva3 := domain.NuevaCueva("C3", "Cueva 3")
	cueva3.X, cueva3.Y = 2, 2
	grafo.AgregarCueva(cueva1)
	grafo.AgregarCueva(cueva2)
	grafo.AgregarCueva(cueva3)

	// Agregar aristas dirigidas
	arista1 := domain.NuevaArista("C1", "C2", 10.0, true)
	arista2 := domain.NuevaArista("C2", "C3", 15.0, true)
	grafo.AgregarArista(arista1)
	grafo.AgregarArista(arista2)

	// Crear servicio de conexión
	servicioConexion := service.NuevoServicioConexion(grafo)

	// Cambiar sentido de múltiples rutas
	solicitudes := []*service.CambiarSentidoRuta{
		{DesdeCuevaID: "C1", HastaCuevaID: "C2"},
		{DesdeCuevaID: "C2", HastaCuevaID: "C3"},
	}

	errores := servicioConexion.CambiarSentidoMultiplesRutas(solicitudes)
	if len(errores) != 0 {
		t.Fatalf("Se esperaban 0 errores, pero se obtuvieron %d: %v", len(errores), errores)
	}

	// Verificar que ambas aristas se invirtieron
	aristas := grafo.Aristas
	if len(aristas) != 2 {
		t.Fatalf("Se esperaban 2 aristas, pero se encontraron %d", len(aristas))
	}

	// Verificar primera arista invertida (C2 -> C1)
	arista1Invertida := false
	arista2Invertida := false

	for _, arista := range aristas {
		if arista.Desde == "C2" && arista.Hasta == "C1" && arista.Distancia == 10.0 {
			arista1Invertida = true
		}
		if arista.Desde == "C3" && arista.Hasta == "C2" && arista.Distancia == 15.0 {
			arista2Invertida = true
		}
	}

	if !arista1Invertida {
		t.Fatal("La primera arista (C1->C2) no se invirtió correctamente")
	}

	if !arista2Invertida {
		t.Fatal("La segunda arista (C2->C3) no se invirtió correctamente")
	}
}

func TestInvertirRutasDesdeCueva(t *testing.T) {
	// Crear un grafo dirigido de prueba
	grafo := domain.NuevoGrafo(true)

	// Agregar cuevas
	cueva1 := domain.NuevaCueva("C1", "Cueva 1")
	cueva1.X, cueva1.Y = 0, 0
	cueva2 := domain.NuevaCueva("C2", "Cueva 2")
	cueva2.X, cueva2.Y = 1, 1
	cueva3 := domain.NuevaCueva("C3", "Cueva 3")
	cueva3.X, cueva3.Y = 2, 2
	grafo.AgregarCueva(cueva1)
	grafo.AgregarCueva(cueva2)
	grafo.AgregarCueva(cueva3)

	// Agregar aristas dirigidas salientes de C1
	arista1 := domain.NuevaArista("C1", "C2", 10.0, true)
	arista2 := domain.NuevaArista("C1", "C3", 15.0, true)
	arista3 := domain.NuevaArista("C2", "C3", 20.0, true) // No debería ser afectada
	grafo.AgregarArista(arista1)
	grafo.AgregarArista(arista2)
	grafo.AgregarArista(arista3)

	// Crear servicio de conexión
	servicioConexion := service.NuevoServicioConexion(grafo)

	// Invertir todas las rutas salientes de C1
	err := servicioConexion.InvertirRutasDesdeCueva("C1")
	if err != nil {
		t.Fatalf("Error al invertir rutas desde cueva: %v", err)
	}

	// Verificar que las rutas se invirtieron correctamente
	aristas := grafo.Aristas
	if len(aristas) != 3 {
		t.Fatalf("Se esperaban 3 aristas, pero se encontraron %d", len(aristas))
	}

	c1aC2Invertida := false
	c1aC3Invertida := false
	c2aC3Intacta := false

	for _, arista := range aristas {
		if arista.Desde == "C2" && arista.Hasta == "C1" && arista.Distancia == 10.0 {
			c1aC2Invertida = true
		}
		if arista.Desde == "C3" && arista.Hasta == "C1" && arista.Distancia == 15.0 {
			c1aC3Invertida = true
		}
		if arista.Desde == "C2" && arista.Hasta == "C3" && arista.Distancia == 20.0 {
			c2aC3Intacta = true
		}
	}

	if !c1aC2Invertida {
		t.Fatal("La arista C1->C2 no se invirtió correctamente")
	}

	if !c1aC3Invertida {
		t.Fatal("La arista C1->C3 no se invirtió correctamente")
	}

	if !c2aC3Intacta {
		t.Fatal("La arista C2->C3 no debería haber sido afectada")
	}
}

func TestInvertirRutasHaciaCueva(t *testing.T) {
	// Crear un grafo dirigido de prueba
	grafo := domain.NuevoGrafo(true)

	// Agregar cuevas
	cueva1 := domain.NuevaCueva("C1", "Cueva 1")
	cueva1.X, cueva1.Y = 0, 0
	cueva2 := domain.NuevaCueva("C2", "Cueva 2")
	cueva2.X, cueva2.Y = 1, 1
	cueva3 := domain.NuevaCueva("C3", "Cueva 3")
	cueva3.X, cueva3.Y = 2, 2
	grafo.AgregarCueva(cueva1)
	grafo.AgregarCueva(cueva2)
	grafo.AgregarCueva(cueva3)

	// Agregar aristas dirigidas entrantes a C3
	arista1 := domain.NuevaArista("C1", "C3", 10.0, true)
	arista2 := domain.NuevaArista("C2", "C3", 15.0, true)
	arista3 := domain.NuevaArista("C1", "C2", 20.0, true) // No debería ser afectada
	grafo.AgregarArista(arista1)
	grafo.AgregarArista(arista2)
	grafo.AgregarArista(arista3)

	// Crear servicio de conexión
	servicioConexion := service.NuevoServicioConexion(grafo)

	// Invertir todas las rutas entrantes a C3
	err := servicioConexion.InvertirRutasHaciaCueva("C3")
	if err != nil {
		t.Fatalf("Error al invertir rutas hacia cueva: %v", err)
	}

	// Verificar que las rutas se invirtieron correctamente
	aristas := grafo.Aristas
	if len(aristas) != 3 {
		t.Fatalf("Se esperaban 3 aristas, pero se encontraron %d", len(aristas))
	}

	c1aC3Invertida := false
	c2aC3Invertida := false
	c1aC2Intacta := false

	for _, arista := range aristas {
		if arista.Desde == "C3" && arista.Hasta == "C1" && arista.Distancia == 10.0 {
			c1aC3Invertida = true
		}
		if arista.Desde == "C3" && arista.Hasta == "C2" && arista.Distancia == 15.0 {
			c2aC3Invertida = true
		}
		if arista.Desde == "C1" && arista.Hasta == "C2" && arista.Distancia == 20.0 {
			c1aC2Intacta = true
		}
	}

	if !c1aC3Invertida {
		t.Fatal("La arista C1->C3 no se invirtió correctamente")
	}

	if !c2aC3Invertida {
		t.Fatal("La arista C2->C3 no se invirtió correctamente")
	}

	if !c1aC2Intacta {
		t.Fatal("La arista C1->C2 no debería haber sido afectada")
	}
}
