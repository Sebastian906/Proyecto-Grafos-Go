package unit

import (
	"os"
	"proyecto-grafos-go/pkg/utils"
	"testing"
)

// Estructura de prueba para JSON/XML
type DatosPrueba struct {
	ID     string  `json:"id" xml:"id"`
	Nombre string  `json:"nombre" xml:"nombre"`
	Valor  float64 `json:"valor" xml:"valor"`
}

func TestFileUtils_ExisteArchivo(t *testing.T) {
	// Crear archivo temporal
	archivo, err := utils.CrearArchivoTemporal("test_existe")
	if err != nil {
		t.Fatalf("Error creando archivo temporal: %v", err)
	}
	defer os.Remove(archivo.Name())
	archivo.Close()

	// Test archivo que existe
	if !utils.ExisteArchivo(archivo.Name()) {
		t.Error("Archivo existente no detectado")
	}

	// Test archivo que no existe
	if utils.ExisteArchivo("archivo_inexistente.txt") {
		t.Error("Archivo inexistente detectado como existente")
	}
}

func TestFileUtils_ExisteDirectorio(t *testing.T) {
	// Test directorio que existe
	if !utils.ExisteDirectorio(".") {
		t.Error("Directorio actual no detectado")
	}

	// Test directorio que no existe
	if utils.ExisteDirectorio("directorio_inexistente") {
		t.Error("Directorio inexistente detectado como existente")
	}
}

func TestFileUtils_CrearDirectorio(t *testing.T) {
	dirPrueba := "test_dir_temp"
	defer os.RemoveAll(dirPrueba)

	// Test crear directorio
	err := utils.CrearDirectorio(dirPrueba)
	if err != nil {
		t.Errorf("Error creando directorio: %v", err)
	}

	if !utils.ExisteDirectorio(dirPrueba) {
		t.Error("Directorio no fue creado")
	}

	// Test crear directorio que ya existe (no debe dar error)
	err = utils.CrearDirectorio(dirPrueba)
	if err != nil {
		t.Errorf("Error al intentar crear directorio existente: %v", err)
	}
}

func TestFileUtils_LeerEscribirArchivo(t *testing.T) {
	archivo := "test_file_temp.txt"
	contenido := []byte("Contenido de prueba")
	defer os.Remove(archivo)

	// Test escribir archivo
	err := utils.EscribirArchivo(archivo, contenido)
	if err != nil {
		t.Errorf("Error escribiendo archivo: %v", err)
	}

	// Test leer archivo
	contenidoLeido, err := utils.LeerArchivo(archivo)
	if err != nil {
		t.Errorf("Error leyendo archivo: %v", err)
	}

	if string(contenidoLeido) != string(contenido) {
		t.Errorf("Contenido leído no coincide. Esperado: %s, Obtenido: %s",
			string(contenido), string(contenidoLeido))
	}

	// Test leer archivo inexistente
	_, err = utils.LeerArchivo("archivo_inexistente.txt")
	if err == nil {
		t.Error("Debería dar error al leer archivo inexistente")
	}
}

func TestFileUtils_CopiarArchivo(t *testing.T) {
	origen := "test_origen.txt"
	destino := "test_destino.txt"
	contenido := []byte("Contenido para copiar")

	defer os.Remove(origen)
	defer os.Remove(destino)

	// Crear archivo origen
	err := utils.EscribirArchivo(origen, contenido)
	if err != nil {
		t.Fatalf("Error creando archivo origen: %v", err)
	}

	// Test copiar archivo
	err = utils.CopiarArchivo(origen, destino)
	if err != nil {
		t.Errorf("Error copiando archivo: %v", err)
	}

	// Verificar que el archivo fue copiado
	contenidoDestino, err := utils.LeerArchivo(destino)
	if err != nil {
		t.Errorf("Error leyendo archivo destino: %v", err)
	}

	if string(contenidoDestino) != string(contenido) {
		t.Error("El contenido copiado no coincide con el original")
	}

	// Test copiar archivo inexistente
	err = utils.CopiarArchivo("inexistente.txt", "destino2.txt")
	if err == nil {
		t.Error("Debería dar error al copiar archivo inexistente")
	}
}

