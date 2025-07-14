package unit

import (
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/pkg/utils"
	"testing"
)

func TestGraphUtils_EsConexo(t *testing.T) {
	// Test grafo conexo
	grafoConexo := domain.NuevoGrafo(false)
	c1 := domain.NuevaCueva("C1", "Cueva 1")
	c2 := domain.NuevaCueva("C2", "Cueva 2")
	c3 := domain.NuevaCueva("C3", "Cueva 3")

	grafoConexo.AgregarCueva(c1)
	grafoConexo.AgregarCueva(c2)
	grafoConexo.AgregarCueva(c3)
	grafoConexo.AgregarConexion("C1", "C2", 10.0)
	grafoConexo.AgregarConexion("C2", "C3", 15.0)

	if !utils.EsConexo(grafoConexo) {
		t.Error("Grafo conexo no reconocido como conexo")
	}

	// Test grafo no conexo
	grafoNoConexo := domain.NuevoGrafo(false)
	c4 := domain.NuevaCueva("C4", "Cueva 4")

	grafoNoConexo.AgregarCueva(c1)
	grafoNoConexo.AgregarCueva(c2)
	grafoNoConexo.AgregarCueva(c3)
	grafoNoConexo.AgregarCueva(c4)
	grafoNoConexo.AgregarConexion("C1", "C2", 10.0)
	// C3 y C4 quedan aislados

	if utils.EsConexo(grafoNoConexo) {
		t.Error("Grafo no conexo reconocido como conexo")
	}
}

func TestGraphUtils_TieneCiclos(t *testing.T) {
	// Test grafo sin ciclos (árbol)
	grafoSinCiclos := domain.NuevoGrafo(false)
	c1 := domain.NuevaCueva("C1", "Cueva 1")
	c2 := domain.NuevaCueva("C2", "Cueva 2")
	c3 := domain.NuevaCueva("C3", "Cueva 3")

	grafoSinCiclos.AgregarCueva(c1)
	grafoSinCiclos.AgregarCueva(c2)
	grafoSinCiclos.AgregarCueva(c3)
	grafoSinCiclos.AgregarConexion("C1", "C2", 10.0)
	grafoSinCiclos.AgregarConexion("C2", "C3", 15.0)

	if utils.TieneCiclos(grafoSinCiclos) {
		t.Error("Grafo sin ciclos reconocido como con ciclos")
	}

	// Test grafo con ciclos
	grafoConCiclos := domain.NuevoGrafo(false)
	grafoConCiclos.AgregarCueva(c1)
	grafoConCiclos.AgregarCueva(c2)
	grafoConCiclos.AgregarCueva(c3)
	grafoConCiclos.AgregarConexion("C1", "C2", 10.0)
	grafoConCiclos.AgregarConexion("C2", "C3", 15.0)
	grafoConCiclos.AgregarConexion("C3", "C1", 20.0) // Forma un ciclo

	if !utils.TieneCiclos(grafoConCiclos) {
		t.Error("Grafo con ciclos no reconocido como con ciclos")
	}
}

func TestGraphUtils_ObtenerGrado(t *testing.T) {
	// Setup
	grafo := domain.NuevoGrafo(false)
	c1 := domain.NuevaCueva("C1", "Cueva 1")
	c2 := domain.NuevaCueva("C2", "Cueva 2")
	c3 := domain.NuevaCueva("C3", "Cueva 3")

	grafo.AgregarCueva(c1)
	grafo.AgregarCueva(c2)
	grafo.AgregarCueva(c3)
	grafo.AgregarConexion("C1", "C2", 10.0)
	grafo.AgregarConexion("C1", "C3", 15.0)

	// Test grado de C1 (conectado a C2 y C3)
	grado := utils.ObtenerGrado(grafo, "C1")
	if grado != 2 {
		t.Errorf("Grado de C1 debería ser 2, obtuvo %d", grado)
	}

	// Test grado de C2 (conectado solo a C1)
	grado = utils.ObtenerGrado(grafo, "C2")
	if grado != 1 {
		t.Errorf("Grado de C2 debería ser 1, obtuvo %d", grado)
	}

	// Test grado de cueva inexistente
	grado = utils.ObtenerGrado(grafo, "C999")
	if grado != 0 {
		t.Errorf("Grado de cueva inexistente debería ser 0, obtuvo %d", grado)
	}
}

func TestGraphUtils_ObtenerVecinos(t *testing.T) {
	// Setup
	grafo := domain.NuevoGrafo(false)
	c1 := domain.NuevaCueva("C1", "Cueva 1")
	c2 := domain.NuevaCueva("C2", "Cueva 2")
	c3 := domain.NuevaCueva("C3", "Cueva 3")
	c4 := domain.NuevaCueva("C4", "Cueva 4")

	grafo.AgregarCueva(c1)
	grafo.AgregarCueva(c2)
	grafo.AgregarCueva(c3)
	grafo.AgregarCueva(c4)
	grafo.AgregarConexion("C1", "C2", 10.0)
	grafo.AgregarConexion("C1", "C3", 15.0)

	// Test vecinos de C1
	vecinos := grafo.ObtenerVecinos("C1")
	if len(vecinos) != 2 {
		t.Errorf("C1 debería tener 2 vecinos, obtuvo %d", len(vecinos))
	}

	// Verificar que C2 y C3 están en la lista de vecinos
	vecinosMap := make(map[string]bool)
	for _, vecino := range vecinos {
		vecinosMap[vecino] = true
	}
	if !vecinosMap["C2"] || !vecinosMap["C3"] {
		t.Error("Vecinos de C1 deberían incluir C2 y C3")
	}

	// Test vecinos de C4 (sin conexiones)
	vecinos = grafo.ObtenerVecinos("C4")
	if len(vecinos) != 0 {
		t.Errorf("C4 debería tener 0 vecinos, obtuvo %d", len(vecinos))
	}
}

