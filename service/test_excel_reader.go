package service

import (
	"io"
	"strconv"

	"github.com/oybek/choguuket/model"
	"github.com/xuri/excelize/v2"
)

type TestExcelReader struct{}

func (r *TestExcelReader) Read(reader io.Reader) ([]model.Medicine, error) {
	file, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	medicines := make([]model.Medicine, 0, len(rows))

	for _, row := range rows {
		if len(rows) < 2 {
			continue
		}
		amount, err := strconv.Atoi(row[1])
		if err != nil {
			continue
		}

		medicine := model.Medicine{
			Name:   row[0],
			Amount: int64(amount),
		}
		if !medicine.IsValid() {
			continue
		}

		medicines = append(medicines, medicine)
	}

	return medicines, nil
}
