package repositories

import (
	"context"
	"errors"

	"github.com/AnatolyKoltun/calculator-storage/database"
	"github.com/AnatolyKoltun/calculator-storage/models"
)

type CalculationRepository struct{}

func (r CalculationRepository) Save(ctx context.Context, calc *models.Calculation) error {
	query := `INSERT INTO calculations (argument1, argument2, operator, result, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := database.DB.QueryRow(ctx, query,
		calc.Argument1, calc.Argument2, calc.Operator,
		calc.Result, calc.CreatedAt).Scan(&calc.ID)

	return err
}

func (r CalculationRepository) GetList(ctx context.Context, filter models.FilterRequest) ([]models.Calculation, error) {
	query := `SELECT id, argument1, argument2, operator, result, created_at FROM calculations WHERE 1=1`
	args := []interface{}{}
	argCount := 1
	calculations := []models.Calculation{}

	if filter.DateFrom != "" {
		query += ` AND created_at >= $` + string(rune('0'+argCount))
		args = append(args, filter.DateFrom+" 00:00:00")
		argCount++
	}

	if filter.DateTo != "" {
		query += ` AND created_at <= $` + string(rune('0'+argCount))
		args = append(args, filter.DateTo+" 23:59:59")
		argCount++
	}

	query += ` ORDER BY created_at DESC`

	rows, err := database.DB.Query(ctx, query, args...)

	if err != nil {
		return calculations, errors.New("Ошибка запроса к БД: " + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var calc models.Calculation

		err := rows.Scan(&calc.ID, &calc.Argument1, &calc.Argument2,
			&calc.Operator, &calc.Result, &calc.CreatedAt)

		if err != nil {
			return calculations, errors.New("Ошибка чтения данных: " + err.Error())
		}

		calculations = append(calculations, calc)
	}

	if rows.Err() != nil {
		return calculations, errors.New("Ошибка при обработке результатов: " + rows.Err().Error())
	}

	return calculations, nil
}
