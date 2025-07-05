package service

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
)

type ServicioCueva struct {
	grafo *domain.Grafo
}

// Crea una nuevo servicio de cuevas
func ServicioNuevaCueva(grafo *domain.Grafo) *ServicioCueva {
	return &ServicioCueva{
		grafo: grafo,
	}
}

// Solicitud para crear una nueva cueva
type SolicitudCrearCrueva struct {
	ID       string         `json:"id"`
	Nombre   string         `json:"nombre"`
	X        float64        `json:"x"`
	Y        float64        `json:"y"`
	Recursos map[string]int `json:"recursos"`
}

// Solicitud para conectar dos cuevas
type SolicitudConectarCuevas struct {
	DesdeCuevaId    string  `json:"desde_cueva_id"`
	HastaCuevaId    string  `json:"hasta_cueva_id"`
	Distancia       float64 `json:"distancia"`
	EsDirigido      bool    `json:"es_dirigido"`
	EsBidireccional bool    `json:"es_bidireccional"`
}

// crea una nueva cueva en el grafo
func (sc *ServicioCueva) CrearCueva(solicitud *SolicitudCrearCrueva) (*domain.Cueva, error) {
	// Validar que el ID no esté vacío
	if solicitud.ID == "" {
		return nil, fmt.Errorf("ID de cueva no puede estar vacío")
	}

	// Validar que el nombre no esté vacío
	if solicitud.Nombre == "" {
		return nil, fmt.Errorf("nombre de la cueva no puede estar vacío")
	}

	// Verificar que la cueva no exista ya
	if _, existe := sc.grafo.ObtenerCueva(solicitud.ID); existe {
		return nil, fmt.Errorf("cueva con el ID %s ya existe", solicitud.ID)
	}

	// Crear la nueva cueva
	cueva := domain.NuevaCueva(solicitud.ID, solicitud.Nombre)
	cueva.X = solicitud.X
	cueva.Y = solicitud.Y

	// Agregar recursos si se proporcionaron
	if solicitud.Recursos != nil {
		for recurso, cantidad := range solicitud.Recursos {
			cueva.AgregarRecurso(recurso, cantidad)
		}
	}

	// Agregar la cueva al grafo
	if err := sc.grafo.AgregarCueva(cueva); err != nil {
		return nil, fmt.Errorf("fallo al agregar la cueva: %v", err)
	}

	return cueva, nil
}

// ConnectCaves conecta dos cuevas existentes
func (sc *ServicioCueva) ConectarCuevas(solicitud *SolicitudConectarCuevas) error {
	// Validar que las cuevas existan
	if _, existe := sc.grafo.ObtenerCueva(solicitud.DesdeCuevaId); !existe {
		return fmt.Errorf("cueva %s no existe", solicitud.DesdeCuevaId)
	}

	if _, existe := sc.grafo.ObtenerCueva(solicitud.HastaCuevaId); !existe {
		return fmt.Errorf("cueva %s no existe", solicitud.HastaCuevaId)
	}

	// Validar que no sea una conexión a sí mismo
	if solicitud.DesdeCuevaId == solicitud.HastaCuevaId {
		return fmt.Errorf("no se puede conectar la cueva hacia si misma")
	}

	// Validar que la distancia sea positiva
	if solicitud.Distancia <= 0 {
		return fmt.Errorf("distancia debe ser un valor positivo")
	}

	// Crear la arista principal
	arista := domain.NuevaArista(solicitud.DesdeCuevaId, solicitud.HastaCuevaId, solicitud.Distancia, solicitud.EsDirigido)
	if err := sc.grafo.AgregarArista(arista); err != nil {
		return fmt.Errorf("fallo al agregar la arista: %v", err)
	}

	// Si es bidireccional, crear la arista inversa
	if solicitud.EsBidireccional {
		aristaInversa := domain.NuevaArista(solicitud.HastaCuevaId, solicitud.DesdeCuevaId, solicitud.Distancia, solicitud.EsDirigido)
		if err := sc.grafo.AgregarArista(aristaInversa); err != nil {
			return fmt.Errorf("fallo al agregar la arista inversa: %v", err)
		}
	}

	return nil
}

// ObtenerCueva obtiene una cueva por su ID
func (sc *ServicioCueva) ObtenerCueva(cuevaID string) (*domain.Cueva, error) {
	cueva, existe := sc.grafo.ObtenerCueva(cuevaID)
	if !existe {
		return nil, fmt.Errorf("cueva %s no encontrada", cuevaID)
	}
	return cueva, nil
}

// lista todas las cuevas
func (sc *ServicioCueva) ListarCuevas() []*domain.Cueva {
	var cuevas []*domain.Cueva
	for _, cueva := range sc.grafo.Cuevas {
		cuevas = append(cuevas, cueva)
	}
	return cuevas
}

// actualiza los datos de una cueva existente
func (sc *ServicioCueva) ActualizarCueva(cuevaID string, actualizaciones map[string]interface{}) error {
	cueva, existe := sc.grafo.ObtenerCueva(cuevaID)
	if !existe {
		return fmt.Errorf("cueva %s no encontrada", cuevaID)
	}

	// Actualizar campos permitidos
	for campo, valor := range actualizaciones {
		switch campo {
		case "nombre":
			if nombre, ok := valor.(string); ok {
				cueva.Nombre = nombre
			}
		case "x":
			if x, ok := valor.(float64); ok {
				cueva.X = x
			}
		case "y":
			if y, ok := valor.(float64); ok {
				cueva.Y = y
			}
		case "recursos":
			if recursos, ok := valor.(map[string]int); ok {
				cueva.Recursos = recursos
			}
		}
	}

	return nil
}

