package handler

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/service"
)

// GraphHandler maneja operaciones relacionadas con el grafo
type GraphHandler struct {
	grafoService *service.ServicioGrafo
}

// NuevoGraphHandler crea una nueva instancia del handler de grafo
func NuevoGraphHandler(grafoService *service.ServicioGrafo) *GraphHandler {
	return &GraphHandler{
		grafoService: grafoService,
	}
}

// CrearGrafo maneja la creación de un nuevo grafo
func (gh *GraphHandler) CrearGrafo(esDirigido bool) (*domain.Grafo, error) {
	if gh.grafoService == nil {
		return nil, fmt.Errorf("servicio de grafo no inicializado")
	}

	grafo := domain.NuevoGrafo(esDirigido)
	return grafo, nil
}

// ObtenerGrafo maneja la obtención del grafo actual
func (gh *GraphHandler) ObtenerGrafo() (*domain.Grafo, error) {
	if gh.grafoService == nil {
		return nil, fmt.Errorf("servicio de grafo no inicializado")
	}

	grafo := gh.grafoService.ObtenerGrafo()
	if grafo == nil {
		return nil, fmt.Errorf("no hay grafo cargado")
	}

	return grafo, nil
}

// CargarGrafoDesdeArchivo maneja la carga de un grafo desde archivo
func (gh *GraphHandler) CargarGrafoDesdeArchivo(nombreArchivo string) error {
	if gh.grafoService == nil {
		return fmt.Errorf("servicio de grafo no inicializado")
	}

	if nombreArchivo == "" {
		return fmt.Errorf("nombre de archivo no puede estar vacío")
	}

	return gh.grafoService.CargarGrafo(nombreArchivo)
}

// GuardarGrafoEnArchivo maneja el guardado de un grafo en archivo
func (gh *GraphHandler) GuardarGrafoEnArchivo(nombreArchivo string) error {
	if gh.grafoService == nil {
		return fmt.Errorf("servicio de grafo no inicializado")
	}

	if nombreArchivo == "" {
		return fmt.Errorf("nombre de archivo no puede estar vacío")
	}

	return gh.grafoService.GuardarGrafo(nombreArchivo)
}

// CambiarTipoGrafo maneja el cambio de tipo de grafo
func (gh *GraphHandler) CambiarTipoGrafo(esDirigido bool) error {
	if gh.grafoService == nil {
		return fmt.Errorf("servicio de grafo no inicializado")
	}

	gh.grafoService.CambiarTipoGrafo(esDirigido)
	return nil
}

// ObtenerEstadisticasGrafo maneja la obtención de estadísticas del grafo
func (gh *GraphHandler) ObtenerEstadisticasGrafo() (map[string]interface{}, error) {
	if gh.grafoService == nil {
		return nil, fmt.Errorf("servicio de grafo no inicializado")
	}

	grafo := gh.grafoService.ObtenerGrafo()
	if grafo == nil {
		return nil, fmt.Errorf("no hay grafo cargado")
	}

	estadisticas := map[string]interface{}{
		"total_cuevas":       len(grafo.Cuevas),
		"total_aristas":      len(grafo.Aristas),
		"es_dirigido":        grafo.EsDirigido,
		"aristas_obstruidas": gh.contarAristasObstruidas(grafo),
	}

	return estadisticas, nil
}

// ValidarGrafo maneja la validación de integridad del grafo
func (gh *GraphHandler) ValidarGrafo() ([]string, error) {
	if gh.grafoService == nil {
		return nil, fmt.Errorf("servicio de grafo no inicializado")
	}

	grafo := gh.grafoService.ObtenerGrafo()
	if grafo == nil {
		return nil, fmt.Errorf("no hay grafo cargado")
	}

	var errores []string

	// Validar que todas las aristas tengan cuevas válidas
	for _, arista := range grafo.Aristas {
		if _, existe := grafo.Cuevas[arista.Desde]; !existe {
			errores = append(errores, fmt.Sprintf("cueva origen '%s' no existe", arista.Desde))
		}
		if _, existe := grafo.Cuevas[arista.Hasta]; !existe {
			errores = append(errores, fmt.Sprintf("cueva destino '%s' no existe", arista.Hasta))
		}
		if arista.Distancia < 0 {
			errores = append(errores, fmt.Sprintf("distancia negativa en arista %s -> %s", arista.Desde, arista.Hasta))
		}
	}

	// Validar que no haya cuevas duplicadas
	idsVistos := make(map[string]bool)
	for id := range grafo.Cuevas {
		if idsVistos[id] {
			errores = append(errores, fmt.Sprintf("ID de cueva duplicado: %s", id))
		}
		idsVistos[id] = true
	}

	return errores, nil
}

