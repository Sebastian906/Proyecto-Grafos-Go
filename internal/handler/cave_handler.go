package handler

import (
	"fmt"
	"proyecto-grafos-go/internal/service"
)

// CaveHandler maneja operaciones relacionadas con las cuevas
type CaveHandler struct {
	cuevaService *service.ServicioCueva
}

// NuevoCaveHandler crea una nueva instancia del handler de cuevas
func NuevoCaveHandler(cuevaService *service.ServicioCueva) *CaveHandler {
	return &CaveHandler{
		cuevaService: cuevaService,
	}
}

// CrearCueva maneja la creación de una nueva cueva
func (ch *CaveHandler) CrearCueva(solicitud service.SolicitudCueva) error {
	if ch.cuevaService == nil {
		return fmt.Errorf("servicio de cueva no inicializado")
	}

	// Validar datos de entrada
	if err := ch.validarSolicitudCueva(solicitud); err != nil {
		return fmt.Errorf("datos de cueva inválidos: %w", err)
	}

	return ch.cuevaService.CrearCueva(solicitud)
}

// ConectarCuevas maneja la conexión entre dos cuevas
func (ch *CaveHandler) ConectarCuevas(desde, hasta string, distancia float64, esDirigido bool) error {
	if ch.cuevaService == nil {
		return fmt.Errorf("servicio de cueva no inicializado")
	}

	// Validar parámetros
	if desde == "" || hasta == "" {
		return fmt.Errorf("IDs de cuevas no pueden estar vacíos")
	}

	if desde == hasta {
		return fmt.Errorf("no se puede conectar una cueva consigo misma")
	}

	if distancia < 0 {
		return fmt.Errorf("distancia no puede ser negativa")
	}

	return ch.cuevaService.ConectarCuevas(desde, hasta, distancia, esDirigido)
}

// EliminarCueva maneja la eliminación de una cueva
func (ch *CaveHandler) EliminarCueva(id string) error {
	if ch.cuevaService == nil {
		return fmt.Errorf("servicio de cueva no inicializado")
	}

	if id == "" {
		return fmt.Errorf("ID de cueva no puede estar vacío")
	}

	return ch.cuevaService.EliminarCueva(id)
}

// ModificarCueva maneja la modificación de una cueva existente
func (ch *CaveHandler) ModificarCueva(id string, solicitud service.SolicitudCueva) error {
	if ch.cuevaService == nil {
		return fmt.Errorf("servicio de cueva no inicializado")
	}

	if id == "" {
		return fmt.Errorf("ID de cueva no puede estar vacío")
	}

	// Validar datos de entrada
	if err := ch.validarSolicitudCueva(solicitud); err != nil {
		return fmt.Errorf("datos de cueva inválidos: %w", err)
	}

	return ch.cuevaService.ModificarCueva(id, solicitud)
}

// ActualizarCueva actualiza una cueva existente
func (ch *CaveHandler) ActualizarCueva(id string, solicitud service.SolicitudCueva) error {
	if ch.cuevaService == nil {
		return fmt.Errorf("servicio de cueva no inicializado")
	}

	if id == "" {
		return fmt.Errorf("ID de cueva no puede estar vacío")
	}

	// Validar datos de entrada
	if err := ch.validarSolicitudCueva(solicitud); err != nil {
		return fmt.Errorf("datos de cueva inválidos: %w", err)
	}

	// Verificar que la cueva existe
	_, existe := ch.cuevaService.ObtenerCueva(id)
	if !existe {
		return fmt.Errorf("cueva con ID %s no encontrada", id)
	}

	// Actualizar la cueva
	return ch.cuevaService.ActualizarCueva(id, solicitud)
}

// ObtenerCueva maneja la obtención de información de una cueva
func (ch *CaveHandler) ObtenerCueva(id string) (*service.DetalleCueva, error) {
	if ch.cuevaService == nil {
		return nil, fmt.Errorf("servicio de cueva no inicializado")
	}

	if id == "" {
		return nil, fmt.Errorf("ID de cueva no puede estar vacío")
	}

	return ch.cuevaService.ObtenerDetalleCueva(id)
}

// ListarCuevas maneja la obtención de todas las cuevas
func (ch *CaveHandler) ListarCuevas() ([]string, error) {
	if ch.cuevaService == nil {
		return nil, fmt.Errorf("servicio de cueva no inicializado")
	}

	cuevas := ch.cuevaService.ListarCuevas()
	return cuevas, nil
}

// AgregarRecurso maneja la adición de recursos a una cueva
func (ch *CaveHandler) AgregarRecurso(idCueva, recurso string, cantidad int) error {
	if ch.cuevaService == nil {
		return fmt.Errorf("servicio de cueva no inicializado")
	}

	if idCueva == "" {
		return fmt.Errorf("ID de cueva no puede estar vacío")
	}

	if recurso == "" {
		return fmt.Errorf("nombre de recurso no puede estar vacío")
	}

	if cantidad <= 0 {
		return fmt.Errorf("cantidad debe ser mayor a cero")
	}

	return ch.cuevaService.AgregarRecurso(idCueva, recurso, cantidad)
}

