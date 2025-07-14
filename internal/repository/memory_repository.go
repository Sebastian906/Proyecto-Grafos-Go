package repository

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"sync"
)

// MemoryRepository implementa un repositorio en memoria para testing
type MemoryRepository struct {
	grafos map[string]*domain.Grafo
	mutex  sync.RWMutex
}

// NuevoMemoryRepository crea una nueva instancia del repositorio en memoria
func NuevoMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		grafos: make(map[string]*domain.Grafo),
	}
}

// Guardar almacena un grafo en memoria
func (mr *MemoryRepository) Guardar(nombre string, grafo *domain.Grafo) error {
	if grafo == nil {
		return fmt.Errorf("grafo no puede ser nil")
	}

	if nombre == "" {
		return fmt.Errorf("nombre no puede estar vacío")
	}

	mr.mutex.Lock()
	defer mr.mutex.Unlock()

	// Crear una copia del grafo para evitar mutaciones
	grafoCopia := mr.copiarGrafo(grafo)
	mr.grafos[nombre] = grafoCopia

	return nil
}

// Cargar recupera un grafo desde memoria
func (mr *MemoryRepository) Cargar(nombre string) (*domain.Grafo, error) {
	if nombre == "" {
		return nil, fmt.Errorf("nombre no puede estar vacío")
	}

	mr.mutex.RLock()
	defer mr.mutex.RUnlock()

	grafo, existe := mr.grafos[nombre]
	if !existe {
		return nil, fmt.Errorf("grafo '%s' no encontrado", nombre)
	}

	// Retornar una copia para evitar mutaciones
	return mr.copiarGrafo(grafo), nil
}

// Listar retorna todos los nombres de grafos almacenados
func (mr *MemoryRepository) Listar() []string {
	mr.mutex.RLock()
	defer mr.mutex.RUnlock()

	nombres := make([]string, 0, len(mr.grafos))
	for nombre := range mr.grafos {
		nombres = append(nombres, nombre)
	}

	return nombres
}

// Eliminar remueve un grafo de memoria
func (mr *MemoryRepository) Eliminar(nombre string) error {
	if nombre == "" {
		return fmt.Errorf("nombre no puede estar vacío")
	}

	mr.mutex.Lock()
	defer mr.mutex.Unlock()

	if _, existe := mr.grafos[nombre]; !existe {
		return fmt.Errorf("grafo '%s' no encontrado", nombre)
	}

	delete(mr.grafos, nombre)
	return nil
}

// Existe verifica si un grafo existe en memoria
func (mr *MemoryRepository) Existe(nombre string) bool {
	mr.mutex.RLock()
	defer mr.mutex.RUnlock()

	_, existe := mr.grafos[nombre]
	return existe
}

// Limpiar remueve todos los grafos de memoria
func (mr *MemoryRepository) Limpiar() {
	mr.mutex.Lock()
	defer mr.mutex.Unlock()

	mr.grafos = make(map[string]*domain.Grafo)
}

// Contar retorna el número de grafos almacenados
func (mr *MemoryRepository) Contar() int {
	mr.mutex.RLock()
	defer mr.mutex.RUnlock()

	return len(mr.grafos)
}

// copiarGrafo crea una copia profunda de un grafo
func (mr *MemoryRepository) copiarGrafo(original *domain.Grafo) *domain.Grafo {
	if original == nil {
		return nil
	}

	// Crear nuevo grafo
	copia := domain.NuevoGrafo(original.EsDirigido)

	// Copiar cuevas
	for id, cueva := range original.Cuevas {
		nuevaCueva := &domain.Cueva{
			ID:       cueva.ID,
			Nombre:   cueva.Nombre,
			X:        cueva.X,
			Y:        cueva.Y,
			Recursos: make(map[string]int),
		}

		// Copiar recursos
		for recurso, cantidad := range cueva.Recursos {
			nuevaCueva.Recursos[recurso] = cantidad
		}

		copia.Cuevas[id] = nuevaCueva
	}

	// Copiar aristas
	for _, arista := range original.Aristas {
		nuevaArista := &domain.Arista{
			Desde:       arista.Desde,
			Hasta:       arista.Hasta,
			Distancia:   arista.Distancia,
			EsDirigido:  arista.EsDirigido,
			EsObstruido: arista.EsObstruido,
		}

		copia.Aristas = append(copia.Aristas, nuevaArista)
	}

	return copia
}

// GuardarJSON guarda un grafo en formato JSON en memoria (simulado)
func (mr *MemoryRepository) GuardarJSON(nombre string, grafo *domain.Grafo) error {
	return mr.Guardar(nombre+".json", grafo)
}

// CargarJSON carga un grafo desde formato JSON en memoria (simulado)
func (mr *MemoryRepository) CargarJSON(nombre string) (*domain.Grafo, error) {
	return mr.Cargar(nombre)
}

// GuardarXML guarda un grafo en formato XML en memoria (simulado)
func (mr *MemoryRepository) GuardarXML(nombre string, grafo *domain.Grafo) error {
	return mr.Guardar(nombre+".xml", grafo)
}

// CargarXML carga un grafo desde formato XML en memoria (simulado)
func (mr *MemoryRepository) CargarXML(nombre string) (*domain.Grafo, error) {
	return mr.Cargar(nombre)
}

// GuardarTXT guarda un grafo en formato TXT en memoria (simulado)
func (mr *MemoryRepository) GuardarTXT(nombre string, grafo *domain.Grafo) error {
	return mr.Guardar(nombre+".txt", grafo)
}

// CargarTXT carga un grafo desde formato TXT en memoria (simulado)
func (mr *MemoryRepository) CargarTXT(nombre string) (*domain.Grafo, error) {
	return mr.Cargar(nombre)
}

// CrearBackup crea un backup del grafo en memoria
func (mr *MemoryRepository) CrearBackup(nombre string) error {
	grafo, err := mr.Cargar(nombre)
	if err != nil {
		return fmt.Errorf("error cargando grafo para backup: %w", err)
	}

	nombreBackup := fmt.Sprintf("%s.backup", nombre)
	return mr.Guardar(nombreBackup, grafo)
}

// RestaurarBackup restaura un grafo desde su backup
func (mr *MemoryRepository) RestaurarBackup(nombre string) error {
	nombreBackup := fmt.Sprintf("%s.backup", nombre)
	grafo, err := mr.Cargar(nombreBackup)
	if err != nil {
		return fmt.Errorf("error cargando backup: %w", err)
	}

	return mr.Guardar(nombre, grafo)
}

// ObtenerEstadisticas retorna estadísticas del repositorio
func (mr *MemoryRepository) ObtenerEstadisticas() map[string]interface{} {
	mr.mutex.RLock()
	defer mr.mutex.RUnlock()

	stats := map[string]interface{}{
		"total_grafos": len(mr.grafos),
		"nombres":      make([]string, 0, len(mr.grafos)),
	}

	for nombre := range mr.grafos {
		stats["nombres"] = append(stats["nombres"].([]string), nombre)
	}

	return stats
}
