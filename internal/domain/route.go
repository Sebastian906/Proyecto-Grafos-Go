package domain

import "fmt"

type Ruta struct {
	ID             string   `json:"id"`
	CuevaIDs       []string `json:"cueva_ids"`
	DistanciaTotal float64  `json:"distancia_total"`
	EstaCompleto   bool     `json:"esta_completo"`
}

// Función para crear una nueva ruta
func NuevaRuta(id string) *Ruta {
	return &Ruta{
		ID:             id,
		CuevaIDs:       make([]string, 0),
		DistanciaTotal: 0.0,
		EstaCompleto:   false,
	}
}

// Función para agregar una nueva ruta
func (r *Ruta) AgregarCueva(cuevaID string, distancia float64) {
	r.CuevaIDs = append(r.CuevaIDs, cuevaID)
	r.DistanciaTotal += distancia
}

// Función para obtener la última cueva en la ruta
func (r *Ruta) UltimaCueva() string {
	if len(r.CuevaIDs) == 0 {
		return ""
	}
	return r.CuevaIDs[len(r.CuevaIDs)-1]
}

// Función para formatear los datos de la ruta
func (r *Ruta) String() string {
	return fmt.Sprintf("Ruta{ID: %s, Cuevas: %v, Distancia: %.2f, Compeltado: %t}",
		r.ID, r.CuevaIDs, r.DistanciaTotal, r.EstaCompleto)
}