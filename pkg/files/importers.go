package files

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"gestione-ordini/pkg/database"
	"strconv"
)

var ErrUnexpectedHeader = errors.New("unexpected CSV header")

var productExpectedHeader = []string{"id", "id_tipologia", "id_fornitore", "id_unita_di_misura", "descrizione", "codice"}

func ImportProductsFromCSV(data []byte) ([]database.Product, error) {
	reader := csv.NewReader(bytes.NewReader(data))

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("%w: expected %v found %v", ErrUnexpectedHeader, productExpectedHeader, header)
	}

	records, err := reader.ReadAll()
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
