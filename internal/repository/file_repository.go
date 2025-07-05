package repository

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"proyecto-grafos-go/internal/domain"
	"strconv"
	"strings"
)

// Carga y guardado de grafos desde archivos
type RepositorioArchivo struct {
	dataDir string
}

// Función para crear un nuevo repositorio de archivos
func NuevoRepositorio(dataDir string) *RepositorioArchivo {
	return &RepositorioArchivo{
		dataDir: dataDir,
	}
}

// Estructura para datos de serialización
type DataGrafo struct {
	XMLNombre  xml.Name         `xml:"grafo" json:"-"`
	Cuevas     []*domain.Cueva  `xml:"cuevas>cueva" json:"cuevas"`
	Aristas    []*domain.Arista `xml:"aristas>arista" json:"aristas"`
	EsDirigido bool             `xml:"es_dirigido" json:"es_dirigido"`
}

// Función para cargar un grafo desde un archivo JSON
func (ra *RepositorioArchivo) CargarJSON(archivo string) (*domain.Grafo, error) {
	dirArchivo := filepath.Join(ra.dataDir, archivo)

	data, err := os.ReadFile(dirArchivo)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}

	var dataGrafo DataGrafo
	if err := json.Unmarshal(data, &dataGrafo); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return ra.construirGrafo(&dataGrafo)
}

// Función para cargar un grafo desde un archivo XML
func (ra *RepositorioArchivo) CargarXML(archivo string) (*domain.Grafo, error) {
	dirArchivo := filepath.Join(ra.dataDir, archivo)

	data, err := os.ReadFile(dirArchivo)
	if err != nil {
		return nil, fmt.Errorf("error reading XML file: %v", err)
	}

	var dataGrafo DataGrafo
	if err := xml.Unmarshal(data, &dataGrafo); err != nil {
		return nil, fmt.Errorf("error parsing XML: %v", err)
	}

	return ra.construirGrafo(&dataGrafo)
}

