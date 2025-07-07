package cli

import (
	"fmt"
	"proyecto-grafos-go/internal/handler"
	"proyecto-grafos-go/internal/service"
)

type MainMenu struct {
	grafoSvc        *service.ServicioGrafo
	cuevaSvc        *service.ServicioCueva
	validacionSvc   *service.ServicioValidacion
	conexionSvc     *service.ServicioConexion
	analysisHandler *handler.AnalysisHandler
}

func NuevoMainMenu(
	grafoSvc *service.ServicioGrafo,
	cuevaSvc *service.ServicioCueva,
	validacionSvc *service.ServicioValidacion,
	conexionSvc *service.ServicioConexion,
	analysisHandler *handler.AnalysisHandler,
) *MainMenu {
	return &MainMenu{
		grafoSvc:        grafoSvc,
		cuevaSvc:        cuevaSvc,
		validacionSvc:   validacionSvc,
		conexionSvc:     conexionSvc,
		analysisHandler: analysisHandler,
	}
}

func (m *MainMenu) Mostrar() {
	for {
		fmt.Println("\n=== Menú Principal ===")
		fmt.Println("1. Cargar grafo desde archivo (1a)")
		fmt.Println("2. Gestión de cuevas y conexiones (1b, 2a)")
		fmt.Println("3. Cambiar tipo de grafo (dirigido/no) (1c)")
		fmt.Println("4. Análisis del grafo (1d-1f)")
		fmt.Println("5. Salir")

		opcion := ObtenerInputInt("Seleccione una opción: ")

		switch opcion {
		case 1:
			m.cargarGrafo()
		case 2:
			m.mostrarMenuCuevas()
		case 3:
			m.cambiarTipoGrafo()
		case 4:
			m.mostrarMenuAnalisis()
		case 5:
			return
		default:
			fmt.Println("Opción inválida")
		}
	}
}

func (m *MainMenu) cargarGrafo() {
	archivo := ObtenerInputString("Nombre del archivo (ej: caves.json): ")
	if err := m.grafoSvc.CargarGrafo(archivo); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Grafo cargado correctamente")
	}
}

func (m *MainMenu) cambiarTipoGrafo() {
	dirigido := ObtenerInputBool("¿Grafo dirigido? (s/n): ")
	m.grafoSvc.CambiarTipoGrafo(dirigido)
	fmt.Println("Tipo de grafo actualizado")
}

func (m *MainMenu) mostrarMenuCuevas() {
	menuCuevas := NuevoMenuCueva(m.cuevaSvc, m.conexionSvc)
	menuCuevas.Mostrar()
}

func (m *MainMenu) mostrarMenuAnalisis() {
	menuAnalisis := NuevoMenuAnalisis(m.validacionSvc, m.grafoSvc, m.conexionSvc, m.analysisHandler)
	menuAnalisis.Mostrar()
}

// MostrarMenuAnalisis expone el menú de análisis públicamente
func (m *MainMenu) MostrarMenuAnalisis() {
	m.mostrarMenuAnalisis()
}
