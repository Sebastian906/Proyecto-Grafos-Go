#!/bin/bash

# Script para ejecutar todas las pruebas del proyecto

echo "=== EJECUCIÓN DE PRUEBAS - SISTEMA DE GESTIÓN DE CUEVAS ==="
echo ""

# Configuración
COVERAGE_DIR="coverage"
COVERAGE_FILE="coverage.out"
COVERAGE_HTML="coverage.html"

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

log_test() {
    echo -e "${BLUE}[TEST]${NC} $1"
}

# Verificar Go
check_go() {
    if ! command -v go &> /dev/null; then
        log_error "Go no está instalado o no está en PATH"
        exit 1
    fi
    log_info "Go version: $(go version)"
}

# Limpiar resultados anteriores
clean_results() {
    log_info "Limpiando resultados anteriores..."
    rm -rf $COVERAGE_DIR
    mkdir -p $COVERAGE_DIR
    log_info "Limpieza completada"
}

# Ejecutar tests unitarios
run_unit_tests() {
    log_test "Ejecutando tests unitarios..."
    
    go test -v ./tests/unit/... -coverprofile=$COVERAGE_DIR/unit.out
    unit_result=$?
    
    if [ $unit_result -eq 0 ]; then
        log_info "Tests unitarios: PASARON"
    else
        log_error "Tests unitarios: FALLARON"
    fi
    
    return $unit_result
}

# Ejecutar tests de integración
run_integration_tests() {
    log_test "Ejecutando tests de integración..."
    
    go test -v ./tests/integration/... -coverprofile=$COVERAGE_DIR/integration.out
    integration_result=$?
    
    if [ $integration_result -eq 0 ]; then
        log_info "Tests de integración: PASARON"
    else
        log_error "Tests de integración: FALLARON"
    fi
    
    return $integration_result
}

# Ejecutar tests de paquetes
run_package_tests() {
    log_test "Ejecutando tests de paquetes..."
    
    go test -v ./pkg/... -coverprofile=$COVERAGE_DIR/pkg.out
    pkg_result=$?
    
    if [ $pkg_result -eq 0 ]; then
        log_info "Tests de paquetes: PASARON"
    else
        log_error "Tests de paquetes: FALLARON"
    fi
    
    return $pkg_result
}

# Ejecutar tests de servicios
run_service_tests() {
    log_test "Ejecutando tests de servicios..."
    
    go test -v ./internal/service/... -coverprofile=$COVERAGE_DIR/service.out
    service_result=$?
    
    if [ $service_result -eq 0 ]; then
        log_info "Tests de servicios: PASARON"
    else
        log_error "Tests de servicios: FALLARON"
    fi
    
    return $service_result
}

# Ejecutar benchmarks
run_benchmarks() {
    log_test "Ejecutando benchmarks..."
    
    go test -bench=. -benchmem ./tests/performance/... > $COVERAGE_DIR/benchmarks.txt
    bench_result=$?
    
    if [ $bench_result -eq 0 ]; then
        log_info "Benchmarks: COMPLETADOS"
        log_info "Resultados guardados en: $COVERAGE_DIR/benchmarks.txt"
    else
        log_error "Benchmarks: FALLARON"
    fi
    
    return $bench_result
}

# Ejecutar race condition tests
run_race_tests() {
    log_test "Ejecutando tests de race conditions..."
    
    go test -race ./... > $COVERAGE_DIR/race.txt 2>&1
    race_result=$?
    
    if [ $race_result -eq 0 ]; then
        log_info "Race tests: PASARON"
    else
        log_error "Race tests: DETECTARON PROBLEMAS"
        log_warn "Revisar: $COVERAGE_DIR/race.txt"
    fi
    
    return $race_result
}

