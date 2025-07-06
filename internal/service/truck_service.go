package service

import (
	"fmt"
	"proyecto-grafos-go/internal/domain"
	"time"
)

// TipoCamion define los tipos de camión disponibles
type TipoCamion string

const (
	CamionPequeno TipoCamion = "PEQUEÑO"
	CamionMediano TipoCamion = "MEDIANO"
	CamionGrande  TipoCamion = "GRANDE"
)

// EstadoCamion define los estados del camión durante la simulación
type EstadoCamion string

const (
	EnAlmacen    EstadoCamion = "EN_ALMACEN"
	EnTransito   EstadoCamion = "EN_TRANSITO"
	Entregando   EstadoCamion = "ENTREGANDO"
	Regresando   EstadoCamion = "REGRESANDO"
	Completado   EstadoCamion = "COMPLETADO"
	Interrumpido EstadoCamion = "INTERRUMPIDO"
)

// Camion representa un camión de entrega
type Camion struct {
	ID                 string         `json:"id"`
	Tipo               TipoCamion     `json:"tipo"`
	CapacidadMaxima    int            `json:"capacidad_maxima"`
	VelocidadPromedio  float64        `json:"velocidad_promedio"` // km/h
	CargaActual        map[string]int `json:"carga_actual"`
	CuevaActual        string         `json:"cueva_actual"`
	Estado             EstadoCamion   `json:"estado"`
	RutaAsignada       *domain.Ruta   `json:"ruta_asignada"`
	TiempoInicio       time.Time      `json:"tiempo_inicio"`
	TiempoFin          time.Time      `json:"tiempo_fin"`
	DistanciaRecorrida float64        `json:"distancia_recorrida"`
}

// SimulacionResultado representa el resultado de una simulación de entrega
type SimulacionResultado struct {
	CamionID            string                    `json:"camion_id"`
	TipoRecorrido       TipoRecorrido             `json:"tipo_recorrido"`
	RutaCompleta        []string                  `json:"ruta_completa"`
	EntregasRealizadas  map[string]map[string]int `json:"entregas_realizadas"` // cueva_id -> recurso -> cantidad
	TiempoTotal         time.Duration             `json:"tiempo_total"`
	DistanciaTotal      float64                   `json:"distancia_total"`
	Exitoso             bool                      `json:"exitoso"`
	Errores             []string                  `json:"errores"`
	EstadisticasEntrega map[string]interface{}    `json:"estadisticas_entrega"`
}

// TruckService proporciona funcionalidades para la simulación de camiones
type TruckService struct {
	traversalService *TraversalService
	graphService     *ServicioGrafo
	camiones         map[string]*Camion
}

// NuevoTruckService crea una nueva instancia del servicio de camiones
func NuevoTruckService(traversalService *TraversalService, graphService *ServicioGrafo) *TruckService {
	return &TruckService{
		traversalService: traversalService,
		graphService:     graphService,
		camiones:         make(map[string]*Camion),
	}
}

// CrearCamion crea un nuevo camión con especificaciones dadas
func (ts *TruckService) CrearCamion(id string, tipo TipoCamion, cuevaOrigen string) (*Camion, error) {
	if _, existe := ts.camiones[id]; existe {
		return nil, fmt.Errorf("camión con ID '%s' ya existe", id)
	}

	// Configurar capacidad y velocidad según el tipo
	var capacidad int
	var velocidad float64

	switch tipo {
	case CamionPequeno:
		capacidad = 100
		velocidad = 60.0 // km/h
	case CamionMediano:
		capacidad = 200
		velocidad = 50.0 // km/h
	case CamionGrande:
		capacidad = 400
		velocidad = 40.0 // km/h
	default:
		return nil, fmt.Errorf("tipo de camión no válido: %s", tipo)
	}

	camion := &Camion{
		ID:                 id,
		Tipo:               tipo,
		CapacidadMaxima:    capacidad,
		VelocidadPromedio:  velocidad,
		CargaActual:        make(map[string]int),
		CuevaActual:        cuevaOrigen,
		Estado:             EnAlmacen,
		DistanciaRecorrida: 0.0,
	}

	ts.camiones[id] = camion
	return camion, nil
}

// CargarInsumos carga insumos en el camión
func (ts *TruckService) CargarInsumos(camionID string, insumos map[string]int) error {
	camion, existe := ts.camiones[camionID]
	if !existe {
		return fmt.Errorf("camión '%s' no encontrado", camionID)
	}

	if camion.Estado != EnAlmacen {
		return fmt.Errorf("camión debe estar en almacén para cargar insumos")
	}

	// Verificar capacidad
	pesoTotal := 0
	for _, cantidad := range insumos {
		pesoTotal += cantidad
	}

	pesoActual := 0
	for _, cantidad := range camion.CargaActual {
		pesoActual += cantidad
	}

	if pesoActual+pesoTotal > camion.CapacidadMaxima {
		return fmt.Errorf("excede capacidad máxima del camión (%d)", camion.CapacidadMaxima)
	}

	// Cargar insumos
	for recurso, cantidad := range insumos {
		camion.CargaActual[recurso] += cantidad
	}

	return nil
}

