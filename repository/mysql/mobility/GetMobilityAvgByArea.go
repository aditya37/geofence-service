package mobility

import (
	"context"

	"github.com/aditya37/geofence-service/entity"
)

func (mr *mobilityRepo) GetMobilityAverageByArea(ctx context.Context, interval, areaid int64) ([]*entity.ResultGetAvgByArea, error) {
	arg := []interface{}{
		&interval,
		&areaid,
	}
	rows, err := mr.db.QueryContext(ctx, mysqlQueryGetAvgMobiltyByArea, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.ResultGetAvgByArea
	for rows.Next() {

		var record entity.ResultGetAvgByArea
		if err := rows.Scan(
			&record.Enter,
			&record.Exit,
			&record.Inside,
			&record.ModifiedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, &entity.ResultGetAvgByArea{
			Exit:       record.Exit,
			Enter:      record.Enter,
			Inside:     record.Inside,
			ModifiedAt: record.ModifiedAt,
		})
	}
	return result, nil
}