func TestFileUtils_EliminarArchivo(t *testing.T) {
	archivo := "test_eliminar.txt"
	contenido := []byte("Para eliminar")

	// Crear archivo
	err := utils.EscribirArchivo(archivo, contenido)
	if err != nil {
		t.Fatalf("Error creando archivo: %v", err)
	}

	// Verificar que existe
	if !utils.ExisteArchivo(archivo) {
		t.Fatal("Archivo no fue creado")
	}

	// Test eliminar archivo
	err = utils.EliminarArchivo(archivo)
	if err != nil {
		t.Errorf("Error eliminando archivo: %v", err)
	}

	// Verificar que fue eliminado
	if utils.ExisteArchivo(archivo) {
		t.Error("Archivo no fue eliminado")
	}

	// Test eliminar archivo inexistente (no debe dar error)
	err = utils.EliminarArchivo("inexistente.txt")
	if err != nil {
		t.Errorf("Error eliminando archivo inexistente: %v", err)
	}
}

func TestFileUtils_ObtenerExtension(t *testing.T) {
	tests := []struct {
		archivo   string
		extension string
	}{
		{"archivo.txt", ".txt"},
		{"imagen.PNG", ".png"},
		{"data.json", ".json"},
		{"sin_extension", ""},
		{"archivo.con.puntos.xml", ".xml"},
	}

	for _, test := range tests {
		ext := utils.ObtenerExtension(test.archivo)
		if ext != test.extension {
			t.Errorf("Para %s esperaba %s, obtuvo %s",
				test.archivo, test.extension, ext)
		}
	}
}

func TestFileUtils_ValidarExtension(t *testing.T) {
	extensionesValidas := []string{".json", ".xml", ".txt"}

	// Test extensiones válidas
	archivosValidos := []string{"datos.json", "config.XML", "readme.TXT"}
	for _, archivo := range archivosValidos {
		if !utils.ValidarExtension(archivo, extensionesValidas) {
			t.Errorf("Archivo válido rechazado: %s", archivo)
		}
	}

	// Test extensiones inválidas
	archivosInvalidos := []string{"imagen.png", "script.js", "sin_extension"}
	for _, archivo := range archivosInvalidos {
		if utils.ValidarExtension(archivo, extensionesValidas) {
			t.Errorf("Archivo inválido aceptado: %s", archivo)
		}
	}
}

func TestFileUtils_ObtenerNombreArchivo(t *testing.T) {
	tests := []struct {
		ruta     string
		esperado string
	}{
		{"/path/to/archivo.txt", "archivo.txt"},
		{"../datos.json", "datos.json"},
		{"simple.xml", "simple.xml"},
		{"/solo/directorio/", "directorio"},
	}

	for _, test := range tests {
		nombre := utils.ObtenerNombreArchivo(test.ruta)
		if nombre != test.esperado {
			t.Errorf("Para %s esperaba %s, obtuvo %s",
				test.ruta, test.esperado, nombre)
		}
	}
}

func TestFileUtils_ObtenerNombreSinExtension(t *testing.T) {
	tests := []struct {
		archivo  string
		esperado string
	}{
		{"archivo.txt", "archivo"},
		{"imagen.PNG", "imagen"},
		{"sin_extension", "sin_extension"},
		{"archivo.con.puntos.xml", "archivo.con.puntos"},
	}

	for _, test := range tests {
		nombre := utils.ObtenerNombreSinExtension(test.archivo)
		if nombre != test.esperado {
			t.Errorf("Para %s esperaba %s, obtuvo %s",
				test.archivo, test.esperado, nombre)
		}
	}
}

