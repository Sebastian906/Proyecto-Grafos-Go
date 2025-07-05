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

// Obstruir o desobstruir una conexión especifica
func (sc *ServicioConexion) ObstruirConexion(solicitud *ObstruirConexion) error {
	// Verificar que las cuevas existan
	if _, existe := sc.grafo.ObtenerCueva(solicitud.DesdeCuevaID); !existe {
		return fmt.Errorf("cueva %s no existe", solicitud.DesdeCuevaID)
	}
	if _, existe := sc.grafo.ObtenerCueva(solicitud.HastaCuevaID); !existe {
		return fmt.Errorf("cueva %s no existe", solicitud.HastaCuevaID)
	}

	// Encontrar y modificar la arista
	aristaEncontrada := false
	for _, arista := range sc.grafo.Aristas {
		if arista.Desde == solicitud.DesdeCuevaID && arista.Hasta == solicitud.HastaCuevaID {
			arista.EsObstruido = solicitud.EsObstruido
			aristaEncontrada = true
		}
	}

	if !aristaEncontrada {
		return fmt.Errorf("conexión desde %s hasta %s no existe", solicitud.DesdeCuevaID, solicitud.HastaCuevaID)
	}

	return nil
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
