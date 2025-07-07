package service

import (
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/pkg/algorithms"
	"testing"
)

// TestMSTServiceCompleto realiza pruebas completas del servicio MST
func TestMSTServiceCompleto(t *testing.T) {
	// Crear grafo de prueba
	grafo := domain.NuevoGrafo(false)

	// Agregar cuevas de prueba
	cuevas := map[string]string{
		"BASE":    "Base Central",
		"NORTE":   "Cueva Norte",
		"SUR":     "Cueva Sur",
		"ESTE":    "Cueva Este",
		"OESTE":   "Cueva Oeste",
		"AISLADA": "Cueva Aislada",
	}

	for id, nombre := range cuevas {
		cueva := &domain.Cueva{
			ID:       id,
			Nombre:   nombre,
			Recursos: map[string]int{"minerales": 10, "agua": 5},
		}
		grafo.Cuevas[id] = cueva
	}

	// Agregar aristas (AISLADA sin conexiones)
	aristas := []*domain.Arista{
		{Desde: "BASE", Hasta: "NORTE", Distancia: 10.0, EsDirigido: false, EsObstruido: false},
		{Desde: "BASE", Hasta: "SUR", Distancia: 8.0, EsDirigido: false, EsObstruido: false},
		{Desde: "BASE", Hasta: "ESTE", Distancia: 12.0, EsDirigido: false, EsObstruido: false},
		{Desde: "BASE", Hasta: "OESTE", Distancia: 15.0, EsDirigido: false, EsObstruido: false},
		{Desde: "NORTE", Hasta: "ESTE", Distancia: 6.0, EsDirigido: false, EsObstruido: false},
		{Desde: "SUR", Hasta: "OESTE", Distancia: 9.0, EsDirigido: false, EsObstruido: false},
		{Desde: "NORTE", Hasta: "SUR", Distancia: 20.0, EsDirigido: false, EsObstruido: false},
	}

	for _, arista := range aristas {
		grafo.Aristas = append(grafo.Aristas, arista)
	}

	// Crear servicio MST
	grafoSvc := NuevoServicioGrafo(grafo, nil)
	mstService := NuevoMSTService(grafoSvc)

	t.Run("MST desde nodo central", func(t *testing.T) {
		resultado, err := mstService.ObtenerMSTDesdeCueva(grafo, "BASE")
		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}

		// Verificaciones
		if resultado.CuevaOrigen != "BASE" {
			t.Errorf("Cueva origen esperada: BASE, obtenida: %s", resultado.CuevaOrigen)
		}

		// Debe alcanzar 5 de 6 cuevas (todas excepto AISLADA)
		expectedAlcanzables := 5
		if resultado.TotalAlcanzables != expectedAlcanzables {
			t.Errorf("Cuevas alcanzables esperadas: %d, obtenidas: %d", expectedAlcanzables, resultado.TotalAlcanzables)
		}

		// Debe tener exactamente 1 cueva no alcanzable (AISLADA)
		if resultado.TotalNoAlcanzable != 1 {
			t.Errorf("Cuevas no alcanzables esperadas: 1, obtenidas: %d", resultado.TotalNoAlcanzable)
		}

		// No debe ser completo por la cueva aislada
		if resultado.EsCompleto {
			t.Error("El MST no debería ser completo debido a cueva aislada")
		}

		// El MST debe tener 4 aristas para 5 nodos conectados
		if resultado.MST.MST.NumAristas != 4 {
			t.Errorf("Aristas en MST esperadas: 4, obtenidas: %d", resultado.MST.MST.NumAristas)
		}
	})

	t.Run("MST desde nodo periférico", func(t *testing.T) {
		resultado, err := mstService.ObtenerMSTDesdeCueva(grafo, "NORTE")
		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}

		// Debe alcanzar las mismas cuevas que desde BASE
		if resultado.TotalAlcanzables != 5 {
			t.Errorf("Cuevas alcanzables esperadas: 5, obtenidas: %d", resultado.TotalAlcanzables)
		}

		// El peso total debe ser el mismo (propiedad del MST)
		pesoEsperado := 33.0
		if resultado.MST.MST.PesoTotal != pesoEsperado {
			t.Errorf("Peso total esperado: %.1f, obtenido: %.1f", pesoEsperado, resultado.MST.MST.PesoTotal)
		}
	})

	t.Run("MST desde nodo aislado", func(t *testing.T) {
		resultado, err := mstService.ObtenerMSTDesdeCueva(grafo, "AISLADA")
		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}

		// Solo debe alcanzar a sí misma
		if resultado.TotalAlcanzables != 1 {
			t.Errorf("Cuevas alcanzables esperadas: 1, obtenidas: %d", resultado.TotalAlcanzables)
		}

		// Debe tener 5 cuevas no alcanzables
		if resultado.TotalNoAlcanzable != 5 {
			t.Errorf("Cuevas no alcanzables esperadas: 5, obtenidas: %d", resultado.TotalNoAlcanzable)
		}

		// El peso debe ser 0 (sin aristas)
		if resultado.MST.MST.PesoTotal != 0.0 {
			t.Errorf("Peso total esperado: 0.0, obtenido: %.1f", resultado.MST.MST.PesoTotal)
		}

		// Sin aristas en el MST
		if resultado.MST.MST.NumAristas != 0 {
			t.Errorf("Aristas esperadas: 0, obtenidas: %d", resultado.MST.MST.NumAristas)
		}
	})

	t.Run("Cueva inexistente", func(t *testing.T) {
		_, err := mstService.ObtenerMSTDesdeCueva(grafo, "INEXISTENTE")
		if err == nil {
			t.Error("Se esperaba error para cueva inexistente")
		}
	})

	t.Run("Verificar rutas mínimas", func(t *testing.T) {
		resultado, err := mstService.ObtenerMSTDesdeCueva(grafo, "BASE")
		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}

		// Verificar que se generaron rutas mínimas
		if len(resultado.RutasMinimas) == 0 {
			t.Error("Se esperaban rutas mínimas generadas")
		}

		// Verificar ruta directa BASE -> NORTE
		if ruta, existe := resultado.RutasMinimas["NORTE"]; !existe {
			t.Error("Se esperaba ruta a NORTE")
		} else if len(ruta) != 2 || ruta[0] != "BASE" || ruta[1] != "NORTE" {
			t.Errorf("Ruta BASE->NORTE esperada: [BASE NORTE], obtenida: %v", ruta)
		}
	})
}

