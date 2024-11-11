package files

import (
	"encoding/csv"
	"errors"
	"fmt"
	"gestione-ordini/pkg/database"
	"io"
	"slices"
	"strconv"
)

var ErrUnexpectedHeader = errors.New("unexpected CSV header")

var productExpectedHeader = []string{"id", "id_tipologia", "id_fornitore", "id_unita_di_misura", "descrizione", "codice"}

func ImportProductsFromCSV(reader io.Reader) ([]database.Product, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = csvComma
	csvReader.LazyQuotes = true

	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	if !slices.Equal(header, productExpectedHeader) {
		return nil, fmt.Errorf("%w: expected %v found %v", ErrUnexpectedHeader, productExpectedHeader, header)
	}

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var products []database.Product
	for _, r := range records {
		id, err := strconv.Atoi(r[0])
		if err != nil {
			return nil, err
		}
		productTypeId, err := strconv.Atoi(r[1])
		if err != nil {
			return nil, err
		}
		supplierId, err := strconv.Atoi(r[2])
		if err != nil {
			return nil, err
		}
		unitOfMeasureId, err := strconv.Atoi(r[3])
		if err != nil {
			return nil, err
		}

		products = append(products, database.Product{
			ID:              id,
			ProductTypeID:   productTypeId,
			SupplierID:      supplierId,
			UnitOfMeasureID: unitOfMeasureId,
			Description:     r[4],
			Code:            r[5],
		})
	}

	return products, nil
}
