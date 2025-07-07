package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Obtener un número entero del usuario
func ObtenerInputInt(prompt string) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		num, err := strconv.Atoi(input)
		if err == nil {
			return num
		}
		fmt.Println("Error: ingrese un número válido")
	}
}

func ObtenerInputString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func ObtenerInputFloat(prompt string) float64 {
	for {
		input := ObtenerInputString(prompt)
		num, err := strconv.ParseFloat(input, 64)
		if err == nil {
			return num
		}
		fmt.Println("Error: ingrese un número válido")
	}
}

func ObtenerInputBool(prompt string) bool {
	input := ObtenerInputString(prompt + " (s/n): ")
	return strings.ToLower(input) == "s"
}

// LeerEntrada es un alias para ObtenerInputString para mantener compatibilidad
func LeerEntrada(prompt string) string {
	return ObtenerInputString(prompt)
}

// Solicitar confirmación del usuario (S/N)
func SolicitarConfirmacion(prompt string) bool {
	for {
		respuesta := ObtenerInputString(prompt + " (S/N): ")
		respuesta = strings.ToLower(strings.TrimSpace(respuesta))

		if respuesta == "s" || respuesta == "si" || respuesta == "sí" || respuesta == "y" || respuesta == "yes" {
			return true
		} else if respuesta == "n" || respuesta == "no" {
			return false
		}

		fmt.Println("❌ Por favor, responda con 'S' para Sí o 'N' para No")
	}
}
