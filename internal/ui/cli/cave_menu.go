package cli

import (
	"fmt"
	"proyecto-grafos-go/internal/service"
)

type MenuCueva struct {
	cuevaSvc    *service.ServicioCueva
	conexionSvc *service.ServicioConexion
}

func NuevoMenuCueva(cuevaSvc *service.ServicioCueva, conexionSvc *service.ServicioConexion) *MenuCueva {
	return &MenuCueva{
		cuevaSvc:    cuevaSvc,
		conexionSvc: conexionSvc,
	}
}

func (m *MenuCueva) Mostrar() {
	for {
		fmt.Println("\n=== Menú de Gestión de Cuevas y Conexiones ===")
		fmt.Println("1. Crear nueva cueva")
		fmt.Println("2. Conectar cuevas")
		fmt.Println("3. Listar cuevas")
		fmt.Println("4. Obstruir una conexión específica")
		fmt.Println("5. Desobstruir una conexión específica")
		fmt.Println("6. Obstruir múltiples conexiones")
		fmt.Println("7. Obstruir todas las conexiones de una cueva")
		fmt.Println("8. Desobstruir todas las conexiones de una cueva")
		fmt.Println("9. Listar conexiones obstruidas")
		fmt.Println("10. Desobstruir todas las conexiones del grafo")
		fmt.Println("11. Cambiar sentido de una ruta específica")
		fmt.Println("12. Cambiar sentido de múltiples rutas")
		fmt.Println("13. Invertir todas las rutas salientes de una cueva")
		fmt.Println("14. Invertir todas las rutas entrantes a una cueva")
		fmt.Println("15. Mostrar estadísticas de conexiones")
		fmt.Println("16. Volver al menú principal")

		opcion := ObtenerInputInt("Seleccione una opción: ")

		switch opcion {
		case 1:
			m.crearCueva()
		case 2:
			m.conectarCuevas()
		case 3:
			m.listarCuevas()
		case 4:
			m.obstruirConexion(true)
		case 5:
			m.obstruirConexion(false)
		case 6:
			m.obstruirMultiplesConexiones()
		case 7:
			m.obstruirTodasConexionesCueva(true)
		case 8:
			m.obstruirTodasConexionesCueva(false)
		case 9:
			m.listarConexionesObstruidas()
		case 10:
			m.desobstruirTodasConexiones()
		case 11:
			m.cambiarSentidoRuta()
		case 12:
			m.cambiarSentidoMultiplesRutas()
		case 13:
			m.invertirRutasDesdeCueva()
		case 14:
			m.invertirRutasHaciaCueva()
		case 15:
			m.mostrarEstadisticasConexiones()
		case 16:
			return
		default:
			fmt.Println("Opción inválida")
		}
	}
}

// ===================== FUNCIONES DE CUEVAS =====================

func (m *MenuCueva) crearCueva() {
	id := ObtenerInputString("ID de la cueva: ")
	nombre := ObtenerInputString("Nombre: ")
	x := ObtenerInputFloat("Coordenada X: ")
	y := ObtenerInputFloat("Coordenada Y: ")

	err := m.cuevaSvc.CrearCueva(service.SolicitudCueva{
		ID:     id,
		Nombre: nombre,
		X:      x,
		Y:      y,
	})

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Cueva creada exitosamente")
	}
}

func (m *MenuCueva) conectarCuevas() {
	desde := ObtenerInputString("ID cueva origen: ")
	hasta := ObtenerInputString("ID cueva destino: ")
	distancia := ObtenerInputFloat("Distancia: ")
	dirigido := ObtenerInputBool("¿Es dirigido?")
	bidireccional := ObtenerInputBool("¿Bidireccional?")

	if err := m.cuevaSvc.Conectar(desde, hasta, distancia, dirigido, bidireccional); err != nil {
		fmt.Println("Error conectando cuevas:", err)
	} else {
		fmt.Println("Cuevas conectadas exitosamente")
	}
}