// TestMSTSimple realiza una prueba simple del algoritmo de Prim
func TestMSTSimple(t *testing.T) {
	grafo := domain.NuevoGrafo(false)

	// Grafo simple de 4 nodos
	cuevas := []string{"A", "B", "C", "D"}
	for _, id := range cuevas {
		cueva := &domain.Cueva{
			ID:       id,
			Nombre:   "Cueva " + id,
			Recursos: map[string]int{"test": 1},
		}
		grafo.Cuevas[id] = cueva
	}

	// Aristas con pesos conocidos
	aristas := []*domain.Arista{
		{Desde: "A", Hasta: "B", Distancia: 5.0, EsDirigido: false, EsObstruido: false},
		{Desde: "A", Hasta: "C", Distancia: 3.0, EsDirigido: false, EsObstruido: false},
		{Desde: "B", Hasta: "C", Distancia: 2.0, EsDirigido: false, EsObstruido: false},
		{Desde: "B", Hasta: "D", Distancia: 4.0, EsDirigido: false, EsObstruido: false},
		{Desde: "C", Hasta: "D", Distancia: 6.0, EsDirigido: false, EsObstruido: false},
	}

	for _, arista := range aristas {
		grafo.Aristas = append(grafo.Aristas, arista)
	}

	// Probar desde cada nodo
	for _, origen := range cuevas {
		t.Run("MST desde "+origen, func(t *testing.T) {
			resultado, err := algorithms.Prim(grafo, origen)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}

			// Verificar propiedades básicas
			if !resultado.EsCompleto {
				t.Error("El MST debería ser completo para grafo conectado")
			}

			if len(resultado.Alcanzables) != 4 {
				t.Errorf("Cuevas alcanzables esperadas: 4, obtenidas: %d", len(resultado.Alcanzables))
			}

			if resultado.MST.NumAristas != 3 {
				t.Errorf("Aristas esperadas: 3, obtenidas: %d", resultado.MST.NumAristas)
			}

			// Peso total debe ser 9.0 para este grafo específico
			pesoEsperado := 9.0
			if resultado.MST.PesoTotal != pesoEsperado {
				t.Errorf("Peso total esperado: %.1f, obtenido: %.1f", pesoEsperado, resultado.MST.PesoTotal)
			}
		})
	}
}