// Función para cargar un grafo desde un archivo de texto
func (ra *RepositorioArchivo) CargarTXT(archivo string) (*domain.Grafo, error) {
	dirArchivo := filepath.Join(ra.dataDir, archivo)

	file, err := os.Open(dirArchivo)
	if err != nil {
		return nil, fmt.Errorf("error opening TXT file: %v", err)
	}
	defer file.Close()

	var dataGrafo DataGrafo
	dataGrafo.Cuevas = make([]*domain.Cueva, 0)
	dataGrafo.Aristas = make([]*domain.Arista, 0)

	scanner := bufio.NewScanner(file)
	section := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignorar líneas vacías y comentarios
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Identificar secciones
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.Trim(line, "[]")
			continue
		}

		switch section {
		case "grafo":
			if err := ra.parseConfigGrafo(line, &dataGrafo); err != nil {
				return nil, fmt.Errorf("error analizando la configuración del grafo: %v", err)
			}
		case "cuevas":
			if err := ra.parseLineaCueva(line, &dataGrafo); err != nil {
				return nil, fmt.Errorf("error analizando la línea de la cueva: %v", err)
			}
		case "aristas":
			if err := ra.parseLineaArista(line, &dataGrafo); err != nil {
				return nil, fmt.Errorf("error analizando la línea de la arista: %v", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error al leer el archivo TXT: %v", err)
	}

	return ra.construirGrafo(&dataGrafo)
}

// parsea la configuración del grafo
func (ra *RepositorioArchivo) parseConfigGrafo(linea string, dataGrafo *DataGrafo) error {
	partes := strings.Split(linea, "=")
	if len(partes) != 2 {
		return fmt.Errorf("configuración de grafo inválida: %s", linea)
	}

	key := strings.TrimSpace(partes[0])
	valor := strings.TrimSpace(partes[1])

	switch key {
	case "dirigido":
		dirigido, err := strconv.ParseBool(valor)
		if err != nil {
			return fmt.Errorf("valor de dirección inválido: %s", valor)
		}
		dataGrafo.EsDirigido = dirigido
	}

	return nil
}

// parsea una línea de cueva del archivo TXT
func (ra *RepositorioArchivo) parseLineaCueva(linea string, dataGrafo *DataGrafo) error {
	// Formato: ID,Name,X,Y,recurso1:cantidad1,recurso2:cantidad2
	partes := strings.Split(linea, ",")
	if len(partes) < 4 {
		return fmt.Errorf("formato inválido de cueva: %s", linea)
	}

	id := strings.TrimSpace(partes[0])
	nombre := strings.TrimSpace(partes[1])

	x, err := strconv.ParseFloat(strings.TrimSpace(partes[2]), 64)
	if err != nil {
		return fmt.Errorf("coordinada X inválida: %s", partes[2])
	}

	y, err := strconv.ParseFloat(strings.TrimSpace(partes[3]), 64)
	if err != nil {
		return fmt.Errorf("coordinada Y inválida: %s", partes[3])
	}

	cueva := domain.NuevaCueva(id, nombre)
	cueva.X = x
	cueva.Y = y

	// Parsear recursos si existen
	for i := 4; i < len(partes); i++ {
		parteRecurso := strings.TrimSpace(partes[i])
		if parteRecurso == "" {
			continue
		}

		partesRecurso := strings.Split(parteRecurso, ":")
		if len(partesRecurso) != 2 {
			return fmt.Errorf("formato de recurso inválido: %s", parteRecurso)
		}

		recurso := strings.TrimSpace(partesRecurso[0])
		cantidad, err := strconv.Atoi(strings.TrimSpace(partesRecurso[1]))
		if err != nil {
			return fmt.Errorf("cantidad de recurso inválido: %s", partesRecurso[1])
		}

		cueva.AgregarRecurso(recurso, cantidad)
	}

	dataGrafo.Cuevas = append(dataGrafo.Cuevas, cueva)
	return nil
}

// parsea una línea de arista del archivo TXT
func (ra *RepositorioArchivo) parseLineaArista(linea string, dataGrafo *DataGrafo) error {
	// Formato: From,To,Distance,IsDirected
	partes := strings.Split(linea, ",")
	if len(partes) < 3 {
		return fmt.Errorf("formado de arista inválido: %s", linea)
	}

	desde := strings.TrimSpace(partes[0])
	hasta := strings.TrimSpace(partes[1])

	distancia, err := strconv.ParseFloat(strings.TrimSpace(partes[2]), 64)
	if err != nil {
		return fmt.Errorf("distancia inválida: %s", partes[2])
	}

	esDirigido := dataGrafo.EsDirigido
	if len(partes) > 3 {
		dirigido, err := strconv.ParseBool(strings.TrimSpace(partes[3]))
		if err != nil {
			return fmt.Errorf("valor de dirección inválido: %s", partes[3])
		}
		esDirigido = dirigido
	}

	arista := domain.NuevaArista(desde, hasta, distancia, esDirigido)
	dataGrafo.Aristas = append(dataGrafo.Aristas, arista)

	return nil
}

// Función para construir el grafo
func (ra *RepositorioArchivo) construirGrafo(dataGrafo *DataGrafo) (*domain.Grafo, error) {
	grafo := domain.NuevoGrafo(dataGrafo.EsDirigido)

	// Agregar cuevas
	for _, cueva := range dataGrafo.Cuevas {
		if err := grafo.AgregarCueva(cueva); err != nil {
			return nil, fmt.Errorf("error agregando cueva %s: %v", cueva.ID, err)
		}
	}

	// Agregar aristas
	for _, arista := range dataGrafo.Aristas {
		if err := grafo.AgregarArista(arista); err != nil {
			return nil, fmt.Errorf("error agregando arista %s->%s: %v", arista.Desde, arista.Hasta, err)
		}
	}

	return grafo, nil
}

// guardar el grafo en un archivo JSON
func (ra *RepositorioArchivo) GuardarJSON(grafo *domain.Grafo, archivo string) error {
	dirArchivo := filepath.Join(ra.dataDir, archivo)

	dataGrafo := ra.extraerDatosGrafo(grafo)

	data, err := json.MarshalIndent(dataGrafo, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	if err := os.WriteFile(dirArchivo, data, 0644); err != nil {
		return fmt.Errorf("error al escribir el archivo JSON: %v", err)
	}

	return nil
}

// guardar el grafo en un archivo XML
func (ra *RepositorioArchivo) GuardarXML(grafo *domain.Grafo, archivo string) error {
    dirArchivo := filepath.Join(ra.dataDir, archivo)
    
    dataGrafo := ra.extraerDatosGrafo(grafo)
    
    data, err := xml.MarshalIndent(dataGrafo, "", "  ")
    if err != nil {
        return fmt.Errorf("error ordenando el XML: %v", err)
    }
    
    // Agregar header XML
    xmlData := []byte(xml.Header + string(data))
    
    if err := os.WriteFile(dirArchivo, xmlData, 0644); err != nil {
        return fmt.Errorf("error al escribir el archivo XML: %v", err)
    }
    
    return nil
}

// extraer los datos del grafo para serialización
func (ra *RepositorioArchivo) extraerDatosGrafo(grafo *domain.Grafo) *DataGrafo {
	dataGrafo := &DataGrafo{
		EsDirigido: grafo.EsDirigido,
		Cuevas:     make([]*domain.Cueva, 0),
		Aristas:    make([]*domain.Arista, 0),
	}

	// Extraer cuevas
	for _, cueva := range grafo.Cuevas {
		dataGrafo.Cuevas = append(dataGrafo.Cuevas, cueva)
	}

	// Extraer aristas
	dataGrafo.Aristas = grafo.ObtenerAristas()

	return dataGrafo
}

// listar los archivos disponibles en el directorio de datos
func (ra *RepositorioArchivo) ListarArchivos() ([]string, error) {
    archivos, err := os.ReadDir(ra.dataDir)
    if err != nil {
        return nil, fmt.Errorf("error reading data directory: %v", err)
    }
    
    var nArchivos []string
    for _, archivo := range archivos {
        if !archivo.IsDir() {
            ext := filepath.Ext(archivo.Name())
            if ext == ".json" || ext == ".xml" || ext == ".txt" {
                nArchivos = append(nArchivos, archivo.Name())
            }
        }
    }
    
    return nArchivos, nil
}