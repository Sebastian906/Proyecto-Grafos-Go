package utils

import (
	"regexp"
	"strings"
)

// ValidarID valida que un ID sea válido (no vacío, formato correcto)
func ValidarID(id string) bool {
	if strings.TrimSpace(id) == "" {
		return false
	}

	if len(id) > 50 {
		return false
	}

	// Verificar que solo contenga caracteres alfanuméricos y guiones
	match, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", id)
	if !match {
		return false
	}

	// Debe contener al menos una letra (no solo números)
	hasLetter, _ := regexp.MatchString("[a-zA-Z]", id)
	return hasLetter
}

// ValidarNombre valida que un nombre sea válido
func ValidarNombre(nombre string) bool {
	if strings.TrimSpace(nombre) == "" {
		return false
	}

	if len(nombre) < 3 {
		return false
	}

	if len(nombre) > 100 {
		return false
	}

	return true
}

// ValidarCoordenadas valida que las coordenadas sean válidas
func ValidarCoordenadas(x, y float64) bool {
	return x >= -1000 && x <= 1000 && y >= -1000 && y <= 1000
}

// ValidarDistancia valida que una distancia sea válida
func ValidarDistancia(distancia float64) bool {
	return distancia > 0 && distancia <= 10000
}

// ValidarArchivoExtension valida que un archivo tenga una extensión válida
func ValidarArchivoExtension(archivo string) bool {
	archivo = strings.ToLower(archivo)
	return strings.HasSuffix(archivo, ".json") ||
		strings.HasSuffix(archivo, ".xml") ||
		strings.HasSuffix(archivo, ".txt")
}

// ValidarTipoCamion valida que el tipo de camión sea válido (A, B, C)
func ValidarTipoCamion(tipo string) bool {
	tipo = strings.ToUpper(strings.TrimSpace(tipo))
	return tipo == "A" || tipo == "B" || tipo == "C"
}
