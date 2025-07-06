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
