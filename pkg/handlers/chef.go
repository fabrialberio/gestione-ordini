package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"strconv"
	"time"
)

func GetChef(w http.ResponseWriter, r *http.Request) {
	productTypes, err := appContext.Database(r).FindAllProductTypes()
	if err != nil {
		ShowError(w, r, err)
		return
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
	amount, _ := strconv.Atoi(r.FormValue(keyOrderAmount))

	selectedProductId, err := strconv.Atoi(r.FormValue(keyOrderProductID))
	if err != nil {
		selectedProductId = 1
	}

	product, err := appContext.Database(r).FindProduct(selectedProductId)
	if err != nil {
		ShowError(w, r, err)
		return
	}

	appContext.ExecuteTemplate(w, r, "input", constructAmountInput(product, strconv.Itoa(amount)))
}

func constructAmountInput(product database.Product, defaultValue string) components.Input {
	return components.Input{
		Label:        "Quantità (" + product.UnitOfMeasure.Symbol + ")",
		Name:         keyOrderAmount,
		Type:         "number",
		DefaultValue: defaultValue,
	}
}