// RemoverRecurso maneja la remoción de recursos de una cueva
func (ch *CaveHandler) RemoverRecurso(idCueva, recurso string, cantidad int) error {
	if ch.cuevaService == nil {
		return fmt.Errorf("servicio de cueva no inicializado")
	}

	if idCueva == "" {
		return fmt.Errorf("ID de cueva no puede estar vacío")
	}

	if recurso == "" {
		return fmt.Errorf("nombre de recurso no puede estar vacío")
	}

	if cantidad <= 0 {
		return fmt.Errorf("cantidad debe ser mayor a cero")
	}

	return ch.cuevaService.RemoverRecurso(idCueva, recurso, cantidad)
}

// VerificarConexion verifica si dos cuevas están conectadas
func (ch *CaveHandler) VerificarConexion(desde, hasta string) (bool, error) {
	if ch.cuevaService == nil {
		return false, fmt.Errorf("servicio de cueva no inicializado")
	}

	if desde == "" || hasta == "" {
		return false, fmt.Errorf("IDs de cuevas no pueden estar vacíos")
	}

	return ch.cuevaService.ExisteConexion(desde, hasta)
}

// ObtenerVecinos obtiene las cuevas vecinas de una cueva
func (ch *CaveHandler) ObtenerVecinos(id string) ([]string, error) {
	if ch.cuevaService == nil {
		return nil, fmt.Errorf("servicio de cueva no inicializado")
	}

	if id == "" {
		return nil, fmt.Errorf("ID de cueva no puede estar vacío")
	}

	vecinos := ch.cuevaService.ObtenerVecinos(id)
	return vecinos, nil
}

// CalcularDistancia calcula la distancia directa entre dos cuevas
func (ch *CaveHandler) CalcularDistancia(desde, hasta string) (float64, error) {
	if ch.cuevaService == nil {
		return 0, fmt.Errorf("servicio de cueva no inicializado")
	}

	if desde == "" || hasta == "" {
		return 0, fmt.Errorf("IDs de cuevas no pueden estar vacíos")
	}

	return ch.cuevaService.CalcularDistanciaDirecta(desde, hasta)
}

// EstablecerUbicacion establece la ubicación de una cueva
func (ch *CaveHandler) EstablecerUbicacion(id string, x, y float64) error {
	if ch.cuevaService == nil {
		return fmt.Errorf("servicio de cueva no inicializado")
	}

	if id == "" {
		return fmt.Errorf("ID de cueva no puede estar vacío")
	}

	return ch.cuevaService.EstablecerUbicacion(id, x, y)
}

// ObtenerEstadisticasCueva obtiene estadísticas de una cueva
func (ch *CaveHandler) ObtenerEstadisticasCueva(id string) (map[string]interface{}, error) {
	if ch.cuevaService == nil {
		return nil, fmt.Errorf("servicio de cueva no inicializado")
	}

	if id == "" {
		return nil, fmt.Errorf("ID de cueva no puede estar vacío")
	}

	return ch.cuevaService.ObtenerEstadisticas(id)
}

// GenerarReporteCueva genera un reporte detallado de una cueva
func (ch *CaveHandler) GenerarReporteCueva(id string) (string, error) {
	if ch.cuevaService == nil {
		return "", fmt.Errorf("servicio de cueva no inicializado")
	}

	if id == "" {
		return "", fmt.Errorf("ID de cueva no puede estar vacío")
	}

	detalle, err := ch.cuevaService.ObtenerDetalleCueva(id)
	if err != nil {
		return "", fmt.Errorf("error obteniendo detalles de cueva: %w", err)
	}

	reporte := fmt.Sprintf("REPORTE DE CUEVA\n")
	reporte += fmt.Sprintf("================\n")
	reporte += fmt.Sprintf("ID: %s\n", detalle.ID)
	reporte += fmt.Sprintf("Nombre: %s\n", detalle.Nombre)
	reporte += fmt.Sprintf("Ubicación: (%.2f, %.2f)\n", detalle.X, detalle.Y)
	reporte += fmt.Sprintf("Conexiones: %d\n", detalle.NumConexiones)

	if len(detalle.Recursos) > 0 {
		reporte += "Recursos:\n"
		for recurso, cantidad := range detalle.Recursos {
			reporte += fmt.Sprintf("  - %s: %d\n", recurso, cantidad)
		}
	} else {
		reporte += "Sin recursos almacenados\n"
	}

	return reporte, nil
}

// Métodos auxiliares privados

func (ch *CaveHandler) validarSolicitudCueva(solicitud service.SolicitudCueva) error {
	if solicitud.ID == "" {
		return fmt.Errorf("ID de cueva no puede estar vacío")
	}

	if solicitud.Nombre == "" {
		return fmt.Errorf("nombre de cueva no puede estar vacío")
	}

	if len(solicitud.ID) > 50 {
		return fmt.Errorf("ID de cueva no puede tener más de 50 caracteres")
	}

	if len(solicitud.Nombre) > 100 {
		return fmt.Errorf("nombre de cueva no puede tener más de 100 caracteres")
	}

	// Validar coordenadas dentro de rangos razonables
	if solicitud.X < -1000 || solicitud.X > 1000 {
		return fmt.Errorf("coordenada X debe estar entre -1000 y 1000")
	}

	if solicitud.Y < -1000 || solicitud.Y > 1000 {
		return fmt.Errorf("coordenada Y debe estar entre -1000 y 1000")
	}

	return nil
}