func TestFileUtils_ObtenerTamanoArchivo(t *testing.T) {
	archivo := "test_tamano.txt"
	contenido := []byte("12345")
	defer os.Remove(archivo)

	// Crear archivo
	err := utils.EscribirArchivo(archivo, contenido)
	if err != nil {
		t.Fatalf("Error creando archivo: %v", err)
	}

	// Test obtener tamaño
	tamano, err := utils.ObtenerTamanoArchivo(archivo)
	if err != nil {
		t.Errorf("Error obteniendo tamaño: %v", err)
	}

	if tamano != int64(len(contenido)) {
		t.Errorf("Tamaño incorrecto. Esperado: %d, Obtenido: %d",
			len(contenido), tamano)
	}

	// Test archivo inexistente
	_, err = utils.ObtenerTamanoArchivo("inexistente.txt")
	if err == nil {
		t.Error("Debería dar error para archivo inexistente")
	}
}

func TestFileUtils_JSON(t *testing.T) {
	archivo := "test_datos.json"
	defer os.Remove(archivo)

	datosPrueba := DatosPrueba{
		ID:     "test123",
		Nombre: "Prueba JSON",
		Valor:  42.5,
	}

	// Test guardar JSON
	err := utils.GuardarJSON(archivo, datosPrueba)
	if err != nil {
		t.Errorf("Error guardando JSON: %v", err)
	}

	// Test cargar JSON
	var datosLeidos DatosPrueba
	err = utils.CargarJSON(archivo, &datosLeidos)
	if err != nil {
		t.Errorf("Error cargando JSON: %v", err)
	}

	// Verificar datos
	if datosLeidos.ID != datosPrueba.ID ||
		datosLeidos.Nombre != datosPrueba.Nombre ||
		datosLeidos.Valor != datosPrueba.Valor {
		t.Error("Datos JSON no coinciden después de cargar")
	}
}

func TestFileUtils_XML(t *testing.T) {
	archivo := "test_datos.xml"
	defer os.Remove(archivo)

	datosPrueba := DatosPrueba{
		ID:     "test456",
		Nombre: "Prueba XML",
		Valor:  33.7,
	}

	// Test guardar XML
	err := utils.GuardarXML(archivo, datosPrueba)
	if err != nil {
		t.Errorf("Error guardando XML: %v", err)
	}

	// Test cargar XML
	var datosLeidos DatosPrueba
	err = utils.CargarXML(archivo, &datosLeidos)
	if err != nil {
		t.Errorf("Error cargando XML: %v", err)
	}

	// Verificar datos
	if datosLeidos.ID != datosPrueba.ID ||
		datosLeidos.Nombre != datosPrueba.Nombre ||
		datosLeidos.Valor != datosPrueba.Valor {
		t.Error("Datos XML no coinciden después de cargar")
	}
}

func TestFileUtils_LimpiarNombreArchivo(t *testing.T) {
	tests := []struct {
		original string
		esperado string
	}{
		{"archivo normal.txt", "archivo normal.txt"},
		{"archivo<>:|?*.txt", "archivo______.txt"},
		{"  archivo con espacios  ", "archivo con espacios"},
		{"", "archivo_sin_nombre"},
		{"///\\\\\\", "archivo_sin_nombre"},
	}

	for _, test := range tests {
		limpio := utils.LimpiarNombreArchivo(test.original)
		if limpio != test.esperado {
			t.Errorf("Para '%s' esperaba '%s', obtuvo '%s'",
				test.original, test.esperado, limpio)
		}
	}
}

func TestFileUtils_ConfiguracionCuevaAcme(t *testing.T) {
	// Test obtener ruta de configuración
	ruta := utils.ObtenerRutaConfiguracionCuevaAcme()
	esperado := "data/caves_cueva_acme_default.json"
	if ruta != esperado {
		t.Errorf("Ruta incorrecta. Esperado: %s, Obtenido: %s", esperado, ruta)
	}

	// Test verificar archivo (puede o no existir dependiendo del entorno)
	existe := utils.VerificarArchivoConfiguracionCuevaAcme()
	t.Logf("Archivo de configuración Cueva Acme existe: %v", existe)

	// Si existe, probar cargar configuración
	if existe {
		rutaCargada, err := utils.CargarConfiguracionCuevaAcmePorDefecto()
		if err != nil {
			t.Errorf("Error cargando configuración por defecto: %v", err)
		}
		if rutaCargada != esperado {
			t.Errorf("Ruta cargada incorrecta. Esperado: %s, Obtenido: %s", esperado, rutaCargada)
		}
	}
}
