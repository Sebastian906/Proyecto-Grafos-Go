package algorithms

import (
	"proyecto-grafos-go/internal/domain"
	"testing"
)

func TestPrim(t *testing.T) {
	// Crear un grafo de prueba simple
	grafo := domain.NuevoGrafo(false)

	// Agregar cuevas
	cuevas := []*domain.Cueva{
		{ID: "A", Nombre: "Cueva A"},
		{ID: "B", Nombre: "Cueva B"},
		{ID: "C", Nombre: "Cueva C"},
		{ID: "D", Nombre: "Cueva D"},
	}

	for _, cueva := range cuevas {
		grafo.Cuevas[cueva.ID] = cueva
	}

	// Agregar aristas
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

	// Ejecutar Prim desde cueva A
	resultado, err := Prim(grafo, "A")

	// Verificar que no hay error
	if err != nil {
		t.Fatalf("Error en algoritmo Prim: %v", err)
	}

	// Verificar que el resultado es válido
	if resultado == nil {
		t.Fatal("Resultado es nil")
	}

	// Verificar cueva origen
	if resultado.CuevaOrigen != "A" {
		t.Errorf("Cueva origen esperada: A, obtenida: %s", resultado.CuevaOrigen)
	}

	// Verificar que es completo (todas las cuevas son alcanzables)
	if !resultado.EsCompleto {
		t.Error("El MST debería ser completo para este grafo conectado")
	}

	// Verificar número de cuevas alcanzables
	if len(resultado.Alcanzables) != 4 {
		t.Errorf("Cuevas alcanzables esperadas: 4, obtenidas: %d", len(resultado.Alcanzables))
	}

	// Verificar que el MST tiene 3 aristas (n-1 para n nodos)
	if resultado.MST.NumAristas != 3 {
		t.Errorf("Aristas en MST esperadas: 3, obtenidas: %d", resultado.MST.NumAristas)
	}

	// Verificar peso total (debería ser 5 + 2 + 4 = 11 para este grafo específico)
	pesoEsperado := 9.0 // A-C(3) + B-C(2) + B-D(4) = 9
	if resultado.MST.PesoTotal != pesoEsperado {
		t.Errorf("Peso total esperado: %.1f, obtenido: %.1f", pesoEsperado, resultado.MST.PesoTotal)
	}
}

func TestPrimCuevaInexistente(t *testing.T) {
	grafo := domain.NuevoGrafo(false)
	grafo.Cuevas["A"] = &domain.Cueva{ID: "A", Nombre: "Cueva A"}

	// Intentar Prim desde cueva que no existe
	_, err := Prim(grafo, "X")

	if err == nil {
		t.Error("Se esperaba error al usar cueva inexistente")
	}
}

func TestPrimGrafoDesconectado(t *testing.T) {
	grafo := domain.NuevoGrafo(false)

	// Crear dos componentes separados
	cuevas := []*domain.Cueva{
		{ID: "A", Nombre: "Cueva A"},
		{ID: "B", Nombre: "Cueva B"},
		{ID: "C", Nombre: "Cueva C"}, // Componente aislado
	}

	for _, cueva := range cuevas {
		grafo.Cuevas[cueva.ID] = cueva
	}

	// Solo conectar A-B, dejando C aislado
	grafo.Aristas = append(grafo.Aristas, &domain.Arista{
		Desde: "A", Hasta: "B", Distancia: 5.0, EsDirigido: false, EsObstruido: false,
	})

	// Ejecutar Prim desde A
	resultado, err := Prim(grafo, "A")

	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	// Verificar que no es completo
	if resultado.EsCompleto {
		t.Error("El MST no debería ser completo para grafo desconectado")
	}

	// Verificar que solo alcanza A y B
	if len(resultado.Alcanzables) != 2 {
		t.Errorf("Cuevas alcanzables esperadas: 2, obtenidas: %d", len(resultado.Alcanzables))
	}

	// Verificar que C no es alcanzable
	if len(resultado.NoAlcanzable) != 1 || resultado.NoAlcanzable[0] != "C" {
		t.Errorf("Cueva C debería estar en no alcanzables: %v", resultado.NoAlcanzable)
	}
}

func TestPrimObtenerRutas(t *testing.T) {
	grafo := domain.NuevoGrafo(false)

	// Crear grafo lineal A-B-C
	cuevas := []*domain.Cueva{
		{ID: "A", Nombre: "Cueva A"},
		{ID: "B", Nombre: "Cueva B"},
		{ID: "C", Nombre: "Cueva C"},
	}

	for _, cueva := range cuevas {
		grafo.Cuevas[cueva.ID] = cueva
	}

	aristas := []*domain.Arista{
		{Desde: "A", Hasta: "B", Distancia: 2.0, EsDirigido: false, EsObstruido: false},
		{Desde: "B", Hasta: "C", Distancia: 3.0, EsDirigido: false, EsObstruido: false},
	}

	for _, arista := range aristas {
		grafo.Aristas = append(grafo.Aristas, arista)
	}

	// Ejecutar Prim desde A
	resultado, err := Prim(grafo, "A")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Obtener rutas
	rutas := resultado.ObtenerRutasDesdeOrigen()

	// Verificar ruta A->A
	if ruta, existe := rutas["A"]; !existe || len(ruta) != 1 || ruta[0] != "A" {
		t.Errorf("Ruta A->A esperada: [A], obtenida: %v", ruta)
	}

	// Verificar ruta A->B
	if ruta, existe := rutas["B"]; !existe || len(ruta) != 2 || ruta[0] != "A" || ruta[1] != "B" {
		t.Errorf("Ruta A->B esperada: [A B], obtenida: %v", ruta)
	}

	// Verificar ruta A->C (debería ser A->B->C)
	if ruta, existe := rutas["C"]; !existe || len(ruta) != 3 || ruta[0] != "A" || ruta[1] != "B" || ruta[2] != "C" {
		t.Errorf("Ruta A->C esperada: [A B C], obtenida: %v", ruta)
	}
}
