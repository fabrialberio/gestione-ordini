package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"net/http"
	"strconv"
	"time"
)

func GetChef(w http.ResponseWriter, r *http.Request) {
	productTypes, err := appContext.Database(r).FindAllProductTypes()
	if err != nil {
		LogError(r, err)
	}

	data := struct {
		ProductInput   components.ProductInput
		AmountInputURL string
		AmountInput    components.Input
		ExpiresAtInput components.Input
	}{
		ProductInput: components.ProductInput{
			ProductSelectName: keyOrderProductID,
			ProductSearchURL:  DestProductSearch,
			SearchInputName:   keyProductSearchQuery,
			ProductTypesName:  keyProductSearchProductTypes,
			ProductTypes:      productTypes,
		},
		AmountInputURL: DestOrderAmountInput,
		ExpiresAtInput: components.Input{
			Label:        "Scadenza",
			Name:         keyOrderRequestedAt,
			Type:         "date",
			DefaultValue: time.Now().Format(dateFormat),
		},
	}

	appContext.ExecuteTemplate(w, r, "chef.html", data)
}

func PostOrderAmountInput(w http.ResponseWriter, r *http.Request) {
	label := "Quantità"

	amount, err := strconv.Atoi(r.FormValue(keyOrderAmount))
	if err != nil {
		id, err := strconv.Atoi(r.FormValue(keyOrderID))
		if err != nil {
			amount = 1
		} else {
			order, _ := appContext.Database(r).FindOrder(id)
			amount = order.Amount
		}
	}

	selectedProductId, err := strconv.Atoi(r.FormValue(keyOrderProductID))
	if err == nil {
		product, err := appContext.Database(r).FindProduct(selectedProductId)
		if err != nil {
			LogError(r, err)
		}

		label += " (" + product.UnitOfMeasure.Symbol + ")"
	}

	appContext.ExecuteTemplate(w, r, "input", components.Input{
		Label:        label,
		Name:         keyOrderAmount,
		Type:         "number",
		DefaultValue: strconv.Itoa(amount),
	})
}