func (m *MenuCueva) listarCuevas() {
	cuevas := m.cuevaSvc.ListarCuevas()

	if len(cuevas) == 0 {
		fmt.Println("No hay cuevas en el grafo")
		return
	}

	fmt.Printf("\n=== Lista de Cuevas (%d) ===\n", len(cuevas))
	for i, id := range cuevas {
		cueva, existe := m.cuevaSvc.ObtenerCueva(id)
		if !existe {
			fmt.Printf("%d. %s (Error: cueva no encontrada)\n", i+1, id)
		} else {
			fmt.Printf("%d. %s - %s (%.2f, %.2f)\n", i+1, cueva.ID, cueva.Nombre, cueva.X, cueva.Y)
		}
	}
}

// ===================== FUNCIONES DE CONEXIONES =====================

func (m *MenuCueva) obstruirConexion(obstruir bool) {
	desde := ObtenerInputString("ID cueva origen: ")
	hasta := ObtenerInputString("ID cueva destino: ")

	solicitud := &service.ObstruirConexion{
		DesdeCuevaID: desde,
		HastaCuevaID: hasta,
		EsObstruido:  obstruir,
	}

	if err := m.conexionSvc.ObstruirConexion(solicitud); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		if obstruir {
			fmt.Printf("Conexión desde %s hasta %s obstruida exitosamente\n", desde, hasta)
		} else {
			fmt.Printf("Conexión desde %s hasta %s desobstruida exitosamente\n", desde, hasta)
		}
	}
}

func (m *MenuCueva) obstruirMultiplesConexiones() {
	fmt.Println("Ingrese las conexiones a obstruir (ingrese 'fin' para terminar):")

	var solicitudes []*service.ObstruirConexion

	for {
		fmt.Printf("\n--- Conexión %d ---\n", len(solicitudes)+1)

		desde := ObtenerInputString("ID cueva origen (o 'fin' para terminar): ")
		if desde == "fin" {
			break
		}

		hasta := ObtenerInputString("ID cueva destino: ")
		obstruir := ObtenerInputBool("¿Obstruir esta conexión?")

		solicitud := &service.ObstruirConexion{
			DesdeCuevaID: desde,
			HastaCuevaID: hasta,
			EsObstruido:  obstruir,
		}

		solicitudes = append(solicitudes, solicitud)
	}

	if len(solicitudes) == 0 {
		fmt.Println("No se ingresaron conexiones")
		return
	}

	errores := m.conexionSvc.ObstruirMultiplesConexiones(solicitudes)

	if len(errores) == 0 {
		fmt.Printf("Todas las %d conexiones fueron procesadas exitosamente\n", len(solicitudes))
	} else {
		fmt.Printf("Se procesaron %d conexiones con %d errores:\n", len(solicitudes), len(errores))
		for _, err := range errores {
			fmt.Printf("- %v\n", err)
		}
	}
}

func (m *MenuCueva) obstruirTodasConexionesCueva(obstruir bool) {
	cuevaID := ObtenerInputString("ID de la cueva: ")

	if err := m.conexionSvc.ObstruirTodasConexionesCueva(cuevaID, obstruir); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		if obstruir {
			fmt.Printf("Todas las conexiones de la cueva %s han sido obstruidas\n", cuevaID)
		} else {
			fmt.Printf("Todas las conexiones de la cueva %s han sido desobstruidas\n", cuevaID)
		}
	}
}

func (m *MenuCueva) listarConexionesObstruidas() {
	conexiones := m.conexionSvc.ListarConexionesObstruidas()

	if len(conexiones) == 0 {
		fmt.Println("No hay conexiones obstruidas en el grafo")
		return
	}

	fmt.Printf("\n=== Conexiones Obstruidas (%d) ===\n", len(conexiones))
	for i, conexion := range conexiones {
		direccion := "↔"
		if conexion["es_dirigido"].(bool) {
			direccion = "→"
		}

		fmt.Printf("%d. %s %s %s (Distancia: %.2f)\n",
			i+1,
			conexion["desde"],
			direccion,
			conexion["hasta"],
			conexion["distancia"],
		)
	}
}

func (m *MenuCueva) desobstruirTodasConexiones() {
	confirmar := ObtenerInputBool("¿Está seguro de que desea desobstruir TODAS las conexiones?")

	if !confirmar {
		fmt.Println("Operación cancelada")
		return
	}

	conexionesDesobstruidas := m.conexionSvc.DesobstruirTodasConexiones()

	if conexionesDesobstruidas == 0 {
		fmt.Println("No había conexiones obstruidas")
	} else {
		fmt.Printf("Se desobstruyeron %d conexiones exitosamente\n", conexionesDesobstruidas)
	}
}

