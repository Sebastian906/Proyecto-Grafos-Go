package performance

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/pkg/algorithms"
	"proyecto-grafos-go/pkg/utils"
	"testing"
)

// BenchmarkDFS prueba el rendimiento del algoritmo DFS
func BenchmarkDFS(b *testing.B) {
	// Crear grafo de prueba con diferentes tamaños
	sizes := []int{10, 50, 100, 500}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("DFS_%d_nodos", size), func(b *testing.B) {
			// Crear grafo con 'size' nodos
			grafo := crearGrafoCompleto(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Ejecutar DFS desde el primer nodo
				visited := make(map[string]bool)
				algorithms.DFS(grafo, "C0", visited)
			}
		})
	}
}

// BenchmarkBFS prueba el rendimiento del algoritmo BFS
func BenchmarkBFS(b *testing.B) {
	sizes := []int{10, 50, 100, 500}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("BFS_%d_nodos", size), func(b *testing.B) {
			grafo := crearGrafoCompleto(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				visited := make(map[string]bool)
				algorithms.BFS(grafo, "C0", visited)
			}
		})
	}
}

// BenchmarkDijkstra prueba el rendimiento del algoritmo de Dijkstra
func BenchmarkDijkstra(b *testing.B) {
	sizes := []int{10, 50, 100, 200}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Dijkstra_%d_nodos", size), func(b *testing.B) {
			grafo := crearGrafoCompleto(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				algorithms.Dijkstra(grafo, "C0", "C1")
			}
		})
	}
}

// BenchmarkPrim prueba el rendimiento del algoritmo de Prim
func BenchmarkPrim(b *testing.B) {
	sizes := []int{10, 50, 100, 200}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Prim_%d_nodos", size), func(b *testing.B) {
			grafo := crearGrafoCompleto(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				algorithms.Prim(grafo)
			}
		})
	}
}

// BenchmarkKruskal prueba el rendimiento del algoritmo de Kruskal
func BenchmarkKruskal(b *testing.B) {
	sizes := []int{10, 50, 100, 200}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Kruskal_%d_nodos", size), func(b *testing.B) {
			grafo := crearGrafoCompleto(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				algorithms.Kruskal(grafo)
			}
		})
	}
}

// BenchmarkValidacionUtils prueba el rendimiento de las utilidades de validación
func BenchmarkValidacionUtils(b *testing.B) {
	b.Run("ValidarID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utils.ValidarID("C123")
		}
	})

	b.Run("ValidarNombre", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utils.ValidarNombre("Cueva de Prueba")
		}
	})

	b.Run("ValidarCoordenadas", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utils.ValidarCoordenadas(10.5, 20.3)
		}
	})
}

// BenchmarkMathUtils prueba el rendimiento de las utilidades matemáticas
func BenchmarkMathUtils(b *testing.B) {
	b.Run("CalcularDistanciaEuclidiana", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utils.CalcularDistanciaEuclidiana(0, 0, 3, 4)
		}
	})

	b.Run("CalcularDistanciaManhattan", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utils.CalcularDistanciaManhattan(0, 0, 3, 4)
		}
	})

	b.Run("Redondear", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utils.Redondear(3.14159, 2)
		}
	})
}

// BenchmarkGraphUtils prueba el rendimiento de las utilidades de grafos
func BenchmarkGraphUtils(b *testing.B) {
	sizes := []int{10, 50, 100}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("EsConexo_%d_nodos", size), func(b *testing.B) {
			grafo := crearGrafoCompleto(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				utils.EsConexo(grafo)
			}
		})

		b.Run(fmt.Sprintf("TieneCiclos_%d_nodos", size), func(b *testing.B) {
			grafo := crearGrafoCompleto(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				utils.TieneCiclos(grafo)
			}
		})

		b.Run(fmt.Sprintf("ObtenerComponentesConexas_%d_nodos", size), func(b *testing.B) {
			grafo := crearGrafoCompleto(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				utils.ObtenerComponentesConexas(grafo)
			}
		})
	}
}

