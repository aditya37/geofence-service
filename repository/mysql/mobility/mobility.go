package mobility

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aditya37/geofence-service/entity"
	"github.com/aditya37/geofence-service/repository"
	"github.com/go-sql-driver/mysql"
)

type mobilityRepo struct {
	db *sql.DB
}

func NewMobilityManager(db *sql.DB) (repository.MobilityManager, error) {
	return &mobilityRepo{
		db: db,
	}, nil
}

// GetLastAggregateFieldValue...
func (mr *mobilityRepo) GetLastAggregateFieldValue(ctx context.Context, field string, geofence_id int64) (int, error) {
	arg := []interface{}{
		&geofence_id,
	}
	queryString := fmt.Sprintf(mysqlQueryGetCurrentMobilityCount, field)
	row := mr.db.QueryRowContext(ctx, queryString, arg...)
	var count int
	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("Field is empty")
		}
		if errDriver, ok := err.(*mysql.MySQLError); ok {
			if errDriver.Number == 1054 {
				return 0, errors.New("unknown column in field list")
			}
		}
		return 0, err
	}
	return count, nil
}

// UpdateAggregateFieldValue....
func (mr *mobilityRepo) UpdateAggregateFieldValue(ctx context.Context, field string, geofence_id int64, value int) error {
	arg := []interface{}{
		&value,
		&geofence_id,
	}
	strQuery := fmt.Sprintf(mysqlQueryUpdateCurrentMobilityCount, field)
	row, err := mr.db.ExecContext(ctx, strQuery, arg...)
	if err != nil {
		return err
	}
	if isAffacted, _ := row.RowsAffected(); isAffacted == 0 {
		return errors.New("Failed update aggregate count")
	}
	return nil
}

// InsertDefaultValueAggregateField
func (mr *mobilityRepo) InsertDefaultValueAggregateField(ctx context.Context, field string, geofence_id int64, value int) error {
	arg := []interface{}{
		&value,
		&geofence_id,
	}
	strQuery := fmt.Sprintf(mysqlQueryInsertDefaultFieldValue, field)
	row, err := mr.db.ExecContext(ctx, strQuery, arg...)
	if err != nil {
		return err
	}
	if isAffacted, _ := row.RowsAffected(); isAffacted == 0 {
		return errors.New("Failed Insert Default Value")
	}
	return nil
}

// GetDailyAverage...
func (mr *mobilityRepo) GetDailyAverage(ctx context.Context, interval int) (*entity.ResultGetDailyAvg, error) {
	arg := []interface{}{
		&interval,
	}
	row := mr.db.QueryRowContext(ctx, mysqlQueryGetDailyAvg, arg...)

	var record entity.ResultGetDailyAvg
	if err := row.Scan(
		&record.Enter,
		&record.Exit,
		&record.Inside,
	); err != nil {
		if err == sql.ErrNoRows {
			// if row nil,return zero data
			return &entity.ResultGetDailyAvg{}, nil
		}
		return nil, err
	}
	return &record, nil
}

// GetAllAreaDailyAverage...
func (mr *mobilityRepo) GetAllAreaDailyAverage(ctx context.Context, interval int) ([]*entity.ResultGetDailyAvg, error) {
	arg := []interface{}{
		&interval,
	}
	rows, err := mr.db.QueryContext(ctx, mysqlQueryGetAllAreaDailyAvg, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*entity.ResultGetDailyAvg
	for rows.Next() {
		var record entity.ResultGetDailyAvg
		if err := rows.Scan(
			&record.GeofencId,
			&record.Enter,
			&record.Exit,
			&record.Inside,
		); err != nil {
			return nil, err
		}
		result = append(result, &entity.ResultGetDailyAvg{
			GeofencId: record.GeofencId,
			Enter:     record.Enter,
			Exit:      record.Exit,
			Inside:    record.Inside,
		})
	}

	return result, nil
}

// Close...
func (mr *mobilityRepo) Close() error {
	return mr.db.Close()
}
