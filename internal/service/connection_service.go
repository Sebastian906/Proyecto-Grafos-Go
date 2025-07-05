package service

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
)

// Operaciones de conexión con grafos
type ServicioConexion struct {
	grafo *domain.Grafo
}

// Nuevo servicio de conexiones
func NuevoServicioConexion(grafo *domain.Grafo) *ServicioConexion {
	return &ServicioConexion{
		grafo: grafo,
	}
}

// Solicitud para cambiar el tipo de grafo
type CambiarTipoGrafo struct {
	EsDirigido bool `json:"es_dirigido"`
}

// Solicitud para obstruir conexión
type ObstruirConexion struct {
	DesdeCuevaID string `json:"desde_cueva_id"`
	HastaCuevaID string `json:"hasta_cueva_id"`
	EsObstruido  bool   `json:"es_obstruido"`
}

// Solicitud para cambiarle la dirección a una conexión
type CambiarDireccion struct {
	DesdeCuevaID   string `json:"desde_cueva_id"`
	HastaCuevaID   string `json:"hasta_cueva_id"`
	NuevaDireccion bool   `json:"nueva_direccion"` // true = dirigido, false = no dirigido
}

// Solicitud para cambiar el sentido de una ruta (invertir dirección)
type CambiarSentidoRuta struct {
	DesdeCuevaID string `json:"desde_cueva_id"`
	HastaCuevaID string `json:"hasta_cueva_id"`
}

// Cambiar el grafo si es dirigido o no dirigido
func (sc *ServicioConexion) CambiarTipoGrafo(solicitud *CambiarTipoGrafo) error {
	// Si el grafo ya es del tipo solicitado, no hacer nada
	if sc.grafo.EsDirigido == solicitud.EsDirigido {
		return nil
	}

	// Cambiar el tipo de grafo
	sc.grafo.EsDirigido = solicitud.EsDirigido

	// Si se cambia a no dirigido, agregar aristas inversas faltantes
	if !solicitud.EsDirigido {
		return sc.convertirANoDirigido()
	}

	// Si se cambia a dirigido, eliminar aristas duplicadas
	return sc.convertirADirigido()
}

// Convertir el grafo a no dirigido
func (sc *ServicioConexion) convertirANoDirigido() error {
	aristasExistentes := make(map[string]bool)
	var nuevasAristas []*domain.Arista

	// Marcar todas las aristas existentes
	for _, arista := range sc.grafo.Aristas {
		clave := fmt.Sprintf("%s-%s", arista.Desde, arista.Hasta)
		aristasExistentes[clave] = true
		nuevasAristas = append(nuevasAristas, arista)
	}

	// Agregar aristas inversas faltantes
	for _, arista := range sc.grafo.Aristas {
		claveInversa := fmt.Sprintf("%s-%s", arista.Hasta, arista.Desde)
		if !aristasExistentes[claveInversa] {
			aristaInversa := domain.NuevaArista(arista.Hasta, arista.Desde, arista.Distancia, false)
			aristaInversa.EsObstruido = arista.EsObstruido
			nuevasAristas = append(nuevasAristas, aristaInversa)
		}
	}

	// Actualizar todas las aristas como no dirigidas
	for _, arista := range nuevasAristas {
		arista.EsDirigido = false
	}

	sc.grafo.Aristas = nuevasAristas
	return nil
}

// Convertir el grafo a dirigido
func (sc *ServicioConexion) convertirADirigido() error {
	// Actualizar todas las aristas como dirigidas
	for _, arista := range sc.grafo.Aristas {
		arista.EsDirigido = true
	}
	return nil
}

// Obstruir o desobstruir una conexión específica
func (sc *ServicioConexion) ObstruirConexion(solicitud *ObstruirConexion) error {
	// Verificar que las cuevas existan
	if _, existe := sc.grafo.ObtenerCueva(solicitud.DesdeCuevaID); !existe {
		return fmt.Errorf("cueva %s no existe", solicitud.DesdeCuevaID)
	}
	if _, existe := sc.grafo.ObtenerCueva(solicitud.HastaCuevaID); !existe {
		return fmt.Errorf("cueva %s no existe", solicitud.HastaCuevaID)
	}

	// Encontrar y modificar la arista
	aristasModificadas := 0
	for _, arista := range sc.grafo.Aristas {
		if arista.Desde == solicitud.DesdeCuevaID && arista.Hasta == solicitud.HastaCuevaID {
			arista.EsObstruido = solicitud.EsObstruido
			aristasModificadas++
		}
		// Si es un grafo no dirigido, también obstruir la arista inversa
		if !sc.grafo.EsDirigido && arista.Desde == solicitud.HastaCuevaID && arista.Hasta == solicitud.DesdeCuevaID {
			arista.EsObstruido = solicitud.EsObstruido
			aristasModificadas++
		}
	}

	if aristasModificadas == 0 {
		return fmt.Errorf("conexión desde %s hasta %s no existe", solicitud.DesdeCuevaID, solicitud.HastaCuevaID)
	}

	return nil
}

