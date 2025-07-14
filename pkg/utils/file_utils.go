package utils

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"proyecto-grafos-go/internal/domain"
	"proyecto-grafos-go/internal/repository"
)

// ExisteArchivo verifica si un archivo existe
func ExisteArchivo(ruta string) bool {
	_, err := os.Stat(ruta)
	return !os.IsNotExist(err)
}

// ExisteDirectorio verifica si un directorio existe
func ExisteDirectorio(ruta string) bool {
	info, err := os.Stat(ruta)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// CrearDirectorio crea un directorio si no existe
func CrearDirectorio(ruta string) error {
	if !ExisteDirectorio(ruta) {
		return os.MkdirAll(ruta, 0755)
	}
	return nil
}

// LeerArchivo lee un archivo completo y devuelve su contenido
func LeerArchivo(ruta string) ([]byte, error) {
	if !ExisteArchivo(ruta) {
		return nil, fmt.Errorf("el archivo %s no existe", ruta)
	}
	return os.ReadFile(ruta)
}

// EscribirArchivo escribe contenido a un archivo
func EscribirArchivo(ruta string, contenido []byte) error {
	// Crear directorio padre si no existe
	dir := filepath.Dir(ruta)
	if err := CrearDirectorio(dir); err != nil {
		return fmt.Errorf("error creando directorio %s: %v", dir, err)
	}

	return os.WriteFile(ruta, contenido, 0644)
}

// CopiarArchivo copia un archivo de origen a destino
func CopiarArchivo(origen, destino string) error {
	if !ExisteArchivo(origen) {
		return fmt.Errorf("el archivo origen %s no existe", origen)
	}

	archivoOrigen, err := os.Open(origen)
	if err != nil {
		return fmt.Errorf("error abriendo archivo origen: %v", err)
	}
	defer archivoOrigen.Close()

	// Crear directorio padre si no existe
	dir := filepath.Dir(destino)
	if err := CrearDirectorio(dir); err != nil {
		return fmt.Errorf("error creando directorio %s: %v", dir, err)
	}

	archivoDestino, err := os.Create(destino)
	if err != nil {
		return fmt.Errorf("error creando archivo destino: %v", err)
	}
	defer archivoDestino.Close()

	_, err = io.Copy(archivoDestino, archivoOrigen)
	if err != nil {
		return fmt.Errorf("error copiando archivo: %v", err)
	}

	return nil
}

// EliminarArchivo elimina un archivo si existe
func EliminarArchivo(ruta string) error {
	if !ExisteArchivo(ruta) {
		return nil // No es error si no existe
	}
	return os.Remove(ruta)
}

// ObtenerExtension obtiene la extensión de un archivo
func ObtenerExtension(nombreArchivo string) string {
	ext := filepath.Ext(nombreArchivo)
	return strings.ToLower(ext)
}

// ValidarExtension verifica si la extensión del archivo es válida
func ValidarExtension(nombreArchivo string, extensionesValidas []string) bool {
	ext := ObtenerExtension(nombreArchivo)
	for _, extValida := range extensionesValidas {
		if ext == strings.ToLower(extValida) {
			return true
		}
	}
	return false
}

// ObtenerNombreArchivo obtiene el nombre del archivo sin la ruta
func ObtenerNombreArchivo(ruta string) string {
	return filepath.Base(ruta)
}

// ObtenerNombreSinExtension obtiene el nombre del archivo sin extensión
func ObtenerNombreSinExtension(nombreArchivo string) string {
	ext := filepath.Ext(nombreArchivo)
	return strings.TrimSuffix(nombreArchivo, ext)
}

// ObtenerTamanoArchivo obtiene el tamaño de un archivo en bytes
func ObtenerTamanoArchivo(ruta string) (int64, error) {
	if !ExisteArchivo(ruta) {
		return 0, fmt.Errorf("el archivo %s no existe", ruta)
	}

	info, err := os.Stat(ruta)
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

// CargarJSON carga un archivo JSON en una estructura
func CargarJSON(ruta string, destino interface{}) error {
	datos, err := LeerArchivo(ruta)
	if err != nil {
		return fmt.Errorf("error leyendo archivo JSON %s: %v", ruta, err)
	}

	if err := json.Unmarshal(datos, destino); err != nil {
		return fmt.Errorf("error deserializando JSON: %v", err)
	}

	return nil
}

// GuardarJSON guarda una estructura como archivo JSON
func GuardarJSON(ruta string, datos interface{}) error {
	jsonData, err := json.MarshalIndent(datos, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando a JSON: %v", err)
	}

	return EscribirArchivo(ruta, jsonData)
}

// CargarXML carga un archivo XML en una estructura
func CargarXML(ruta string, destino interface{}) error {
	datos, err := LeerArchivo(ruta)
	if err != nil {
		return fmt.Errorf("error leyendo archivo XML %s: %v", ruta, err)
	}

	if err := xml.Unmarshal(datos, destino); err != nil {
		return fmt.Errorf("error deserializando XML: %v", err)
	}

	return nil
}

// GuardarXML guarda una estructura como archivo XML
func GuardarXML(ruta string, datos interface{}) error {
	xmlData, err := xml.MarshalIndent(datos, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando a XML: %v", err)
	}

	// Agregar header XML
	contenido := []byte(xml.Header + string(xmlData))
	return EscribirArchivo(ruta, contenido)
}

// ListarArchivos lista todos los archivos en un directorio
func ListarArchivos(directorio string) ([]string, error) {
	if !ExisteDirectorio(directorio) {
		return nil, fmt.Errorf("el directorio %s no existe", directorio)
	}

	archivos := []string{}
	err := filepath.Walk(directorio, func(ruta string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			archivos = append(archivos, ruta)
		}
		return nil
	})

	return archivos, err
}

// ListarArchivosPorExtension lista archivos con una extensión específica
func ListarArchivosPorExtension(directorio, extension string) ([]string, error) {
	archivos, err := ListarArchivos(directorio)
	if err != nil {
		return nil, err
	}

	var archivosFiltrados []string
	extension = strings.ToLower(extension)
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}

	for _, archivo := range archivos {
		if ObtenerExtension(archivo) == extension {
			archivosFiltrados = append(archivosFiltrados, archivo)
		}
	}

	return archivosFiltrados, nil
}

