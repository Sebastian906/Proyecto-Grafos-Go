package handler

import (
	"encoding/json"
	"fmt"
	"proyecto-grafos-go/internal/service"
)

// Controlador para operaciones de conexiones
type ControladorConexion struct {
	servicioConexion *service.ServicioConexion
}

// Nuevo controlador de conexiones
func NuevoControladorConexion(servicioConexion *service.ServicioConexion) *ControladorConexion {
	return &ControladorConexion{
		servicioConexion: servicioConexion,
	}
}

// Manejar cambio de tipo de grafo (dirigido/no dirigido)
func (cc *ControladorConexion) ManejarCambiarTipoGrafo(datos []byte) (string, error) {
	var solicitud service.CambiarTipoGrafo
	if err := json.Unmarshal(datos, &solicitud); err != nil {
		return "", fmt.Errorf("error al parsear datos: %v", err)
	}

	if err := cc.servicioConexion.CambiarTipoGrafo(&solicitud); err != nil {
		return "", err
	}

	tipoGrafo := "no dirigido"
	if solicitud.EsDirigido {
		tipoGrafo = "dirigido"
	}

	return fmt.Sprintf("Tipo de grafo cambiado exitosamente a: %s", tipoGrafo), nil
}

// Manejar obstrucción de conexiones
func (cc *ControladorConexion) ManejarObstruirConexion(datos []byte) (string, error) {
	var solicitud service.ObstruirConexion
	if err := json.Unmarshal(datos, &solicitud); err != nil {
		return "", fmt.Errorf("error al parsear datos: %v", err)
	}

	if err := cc.servicioConexion.ObstruirConexion(&solicitud); err != nil {
		return "", err
	}

	estado := "obstruida"
	if !solicitud.EsObstruido {
		estado = "desobstruida"
	}

	return fmt.Sprintf("Conexión desde %s hasta %s %s exitosamente",
		solicitud.DesdeCuevaID, solicitud.HastaCuevaID, estado), nil
}

// Manejar cambio de dirección de conexiones
func (cc *ControladorConexion) ManejarCambiarDireccionConexion(datos []byte) (string, error) {
	var solicitud service.CambiarDireccion
	if err := json.Unmarshal(datos, &solicitud); err != nil {
		return "", fmt.Errorf("error al parsear datos: %v", err)
	}

	if err := cc.servicioConexion.CambiarDireccionConexion(&solicitud); err != nil {
		return "", err
	}

	tipo := "no dirigida"
	if solicitud.NuevaDireccion {
		tipo = "dirigida"
	}

	return fmt.Sprintf("Conexión desde %s hasta %s cambiada a %s exitosamente",
		solicitud.DesdeCuevaID, solicitud.HastaCuevaID, tipo), nil
}

// Manejar cambio de sentido de ruta
func (cc *ControladorConexion) ManejarCambiarSentidoRuta(datos []byte) (string, error) {
	var solicitud service.CambiarSentidoRuta
	if err := json.Unmarshal(datos, &solicitud); err != nil {
		return "", fmt.Errorf("error al parsear datos: %v", err)
	}

	if err := cc.servicioConexion.CambiarSentidoRuta(&solicitud); err != nil {
		return "", err
	}

	return fmt.Sprintf("Sentido de ruta cambiado exitosamente: ahora va desde %s hasta %s",
		solicitud.HastaCuevaID, solicitud.DesdeCuevaID), nil
}

// Manejar cambio de sentido de múltiples rutas
func (cc *ControladorConexion) ManejarCambiarSentidoMultiplesRutas(datos []byte) (string, error) {
	var solicitudes []*service.CambiarSentidoRuta
	if err := json.Unmarshal(datos, &solicitudes); err != nil {
		return "", fmt.Errorf("error al parsear datos: %v", err)
	}

	errores := cc.servicioConexion.CambiarSentidoMultiplesRutas(solicitudes)

	if len(errores) > 0 {
		mensaje := fmt.Sprintf("Se procesaron %d rutas con %d errores:\n", len(solicitudes), len(errores))
		for _, err := range errores {
			mensaje += fmt.Sprintf("- %v\n", err)
		}
		return mensaje, nil
	}

	return fmt.Sprintf("Sentido de %d rutas cambiado exitosamente", len(solicitudes)), nil
}

// Manejar inversión de rutas desde una cueva
func (cc *ControladorConexion) ManejarInvertirRutasDesdeCueva(cuevaID string) (string, error) {
	if err := cc.servicioConexion.InvertirRutasDesdeCueva(cuevaID); err != nil {
		return "", err
	}

	return fmt.Sprintf("Todas las rutas dirigidas salientes de la cueva %s han sido invertidas", cuevaID), nil
}

// Manejar inversión de rutas hacia una cueva
func (cc *ControladorConexion) ManejarInvertirRutasHaciaCueva(cuevaID string) (string, error) {
	if err := cc.servicioConexion.InvertirRutasHaciaCueva(cuevaID); err != nil {
		return "", err
	}

	return fmt.Sprintf("Todas las rutas dirigidas entrantes a la cueva %s han sido invertidas", cuevaID), nil
}

// Listar todas las conexiones
func (cc *ControladorConexion) ListarConexiones() ([]byte, error) {
	conexiones := cc.servicioConexion.ListarConexiones()
	return json.MarshalIndent(conexiones, "", "  ")
}

// Listar conexiones activas (no obstruidas)
func (cc *ControladorConexion) ListarConexionesActivas() ([]byte, error) {
	conexiones := cc.servicioConexion.ObtenerConexionesActivas()
	return json.MarshalIndent(conexiones, "", "  ")
}

// Listar conexiones obstruidas
func (cc *ControladorConexion) ListarConexionesObstruidas() ([]byte, error) {
	conexiones := cc.servicioConexion.ListarConexionesObstruidas()
	return json.MarshalIndent(conexiones, "", "  ")
}

// Obtener estadísticas de conexiones
func (cc *ControladorConexion) ObtenerEstadisticasConexiones() ([]byte, error) {
	estadisticas := cc.servicioConexion.EstadisticasConexiones()
	return json.MarshalIndent(estadisticas, "", "  ")
}

// Eliminar conexión específica
func (cc *ControladorConexion) ManejarEliminarConexion(desdeCuevaID, hastaCuevaID string) (string, error) {
	if err := cc.servicioConexion.EliminarConexion(desdeCuevaID, hastaCuevaID); err != nil {
		return "", err
	}

	return fmt.Sprintf("Conexión desde %s hasta %s eliminada exitosamente", desdeCuevaID, hastaCuevaID), nil
}

// Desobstruir todas las conexiones
func (cc *ControladorConexion) ManejarDesobstruirTodasConexiones() (string, error) {
	conexionesDesobstruidas := cc.servicioConexion.DesobstruirTodasConexiones()
	return fmt.Sprintf("%d conexiones desobstruidas exitosamente", conexionesDesobstruidas), nil
}
