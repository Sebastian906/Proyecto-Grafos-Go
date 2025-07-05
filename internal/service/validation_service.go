package service

import (
	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/pkg/algorithms"
)

type ServicioValidacion struct {
	grafo *domain.Grafo
}

func NuevoServicioValidacion(grafo *domain.Grafo) *ServicioValidacion {
	return &ServicioValidacion{grafo: grafo}
}

// 1d: Conectividad fuerte
func (sv *ServicioValidacion) EsFuertementeConectado() bool {
	return algorithms.VerificarConectividadFuerte(sv.grafo)
}

// 1e: Detección de pozos
func (sv *ServicioValidacion) DetectarPozos() []string {
	var pozos []string
	for id := range sv.grafo.Cuevas {
		if len(sv.grafo.AristasSalientes(id)) == 0 &&
			len(sv.grafo.ProximasAristas(id)) == 0 {
			pozos = append(pozos, id)
		}
	}
	return pozos
}
