package integration

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/repository"
	"proyecto-grafos-go/internal/service"
	"testing"
)

// TestRequirement2a prueba el requisito 2a: Obstrucción de caminos
func TestRequirement2a(t *testing.T) {
	fmt.Println("=== Prueba del Requisito 2a: Obstrucción de Caminos ===")

	// Inicialización
	grafo := domain.NuevoGrafo(false)
	repo := repository.NuevoRepositorio("../../data/")

	// Servicios
	grafoSvc := service.NuevoServicioGrafo(grafo, repo)
	conexionSvc := service.NuevoServicioConexion(grafo)

	// Cargar el grafo
	fmt.Println("\n1. Cargando grafo desde archivo...")
	if err := grafoSvc.CargarGrafo("caves_with_obstacles.json"); err != nil {
		t.Fatalf("Error cargando archivo: %v", err)
	}
	fmt.Println("Grafo cargado exitosamente")

	// Mostrar estadísticas iniciales
	fmt.Println("\n2. Estadísticas iniciales:")
	stats := conexionSvc.EstadisticasConexiones()
	fmt.Printf("- Total conexiones: %v\n", stats["total_conexiones"])
	fmt.Printf("- Conexiones activas: %v\n", stats["conexiones_activas"])
	fmt.Printf("- Conexiones obstruidas: %v\n", stats["conexiones_obstruidas"])

	// Verificar que hay al menos una conexión obstruida pre-cargada
	if stats["conexiones_obstruidas"].(int) == 0 {
		t.Error("Se esperaba al menos una conexión obstruida en el archivo de prueba")
	}

	// Listar conexiones obstruidas
	fmt.Println("\n3. Conexiones obstruidas pre-cargadas:")
	conexionesObstruidas := conexionSvc.ListarConexionesObstruidas()
	if len(conexionesObstruidas) == 0 {
		fmt.Println("   No hay conexiones obstruidas")
	} else {
		for i, conn := range conexionesObstruidas {
			fmt.Printf("   %d. %s ↔ %s (Distancia: %.2f)\n",
				i+1, conn["desde"], conn["hasta"], conn["distancia"])
		}
	}

	// Obstruir una nueva conexión
	fmt.Println("\n4. Obstruyendo la conexión CR → Silvestre...")
	solicitud := &service.ObstruirConexion{
		DesdeCuevaID: "CR",
		HastaCuevaID: "Silvestre",
		EsObstruido:  true,
	}

	if err := conexionSvc.ObstruirConexion(solicitud); err != nil {
		t.Errorf("Error obstruyendo conexión: %v", err)
	} else {
		fmt.Println("   Conexión obstruida exitosamente")
	}

	// Verificar que la obstrucción funcionó
	statsIntermedias := conexionSvc.EstadisticasConexiones()
	conexionesObstruidasAnteriores := stats["conexiones_obstruidas"].(int)
	conexionesObstruidasActuales := statsIntermedias["conexiones_obstruidas"].(int)

	// En grafos no dirigidos, obstruir una conexión obstruye ambas direcciones
	expectedIncrease := 2 // Se obstruyen ambas direcciones
	if conexionesObstruidasActuales != conexionesObstruidasAnteriores+expectedIncrease {
		t.Errorf("Se esperaba %d conexiones obstruidas, se encontraron %d (grafo no dirigido obstruye ambas direcciones)",
			conexionesObstruidasAnteriores+expectedIncrease, conexionesObstruidasActuales)
	}

	// Mostrar estadísticas después de obstruir
	fmt.Println("\n5. Estadísticas después de obstruir:")
	fmt.Printf("- Total conexiones: %v\n", statsIntermedias["total_conexiones"])
	fmt.Printf("- Conexiones activas: %v\n", statsIntermedias["conexiones_activas"])
	fmt.Printf("- Conexiones obstruidas: %v\n", statsIntermedias["conexiones_obstruidas"])

	// Probar obstrucción múltiple
	fmt.Println("\n6. Probando obstrucción múltiple...")
	solicitudesMultiples := []*service.ObstruirConexion{
		{DesdeCuevaID: "Tazmania", HastaCuevaID: "Coyote", EsObstruido: true},
		{DesdeCuevaID: "Bunny", HastaCuevaID: "Marvin", EsObstruido: true},
	}

	errores := conexionSvc.ObstruirMultiplesConexiones(solicitudesMultiples)
	if len(errores) == 0 {
		fmt.Println("   Todas las conexiones múltiples fueron obstruidas")
	} else {
		t.Errorf("Se encontraron %d errores en obstrucción múltiple:", len(errores))
		for _, err := range errores {
			t.Errorf("   - %v", err)
		}
	}

	// Probar desobstrucción
	fmt.Println("\n7. Probando desobstrucción...")
	solicitudDesobstruir := &service.ObstruirConexion{
		DesdeCuevaID: "CR",
		HastaCuevaID: "Silvestre",
		EsObstruido:  false,
	}

	if err := conexionSvc.ObstruirConexion(solicitudDesobstruir); err != nil {
		t.Errorf("Error desobstruyendo conexión: %v", err)
	} else {
		fmt.Println("   Conexión desobstruida exitosamente")
	}

	// Estadísticas finales
	fmt.Println("\n8. Estadísticas finales:")
	statsFinals := conexionSvc.EstadisticasConexiones()
	fmt.Printf("- Conexiones obstruidas: %v\n", statsFinals["conexiones_obstruidas"])

	// Listar todas las conexiones obstruidas finales
	fmt.Println("\n9. Conexiones obstruidas finales:")
	conexionesObstruidasFinales := conexionSvc.ListarConexionesObstruidas()
	for i, conn := range conexionesObstruidasFinales {
		fmt.Printf("   %d. %s ↔ %s (Distancia: %.2f)\n",
			i+1, conn["desde"], conn["hasta"], conn["distancia"])
	}

	fmt.Println("\n=== Prueba del Requisito 2a COMPLETADA ===")
	fmt.Println("Se puede obstruir conexiones individuales")
	fmt.Println("Se puede obstruir múltiples conexiones")
	fmt.Println("Se puede desobstruir conexiones")
	fmt.Println("Se pueden ver estadísticas de obstrucción")
	fmt.Println("Se pueden listar conexiones obstruidas")
}

// BenchmarkObstruccionConexiones benchmarks para medir el rendimiento
func BenchmarkObstruccionConexiones(b *testing.B) {
	// Inicialización
	grafo := domain.NuevoGrafo(false)
	repo := repository.NuevoRepositorio("../../data/")

	// Servicios
	grafoSvc := service.NuevoServicioGrafo(grafo, repo)
	conexionSvc := service.NuevoServicioConexion(grafo)

	// Cargar el grafo
	if err := grafoSvc.CargarGrafo("caves_with_obstacles.json"); err != nil {
		b.Fatalf("Error cargando archivo: %v", err)
	}

	// Benchmark de obstrucción de conexiones
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solicitud := &service.ObstruirConexion{
			DesdeCuevaID: "CR",
			HastaCuevaID: "Silvestre",
			EsObstruido:  i%2 == 0, // Alternar entre obstruir y desobstruir
		}
		conexionSvc.ObstruirConexion(solicitud)
	}
}
