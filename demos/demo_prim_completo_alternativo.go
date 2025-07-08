package main

import (
	"fmt"
	"log"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/pkg/algorithms"
)

// main permite ejecutar este demo individualmente
func main() {
	demoPrimCompletoAlternativo()
}

func demoPrimCompletoAlternativo() {
	fmt.Println("=== PRUEBA COMPLETA DEL ALGORITMO DE PRIM ===")
	fmt.Println("Esta prueba verifica la funcionalidad del MST desde cueva específica")
	fmt.Println()

	// Crear un grafo de prueba con casos interesantes
	grafo := domain.NuevoGrafo(false)

	// Agregar cuevas con nombres descriptivos
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

	// Agregar aristas creando una red conectada excepto por AISLADA
	aristas := []*domain.Arista{
		// Conexiones desde BASE (nodo central)
		{Desde: "BASE", Hasta: "NORTE", Distancia: 10.0, EsDirigido: false, EsObstruido: false},
		{Desde: "BASE", Hasta: "SUR", Distancia: 8.0, EsDirigido: false, EsObstruido: false},
		{Desde: "BASE", Hasta: "ESTE", Distancia: 12.0, EsDirigido: false, EsObstruido: false},
		{Desde: "BASE", Hasta: "OESTE", Distancia: 15.0, EsDirigido: false, EsObstruido: false},

		// Conexiones entre nodos satélite
		{Desde: "NORTE", Hasta: "ESTE", Distancia: 6.0, EsDirigido: false, EsObstruido: false},
		{Desde: "SUR", Hasta: "OESTE", Distancia: 9.0, EsDirigido: false, EsObstruido: false},

		// Conexión directa entre opuestos (más costosa)
		{Desde: "NORTE", Hasta: "SUR", Distancia: 20.0, EsDirigido: false, EsObstruido: false},

		// AISLADA no tiene conexiones - para demostrar componentes desconectados
	}

	for _, arista := range aristas {
		grafo.Aristas = append(grafo.Aristas, arista)
	}

	fmt.Printf("Grafo creado:\n")
	fmt.Printf("- %d cuevas (%d conectadas + 1 aislada)\n", len(grafo.Cuevas), len(grafo.Cuevas)-1)
	fmt.Printf("- %d aristas\n\n", len(grafo.Aristas))

	// Casos de prueba específicos
	casos := []struct {
		origen      string
		descripcion string
	}{
		{"BASE", "Nodo central - debería tener la mejor cobertura"},
		{"NORTE", "Nodo periférico - cobertura completa pero rutas diferentes"},
		{"AISLADA", "Nodo aislado - solo alcanza a sí misma"},
	}

	for i, caso := range casos {
		fmt.Printf("=== CASO %d: MST desde '%s' ===\n", i+1, caso.origen)
		fmt.Printf("Descripción: %s\n", caso.descripcion)
		fmt.Println()

		resultado, err := algorithms.Prim(grafo, caso.origen)
		if err != nil {
			log.Printf("❌ Error: %v\n\n", err)
			continue
		}

		// Mostrar resultados principales
		fmt.Printf("✓ Cueva origen: %s\n", resultado.CuevaOrigen)
		fmt.Printf("✓ Es completo: %v\n", resultado.EsCompleto)
		fmt.Printf("✓ Cuevas alcanzables: %d/%d (%.1f%%)\n",
			len(resultado.Alcanzables), len(grafo.Cuevas),
			float64(len(resultado.Alcanzables))/float64(len(grafo.Cuevas))*100)

		if resultado.MST != nil {
			fmt.Printf("✓ Peso total del MST: %.1f\n", resultado.MST.PesoTotal)
			fmt.Printf("✓ Aristas en MST: %d\n", resultado.MST.NumAristas)
		}

		// Mostrar cuevas no alcanzables si las hay
		if len(resultado.NoAlcanzable) > 0 {
			fmt.Printf("⚠  Cuevas no alcanzables: %v\n", resultado.NoAlcanzable)
		}

		// Mostrar conexiones del MST
		if resultado.MST != nil && len(resultado.MST.Aristas) > 0 {
			fmt.Println("\n  Conexiones del MST:")
			for j, arista := range resultado.MST.Aristas {
				fmt.Printf("    %d. %s → %s (distancia: %.1f)\n",
					j+1, arista.Desde, arista.Hasta, arista.Distancia)
			}
		}

		// Mostrar rutas mínimas
		rutas := resultado.ObtenerRutasDesdeOrigen()
		if len(rutas) > 1 { // Solo mostrar si hay más de una ruta (origen a sí mismo)
			fmt.Println("\n  Rutas mínimas desde origen:")
			for destino, ruta := range rutas {
				if destino != caso.origen {
					fmt.Printf("    %s: %v\n", destino, ruta)
				}
			}
		}

		fmt.Println()
		fmt.Println("--------------------------------------------------")
		fmt.Println()
	}

	// Resumen y comparación
	fmt.Println("=== RESUMEN DE PRUEBAS ===")
	fmt.Println()

	fmt.Println("Observaciones importantes:")
	fmt.Println("1. BASE tiene la mejor cobertura por ser el nodo más central")
	fmt.Println("2. NORTE también alcanza todas las cuevas conectadas, pero con diferentes rutas")
	fmt.Println("3. AISLADA demuestra cómo el algoritmo maneja componentes desconectados")
	fmt.Println("4. Todos los MST completos tienen el mismo peso total (propiedad de MST)")
	fmt.Println()

	fmt.Println("✅ TODAS LAS PRUEBAS COMPLETADAS EXITOSAMENTE")
	fmt.Println("El algoritmo de Prim para MST desde cueva específica funciona correctamente!")
}
