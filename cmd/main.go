package main

import (
    "proyecto-grafos-go/internal/domain"
    "proyecto-grafos-go/internal/repository"
    "proyecto-grafos-go/internal/service"
    "proyecto-grafos-go/internal/ui/cli"
)

func main() {
    // Inicialización
    grafo := domain.NuevoGrafo(false)
    repo := repository.NuevoRepositorio("data/")
    
    // Servicios
    grafoSvc := service.NuevoServicioGrafo(grafo, repo)
    cuevaSvc := service.ServicioNuevaCueva(grafo)
    validacionSvc := service.NuevoServicioValidacion(grafo)

    // Menú principal
    menu := cli.NuevoMainMenu(grafoSvc, cuevaSvc, validacionSvc)
    menu.Mostrar()
}