package domain

import "fmt"

// Representación de cueva en el sistema
type Cueva struct {
	ID       string         `json:"id" xml:"id"`
	Nombre   string         `json:"nombre" xml:"nombre"`
	Recursos map[string]int `json:"recursos" xml:"recursos"`
	X        float64        `json:"x" xml:"x"` // Coordenada x
	Y        float64        `json:"y" xml:"y"` // Coordenada y
}

// Función para crear una nueva cueva
func NuevaCueva(id, nombre string) *Cueva {
	return &Cueva{
		ID:       id,
		Nombre:   nombre,
		Recursos: make(map[string]int),
		X:        0,
		Y:        0,
	}
}

// Función de agregar un recurso a la cueva
func (c *Cueva) AgregarRecurso(recurso string, cantidad int) {
	c.Recursos[recurso] = cantidad
}

// Función de obtener cantidad de un recurso específico
func (c *Cueva) ObtenerRecurso(recurso string) int {
	return c.Recursos[recurso]
}

// Función para formatear los datos de la cueva
func (c *Cueva) String() string {
	return fmt.Sprintf("Cueva{ID: %s, Nombre: %s, Recursos: %v}", c.ID, c.Nombre, c.Recursos)
}