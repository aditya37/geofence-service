package geofence

import (
	"context"

	"github.com/aditya37/geofence-service/entity"
)

const queryGetTypes = `SELECT id,type_name FROM mst_geofence_type`

func (gm *geofenceManager) GetGeofenceTypes(ctx context.Context) ([]*entity.GeofenceType, error) {
	rows, err := gm.db.QueryContext(ctx, queryGetTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var item []*entity.GeofenceType
	for rows.Next() {
		var record entity.GeofenceType
		if err := rows.Scan(
			&record.Id,
			&record.TypeName,
		); err != nil {
			return nil, err
		}
		item = append(item, &entity.GeofenceType{
			Id:       record.Id,
			TypeName: record.TypeName,
		})
	}

	return item, nil
}
