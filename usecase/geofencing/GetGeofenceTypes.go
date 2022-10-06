package geofencing

import (
	"context"
	"net/http"

	"github.com/aditya37/geofence-service/entity"
	"github.com/aditya37/geofence-service/util"
)

func (gu *GeofencingUsecase) GetGeofenceTypes(ctx context.Context) ([]entity.GeofenceType, error) {
	resp, err := gu.geofenceManager.GetGeofenceTypes(ctx)
	if err != nil {
		return []entity.GeofenceType{}, &util.ErrorMsg{
			HttpRespCode: http.StatusUnprocessableEntity,
			Description:  err.Error(),
		}
	}
	var item []entity.GeofenceType
	for _, v := range resp {
		item = append(item, entity.GeofenceType{
			Id:       v.Id,
			TypeName: v.TypeName,
		})
	}
	return item, nil
}