func TestGraphUtils_EsAciclico(t *testing.T) {
	// Test grafo acíclico (árbol)
	grafoAciclico := domain.NuevoGrafo(false)
	c1 := domain.NuevaCueva("C1", "Cueva 1")
	c2 := domain.NuevaCueva("C2", "Cueva 2")
	c3 := domain.NuevaCueva("C3", "Cueva 3")

	grafoAciclico.AgregarCueva(c1)
	grafoAciclico.AgregarCueva(c2)
	grafoAciclico.AgregarCueva(c3)
	grafoAciclico.AgregarConexion("C1", "C2", 10.0)
	grafoAciclico.AgregarConexion("C2", "C3", 15.0)

	if !utils.EsAciclico(grafoAciclico) {
		t.Error("Grafo acíclico no reconocido como acíclico")
	}

	// Test grafo con ciclos
	grafoConCiclos := domain.NuevoGrafo(false)
	grafoConCiclos.AgregarCueva(c1)
	grafoConCiclos.AgregarCueva(c2)
	grafoConCiclos.AgregarCueva(c3)
	grafoConCiclos.AgregarConexion("C1", "C2", 10.0)
	grafoConCiclos.AgregarConexion("C2", "C3", 15.0)
	grafoConCiclos.AgregarConexion("C3", "C1", 20.0) // Forma un ciclo

	if utils.EsAciclico(grafoConCiclos) {
		t.Error("Grafo con ciclos reconocido como acíclico")
	}
}

func TestGraphUtils_ObtenerComponentesConexas(t *testing.T) {
	// Test con múltiples componentes conexas
	grafo := domain.NuevoGrafo(false)

	// Componente 1: C1-C2
	c1 := domain.NuevaCueva("C1", "Cueva 1")
	c2 := domain.NuevaCueva("C2", "Cueva 2")
	grafo.AgregarCueva(c1)
	grafo.AgregarCueva(c2)
	grafo.AgregarConexion("C1", "C2", 10.0)

	// Componente 2: C3-C4
	c3 := domain.NuevaCueva("C3", "Cueva 3")
	c4 := domain.NuevaCueva("C4", "Cueva 4")
	grafo.AgregarCueva(c3)
	grafo.AgregarCueva(c4)
	grafo.AgregarConexion("C3", "C4", 15.0)

	// Componente 3: C5 (aislado)
	c5 := domain.NuevaCueva("C5", "Cueva 5")
	grafo.AgregarCueva(c5)

	componentes := utils.ObtenerComponentesConexas(grafo)
	if len(componentes) != 3 {
		t.Errorf("Debería haber 3 componentes conexas, obtuvo %d", len(componentes))
	}

	// Verificar que cada componente tiene el tamaño correcto
	tamanos := make([]int, len(componentes))
	for i, comp := range componentes {
		tamanos[i] = len(comp)
	}

	// Ordenar tamaños para comparación
	for i := 0; i < len(tamanos)-1; i++ {
		for j := i + 1; j < len(tamanos); j++ {
			if tamanos[i] > tamanos[j] {
				tamanos[i], tamanos[j] = tamanos[j], tamanos[i]
			}
		}
	}

	expectedTamanos := []int{1, 2, 2} // Un componente de tamaño 1 y dos de tamaño 2
	for i, expected := range expectedTamanos {
		if tamanos[i] != expected {
			t.Errorf("Tamaño del componente %d debería ser %d, obtuvo %d", i, expected, tamanos[i])
		}
	}
}

func TestGraphUtils_EsCompleto(t *testing.T) {
	// Test grafo completo con 3 nodos
	grafoCompleto := domain.NuevoGrafo(false)
	c1 := domain.NuevaCueva("C1", "Cueva 1")
	c2 := domain.NuevaCueva("C2", "Cueva 2")
	c3 := domain.NuevaCueva("C3", "Cueva 3")

	grafoCompleto.AgregarCueva(c1)
	grafoCompleto.AgregarCueva(c2)
	grafoCompleto.AgregarCueva(c3)
	grafoCompleto.AgregarConexion("C1", "C2", 10.0)
	grafoCompleto.AgregarConexion("C1", "C3", 15.0)
	grafoCompleto.AgregarConexion("C2", "C3", 20.0)

	if !utils.EsCompleto(grafoCompleto) {
		t.Error("Grafo completo no reconocido como completo")
	}

	// Test grafo incompleto
	grafoIncompleto := domain.NuevoGrafo(false)
	grafoIncompleto.AgregarCueva(c1)
	grafoIncompleto.AgregarCueva(c2)
	grafoIncompleto.AgregarCueva(c3)
	grafoIncompleto.AgregarConexion("C1", "C2", 10.0)
	// Falta conexión C1-C3 y C2-C3

	if utils.EsCompleto(grafoIncompleto) {
		t.Error("Grafo incompleto reconocido como completo")
	}
}
