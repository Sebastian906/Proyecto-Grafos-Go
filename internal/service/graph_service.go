package service

import (
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/repository"
	"strings"
)

type ServicioGrafo struct {
	grafo       *domain.Grafo
	repositorio *repository.RepositorioArchivo
}

func NuevoServicioGrafo(grafo *domain.Grafo, repositorio *repository.RepositorioArchivo) *ServicioGrafo {
	return &ServicioGrafo{grafo: grafo, repositorio: repositorio}
}

// 1a: Cargar grafo desde archivo
func (sg *ServicioGrafo) CargarGrafo(archivo string) error {
	// Detectar tipo de archivo por extensión
	if strings.HasSuffix(strings.ToLower(archivo), ".xml") {
		grafo, err := sg.repositorio.CargarXML(archivo)
		if err != nil {
			return err
		}
		*sg.grafo = *grafo
		return nil
	} else if strings.HasSuffix(strings.ToLower(archivo), ".txt") {
		grafo, err := sg.repositorio.CargarTXT(archivo)
		if err != nil {
			return err
		}
		*sg.grafo = *grafo
		return nil
	}

	// Por defecto, cargar como JSON
	grafo, err := sg.repositorio.CargarJSON(archivo)
	if err != nil {
		return err
	}
	*sg.grafo = *grafo // Actualizar el grafo existente
	return nil
}

// 1c: Cambiar tipo de grafo (dirigido/no dirigido)
func (sg *ServicioGrafo) CambiarTipoGrafo(esDirigido bool) {
	sg.grafo.EsDirigido = esDirigido
	// Reconectar aristas si es necesario
}

// 1f: Grados de vértices
func (sg *ServicioGrafo) ObtenerGradosVertices() map[string]map[string]int {
	grados := make(map[string]map[string]int)
	for id := range sg.grafo.Cuevas {
		grados[id] = map[string]int{
			"entrante": len(sg.grafo.ProximasAristas(id)),
			"saliente": len(sg.grafo.AristasSalientes(id)),
			"total":    len(sg.grafo.ProximasAristas(id)) + len(sg.grafo.AristasSalientes(id)),
		}
	}
	return grados
}

// ObtenerGrafo devuelve el grafo actual
func (sg *ServicioGrafo) ObtenerGrafo() *domain.Grafo {
	return sg.grafo
}
