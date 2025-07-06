package handler

import (
	"encoding/json"
	"fmt"
	"proyecto-grafos-go/internal/service"
	"strings"
)

// Controlador para operaciones de conexiones
type ControladorConexion struct {
	servicioConexion   *service.ServicioConexion
	servicioValidacion *service.ServicioValidacion
}

// Nuevo controlador de conexiones
func NuevoControladorConexion(servicioConexion *service.ServicioConexion, servicioValidacion *service.ServicioValidacion) *ControladorConexion {
	return &ControladorConexion{
		servicioConexion:   servicioConexion,
		servicioValidacion: servicioValidacion,
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

// ManejarDetectarCuevasInaccesibles detecta y reporta cuevas inaccesibles tras cambios
func (cc *ControladorConexion) ManejarDetectarCuevasInaccesibles() (string, error) {
	resultado := cc.servicioValidacion.DetectarCuevasInaccesiblesTrasChanged()

	var reporte strings.Builder
	reporte.WriteString("=== ANALISIS DE ACCESIBILIDAD ===\n")
	reporte.WriteString(fmt.Sprintf("Total de cuevas: %d\n", resultado.TotalCuevas))
	reporte.WriteString(fmt.Sprintf("Cuevas accesibles: %d\n", resultado.CuevasAccesibles))
	reporte.WriteString(fmt.Sprintf("Cuevas inaccesibles: %d\n", len(resultado.CuevasInaccesibles)))

	if len(resultado.CuevasInaccesibles) > 0 {
		reporte.WriteString("\nCUEVAS INACCESIBLES:\n")
		for i, cueva := range resultado.CuevasInaccesibles {
			reporte.WriteString(fmt.Sprintf("%d. %s\n", i+1, cueva))
		}
	}

	reporte.WriteString("\nSOLUCIONES PROPUESTAS:\n")
	for _, solucion := range resultado.Soluciones {
		reporte.WriteString(solucion + "\n")
	}

	return reporte.String(), nil
}

// ManejarAnalizarAccesibilidadDesde analiza accesibilidad desde una cueva específica
func (cc *ControladorConexion) ManejarAnalizarAccesibilidadDesde(cuevaInicio string) (string, error) {
	resultado := cc.servicioValidacion.AnalizarAccesibilidad(cuevaInicio)

	var reporte strings.Builder
	reporte.WriteString(fmt.Sprintf("=== ANALISIS DE ACCESIBILIDAD DESDE '%s' ===\n", cuevaInicio))
	reporte.WriteString(fmt.Sprintf("Total de cuevas: %d\n", resultado.TotalCuevas))
	reporte.WriteString(fmt.Sprintf("Cuevas accesibles: %d\n", resultado.CuevasAccesibles))
	reporte.WriteString(fmt.Sprintf("Cuevas inaccesibles: %d\n", len(resultado.CuevasInaccesibles)))

	if len(resultado.CuevasInaccesibles) > 0 {
		reporte.WriteString("\nCUEVAS INACCESIBLES:\n")
		for i, cueva := range resultado.CuevasInaccesibles {
			reporte.WriteString(fmt.Sprintf("%d. %s\n", i+1, cueva))
		}
	}

	reporte.WriteString("\nSOLUCIONES PROPUESTAS:\n")
	for _, solucion := range resultado.Soluciones {
		reporte.WriteString(solucion + "\n")
	}

	return reporte.String(), nil
}

// ManejarCambiarDireccionConConAnalisis cambia dirección y analiza impacto en accesibilidad
func (cc *ControladorConexion) ManejarCambiarDireccionConConAnalisis(datos []byte) (string, error) {
	var solicitud service.CambiarDireccion
	if err := json.Unmarshal(datos, &solicitud); err != nil {
		return "", fmt.Errorf("error al parsear datos: %v", err)
	}

	// Realizar el cambio
	if err := cc.servicioConexion.CambiarDireccionConexion(&solicitud); err != nil {
		return "", err
	}

	tipo := "no dirigida"
	if solicitud.NuevaDireccion {
		tipo = "dirigida"
	}

	// Analizar impacto en accesibilidad
	resultado := cc.servicioValidacion.DetectarCuevasInaccesiblesTrasChanged()

	var reporte strings.Builder
	reporte.WriteString(fmt.Sprintf("Conexion desde %s hasta %s cambiada a %s exitosamente\n\n",
		solicitud.DesdeCuevaID, solicitud.HastaCuevaID, tipo))

	if len(resultado.CuevasInaccesibles) > 0 {
		reporte.WriteString("ATENCION: El cambio ha resultado en cuevas inaccesibles:\n")
		for _, cueva := range resultado.CuevasInaccesibles {
			reporte.WriteString(fmt.Sprintf("- %s\n", cueva))
		}
		reporte.WriteString("\nSOLUCIONES RECOMENDADAS:\n")
		for _, solucion := range resultado.Soluciones {
			reporte.WriteString(solucion + "\n")
		}
	} else {
		reporte.WriteString("Excelente: Todas las cuevas siguen siendo accesibles tras el cambio.")
	}

	return reporte.String(), nil
}