// Obstruir múltiples conexiones en una sola operación
func (sc *ServicioConexion) ObstruirMultiplesConexiones(solicitudes []*ObstruirConexion) []error {
	var errores []error

	for i, solicitud := range solicitudes {
		if err := sc.ObstruirConexion(solicitud); err != nil {
			errores = append(errores, fmt.Errorf("error en conexión %d: %v", i+1, err))
		}
	}

	return errores
}

// Obstruir todas las conexiones de una cueva específica
func (sc *ServicioConexion) ObstruirTodasConexionesCueva(cuevaID string, esObstruido bool) error {
	// Verificar que la cueva exista
	if _, existe := sc.grafo.ObtenerCueva(cuevaID); !existe {
		return fmt.Errorf("cueva %s no existe", cuevaID)
	}

	conexionesModificadas := 0
	for _, arista := range sc.grafo.Aristas {
		if arista.Desde == cuevaID || arista.Hasta == cuevaID {
			arista.EsObstruido = esObstruido
			conexionesModificadas++
		}
	}

	if conexionesModificadas == 0 {
		return fmt.Errorf("la cueva %s no tiene conexiones", cuevaID)
	}

	return nil
}

// Listar todas las conexiones obstruidas
func (sc *ServicioConexion) ListarConexionesObstruidas() []map[string]interface{} {
	var conexionesObstruidas []map[string]interface{}

	for _, arista := range sc.grafo.Aristas {
		if arista.EsObstruido {
			conexion := map[string]interface{}{
				"desde":       arista.Desde,
				"hasta":       arista.Hasta,
				"distancia":   arista.Distancia,
				"es_dirigido": arista.EsDirigido,
			}
			conexionesObstruidas = append(conexionesObstruidas, conexion)
		}
	}

	return conexionesObstruidas
}

// Desobstruir todas las conexiones del grafo
func (sc *ServicioConexion) DesobstruirTodasConexiones() int {
	conexionesDesobstruidas := 0

	for _, arista := range sc.grafo.Aristas {
		if arista.EsObstruido {
			arista.EsObstruido = false
			conexionesDesobstruidas++
		}
	}

	return conexionesDesobstruidas
}

// Cambiar dirección de una conexión específica
func (sc *ServicioConexion) CambiarDireccionConexion(solicitud *CambiarDireccion) error {
	// Verificar que las cuevas existan
	if _, existe := sc.grafo.ObtenerCueva(solicitud.DesdeCuevaID); !existe {
		return fmt.Errorf("cueva %s no existe", solicitud.DesdeCuevaID)
	}
	if _, existe := sc.grafo.ObtenerCueva(solicitud.HastaCuevaID); !existe {
		return fmt.Errorf("cueva %s no existe", solicitud.HastaCuevaID)
	}

	// Encontrar la arista
	var aristaEncontrada *domain.Arista
	for _, arista := range sc.grafo.Aristas {
		if arista.Desde == solicitud.DesdeCuevaID && arista.Hasta == solicitud.HastaCuevaID {
			aristaEncontrada = arista
			break
		}
	}

	if aristaEncontrada == nil {
		return fmt.Errorf("conexión desde %s hasta %s no existe", solicitud.DesdeCuevaID, solicitud.HastaCuevaID)
	}

	// Cambiar la dirección
	aristaEncontrada.EsDirigido = solicitud.NuevaDireccion

	return nil
}

