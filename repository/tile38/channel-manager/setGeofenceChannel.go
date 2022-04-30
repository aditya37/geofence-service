package channelmanager

import (
	"errors"

	tileEntity "github.com/aditya37/geofence-service/entity/tile38"
	"github.com/aditya37/geofence-service/repository"
	"github.com/xjem/t38c"
)

var (
	ErrorWrongDetectType = errors.New("Wrong detect type, please use WITHIN etc.")
)

func (cm *chanManager) SetGeofenceChannel(param tileEntity.Geofence) error {
	detect, err := cm.geofenceQuery(param)
	if err != nil {
		return err
	}
	if err := cm.tile.Channels.SetChan(param.Name, detect).Do(); err != nil {
		return err
	}
	return nil
}

// Mapping detect (Nearby,etc...)
func (cm *chanManager) geofenceQuery(param tileEntity.Geofence) (t38c.GeofenceQueryBuilder, error) {
	geofence := &t38c.Geofence{}
	query := t38c.GeofenceQueryBuilder{}
	if param.Action == repository.Within.ToString() {
		actions := cm.mappingGeofenceActions(param)
		query = geofence.Within(param.Key).
			Feature(param.Feature).
			Actions(actions...)
		return query, nil
	}
	return query, ErrorWrongDetectType
}

// mappingGeofenceActions
func (cm *chanManager) mappingGeofenceActions(param tileEntity.Geofence) []t38c.DetectAction {
	actions := []t38c.DetectAction{}
	for _, val := range param.Detect {
		actions = append(actions, t38c.DetectAction(val))
	}
	return actions
}
