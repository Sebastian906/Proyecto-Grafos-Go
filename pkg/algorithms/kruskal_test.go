package algorithms

import (
	"proyecto-grafos-go/internal/domain"
	"testing"
)

func TestKruskalAlgorithm(t *testing.T) {
	// Test caso 1: Grafo simple conexo
	t.Run("Grafo simple conexo", func(t *testing.T) {
		grafo := domain.NuevoGrafo(false)

		// Agregar cuevas
		grafo.AgregarCueva(&domain.Cueva{ID: "A", Nombre: "Cueva A"})
		grafo.AgregarCueva(&domain.Cueva{ID: "B", Nombre: "Cueva B"})
		grafo.AgregarCueva(&domain.Cueva{ID: "C", Nombre: "Cueva C"})

		// Agregar aristas
		grafo.AgregarArista(&domain.Arista{Desde: "A", Hasta: "B", Distancia: 1.0, EsDirigido: false})
		grafo.AgregarArista(&domain.Arista{Desde: "B", Hasta: "C", Distancia: 2.0, EsDirigido: false})
		grafo.AgregarArista(&domain.Arista{Desde: "A", Hasta: "C", Distancia: 3.0, EsDirigido: false})

		mst, err := Kruskal(grafo)

		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}

		if mst.NumAristas != 2 {
			t.Errorf("Esperaba 2 aristas en MST, obtuvo %d", mst.NumAristas)
		}

		if mst.PesoTotal != 3.0 {
			t.Errorf("Esperaba peso total 3.0, obtuvo %.2f", mst.PesoTotal)
		}

		if mst.NumNodos != 3 {
			t.Errorf("Esperaba 3 nodos, obtuvo %d", mst.NumNodos)
		}
	})

	// Test caso 2: Grafo con aristas obstruidas
	t.Run("Grafo con aristas obstruidas", func(t *testing.T) {
		grafo := domain.NuevoGrafo(false)

		// Agregar cuevas
		grafo.AgregarCueva(&domain.Cueva{ID: "A", Nombre: "Cueva A"})
		grafo.AgregarCueva(&domain.Cueva{ID: "B", Nombre: "Cueva B"})
		grafo.AgregarCueva(&domain.Cueva{ID: "C", Nombre: "Cueva C"})

		// Agregar aristas (una obstruida)
		grafo.AgregarArista(&domain.Arista{Desde: "A", Hasta: "B", Distancia: 1.0, EsDirigido: false, EsObstruido: true})
		grafo.AgregarArista(&domain.Arista{Desde: "B", Hasta: "C", Distancia: 2.0, EsDirigido: false})
		grafo.AgregarArista(&domain.Arista{Desde: "A", Hasta: "C", Distancia: 3.0, EsDirigido: false})

		mst, err := Kruskal(grafo)

		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}

		// El MST debe usar solo las aristas no obstruidas
		if mst.PesoTotal != 5.0 {
			t.Errorf("Esperaba peso total 5.0 (ignorando arista obstruida), obtuvo %.2f", mst.PesoTotal)
		}
	})

	// Test caso 3: Grafo dirigido
	t.Run("Grafo dirigido", func(t *testing.T) {
		grafo := domain.NuevoGrafo(true)

		// Agregar cuevas
		grafo.AgregarCueva(&domain.Cueva{ID: "A", Nombre: "Cueva A"})
		grafo.AgregarCueva(&domain.Cueva{ID: "B", Nombre: "Cueva B"})
		grafo.AgregarCueva(&domain.Cueva{ID: "C", Nombre: "Cueva C"})

		// Agregar aristas dirigidas
		grafo.AgregarArista(&domain.Arista{Desde: "A", Hasta: "B", Distancia: 1.0, EsDirigido: true})
		grafo.AgregarArista(&domain.Arista{Desde: "B", Hasta: "A", Distancia: 1.0, EsDirigido: true})
		grafo.AgregarArista(&domain.Arista{Desde: "B", Hasta: "C", Distancia: 2.0, EsDirigido: true})
		grafo.AgregarArista(&domain.Arista{Desde: "C", Hasta: "B", Distancia: 2.0, EsDirigido: true})
		grafo.AgregarArista(&domain.Arista{Desde: "A", Hasta: "C", Distancia: 3.0, EsDirigido: true})
		grafo.AgregarArista(&domain.Arista{Desde: "C", Hasta: "A", Distancia: 3.0, EsDirigido: true})

		mst, err := Kruskal(grafo)

		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}

		// Para grafos dirigidos, Kruskal debe crear un MST no dirigido
		if mst.NumAristas != 2 {
			t.Errorf("Esperaba 2 aristas en MST, obtuvo %d", mst.NumAristas)
		}
	})

	// Test caso 4: Grafo vacio
	t.Run("Grafo vacio", func(t *testing.T) {
		grafo := domain.NuevoGrafo(false)

		mst, err := Kruskal(grafo)

		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}

		if len(mst.Aristas) != 0 {
			t.Errorf("Esperaba 0 aristas para grafo vacio, obtuvo %d", len(mst.Aristas))
		}

		if mst.PesoTotal != 0.0 {
			t.Errorf("Esperaba peso total 0.0 para grafo vacio, obtuvo %.2f", mst.PesoTotal)
		}
	})
}