// Invertir dirección de una conexión (cambiar sentido de la ruta)
func (sc *ServicioConexion) InvertirDireccionConexion(solicitud *CambiarSentidoRuta) error {
	// Verificar que las cuevas existan
	if _, existe := sc.grafo.ObtenerCueva(solicitud.DesdeCuevaID); !existe {
		return fmt.Errorf("cueva %s no existe", solicitud.DesdeCuevaID)
	}
	if _, existe := sc.grafo.ObtenerCueva(solicitud.HastaCuevaID); !existe {
		return fmt.Errorf("cueva %s no existe", solicitud.HastaCuevaID)
	}

	// Encontrar la arista
	var aristaEncontrada *domain.Arista
	for _, arista := range sc.grafo.Aristas {
		if arista.Desde == solicitud.DesdeCuevaID && arista.Hasta == solicitud.HastaCuevaID {
			aristaEncontrada = arista
			break
		}
	}

	if aristaEncontrada == nil {
		return fmt.Errorf("conexión desde %s hasta %s no existe", solicitud.DesdeCuevaID, solicitud.HastaCuevaID)
	}

	// Cambiar el sentido de la arista
	aristaEncontrada.Desde, aristaEncontrada.Hasta = aristaEncontrada.Hasta, aristaEncontrada.Desde

	return nil
}

// Listar todas las conexiones en el grafo
func (sc *ServicioConexion) ListarConexiones() []map[string]interface{} {
	var conexiones []map[string]interface{}

	for _, arista := range sc.grafo.Aristas {
		conexion := map[string]interface{}{
			"desde":        arista.Desde,
			"hasta":        arista.Hasta,
			"distancia":    arista.Distancia,
			"es_dirigido":  arista.EsDirigido,
			"es_obstruido": arista.EsObstruido,
		}
		conexiones = append(conexiones, conexion)
	}

	return conexiones
}

// Obtener conexiones que no estén obstruidas
func (sc *ServicioConexion) ObtenerConexionesActivas() []map[string]interface{} {
	var conexiones []map[string]interface{}

	for _, arista := range sc.grafo.Aristas {
		if !arista.EsObstruido {
			conexion := map[string]interface{}{
				"desde":       arista.Desde,
				"hasta":       arista.Hasta,
				"distancia":   arista.Distancia,
				"es_dirigido": arista.EsDirigido,
			}
			conexiones = append(conexiones, conexion)
		}
	}

	return conexiones
}

// Eliminar conexiones específicas
func (sc *ServicioConexion) EliminarConexion(desdeCuevaID, hastaCuevaID string) error {
	// Verificar que las cuevas existan
	if _, existe := sc.grafo.ObtenerCueva(desdeCuevaID); !existe {
		return fmt.Errorf("cueva %s no existe", desdeCuevaID)
	}
	if _, existe := sc.grafo.ObtenerCueva(hastaCuevaID); !existe {
		return fmt.Errorf("cueva %s no existe", hastaCuevaID)
	}

	// Filtrar la arista a eliminar
	var aristasActualizadas []*domain.Arista
	aristaEliminada := false

	for _, arista := range sc.grafo.Aristas {
		if arista.Desde == desdeCuevaID && arista.Hasta == hastaCuevaID {
			aristaEliminada = true
			continue
		}
		aristasActualizadas = append(aristasActualizadas, arista)
	}

	if !aristaEliminada {
		return fmt.Errorf("conexión desde %s hasta %s no existe", desdeCuevaID, hastaCuevaID)
	}

	sc.grafo.Aristas = aristasActualizadas
	return nil
}

// Estadísticas de conexiones
func (sc *ServicioConexion) EstadisticasConexiones() map[string]interface{} {
	totalConexiones := len(sc.grafo.Aristas)
	conexionesActivas := 0
	conexionesObstruidas := 0
	conexionesDirigidas := 0
	conexionesNoDirigidas := 0

	for _, arista := range sc.grafo.Aristas {
		if arista.EsObstruido {
			conexionesObstruidas++
		} else {
			conexionesActivas++
		}

		if arista.EsDirigido {
			conexionesDirigidas++
		} else {
			conexionesNoDirigidas++
		}
	}

	return map[string]interface{}{
		"total_conexiones":        totalConexiones,
		"conexiones_activas":      conexionesActivas,
		"conexiones_obstruidas":   conexionesObstruidas,
		"conexiones_dirigidas":    conexionesDirigidas,
		"conexiones_no_dirigidas": conexionesNoDirigidas,
		"tipo_grafo":              sc.grafo.EsDirigido,
	}
}