# Generar reporte de cobertura
generate_coverage() {
    log_info "Generando reporte de cobertura..."
    
    # Combinar archivos de cobertura
    echo "mode: atomic" > $COVERAGE_DIR/$COVERAGE_FILE
    
    for coverage_file in $COVERAGE_DIR/*.out; do
        if [ -f "$coverage_file" ]; then
            tail -n +2 "$coverage_file" >> $COVERAGE_DIR/$COVERAGE_FILE
        fi
    done
    
    # Generar reporte HTML
    go tool cover -html=$COVERAGE_DIR/$COVERAGE_FILE -o $COVERAGE_DIR/$COVERAGE_HTML
    
    # Mostrar resumen de cobertura
    coverage_percent=$(go tool cover -func=$COVERAGE_DIR/$COVERAGE_FILE | grep total | awk '{print $3}')
    log_info "Cobertura total: $coverage_percent"
    log_info "Reporte HTML: $COVERAGE_DIR/$COVERAGE_HTML"
}

# Ejecutar tests específicos por requisito
run_requirement_tests() {
    log_test "Ejecutando tests por requisito..."
    
    echo "=== REQUISITO 1: Gestión de Grafos ===" > $COVERAGE_DIR/requirements.txt
    go test -v ./tests/unit/ -run "TestRequisito1" >> $COVERAGE_DIR/requirements.txt 2>&1
    
    echo "=== REQUISITO 2: Gestión de Conexiones ===" >> $COVERAGE_DIR/requirements.txt
    go test -v ./tests/unit/ -run "TestRequisito2" >> $COVERAGE_DIR/requirements.txt 2>&1
    
    echo "=== REQUISITO 3: MST ===" >> $COVERAGE_DIR/requirements.txt
    go test -v ./tests/unit/ -run "TestRequisito3" >> $COVERAGE_DIR/requirements.txt 2>&1
    
    log_info "Tests por requisito completados: $COVERAGE_DIR/requirements.txt"
}

# Validar calidad del código
validate_code_quality() {
    log_test "Validando calidad del código..."
    
    # Verificar formato
    log_info "Verificando formato del código..."
    gofmt -l . > $COVERAGE_DIR/format.txt
    if [ -s $COVERAGE_DIR/format.txt ]; then
        log_warn "Archivos con formato incorrecto encontrados"
        cat $COVERAGE_DIR/format.txt
    else
        log_info "Formato del código: OK"
    fi
    
    # Verificar imports
    if command -v goimports &> /dev/null; then
        log_info "Verificando imports..."
        goimports -l . > $COVERAGE_DIR/imports.txt
        if [ -s $COVERAGE_DIR/imports.txt ]; then
            log_warn "Imports incorrectos encontrados"
        else
            log_info "Imports: OK"
        fi
    fi
    
    # Verificar vet
    log_info "Ejecutando go vet..."
    go vet ./... > $COVERAGE_DIR/vet.txt 2>&1
    if [ $? -eq 0 ]; then
        log_info "Go vet: OK"
    else
        log_warn "Go vet encontró problemas"
    fi
}

# Mostrar resumen final
show_summary() {
    echo ""
    echo "=== RESUMEN DE PRUEBAS ==="
    echo "Resultados guardados en: $COVERAGE_DIR/"
    echo ""
    
    if [ -f "$COVERAGE_DIR/$COVERAGE_FILE" ]; then
        coverage_percent=$(go tool cover -func=$COVERAGE_DIR/$COVERAGE_FILE | grep total | awk '{print $3}')
        log_info "Cobertura total: $coverage_percent"
    fi
    
    echo "Archivos generados:"
    ls -la $COVERAGE_DIR/
    
    echo ""
    echo "Para ver el reporte de cobertura HTML:"
    echo "  open $COVERAGE_DIR/$COVERAGE_HTML"
    echo ""
}

# Mostrar ayuda
show_help() {
    echo "Uso: $0 [opciones]"
    echo ""
    echo "Opciones:"
    echo "  -h, --help          Mostrar esta ayuda"
    echo "  -u, --unit          Solo tests unitarios"
    echo "  -i, --integration   Solo tests de integración"
    echo "  -p, --package       Solo tests de paquetes"
    echo "  -s, --service       Solo tests de servicios"
    echo "  -b, --bench         Solo benchmarks"
    echo "  -r, --race          Solo race condition tests"
    echo "  -c, --coverage      Generar solo reporte de cobertura"
    echo "  -q, --quality       Validar solo calidad del código"
    echo "  -a, --all           Ejecutar todas las pruebas (default)"
    echo ""
    echo "Ejemplos:"
    echo "  $0                  Ejecutar todas las pruebas"
    echo "  $0 -u -c           Ejecutar tests unitarios y generar cobertura"
    echo "  $0 -b              Ejecutar solo benchmarks"
}

# Procesar argumentos
UNIT=false
INTEGRATION=false
PACKAGE=false
SERVICE=false
BENCH=false
RACE=false
COVERAGE=false
QUALITY=false
ALL=true

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -u|--unit)
            UNIT=true
            ALL=false
            shift
            ;;
        -i|--integration)
            INTEGRATION=true
            ALL=false
            shift
            ;;
        -p|--package)
            PACKAGE=true
            ALL=false
            shift
            ;;
        -s|--service)
            SERVICE=true
            ALL=false
            shift
            ;;
        -b|--bench)
            BENCH=true
            ALL=false
            shift
            ;;
        -r|--race)
            RACE=true
            ALL=false
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            ALL=false
            shift
            ;;
        -q|--quality)
            QUALITY=true
            ALL=false
            shift
            ;;
        -a|--all)
            ALL=true
            shift
            ;;
        *)
            log_error "Opción desconocida: $1"
            show_help
            exit 1
            ;;
    esac
done

# Ejecutar proceso de pruebas
main() {
    log_info "Iniciando proceso de pruebas..."
    
    check_go
    clean_results
    
    total_failures=0
    
    if [ "$ALL" = true ] || [ "$UNIT" = true ]; then
        run_unit_tests
        total_failures=$((total_failures + $?))
    fi
    
    if [ "$ALL" = true ] || [ "$INTEGRATION" = true ]; then
        run_integration_tests
        total_failures=$((total_failures + $?))
    fi
    
    if [ "$ALL" = true ] || [ "$PACKAGE" = true ]; then
        run_package_tests
        total_failures=$((total_failures + $?))
    fi
    
    if [ "$ALL" = true ] || [ "$SERVICE" = true ]; then
        run_service_tests
        total_failures=$((total_failures + $?))
    fi
    
    if [ "$ALL" = true ] || [ "$BENCH" = true ]; then
        run_benchmarks
    fi
    
    if [ "$ALL" = true ] || [ "$RACE" = true ]; then
        run_race_tests
    fi
    
    if [ "$ALL" = true ]; then
        run_requirement_tests
    fi
    
    if [ "$ALL" = true ] || [ "$COVERAGE" = true ]; then
        generate_coverage
    fi
    
    if [ "$ALL" = true ] || [ "$QUALITY" = true ]; then
        validate_code_quality
    fi
    
    show_summary
    
    if [ $total_failures -eq 0 ]; then
        log_info "Todas las pruebas completadas exitosamente"
        exit 0
    else
        log_error "Algunas pruebas fallaron"
        exit 1
    fi
}

# Ejecutar función principal
main
