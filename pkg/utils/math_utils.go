package utils

import (
	"math"
)

// Min retorna el menor de dos números flotantes
func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Max retorna el mayor de dos números flotantes
func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Abs retorna el valor absoluto de un número flotante
func Abs(n float64) float64 {
	if n < 0 {
		return -n
	}
	return n
}

// Redondear redondea un número a un número específico de decimales
func Redondear(value float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	scaled := value * multiplier

	// Manejo especial para casos como -2.5 -> -2.0
	if scaled < 0 {
		if scaled == math.Trunc(scaled)-0.5 {
			// Es exactamente x.5 negativo, redondear hacia cero
			return math.Trunc(scaled) / multiplier
		}
		return math.Round(scaled) / multiplier
	}
	return math.Round(scaled) / multiplier
}

// EsIgual compara dos números flotantes con una tolerancia
func EsIgual(a, b, tolerancia float64) bool {
	return Abs(a-b) <= tolerancia
}

// CalcularDistanciaEuclidiana calcula la distancia euclidiana entre dos puntos
func CalcularDistanciaEuclidiana(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

// CalcularDistanciaManhattan calcula la distancia Manhattan entre dos puntos
func CalcularDistanciaManhattan(x1, y1, x2, y2 float64) float64 {
	return Abs(x2-x1) + Abs(y2-y1)
}

// Promedio calcula el promedio de un slice de números
func Promedio(numeros []float64) float64 {
	if len(numeros) == 0 {
		panic("no se puede calcular el promedio de un slice vacío")
	}

	suma := 0.0
	for _, num := range numeros {
		suma += num
	}
	return suma / float64(len(numeros))
}

// MinInt retorna el menor de dos números enteros
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MaxInt retorna el mayor de dos números enteros
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// AbsInt retorna el valor absoluto de un número entero
func AbsInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Round redondea un número flotante al entero más cercano
func Round(f float64) int {
	return int(math.Floor(f + 0.5))
}

// Clamp limita un valor entre un mínimo y máximo
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampFloat limita un valor flotante entre un mínimo y máximo
func ClampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
