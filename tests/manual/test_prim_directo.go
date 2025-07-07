package main

import (
	"fmt"
	"log"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/pkg/algorithms"
)

// testPrimDirecto ejecuta una prueba directa y simple del algoritmo
func testPrimDirecto() {
	fmt.Println("=== PRUEBA DIRECTA DEL ALGORITMO DE PRIM ===")

	// Crear un grafo de prueba simple
	grafo := domain.NuevoGrafo(false)

	// Agregar cuevas
	cuevas := map[string]string{
		"A": "Cueva Alpha",
		"B": "Cueva Beta",
		"C": "Cueva Gamma",
		"D": "Cueva Delta",
	}

	for id, nombre := range cuevas {
		cueva := &domain.Cueva{
			ID:       id,
			Nombre:   nombre,
			Recursos: map[string]int{"test": 1},
		}
		grafo.Cuevas[id] = cueva
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

	fmt.Printf("Grafo creado con %d cuevas y %d aristas\n\n", len(grafo.Cuevas), len(grafo.Aristas))

	// Probar MST desde cada cueva
	for origen := range cuevas {
		fmt.Printf("--- MST desde cueva '%s' ---\n", origen)

		resultado, err := algorithms.Prim(grafo, origen)
		if err != nil {
			log.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("✓ Cueva origen: %s\n", resultado.CuevaOrigen)
		fmt.Printf("✓ Es completo: %v\n", resultado.EsCompleto)
		fmt.Printf("✓ Cuevas alcanzables: %d/%d\n", len(resultado.Alcanzables), len(grafo.Cuevas))
		fmt.Printf("✓ Peso total: %.1f\n", resultado.MST.PesoTotal)
		fmt.Printf("✓ Aristas en MST: %d\n", resultado.MST.NumAristas)

		if len(resultado.MST.Aristas) > 0 {
			fmt.Println("  Conexiones:")
			for i, arista := range resultado.MST.Aristas {
				fmt.Printf("    %d. %s -> %s (%.1f)\n", i+1, arista.Desde, arista.Hasta, arista.Distancia)
			}
		}

		// Obtener rutas
		rutas := resultado.ObtenerRutasDesdeOrigen()
		fmt.Println("  Rutas desde origen:")
		for destino, ruta := range rutas {
			if destino != origen {
				fmt.Printf("    %s: %v\n", destino, ruta)
			}
		}

		fmt.Println()
	}

	fmt.Println("=== PRUEBA COMPLETADA ===")
	fmt.Println("Si ves este mensaje sin errores, el algoritmo funciona correctamente!")
}

// Esta función ya no es necesaria, la funcionalidad está en test_prim_completo.go
// Manteniendo solo como referencia
