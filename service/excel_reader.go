package service

import (
	"io"

	"github.com/oybek/choguuket/model"
)

type ExcelReader interface {
	Read(io.Reader) ([]model.Medicine, error)
}
