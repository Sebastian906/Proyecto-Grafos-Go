package main

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/repository"
	"proyecto-grafos-go/internal/service"
)

// Programa de prueba rápida para verificar funcionalidad con datos de 9 cuevas
func main() {
	fmt.Println("=== PRUEBA RAPIDA: DETECCION DE CUEVAS INACCESIBLES ===")
	fmt.Println("")

	// Inicializar sistema
	grafo := domain.NuevoGrafo(false)
	repo := repository.NuevoRepositorio("../../data/")
	grafoSvc := service.NuevoServicioGrafo(grafo, repo)
	validacionSvc := service.NuevoServicioValidacion(grafo)

	// Cargar datos de 9 cuevas
	fmt.Println("1. Cargando datos de ejemplo (9 cuevas)...")
	if err := grafoSvc.CargarGrafo("caves_directed_example.json"); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	fmt.Printf("✓ Cargadas %d cuevas y %d conexiones\n", len(grafo.Cuevas), len(grafo.Aristas))

	// Mostrar cuevas cargadas
	fmt.Println("\n2. Cuevas del sistema:")
	for id, cueva := range grafo.Cuevas {
		fmt.Printf("   - %s: %s\n", id, cueva.Nombre)
	}

	// Análisis desde diferentes puntos
	cuevasPrueba := []string{"Silvestre", "Tazmania", "Correcaminos", "Popeye"}

	for _, cuevaInicio := range cuevasPrueba {
		fmt.Printf("\n3. Análisis de accesibilidad desde '%s':\n", cuevaInicio)
		resultado := validacionSvc.AnalizarAccesibilidad(cuevaInicio)

		fmt.Printf("   Total cuevas: %d\n", resultado.TotalCuevas)
		fmt.Printf("   Accesibles: %d\n", resultado.CuevasAccesibles)
		fmt.Printf("   Inaccesibles: %d\n", len(resultado.CuevasInaccesibles))

		if len(resultado.CuevasInaccesibles) > 0 {
			fmt.Printf("   Cuevas inaccesibles: %v\n", resultado.CuevasInaccesibles)
			fmt.Println("   Primeras 3 soluciones:")
			for i, solucion := range resultado.Soluciones {
				if i >= 3 {
					break
				}
				fmt.Printf("   - %s\n", solucion)
			}
		} else {
			fmt.Println("   ✓ Todas las cuevas son accesibles")
		}
	}

	// Probar escenario con obstrucción
	fmt.Println("\n4. Probando escenario con obstrucción...")
	conexionSvc := service.NuevoServicioConexion(grafo)

	// Obstruir conexión crítica
	solicitud := &service.ObstruirConexion{
		DesdeCuevaID: "Silvestre",
		HastaCuevaID: "Tazmania",
		EsObstruido:  true,
	}

	if err := conexionSvc.ObstruirConexion(solicitud); err != nil {
		fmt.Printf("ERROR al obstruir: %v\n", err)
	} else {
		fmt.Println("   ✓ Conexión Silvestre → Tazmania obstruida")

		// Analizar impacto
		resultado := validacionSvc.AnalizarAccesibilidad("Silvestre")
		fmt.Printf("   Ahora desde Silvestre hay %d cuevas inaccesibles: %v\n",
			len(resultado.CuevasInaccesibles), resultado.CuevasInaccesibles)
	}

	fmt.Println("\n=== PRUEBA COMPLETADA EXITOSAMENTE ===")
	fmt.Println("El sistema está listo para usar con las 9 cuevas!")
}
