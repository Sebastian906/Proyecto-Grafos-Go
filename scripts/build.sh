#!/bin/bash

# Script de compilación para el proyecto de grafos

echo "=== COMPILACIÓN DEL SISTEMA DE GESTIÓN DE CUEVAS ==="
echo ""

# Configuración
APP_NAME="sistema-cuevas"
MAIN_FILE="cmd/main.go"
BUILD_DIR="build"
VERSION=$(date +%Y%m%d-%H%M%S)

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Funciones auxiliares
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verificar Go
check_go() {
    if ! command -v go &> /dev/null; then
        log_error "Go no está instalado o no está en PATH"
        exit 1
    fi
    log_info "Go version: $(go version)"
}

# Limpiar builds anteriores
clean_build() {
    log_info "Limpiando builds anteriores..."
    rm -rf $BUILD_DIR
    rm -f $APP_NAME $APP_NAME.exe
    log_info "Limpieza completada"
}

# Descargar dependencias
download_deps() {
    log_info "Descargando dependencias..."
    go mod download
    if [ $? -ne 0 ]; then
        log_error "Error descargando dependencias"
        exit 1
    fi
    log_info "Dependencias descargadas"
}

# Ejecutar tests
run_tests() {
    log_info "Ejecutando tests..."
    go test ./...
    if [ $? -ne 0 ]; then
        log_warn "Algunos tests fallaron, pero continuando con la compilación"
    else
        log_info "Todos los tests pasaron"
    fi
}

# Compilar aplicación
build_app() {
    log_info "Compilando aplicación..."
    
    # Crear directorio de build
    mkdir -p $BUILD_DIR
    
    # Flags de compilación
    LDFLAGS="-X main.version=$VERSION -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    
    # Compilar para el sistema actual
    go build -ldflags "$LDFLAGS" -o $BUILD_DIR/$APP_NAME $MAIN_FILE
    if [ $? -ne 0 ]; then
        log_error "Error en la compilación"
        exit 1
    fi
    
    log_info "Aplicación compilada exitosamente: $BUILD_DIR/$APP_NAME"
}

# Compilar para múltiples plataformas
build_cross_platform() {
    log_info "Compilando para múltiples plataformas..."
    
    # Plataformas objetivo
    platforms=(
        "linux/amd64"
        "linux/386"
        "windows/amd64"
        "windows/386"
        "darwin/amd64"
        "darwin/arm64"
    )
    
    for platform in "${platforms[@]}"; do
        IFS='/' read -r -a platform_split <<< "$platform"
        GOOS="${platform_split[0]}"
        GOARCH="${platform_split[1]}"
        
        output_name="$BUILD_DIR/${APP_NAME}_${GOOS}_${GOARCH}"
        if [ $GOOS = "windows" ]; then
            output_name="${output_name}.exe"
        fi
        
        log_info "Compilando para $GOOS/$GOARCH..."
        env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "$LDFLAGS" -o $output_name $MAIN_FILE
        if [ $? -ne 0 ]; then
            log_error "Error compilando para $GOOS/$GOARCH"
        else
            log_info "Compilado: $output_name"
        fi
    done
}

# Crear paquete de distribución
create_package() {
    log_info "Creando paquete de distribución..."
    
    # Copiar archivos necesarios
    cp -r data $BUILD_DIR/
    cp -r configs $BUILD_DIR/
    cp docs/README.md $BUILD_DIR/
    
    # Crear directorio de logs
    mkdir -p $BUILD_DIR/logs
    
    log_info "Paquete creado en $BUILD_DIR/"
}

# Mostrar ayuda
show_help() {
    echo "Uso: $0 [opciones]"
    echo ""
    echo "Opciones:"
    echo "  -h, --help          Mostrar esta ayuda"
    echo "  -c, --clean         Limpiar builds anteriores"
    echo "  -t, --test          Ejecutar tests antes de compilar"
    echo "  -x, --cross         Compilar para múltiples plataformas"
    echo "  -p, --package       Crear paquete de distribución"
    echo "  -a, --all           Ejecutar todas las opciones"
    echo ""
    echo "Ejemplos:"
    echo "  $0                  Compilación básica"
    echo "  $0 -t               Compilar después de ejecutar tests"
    echo "  $0 -x -p           Compilar multiplataforma y crear paquete"
    echo "  $0 -a              Ejecutar proceso completo"
}

# Procesar argumentos
CLEAN=false
TEST=false
CROSS=false
PACKAGE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -c|--clean)
            CLEAN=true
            shift
            ;;
        -t|--test)
            TEST=true
            shift
            ;;
        -x|--cross)
            CROSS=true
            shift
            ;;
        -p|--package)
            PACKAGE=true
            shift
            ;;
        -a|--all)
            CLEAN=true
            TEST=true
            CROSS=true
            PACKAGE=true
            shift
            ;;
        *)
            log_error "Opción desconocida: $1"
            show_help
            exit 1
            ;;
    esac
done

# Ejecutar proceso de compilación
main() {
    log_info "Iniciando proceso de compilación..."
    
    check_go
    
    if [ "$CLEAN" = true ]; then
        clean_build
    fi
    
    download_deps
    
    if [ "$TEST" = true ]; then
        run_tests
    fi
    
    build_app
    
    if [ "$CROSS" = true ]; then
        build_cross_platform
    fi
    
    if [ "$PACKAGE" = true ]; then
        create_package
    fi
    
    log_info "Proceso de compilación completado exitosamente"
    log_info "Ejecutable disponible en: $BUILD_DIR/$APP_NAME"
}

# Ejecutar función principal
main