func TestUnionFind(t *testing.T) {
	t.Run("Union Find básico", func(t *testing.T) {
		uf := NuevoUnionFind()

		// Crear conjuntos
		uf.MakeSet("A")
		uf.MakeSet("B")
		uf.MakeSet("C")

		// Verificar que están en conjuntos separados
		if uf.Find("A") == uf.Find("B") {
			t.Error("A y B no deberían estar en el mismo conjunto inicialmente")
		}

		// Unir A y B
		if !uf.Union("A", "B") {
			t.Error("Union debería retornar true para conjuntos separados")
		}

		// Verificar que ahora están unidos
		if uf.Find("A") != uf.Find("B") {
			t.Error("A y B deberían estar en el mismo conjunto después de Union")
		}

		// Intentar unir A y B nuevamente
		if uf.Union("A", "B") {
			t.Error("Union debería retornar false para elementos ya unidos")
		}
	})
}

func TestEsConexo(t *testing.T) {
	t.Run("Grafo conexo", func(t *testing.T) {
		grafo := domain.NuevoGrafo(false)

		// Agregar cuevas
		grafo.AgregarCueva(&domain.Cueva{ID: "A", Nombre: "Cueva A"})
		grafo.AgregarCueva(&domain.Cueva{ID: "B", Nombre: "Cueva B"})
		grafo.AgregarCueva(&domain.Cueva{ID: "C", Nombre: "Cueva C"})

		// Agregar aristas que conectan todo
		grafo.AgregarArista(&domain.Arista{Desde: "A", Hasta: "B", Distancia: 1.0})
		grafo.AgregarArista(&domain.Arista{Desde: "B", Hasta: "C", Distancia: 1.0})

		if !EsConexo(grafo) {
			t.Error("El grafo debería ser conexo")
		}
	})

	t.Run("Grafo desconexo", func(t *testing.T) {
		grafo := domain.NuevoGrafo(false)

		// Agregar cuevas
		grafo.AgregarCueva(&domain.Cueva{ID: "A", Nombre: "Cueva A"})
		grafo.AgregarCueva(&domain.Cueva{ID: "B", Nombre: "Cueva B"})
		grafo.AgregarCueva(&domain.Cueva{ID: "C", Nombre: "Cueva C"})

		// Solo conectar A y B, dejando C aislado
		grafo.AgregarArista(&domain.Arista{Desde: "A", Hasta: "B", Distancia: 1.0})

		if EsConexo(grafo) {
			t.Error("El grafo no debería ser conexo")
		}
	})

	t.Run("Grafo con aristas obstruidas", func(t *testing.T) {
		grafo := domain.NuevoGrafo(false)

		// Agregar cuevas
		grafo.AgregarCueva(&domain.Cueva{ID: "A", Nombre: "Cueva A"})
		grafo.AgregarCueva(&domain.Cueva{ID: "B", Nombre: "Cueva B"})

		// Agregar arista obstruida
		grafo.AgregarArista(&domain.Arista{Desde: "A", Hasta: "B", Distancia: 1.0, EsObstruido: true})

		if EsConexo(grafo) {
			t.Error("El grafo no debería ser conexo con aristas obstruidas")
		}
	})
}