// CrearArchivoTemporal crea un archivo temporal con un prefijo
func CrearArchivoTemporal(prefijo string) (*os.File, error) {
	return os.CreateTemp("", prefijo+"_*.tmp")
}

// LimpiarNombreArchivo limpia un nombre de archivo de caracteres inválidos
func LimpiarNombreArchivo(nombre string) string {
	// Caracteres no permitidos en nombres de archivo
	caracteresInvalidos := []string{"<", ">", ":", "\"", "|", "?", "*", "/", "\\"}

	for _, char := range caracteresInvalidos {
		nombre = strings.ReplaceAll(nombre, char, "_")
	}

	// Remover espacios al inicio y final
	nombre = strings.TrimSpace(nombre)

	// Si queda vacío, usar un nombre por defecto
	if nombre == "" || nombre == "_" || strings.Trim(nombre, "_") == "" {
		nombre = "archivo_sin_nombre"
	}

	return nombre
}

// CargarConfiguracionCuevaAcmePorDefecto carga la configuración por defecto de Cueva Acme
func CargarConfiguracionCuevaAcmePorDefecto() (*domain.Grafo, error) {
	rutaDefecto := "caves_cueva_acme_default.json"
	rutaCompleta := "data/" + rutaDefecto

	// Verificar si el archivo existe
	if !ExisteArchivo(rutaCompleta) {
		return nil, fmt.Errorf("archivo de configuración por defecto no encontrado: %s", rutaCompleta)
	}

	// Crear un repositorio de archivos para cargar el JSON
	repo := repository.NuevoRepositorio("data")

	// Cargar el grafo usando el repositorio
	grafo, err := repo.CargarJSON(rutaDefecto)
	if err != nil {
		return nil, fmt.Errorf("error cargando configuración por defecto: %v", err)
	}

	return grafo, nil
}

// ObtenerRutaConfiguracionCuevaAcmePorDefecto obtiene solo la ruta del archivo de configuración por defecto
func ObtenerRutaConfiguracionCuevaAcmePorDefecto() (string, error) {
	rutaDefecto := "caves_cueva_acme_default.json"
	rutaCompleta := "data/" + rutaDefecto

	// Verificar si el archivo existe
	if !ExisteArchivo(rutaCompleta) {
		return "", fmt.Errorf("archivo de configuración por defecto no encontrado: %s", rutaCompleta)
	}

	return rutaDefecto, nil
}

// VerificarArchivoConfiguracionCuevaAcme verifica si existe el archivo de configuración
func VerificarArchivoConfiguracionCuevaAcme() bool {
	rutaDefecto := "data/caves_cueva_acme_default.json"
	return ExisteArchivo(rutaDefecto)
}

// ObtenerRutaConfiguracionCuevaAcme obtiene la ruta del archivo de configuración por defecto
func ObtenerRutaConfiguracionCuevaAcme() string {
	return "data/caves_cueva_acme_default.json"
}
