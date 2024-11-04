package exporters

import (
	"bytes"
	"encoding/csv"
	"gestione-ordini/pkg/database"
	"strconv"
)

func ExportToList(orders []database.Order) []byte {
	builder := bytes.Buffer{}

	for _, order := range orders {
		builder.WriteString(order.Product.Name)
		builder.WriteString(" - ")
		builder.WriteString(order.AmountString)
		builder.WriteString("\n")
	}
	builder.WriteString("\n")

	return builder.Bytes()
}

func ExportToCSV(orders []database.Order) []byte {
	builder := bytes.Buffer{}
	writer := csv.NewWriter(&builder)
	writer.Comma = ';'

	writer.Write([]string{"Prodotto", "Quantità", "Unità di misura", "Fornitore"})
	for _, order := range orders {
		writer.Write([]string{
			order.Product.Name,
			strconv.Itoa(order.Amount),
			order.Product.UnitOfMeasure.Symbol,
			order.Product.Supplier.Name,
		})
	}
	writer.Flush()

	return builder.Bytes()
}
