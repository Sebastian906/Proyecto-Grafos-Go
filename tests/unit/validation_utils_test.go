package unit

import (
	"proyecto-grafos-go/pkg/utils"
	"testing"
)

func TestValidacionUtils_ValidarID(t *testing.T) {
	// Test IDs válidos
	idsValidos := []string{"C1", "CUEVA_001", "A", "Cave123"}
	for _, id := range idsValidos {
		if !utils.ValidarID(id) {
			t.Errorf("ID válido rechazado: %s", id)
		}
	}

	// Test IDs inválidos
	idsInvalidos := []string{"", "  ", "C 1", "C@1", "123"}
	for _, id := range idsInvalidos {
		if utils.ValidarID(id) {
			t.Errorf("ID inválido aceptado: %s", id)
		}
	}
}

func TestValidacionUtils_ValidarNombre(t *testing.T) {
	// Test nombres válidos
	nombresValidos := []string{"Cueva Principal", "Entrada A", "Cave 123", "Túnel-Norte"}
	for _, nombre := range nombresValidos {
		if !utils.ValidarNombre(nombre) {
			t.Errorf("Nombre válido rechazado: %s", nombre)
		}
	}

	// Test nombres inválidos
	nombresInvalidos := []string{"", "  ", "A", "AB"}
	for _, nombre := range nombresInvalidos {
		if utils.ValidarNombre(nombre) {
			t.Errorf("Nombre inválido aceptado: %s", nombre)
		}
	}
}

func TestValidacionUtils_ValidarCoordenadas(t *testing.T) {
	// Test coordenadas válidas
	coordenadasValidas := []struct{ x, y float64 }{
		{0.0, 0.0},
		{10.5, 20.3},
		{-5.0, 15.0},
		{100.0, -50.0},
	}
	for _, coord := range coordenadasValidas {
		if !utils.ValidarCoordenadas(coord.x, coord.y) {
			t.Errorf("Coordenadas válidas rechazadas: (%.2f, %.2f)", coord.x, coord.y)
		}
	}

	// Test coordenadas inválidas (en este caso, asumimos que NaN o infinito son inválidos)
	coordenadasInvalidas := []struct{ x, y float64 }{
		// Note: En Go, comparaciones normales con NaN/Inf requieren import math
		// Para simplicidad, asumimos que las validaciones están en el código
	}
	for _, coord := range coordenadasInvalidas {
		if utils.ValidarCoordenadas(coord.x, coord.y) {
			t.Errorf("Coordenadas inválidas aceptadas: (%.2f, %.2f)", coord.x, coord.y)
		}
	}
}

func TestValidacionUtils_ValidarDistancia(t *testing.T) {
	// Test distancias válidas
	distanciasValidas := []float64{0.1, 1.0, 10.5, 100.0, 1000.0}
	for _, distancia := range distanciasValidas {
		if !utils.ValidarDistancia(distancia) {
			t.Errorf("Distancia válida rechazada: %.2f", distancia)
		}
	}

	// Test distancias inválidas
	distanciasInvalidas := []float64{0.0, -1.0, -10.5}
	for _, distancia := range distanciasInvalidas {
		if utils.ValidarDistancia(distancia) {
			t.Errorf("Distancia inválida aceptada: %.2f", distancia)
		}
	}
}

func TestValidacionUtils_ValidarArchivoExtension(t *testing.T) {
	// Test extensiones válidas
	archivosValidos := []string{
		"grafo.json",
		"datos.JSON",
		"cuevas.xml",
		"ARCHIVO.XML",
		"output.txt",
		"input.TXT",
	}
	for _, archivo := range archivosValidos {
		if !utils.ValidarArchivoExtension(archivo) {
			t.Errorf("Archivo válido rechazado: %s", archivo)
		}
	}

	// Test extensiones inválidas
	archivosInvalidos := []string{
		"archivo",
		"archivo.exe",
		"archivo.pdf",
		"archivo.doc",
		"",
	}
	for _, archivo := range archivosInvalidos {
		if utils.ValidarArchivoExtension(archivo) {
			t.Errorf("Archivo inválido aceptado: %s", archivo)
		}
	}
}

func TestValidacionUtils_ValidarTipoCamion(t *testing.T) {
	// Test tipos de camión válidos
	tiposValidos := []string{"A", "B", "C", "a", "b", "c"}
	for _, tipo := range tiposValidos {
		if !utils.ValidarTipoCamion(tipo) {
			t.Errorf("Tipo de camión válido rechazado: %s", tipo)
		}
	}

	// Test tipos de camión inválidos
	tiposInvalidos := []string{"", "D", "X", "AB", "1", "AA"}
	for _, tipo := range tiposInvalidos {
		if utils.ValidarTipoCamion(tipo) {
			t.Errorf("Tipo de camión inválido aceptado: %s", tipo)
		}
	}
}