// SimularEntregaDFS simula la entrega de insumos usando recorrido DFS
func (ts *TruckService) SimularEntregaDFS(grafo *domain.Grafo, camionID string, cuevaOrigen string) (*SimulacionResultado, error) {
	return ts.simularEntrega(grafo, camionID, cuevaOrigen, DFS)
}

// SimularEntregaBFS simula la entrega de insumos usando recorrido BFS
func (ts *TruckService) SimularEntregaBFS(grafo *domain.Grafo, camionID string, cuevaOrigen string) (*SimulacionResultado, error) {
	return ts.simularEntrega(grafo, camionID, cuevaOrigen, BFS)
}

// simularEntrega realiza la simulación de entrega con el algoritmo especificado
func (ts *TruckService) simularEntrega(grafo *domain.Grafo, camionID string, cuevaOrigen string, tipoRecorrido TipoRecorrido) (*SimulacionResultado, error) {
	camion, existe := ts.camiones[camionID]
	if !existe {
		return nil, fmt.Errorf("camión '%s' no encontrado", camionID)
	}

	// Verificar que el camión tenga carga
	if len(camion.CargaActual) == 0 {
		return nil, fmt.Errorf("camión no tiene carga para entregar")
	}

	// Inicializar simulación
	camion.Estado = EnTransito
	camion.TiempoInicio = time.Now()
	camion.CuevaActual = cuevaOrigen
	camion.DistanciaRecorrida = 0.0

	resultado := &SimulacionResultado{
		CamionID:            camionID,
		TipoRecorrido:       tipoRecorrido,
		RutaCompleta:        make([]string, 0),
		EntregasRealizadas:  make(map[string]map[string]int),
		Errores:             make([]string, 0),
		EstadisticasEntrega: make(map[string]interface{}),
	}

	// Obtener ruta usando el algoritmo especificado
	var recorrido *RecorridoResultado
	var err error

	switch tipoRecorrido {
	case DFS:
		recorrido, err = ts.traversalService.RealizarRecorridoDFS(grafo, cuevaOrigen)
	case BFS:
		recorrido, err = ts.traversalService.RealizarRecorridoBFS(grafo, cuevaOrigen)
	default:
		return nil, fmt.Errorf("tipo de recorrido no válido: %s", tipoRecorrido)
	}

	if err != nil {
		resultado.Errores = append(resultado.Errores, fmt.Sprintf("Error en recorrido: %s", err.Error()))
		resultado.Exitoso = false
		return resultado, err
	}

	// Crear ruta para el camión
	ruta := domain.NuevaRuta(fmt.Sprintf("ruta_%s_%s", camionID, tipoRecorrido))

	// Simular entrega en cada cueva
	cargaOriginal := make(map[string]int)
	for recurso, cantidad := range camion.CargaActual {
		cargaOriginal[recurso] = cantidad
	}

	entregasExitosas := 0
	for _, cuevaID := range recorrido.CuevasVisitas {
		// Agregar cueva a la ruta
		distancia := 0.0
		if len(ruta.CuevaIDs) > 0 {
			distancia = ts.obtenerDistanciaEntreAristas(grafo, ruta.UltimaCueva(), cuevaID)
		}
		ruta.AgregarCueva(cuevaID, distancia)

		// Actualizar posición del camión
		camion.CuevaActual = cuevaID
		camion.DistanciaRecorrida += distancia

		// Obtener cueva para verificar necesidades
		cueva, existe := grafo.ObtenerCueva(cuevaID)
		if !existe {
			resultado.Errores = append(resultado.Errores, fmt.Sprintf("Cueva '%s' no encontrada", cuevaID))
			continue
		}

		// Simular entrega de insumos
		camion.Estado = Entregando
		entregaEnCueva := make(map[string]int)

		// Determinar qué entregar basándose en las necesidades de la cueva
		for recurso, cantidadDisponible := range camion.CargaActual {
			if cantidadDisponible > 0 {
				// Simplificación: entregar cantidad proporcional
				cantidadAEntregar := cantidadDisponible / (len(recorrido.CuevasVisitas) - entregasExitosas)
				if cantidadAEntregar > 0 {
					entregaEnCueva[recurso] = cantidadAEntregar
					camion.CargaActual[recurso] -= cantidadAEntregar

					// Actualizar recursos de la cueva
					cueva.AgregarRecurso(recurso, cueva.ObtenerRecurso(recurso)+cantidadAEntregar)
				}
			}
		}

		resultado.EntregasRealizadas[cuevaID] = entregaEnCueva
		entregasExitosas++

		// Simular tiempo de entrega (basado en distancia y velocidad)
		if distancia > 0 {
			tiempoViaje := time.Duration(distancia/camion.VelocidadPromedio*3600) * time.Second
			time.Sleep(tiempoViaje / 1000) // Simulación acelerada
		}
	}

	// Finalizar simulación
	camion.Estado = Completado
	camion.TiempoFin = time.Now()
	camion.RutaAsignada = ruta

	resultado.RutaCompleta = recorrido.CuevasVisitas
	resultado.TiempoTotal = camion.TiempoFin.Sub(camion.TiempoInicio)
	resultado.DistanciaTotal = camion.DistanciaRecorrida
	resultado.Exitoso = len(resultado.Errores) == 0 && entregasExitosas > 0

	// Generar estadísticas
	resultado.EstadisticasEntrega["entregas_exitosas"] = entregasExitosas
	resultado.EstadisticasEntrega["carga_original"] = cargaOriginal
	resultado.EstadisticasEntrega["carga_restante"] = camion.CargaActual
	resultado.EstadisticasEntrega["eficiencia_entrega"] = float64(entregasExitosas) / float64(len(recorrido.CuevasVisitas)) * 100

	return resultado, nil
}

