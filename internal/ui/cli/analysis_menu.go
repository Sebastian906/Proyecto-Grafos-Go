package cli

import (
	"fmt"
	"proyecto-grafos-go/internal/service"
)

type MenuAnalisis struct {
    validacionSvc *service.ServicioValidacion
    grafoSvc      *service.ServicioGrafo
}

func NuevoMenuAnalisis(
    validacionSvc *service.ServicioValidacion,
    grafoSvc *service.ServicioGrafo,
) *MenuAnalisis {
    return &MenuAnalisis{
        validacionSvc: validacionSvc,
        grafoSvc:      grafoSvc,
    }
}

func (m *MenuAnalisis) Mostrar() {
    for {
        fmt.Println("\n=== Menú de Análisis ===")
        fmt.Println("1. Verificar conectividad fuerte")
        fmt.Println("2. Detectar pozos")
        fmt.Println("3. Mostrar grados de vértices")
        fmt.Println("4. Salir")

        opcion := ObtenerInputInt("Seleccione una opción: ")

        switch opcion {
        case 1:
            m.mostrarConectividad()
        case 2:
            m.mostrarPozos()
        case 3:
            m.mostrarGrados()
        case 4:
            return
        default:
            fmt.Println("Opción inválida")
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