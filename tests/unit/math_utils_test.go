package unit

import (
	"math"
	"proyecto-grafos-go/pkg/utils"
	"testing"
)

func TestMathUtils_CalcularDistanciaEuclidiana(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2 float64
		expected       float64
	}{
		{0, 0, 3, 4, 5.0},          // Triángulo 3-4-5
		{0, 0, 0, 0, 0.0},          // Mismo punto
		{1, 1, 4, 5, 5.0},          // Otro triángulo 3-4-5
		{-1, -1, 2, 3, 5.0},        // Con coordenadas negativas
		{0, 0, 1, 1, math.Sqrt(2)}, // Diagonal unitaria
	}

	for _, test := range tests {
		result := utils.CalcularDistanciaEuclidiana(test.x1, test.y1, test.x2, test.y2)
		if math.Abs(result-test.expected) > 1e-10 {
			t.Errorf("CalcularDistanciaEuclidiana(%.2f, %.2f, %.2f, %.2f) = %.6f, esperaba %.6f",
				test.x1, test.y1, test.x2, test.y2, result, test.expected)
		}
	}
}

func TestMathUtils_CalcularDistanciaManhattan(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2 float64
		expected       float64
	}{
		{0, 0, 3, 4, 7.0},   // |3-0| + |4-0| = 7
		{0, 0, 0, 0, 0.0},   // Mismo punto
		{1, 1, 4, 5, 7.0},   // |4-1| + |5-1| = 7
		{-1, -1, 2, 3, 7.0}, // |2-(-1)| + |3-(-1)| = 7
		{5, 3, 1, 1, 6.0},   // |1-5| + |1-3| = 6
	}

	for _, test := range tests {
		result := utils.CalcularDistanciaManhattan(test.x1, test.y1, test.x2, test.y2)
		if math.Abs(result-test.expected) > 1e-10 {
			t.Errorf("CalcularDistanciaManhattan(%.2f, %.2f, %.2f, %.2f) = %.6f, esperaba %.6f",
				test.x1, test.y1, test.x2, test.y2, result, test.expected)
		}
	}
}

func TestMathUtils_Redondear(t *testing.T) {
	tests := []struct {
		value    float64
		decimals int
		expected float64
	}{
		{3.14159, 2, 3.14},
		{3.14159, 4, 3.1416},
		{3.14159, 0, 3.0},
		{2.5, 0, 3.0},   // Redondeo hacia arriba
		{2.4, 0, 2.0},   // Redondeo hacia abajo
		{-2.5, 0, -2.0}, // Redondeo negativo
	}

	for _, test := range tests {
		result := utils.Redondear(test.value, test.decimals)
		if math.Abs(result-test.expected) > 1e-10 {
			t.Errorf("Redondear(%.6f, %d) = %.6f, esperaba %.6f",
				test.value, test.decimals, result, test.expected)
		}
	}
}

func TestMathUtils_Min(t *testing.T) {
	tests := []struct {
		a, b     float64
		expected float64
	}{
		{1.0, 2.0, 1.0},
		{5.0, 3.0, 3.0},
		{-1.0, -2.0, -2.0},
		{0.0, 0.0, 0.0},
		{1.5, 1.5, 1.5},
	}

	for _, test := range tests {
		result := utils.Min(test.a, test.b)
		if result != test.expected {
			t.Errorf("Min(%.2f, %.2f) = %.2f, esperaba %.2f",
				test.a, test.b, result, test.expected)
		}
	}
}

func TestMathUtils_Max(t *testing.T) {
	tests := []struct {
		a, b     float64
		expected float64
	}{
		{1.0, 2.0, 2.0},
		{5.0, 3.0, 5.0},
		{-1.0, -2.0, -1.0},
		{0.0, 0.0, 0.0},
		{1.5, 1.5, 1.5},
	}

	for _, test := range tests {
		result := utils.Max(test.a, test.b)
		if result != test.expected {
			t.Errorf("Max(%.2f, %.2f) = %.2f, esperaba %.2f",
				test.a, test.b, result, test.expected)
		}
	}
}

func TestMathUtils_Abs(t *testing.T) {
	tests := []struct {
		value    float64
		expected float64
	}{
		{5.0, 5.0},
		{-3.0, 3.0},
		{0.0, 0.0},
		{-0.0, 0.0},
		{1.5, 1.5},
		{-2.7, 2.7},
	}

	for _, test := range tests {
		result := utils.Abs(test.value)
		if result != test.expected {
			t.Errorf("Abs(%.2f) = %.2f, esperaba %.2f",
				test.value, result, test.expected)
		}
	}
}

func TestMathUtils_EsIgual(t *testing.T) {
	tests := []struct {
		a, b       float64
		tolerancia float64
		expected   bool
	}{
		{1.0, 1.0, 0.01, true},
		{1.0, 1.005, 0.01, true},
		{1.0, 1.02, 0.01, false},
		{0.0, 0.0, 0.001, true},
		{3.14159, 3.14160, 0.0001, true},
		{3.14159, 3.14200, 0.0001, false},
	}

	for _, test := range tests {
		result := utils.EsIgual(test.a, test.b, test.tolerancia)
		if result != test.expected {
			t.Errorf("EsIgual(%.6f, %.6f, %.6f) = %v, esperaba %v",
				test.a, test.b, test.tolerancia, result, test.expected)
		}
	}
}

func TestMathUtils_Promedio(t *testing.T) {
	tests := []struct {
		values   []float64
		expected float64
	}{
		{[]float64{1.0, 2.0, 3.0}, 2.0},
		{[]float64{5.0}, 5.0},
		{[]float64{-1.0, 1.0}, 0.0},
		{[]float64{2.5, 3.5, 4.0}, 3.333333333333333},
	}

	for _, test := range tests {
		result := utils.Promedio(test.values)
		if math.Abs(result-test.expected) > 1e-10 {
			t.Errorf("Promedio(%v) = %.6f, esperaba %.6f",
				test.values, result, test.expected)
		}
	}

	// Test con slice vacío
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Promedio([]) debería generar pánico")
		}
	}()
	utils.Promedio([]float64{})
}