// obtenerDistanciaEntreAristas obtiene la distancia entre dos cuevas conectadas
func (ts *TruckService) obtenerDistanciaEntreAristas(grafo *domain.Grafo, desde, hasta string) float64 {
	for _, arista := range grafo.Aristas {
		if arista.Desde == desde && arista.Hasta == hasta {
			return arista.Distancia
		}
		// Si el grafo no es dirigido, también buscar en dirección inversa
		if !grafo.EsDirigido && arista.Desde == hasta && arista.Hasta == desde {
			return arista.Distancia
		}
	}
	return 0.0
}

// ObtenerCamion obtiene un camión por su ID
func (ts *TruckService) ObtenerCamion(camionID string) (*Camion, error) {
	camion, existe := ts.camiones[camionID]
	if !existe {
		return nil, fmt.Errorf("camión '%s' no encontrado", camionID)
	}
	return camion, nil
}

// ListarCamiones obtiene todos los camiones registrados
func (ts *TruckService) ListarCamiones() map[string]*Camion {
	return ts.camiones
}

// ReiniciarCamion reinicia un camión al estado inicial
func (ts *TruckService) ReiniciarCamion(camionID string, cuevaOrigen string) error {
	camion, existe := ts.camiones[camionID]
	if !existe {
		return fmt.Errorf("camión '%s' no encontrado", camionID)
	}

	camion.CargaActual = make(map[string]int)
	camion.CuevaActual = cuevaOrigen
	camion.Estado = EnAlmacen
	camion.RutaAsignada = nil
	camion.DistanciaRecorrida = 0.0
	camion.TiempoInicio = time.Time{}
	camion.TiempoFin = time.Time{}

	return nil
}

// EliminarCamion elimina un camión del sistema
func (ts *TruckService) EliminarCamion(camionID string) error {
	if _, existe := ts.camiones[camionID]; !existe {
		return fmt.Errorf("camión '%s' no encontrado", camionID)
	}

	delete(ts.camiones, camionID)
	return nil
}

// GenerarReporteSimulacion genera un reporte detallado de la simulación
func (ts *TruckService) GenerarReporteSimulacion(resultado *SimulacionResultado) string {
	reporte := "=== REPORTE DE SIMULACIÓN DE ENTREGA ===\n"
	reporte += fmt.Sprintf("Camión: %s\n", resultado.CamionID)
	reporte += fmt.Sprintf("Tipo de Recorrido: %s\n", resultado.TipoRecorrido)
	reporte += fmt.Sprintf("Éxito: %t\n", resultado.Exitoso)
	reporte += fmt.Sprintf("Tiempo Total: %v\n", resultado.TiempoTotal)
	reporte += fmt.Sprintf("Distancia Total: %.2f km\n", resultado.DistanciaTotal)
	reporte += fmt.Sprintf("Cuevas Visitadas: %d\n", len(resultado.RutaCompleta))

	reporte += "\n--- RUTA SEGUIDA ---\n"
	for i, cueva := range resultado.RutaCompleta {
		reporte += fmt.Sprintf("%d. %s\n", i+1, cueva)
	}

	reporte += "\n--- ENTREGAS REALIZADAS ---\n"
	for cueva, entregas := range resultado.EntregasRealizadas {
		if len(entregas) > 0 {
			reporte += fmt.Sprintf("Cueva %s:\n", cueva)
			for recurso, cantidad := range entregas {
				reporte += fmt.Sprintf("  - %s: %d unidades\n", recurso, cantidad)
			}
		}
	}

	if len(resultado.Errores) > 0 {
		reporte += "\n--- ERRORES ---\n"
		for _, error := range resultado.Errores {
			reporte += fmt.Sprintf("- %s\n", error)
		}
	}

	reporte += "\n--- ESTADÍSTICAS ---\n"
	for clave, valor := range resultado.EstadisticasEntrega {
		reporte += fmt.Sprintf("%s: %v\n", clave, valor)
	}

	return reporte
}
