package main

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/repository"
	"proyecto-grafos-go/internal/service"
)

func main() {
	fmt.Println("=== Prueba del Archivo caves_directed_example.json ===")

	// Inicializar repositorio y grafo
	grafo := domain.NuevoGrafo(true)
	repo := repository.NuevoRepositorio("../../../../data/")

	// Cargar el archivo
	fmt.Println("\n1. Cargando archivo caves_directed_example.json...")
	grafoNuevo, err := repo.CargarJSON("caves_directed_example.json")
	if err != nil {
		fmt.Printf("Error al cargar el archivo: %v\n", err)
		return
	}
	grafo = grafoNuevo
	fmt.Println("Archivo cargado exitosamente!")

	// Mostrar información básica del grafo
	fmt.Printf("\n2. Información del grafo:")
	fmt.Printf("\n   - Tipo: %s", func() string {
		if grafo.EsDirigido {
			return "Dirigido"
		}
		return "No dirigido"
	}())
	fmt.Printf("\n   - Número de cuevas: %d", len(grafo.Cuevas))
	fmt.Printf("\n   - Número de conexiones: %d", len(grafo.Aristas))

	// Listar todas las cuevas
	fmt.Println("\n\n3. Cuevas en el grafo:")
	for id, cueva := range grafo.Cuevas {
		fmt.Printf("   - %s: %s (%.1f, %.1f)\n", id, cueva.Nombre, cueva.X, cueva.Y)
	}

	// Crear servicio de conexión
	conexionSvc := service.NuevoServicioConexion(grafo)

	// Mostrar todas las conexiones
	fmt.Println("\n4. Conexiones originales:")
	conexiones := conexionSvc.ListarConexiones()
	for i, conn := range conexiones {
		fmt.Printf("   %d. %s → %s (distancia: %.0f)\n",
			i+1, conn["desde"], conn["hasta"], conn["distancia"])
	}

	// Probar cambio de sentido de una ruta
	fmt.Println("\n5. Probando cambio de sentido de ruta...")
	fmt.Println("   Cambiando: Silvestre → Tasmania  a  Tasmania → Silvestre")

	solicitud := &service.CambiarSentidoRuta{
		DesdeCuevaID: "Silvestre",
		HastaCuevaID: "Tasmania",
	}

	err = conexionSvc.CambiarSentidoRuta(solicitud)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Println("   Cambio exitoso!")

		// Mostrar conexiones después del cambio
		fmt.Println("\n6. Conexiones después del cambio:")
		conexionesNuevas := conexionSvc.ListarConexiones()
		for i, conn := range conexionesNuevas {
			prefix := "   "
			if conn["desde"] == "Tasmania" && conn["hasta"] == "Silvestre" {
				prefix = "   > "
			}
			fmt.Printf("%s%d. %s → %s (distancia: %.0f)\n",
				prefix, i+1, conn["desde"], conn["hasta"], conn["distancia"])
		}
	}

	// Probar inversión de múltiples rutas desde Bunny
	fmt.Println("\n7. Probando inversión de rutas salientes de Bunny...")
	err = conexionSvc.InvertirRutasDesdeCueva("Bunny")
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Println("   Inversión exitosa!")
		fmt.Println("   Las rutas desde Bunny ahora van hacia Bunny")
	}

	// Mostrar estadísticas finales
	fmt.Println("\n8. Estadísticas finales:")
	stats := conexionSvc.EstadisticasConexiones()
	fmt.Printf("   - Total conexiones: %v\n", stats["total_conexiones"])
	fmt.Printf("   - Conexiones activas: %v\n", stats["conexiones_activas"])
	fmt.Printf("   - Conexiones dirigidas: %v\n", stats["conexiones_dirigidas"])

	fmt.Println("\n=== Prueba completada ===")
}
