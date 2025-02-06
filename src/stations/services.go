package station

import (
	"context"
	"fmt"

	"github.com/V-enekoder/trenes/config"
)

// Servicio para crear una estación
func CreateStationService(ctx context.Context, estacion Station) error {
	err := createStationRepository(ctx, estacion)
	if err != nil {
		return fmt.Errorf("error en el servicio de creación de estación: %w", err)
	}
	return nil
}

// Servicio para obtener una estación por ID
func GetStationByIDService(ctx context.Context, id int64) (*Station, error) {
	estacion, err := getStationByIdRepository(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error en el servicio de obtención de estación por ID: %w", err)
	}
	return estacion, nil
}

// Servicio para obtener todas las estaciones
func GetAllStationsService(ctx context.Context) ([]*Station, error) {
	estaciones, err := getAllStationsRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("error en el servicio de obtención de todas las estaciones: %w", err)
	}
	return estaciones, nil
}

// Servicio para actualizar una estación
func UpdateStationService(ctx context.Context, estacion Station) error {
	err := UpdateStationRepository(ctx, estacion)
	if err != nil {
		return fmt.Errorf("error en el servicio de actualización de estación: %w", err)
	}
	return nil
}

// Servicio para eliminar una estación por ID
func DeleteStationService(ctx context.Context, id int64) error {
	err := DeleteStationRepository(ctx, id)
	if err != nil {
		return fmt.Errorf("error en el servicio de eliminación de estación: %w", err)
	}
	return nil
}

func FindOptimalRoadService(ctx context.Context, startID, endID int64) ([]interface{}, float64, error) {
	session, err := config.GetDatabaseConnection(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("error obteniendo la conexión: %w", err)
	}
	defer session.Close(ctx)

	result, err := session.Run(ctx, `
		MATCH (start:Estacion {Id: $startID})
		MATCH (end:Estacion {Id: $endID})
		CALL apoc.algo.dijkstra(start, end, 'CONNECTS_TO', 'distance')
		YIELD path, weight
		RETURN path, weight
	`, map[string]interface{}{"startID": startID, "endID": endID})
	if err != nil {
		return nil, 0, fmt.Errorf("error ejecutando Dijkstra: %w", err)
	}

	if result.Next(ctx) {
		record := result.Record()
		path := record.Values[0].([]interface{}) // path es un slice de interfaces
		weight := record.Values[1].(float64)
		return path, weight, nil
	}

	return nil, 0, fmt.Errorf("no se encontró ruta entre las estaciones %d y %d", startID, endID)
}