// BenchmarkOperacionesGrafo prueba operaciones básicas del grafo
func BenchmarkOperacionesGrafo(b *testing.B) {
	b.Run("AgregarCueva", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			grafo := domain.NuevoGrafo(false)
			cueva := domain.NuevaCueva(fmt.Sprintf("C%d", i), fmt.Sprintf("Cueva %d", i))
			grafo.AgregarCueva(cueva)
		}
	})

	b.Run("AgregarConexion", func(b *testing.B) {
		grafo := domain.NuevoGrafo(false)
		// Agregar algunas cuevas
		for i := 0; i < 100; i++ {
			cueva := domain.NuevaCueva(fmt.Sprintf("C%d", i), fmt.Sprintf("Cueva %d", i))
			grafo.AgregarCueva(cueva)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			from := fmt.Sprintf("C%d", i%100)
			to := fmt.Sprintf("C%d", (i+1)%100)
			grafo.AgregarConexion(from, to, 10.0)
		}
	})

	b.Run("ObtenerCueva", func(b *testing.B) {
		grafo := crearGrafoCompleto(100)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			id := fmt.Sprintf("C%d", i%100)
			grafo.ObtenerCueva(id)
		}
	})
}

// Función auxiliar para crear un grafo completo con n nodos
func crearGrafoCompleto(n int) *domain.Grafo {
	grafo := domain.NuevoGrafo(false)

	// Agregar nodos
	for i := 0; i < n; i++ {
		cueva := domain.NuevaCueva(fmt.Sprintf("C%d", i), fmt.Sprintf("Cueva %d", i))
		cueva.X = float64(i % 10)
		cueva.Y = float64(i / 10)
		grafo.AgregarCueva(cueva)
	}

	// Agregar conexiones (cada nodo conectado a los siguientes 3 nodos)
	for i := 0; i < n; i++ {
		for j := 1; j <= 3 && i+j < n; j++ {
			from := fmt.Sprintf("C%d", i)
			to := fmt.Sprintf("C%d", i+j)
			distancia := utils.CalcularDistanciaEuclidiana(
				float64(i%10), float64(i/10),
				float64((i+j)%10), float64((i+j)/10),
			)
			grafo.AgregarConexion(from, to, distancia)
		}
	}

	return grafo
}

// Función auxiliar para crear un grafo lineal (como una cadena)
func crearGrafoLineal(n int) *domain.Grafo {
	grafo := domain.NuevoGrafo(false)

	// Agregar nodos
	for i := 0; i < n; i++ {
		cueva := domain.NuevaCueva(fmt.Sprintf("C%d", i), fmt.Sprintf("Cueva %d", i))
		cueva.X = float64(i)
		cueva.Y = 0.0
		grafo.AgregarCueva(cueva)
	}

	// Agregar conexiones lineales
	for i := 0; i < n-1; i++ {
		from := fmt.Sprintf("C%d", i)
		to := fmt.Sprintf("C%d", i+1)
		grafo.AgregarConexion(from, to, 1.0)
	}

	return grafo
}

// Función auxiliar para crear un grafo con componentes desconectadas
func crearGrafoDesconectado(n int) *domain.Grafo {
	grafo := domain.NuevoGrafo(false)

	// Agregar nodos
	for i := 0; i < n; i++ {
		cueva := domain.NuevaCueva(fmt.Sprintf("C%d", i), fmt.Sprintf("Cueva %d", i))
		cueva.X = float64(i % 10)
		cueva.Y = float64(i / 10)
		grafo.AgregarCueva(cueva)
	}

	// Crear componentes pequeñas desconectadas
	for i := 0; i < n-1; i += 3 {
		if i+1 < n {
			grafo.AgregarConexion(fmt.Sprintf("C%d", i), fmt.Sprintf("C%d", i+1), 1.0)
		}
		if i+2 < n {
			grafo.AgregarConexion(fmt.Sprintf("C%d", i+1), fmt.Sprintf("C%d", i+2), 1.0)
		}
	}

	return grafo
}