// LimpiarGrafo maneja la limpieza del grafo actual
func (gh *GraphHandler) LimpiarGrafo() error {
	if gh.grafoService == nil {
		return fmt.Errorf("servicio de grafo no inicializado")
	}

	grafo := gh.grafoService.ObtenerGrafo()
	if grafo == nil {
		return fmt.Errorf("no hay grafo cargado")
	}

	// Limpiar cuevas y aristas
	grafo.Cuevas = make(map[string]*domain.Cueva)
	grafo.Aristas = []*domain.Arista{}

	return nil
}

// ObtenerInformacionCompleta maneja la obtención de información completa del grafo
func (gh *GraphHandler) ObtenerInformacionCompleta() (string, error) {
	if gh.grafoService == nil {
		return "", fmt.Errorf("servicio de grafo no inicializado")
	}

	grafo := gh.grafoService.ObtenerGrafo()
	if grafo == nil {
		return "", fmt.Errorf("no hay grafo cargado")
	}

	info := fmt.Sprintf("INFORMACIÓN DEL GRAFO:\n")
	info += fmt.Sprintf("- Tipo: %s\n", gh.obtenerTipoGrafo(grafo))
	info += fmt.Sprintf("- Cuevas: %d\n", len(grafo.Cuevas))
	info += fmt.Sprintf("- Aristas: %d\n", len(grafo.Aristas))
	info += fmt.Sprintf("- Aristas obstruidas: %d\n", gh.contarAristasObstruidas(grafo))

	info += "\nCUEVAS:\n"
	for id, cueva := range grafo.Cuevas {
		info += fmt.Sprintf("- %s: %s (%.2f, %.2f)\n", id, cueva.Nombre, cueva.X, cueva.Y)
	}

	info += "\nCONEXIONES:\n"
	for _, arista := range grafo.Aristas {
		estado := "OK"
		if arista.EsObstruido {
			estado = "OBSTRUIDO"
		}
		info += fmt.Sprintf("- %s -> %s: %.2f km [%s]\n", arista.Desde, arista.Hasta, arista.Distancia, estado)
	}

	return info, nil
}

// ObtenerEstadisticas obtiene estadísticas del grafo
func (gh *GraphHandler) ObtenerEstadisticas() (*service.EstadisticasGrafo, error) {
	if gh.grafoService == nil {
		return nil, fmt.Errorf("servicio de grafo no inicializado")
	}

	return gh.grafoService.ObtenerEstadisticas()
}

// CargarGrafo carga un grafo desde archivo
func (gh *GraphHandler) CargarGrafo(archivo string) error {
	if gh.grafoService == nil {
		return fmt.Errorf("servicio de grafo no inicializado")
	}

	return gh.grafoService.CargarGrafo(archivo)
}

// GuardarGrafo guarda el grafo en un archivo
func (gh *GraphHandler) GuardarGrafo(archivo string) error {
	if gh.grafoService == nil {
		return fmt.Errorf("servicio de grafo no inicializado")
	}

	return gh.grafoService.GuardarGrafo(archivo)
}

// Métodos auxiliares privados

func (gh *GraphHandler) contarAristasObstruidas(grafo *domain.Grafo) int {
	contador := 0
	for _, arista := range grafo.Aristas {
		if arista.EsObstruido {
			contador++
		}
	}
	return contador
}

func (gh *GraphHandler) obtenerTipoGrafo(grafo *domain.Grafo) string {
	if grafo.EsDirigido {
		return "Dirigido"
	}
	return "No dirigido"
}

// CrearGrafoVacio maneja la creación de un grafo vacío
func (gh *GraphHandler) CrearGrafoVacio(esDirigido bool) error {
	if gh.grafoService == nil {
		return fmt.Errorf("servicio de grafo no inicializado")
	}

	nuevoGrafo := domain.NuevoGrafo(esDirigido)
	gh.grafoService.ActualizarGrafo(nuevoGrafo)

	return nil
}

// ExportarGrafo maneja la exportación del grafo a diferentes formatos
func (gh *GraphHandler) ExportarGrafo(nombreArchivo, formato string) error {
	if gh.grafoService == nil {
		return fmt.Errorf("servicio de grafo no inicializado")
	}

	if nombreArchivo == "" {
		return fmt.Errorf("nombre de archivo no puede estar vacío")
	}

	formatosValidos := []string{"json", "xml", "txt"}
	formatoValido := false
	for _, f := range formatosValidos {
		if formato == f {
			formatoValido = true
			break
		}
	}

	if !formatoValido {
		return fmt.Errorf("formato no válido. Formatos soportados: json, xml, txt")
	}

	// Agregar extensión si no la tiene
	if !gh.tieneExtension(nombreArchivo, formato) {
		nombreArchivo += "." + formato
	}

	return gh.grafoService.GuardarGrafo(nombreArchivo)
}

func (gh *GraphHandler) tieneExtension(nombre, extension string) bool {
	longitud := len(nombre)
	longitudExt := len(extension)

	if longitud <= longitudExt {
		return false
	}

	return nombre[longitud-longitudExt:] == extension
}