// agrega un recurso a una cueva
func (sc *ServicioCueva) AgregarRecurso(cuevaID, recurso string, cantidad int) error {
	cueva, existe := sc.grafo.ObtenerCueva(cuevaID)
	if !existe {
		return fmt.Errorf("cueva %s no encontrada", cuevaID)
	}

	cueva.AgregarRecurso(recurso, cantidad)
	return nil
}

// obtiene las conexiones de una cueva
func (sc *ServicioCueva) ObtenerConexiones(cuevaID string) (map[string]interface{}, error) {
	// Verificar que la cueva existe
	if _, existe := sc.grafo.ObtenerCueva(cuevaID); !existe {
		return nil, fmt.Errorf("cueva %s no encontrada", cuevaID)
	}

	proximasAristas := sc.grafo.ProximasAristas(cuevaID)
	aristasSalientes := sc.grafo.AristasSalientes(cuevaID)

	conexiones := map[string]interface{}{
		"cueva_id":  cuevaID,
		"proximas":  make([]map[string]interface{}, 0),
		"salientes": make([]map[string]interface{}, 0),
		"vecinos":   sc.grafo.ObtenerVecinos(cuevaID),
	}

	// Procesar conexiones entrantes
	for _, arista := range proximasAristas {
		conexiones["proximas"] = append(conexiones["proximas"].([]map[string]interface{}), map[string]interface{}{
			"desde":        arista.Desde,
			"distancia":    arista.Distancia,
			"es_dirigida":  arista.EsDirigido,
			"es_obstruida": arista.EsObstruido,
		})
	}

	// Procesar conexiones salientes
	for _, arista := range aristasSalientes {
		conexiones["salientes"] = append(conexiones["salientes"].([]map[string]interface{}), map[string]interface{}{
			"hasta":        arista.Hasta,
			"distancia":    arista.Distancia,
			"es_dirigida":  arista.EsDirigido,
			"es_obstruida": arista.EsObstruido,
		})
	}

	return conexiones, nil
}

// elimina una cueva y todas sus conexiones
func (sc *ServicioCueva) EliminarCueva(cuevaID string) error {
	// Verificar que la cueva existe
	if _, existe := sc.grafo.ObtenerCueva(cuevaID); !existe {
		return fmt.Errorf("cueva %s no encontrada", cuevaID)
	}

	// Eliminar todas las aristas relacionadas con esta cueva
	var aristasRelacionadas []*domain.Arista
	for _, arista := range sc.grafo.ObtenerAristas() {
		if arista.Desde != cuevaID && arista.Hasta != cuevaID {
			aristasRelacionadas = append(aristasRelacionadas, arista)
		}
	}

	// Actualizar las aristas en el grafo
	sc.grafo.Aristas = aristasRelacionadas

	// Eliminar la cueva del mapa
	delete(sc.grafo.Cuevas, cuevaID)

	return nil
}

// obtiene estadísticas de una cueva
func (sc *ServicioCueva) CuevasEstadisticas(cuevaID string) (map[string]interface{}, error) {
	cueva, existe := sc.grafo.ObtenerCueva(cuevaID)
	if !existe {
		return nil, fmt.Errorf("cueva %s no encontrada", cuevaID)
	}

	proximasAristas := sc.grafo.ProximasAristas(cuevaID)
	aristasSalientes := sc.grafo.AristasSalientes(cuevaID)

	// Calcular distancia total de conexiones
	distanciaTotalC := 0.0
	for _, arista := range proximasAristas {
		distanciaTotalC += arista.Distancia
	}

	distanciaTotalS := 0.0
	for _, arista := range aristasSalientes {
		distanciaTotalS += arista.Distancia
	}

	// Contar recursos
	totalRecursos := 0
	for _, cantidad := range cueva.Recursos {
		totalRecursos += cantidad
	}

	estadistica := map[string]interface{}{
		"cueva_id": cuevaID,
		"nombre":   cueva.Nombre,
		"posicion": map[string]float64{
			"x": cueva.X,
			"y": cueva.Y,
		},
		"conexiones": map[string]interface{}{
			"cuenta_proximos":    len(proximasAristas),
			"cuenta_salientes":   len(aristasSalientes),
			"cuenta_total":       len(proximasAristas) + len(aristasSalientes),
			"proxima_distancia":  distanciaTotalC,
			"saliente_distancia": distanciaTotalS,
			"distancia_total":    distanciaTotalC + distanciaTotalS,
		},
		"recursos": map[string]interface{}{
			"tipos":          len(cueva.Recursos),
			"cantidad_total": totalRecursos,
			"detalles":       cueva.Recursos,
		},
	}

	return estadistica, nil
}

// valida la red de cuevas
func (sc *ServicioCueva) ValidarRed() []string {
	var errores []string

	// Verificar cuevas sin nombre
	for _, cueva := range sc.grafo.Cuevas {
		if cueva.Nombre == "" {
			errores = append(errores, fmt.Sprintf("Cueva %s no tiene nombre", cueva.ID))
		}
	}

	// Verificar cuevas duplicadas por posición
	posiciones := make(map[string][]string)
	for _, cueva := range sc.grafo.Cuevas {
		posKey := fmt.Sprintf("%.2f,%.2f", cueva.X, cueva.Y)
		posiciones[posKey] = append(posiciones[posKey], cueva.ID)
	}

	for pos, cuevas := range posiciones {
		if len(cuevas) > 1 {
			errores = append(errores, fmt.Sprintf("Multiples cuevas en la posicion %s: %v", pos, cuevas))
		}
	}

	return errores
}
