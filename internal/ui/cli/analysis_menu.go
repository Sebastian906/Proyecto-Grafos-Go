package cli

import (
	"fmt"
	"proyecto-grafos-go/internal/service"
)

type MenuAnalisis struct {
	validacionSvc *service.ServicioValidacion
	grafoSvc      *service.ServicioGrafo
	conexionSvc   *service.ServicioConexion
}

func NuevoMenuAnalisis(
	validacionSvc *service.ServicioValidacion,
	grafoSvc *service.ServicioGrafo,
	conexionSvc *service.ServicioConexion,
) *MenuAnalisis {
	return &MenuAnalisis{
		validacionSvc: validacionSvc,
		grafoSvc:      grafoSvc,
		conexionSvc:   conexionSvc,
	}
}

func (m *MenuAnalisis) Mostrar() {
	for {
		fmt.Println("\n=== Menu de Analisis ===")
		fmt.Println("1. Verificar conectividad fuerte")
		fmt.Println("2. Detectar pozos")
		fmt.Println("3. Mostrar grados de vertices")
		fmt.Println("4. Detectar cuevas inaccesibles")
		fmt.Println("5. Analizar accesibilidad desde cueva especifica")
		fmt.Println("6. Salir")

		opcion := ObtenerInputInt("Seleccione una opcion: ")

		switch opcion {
		case 1:
			m.mostrarConectividad()
		case 2:
			m.mostrarPozos()
		case 3:
			m.mostrarGrados()
		case 4:
			m.detectarCuevasInaccesibles()
		case 5:
			m.analizarAccesibilidadEspecifica()
		case 6:
			return
		default:
			fmt.Println("Opcion invalida")
		}
	}
}

func (m *MenuAnalisis) mostrarConectividad() {
	if m.validacionSvc.EsFuertementeConectado() {
		fmt.Println("El grafo es fuertemente conectado")
	} else {
		fmt.Println("El grafo NO es fuertemente conectado")
	}
}

func (m *MenuAnalisis) mostrarPozos() {
	pozos := m.validacionSvc.DetectarPozos()
	if len(pozos) == 0 {
		fmt.Println("No hay pozos en el grafo")
	} else {
		fmt.Println("Pozos encontrados:")
		for _, p := range pozos {
			fmt.Println("-", p)
		}
	}
}

func (m *MenuAnalisis) mostrarGrados() {
	grados := m.grafoSvc.ObtenerGradosVertices()
	for id, g := range grados {
		fmt.Printf("%s: Entrantes=%d, Salientes=%d, Total=%d\n",
			id, g["entrante"], g["saliente"], g["total"])
	}
}

func (m *MenuAnalisis) detectarCuevasInaccesibles() {
	resultado := m.validacionSvc.DetectarCuevasInaccesiblesTrasChanged()

	fmt.Println("\n=== ANALISIS DE ACCESIBILIDAD ===")
	fmt.Printf("Total de cuevas: %d\n", resultado.TotalCuevas)
	fmt.Printf("Cuevas accesibles: %d\n", resultado.CuevasAccesibles)
	fmt.Printf("Cuevas inaccesibles: %d\n", len(resultado.CuevasInaccesibles))

	if len(resultado.CuevasInaccesibles) > 0 {
		fmt.Println("\nCUEVAS INACCESIBLES:")
		for i, cueva := range resultado.CuevasInaccesibles {
			fmt.Printf("%d. %s\n", i+1, cueva)
		}
	}

	fmt.Println("\nSOLUCIONES PROPUESTAS:")
	for _, solucion := range resultado.Soluciones {
		fmt.Println(solucion)
	}

	fmt.Println("\nPresione Enter para continuar...")
	ObtenerInputString("")
}

func (m *MenuAnalisis) analizarAccesibilidadEspecifica() {
	cuevaInicio := ObtenerInputString("Ingrese el ID de la cueva de inicio: ")
	if cuevaInicio == "" {
		fmt.Println("ID de cueva no puede estar vacio")
		return
	}

	resultado := m.validacionSvc.AnalizarAccesibilidad(cuevaInicio)

	fmt.Printf("\n=== ANALISIS DE ACCESIBILIDAD DESDE '%s' ===\n", cuevaInicio)
	fmt.Printf("Total de cuevas: %d\n", resultado.TotalCuevas)
	fmt.Printf("Cuevas accesibles: %d\n", resultado.CuevasAccesibles)
	fmt.Printf("Cuevas inaccesibles: %d\n", len(resultado.CuevasInaccesibles))

	if len(resultado.CuevasInaccesibles) > 0 {
		fmt.Println("\nCUEVAS INACCESIBLES:")
		for i, cueva := range resultado.CuevasInaccesibles {
			fmt.Printf("%d. %s\n", i+1, cueva)
		}
	}

	fmt.Println("\nSOLUCIONES PROPUESTAS:")
	for _, solucion := range resultado.Soluciones {
		fmt.Println(solucion)
	}

	fmt.Println("\nPresione Enter para continuar...")
	ObtenerInputString("")
}
