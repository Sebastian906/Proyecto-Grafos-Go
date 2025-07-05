package cli

import (
	"fmt"
	"proyecto-grafos-go/internal/service"
)

type MainMenu struct {
	grafoSvc      *service.ServicioGrafo
	cuevaSvc      *service.ServicioCueva
	validacionSvc *service.ServicioValidacion
}

func NuevoMainMenu(
	grafoSvc *service.ServicioGrafo,
	cuevaSvc *service.ServicioCueva,
	validacionSvc *service.ServicioValidacion,
) *MainMenu {
	return &MainMenu{
		grafoSvc:      grafoSvc,
		cuevaSvc:      cuevaSvc,
		validacionSvc: validacionSvc,
	}
}

func (m *MainMenu) Mostrar() {
	for {
		fmt.Println("\n=== Menú Principal ===")
		fmt.Println("1. Cargar grafo desde archivo (1a)")
		fmt.Println("2. Crear nueva cueva (1b)")
		fmt.Println("3. Conectar cuevas (1b/1c)")
		fmt.Println("4. Cambiar tipo de grafo (dirigido/no) (1c)")
		fmt.Println("5. Análisis del grafo (1d-1f)")
		fmt.Println("6. Salir")

		opcion := ObtenerInputInt("Seleccione una opción: ")

		switch opcion {
		case 1:
			m.cargarGrafo()
		case 2:
			m.crearCueva()
		case 3:
			m.conectarCuevas()
		case 4:
			m.cambiarTipoGrafo()
		case 5:
			m.mostrarMenuAnalisis()
		case 6:
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

func (m *MainMenu) crearCueva() {
	id := ObtenerInputString("ID de la cueva: ")
	nombre := ObtenerInputString("Nombre: ")
	err := m.cuevaSvc.CrearCueva(service.SolicitudCueva{
		ID:     id,
		Nombre: nombre,
		X:      0, // Valor por defecto
		Y:      0, // Valor por defecto
	})
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Cueva creada")
	}
}

func (m *MainMenu) conectarCuevas() {
	desde := ObtenerInputString("ID cueva origen: ")
	hasta := ObtenerInputString("ID cueva destino: ")
	distancia := ObtenerInputFloat("Distancia: ")
	dirigido := ObtenerInputBool("¿Es dirigido? (s/n): ")
	bidireccional := ObtenerInputBool("¿Bidireccional? (s/n): ")

	if err := m.cuevaSvc.Conectar(desde, hasta, distancia, dirigido, bidireccional); err != nil {
		fmt.Println("Error conectando cuevas:", err)
	} else {
		fmt.Println("Cuevas conectadas")
	}
}

func (m *MainMenu) cambiarTipoGrafo() {
	dirigido := ObtenerInputBool("¿Grafo dirigido? (s/n): ")
	m.grafoSvc.CambiarTipoGrafo(dirigido)
	fmt.Println("Tipo de grafo actualizado")
}

func (m *MainMenu) mostrarMenuAnalisis() {
	menuAnalisis := NuevoMenuAnalisis(m.validacionSvc, m.grafoSvc)
	menuAnalisis.Mostrar()
}
