package exporters

import (
	"bytes"
	"encoding/csv"
	"gestione-ordini/pkg/database"
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

	writer.Write([]string{"Prodotto", "Quantità", "Fornitore"})
	for _, order := range orders {
		writer.Write([]string{order.Product.Name, order.AmountString, order.Product.Supplier.Name})
	}
	writer.Flush()

	return builder.Bytes()
}