// Cambiar el sentido de una ruta específica (invertir dirección)
func (sc *ServicioConexion) CambiarSentidoRuta(solicitud *CambiarSentidoRuta) error {
	// Verificar que las cuevas existan
	if _, existe := sc.grafo.ObtenerCueva(solicitud.DesdeCuevaID); !existe {
		return fmt.Errorf("cueva %s no existe", solicitud.DesdeCuevaID)
	}
	if _, existe := sc.grafo.ObtenerCueva(solicitud.HastaCuevaID); !existe {
		return fmt.Errorf("cueva %s no existe", solicitud.HastaCuevaID)
	}

	// Buscar la arista original
	var aristaEncontrada *domain.Arista
	var indiceArista int = -1

	for i, arista := range sc.grafo.Aristas {
		if arista.Desde == solicitud.DesdeCuevaID && arista.Hasta == solicitud.HastaCuevaID {
			aristaEncontrada = arista
			indiceArista = i
			break
		}
	}

	if aristaEncontrada == nil {
		return fmt.Errorf("conexión desde %s hasta %s no existe", solicitud.DesdeCuevaID, solicitud.HastaCuevaID)
	}

	// Verificar si la conexión es dirigida
	if !aristaEncontrada.EsDirigido {
		return fmt.Errorf("no se puede cambiar el sentido de una conexión no dirigida")
	}

	// Crear la nueva arista con dirección invertida
	nuevaArista := domain.NuevaArista(
		aristaEncontrada.Hasta,
		aristaEncontrada.Desde,
		aristaEncontrada.Distancia,
		aristaEncontrada.EsDirigido,
	)
	nuevaArista.EsObstruido = aristaEncontrada.EsObstruido

	// Reemplazar la arista original con la nueva arista invertida
	sc.grafo.Aristas[indiceArista] = nuevaArista

	return nil
}

// Cambiar el sentido de múltiples rutas en una sola operación
func (sc *ServicioConexion) CambiarSentidoMultiplesRutas(solicitudes []*CambiarSentidoRuta) []error {
	var errores []error

	for i, solicitud := range solicitudes {
		if err := sc.CambiarSentidoRuta(solicitud); err != nil {
			errores = append(errores, fmt.Errorf("error en ruta %d: %v", i+1, err))
		}
	}

	return errores
}

// Invertir todas las rutas dirigidas que salen de una cueva específica
func (sc *ServicioConexion) InvertirRutasDesdeCueva(cuevaID string) error {
	// Verificar que la cueva exista
	if _, existe := sc.grafo.ObtenerCueva(cuevaID); !existe {
		return fmt.Errorf("cueva %s no existe", cuevaID)
	}

	rutasInvertidas := 0

	// Buscar todas las aristas que salen de la cueva especificada
	for i, arista := range sc.grafo.Aristas {
		if arista.Desde == cuevaID && arista.EsDirigido {
			// Crear nueva arista invertida
			nuevaArista := domain.NuevaArista(
				arista.Hasta,
				arista.Desde,
				arista.Distancia,
				arista.EsDirigido,
			)
			nuevaArista.EsObstruido = arista.EsObstruido

			// Reemplazar la arista original
			sc.grafo.Aristas[i] = nuevaArista
			rutasInvertidas++
		}
	}

	if rutasInvertidas == 0 {
		return fmt.Errorf("la cueva %s no tiene rutas dirigidas salientes para invertir", cuevaID)
	}

	return nil
}

// Invertir todas las rutas dirigidas que llegan a una cueva específica
func (sc *ServicioConexion) InvertirRutasHaciaCueva(cuevaID string) error {
	// Verificar que la cueva exista
	if _, existe := sc.grafo.ObtenerCueva(cuevaID); !existe {
		return fmt.Errorf("cueva %s no existe", cuevaID)
	}

	rutasInvertidas := 0

	// Buscar todas las aristas que llegan a la cueva especificada
	for i, arista := range sc.grafo.Aristas {
		if arista.Hasta == cuevaID && arista.EsDirigido {
			// Crear nueva arista invertida
			nuevaArista := domain.NuevaArista(
				arista.Hasta,
				arista.Desde,
				arista.Distancia,
				arista.EsDirigido,
			)
			nuevaArista.EsObstruido = arista.EsObstruido

			// Reemplazar la arista original
			sc.grafo.Aristas[i] = nuevaArista
			rutasInvertidas++
		}
	}

	if rutasInvertidas == 0 {
		return fmt.Errorf("la cueva %s no tiene rutas dirigidas entrantes para invertir", cuevaID)
	}

	return nil
}