func (m *MenuCueva) mostrarEstadisticasConexiones() {
	stats := m.conexionSvc.EstadisticasConexiones()

	fmt.Println("\n=== Estadísticas de Conexiones ===")
	fmt.Printf("Total de conexiones: %v\n", stats["total_conexiones"])
	fmt.Printf("Conexiones activas: %v\n", stats["conexiones_activas"])
	fmt.Printf("Conexiones obstruidas: %v\n", stats["conexiones_obstruidas"])
	fmt.Printf("Conexiones dirigidas: %v\n", stats["conexiones_dirigidas"])
	fmt.Printf("Conexiones no dirigidas: %v\n", stats["conexiones_no_dirigidas"])
	fmt.Printf("Tipo de grafo: ")
	if stats["tipo_grafo"].(bool) {
		fmt.Println("Dirigido")
	} else {
		fmt.Println("No dirigido")
	}
}

// ===================== FUNCIONES DE CAMBIO DE SENTIDO DE RUTAS =====================

func (m *MenuCueva) cambiarSentidoRuta() {
	fmt.Println("\n=== Cambiar Sentido de Ruta ===")

	desdeCuevaID := ObtenerInputString("ID de la cueva origen: ")
	hastaCuevaID := ObtenerInputString("ID de la cueva destino: ")

	solicitud := &service.CambiarSentidoRuta{
		DesdeCuevaID: desdeCuevaID,
		HastaCuevaID: hastaCuevaID,
	}

	err := m.conexionSvc.CambiarSentidoRuta(solicitud)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Sentido de ruta cambiado exitosamente: ahora va desde %s hasta %s\n",
			hastaCuevaID, desdeCuevaID)
	}
}

func (m *MenuCueva) cambiarSentidoMultiplesRutas() {
	fmt.Println("\n=== Cambiar Sentido de Múltiples Rutas ===")

	numRutas := ObtenerInputInt("¿Cuántas rutas desea cambiar?: ")
	if numRutas <= 0 {
		fmt.Println("Número inválido de rutas")
		return
	}

	var solicitudes []*service.CambiarSentidoRuta

	for i := 0; i < numRutas; i++ {
		fmt.Printf("\n--- Ruta %d ---\n", i+1)
		desdeCuevaID := ObtenerInputString("ID de la cueva origen: ")
		hastaCuevaID := ObtenerInputString("ID de la cueva destino: ")

		solicitud := &service.CambiarSentidoRuta{
			DesdeCuevaID: desdeCuevaID,
			HastaCuevaID: hastaCuevaID,
		}
		solicitudes = append(solicitudes, solicitud)
	}

	errores := m.conexionSvc.CambiarSentidoMultiplesRutas(solicitudes)

	if len(errores) > 0 {
		fmt.Printf("Se procesaron %d rutas con %d errores:\n", len(solicitudes), len(errores))
		for _, err := range errores {
			fmt.Printf("- %v\n", err)
		}
	} else {
		fmt.Printf("Sentido de %d rutas cambiado exitosamente\n", len(solicitudes))
	}
}

func (m *MenuCueva) invertirRutasDesdeCueva() {
	fmt.Println("\n=== Invertir Rutas Salientes de una Cueva ===")

	cuevaID := ObtenerInputString("ID de la cueva: ")

	err := m.conexionSvc.InvertirRutasDesdeCueva(cuevaID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Todas las rutas dirigidas salientes de la cueva %s han sido invertidas\n", cuevaID)
	}
}

func (m *MenuCueva) invertirRutasHaciaCueva() {
	fmt.Println("\n=== Invertir Rutas Entrantes a una Cueva ===")

	cuevaID := ObtenerInputString("ID de la cueva: ")

	err := m.conexionSvc.InvertirRutasHaciaCueva(cuevaID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Todas las rutas dirigidas entrantes a la cueva %s han sido invertidas\n", cuevaID)
	}
}
