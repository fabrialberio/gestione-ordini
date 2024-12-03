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
		ShowDatabaseQueryError(w, r, err)
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
			Label:        "Data di consegna richiesta",
			Name:         keyOrderRequestedAt,
			DefaultValue: time.Now().Format(dateFormat),
		},
	}

	appContext.ExecuteTemplate(w, r, "chef.html", data)
}

func PostOrderAmountInput(w http.ResponseWriter, r *http.Request) {
	amount := 1
	productId := 0

	formOrder, err := parseOrderForm(r)
	formOrderIsValid := err == nil && formOrder.ProductID != 0

	if formOrderIsValid {
		amount = formOrder.Amount
		productId = formOrder.ProductID
	} else if r.FormValue("isNew") != "true" {
		dbOrder, err := appContext.Database(r).FindOrder(formOrder.ID)
		if err != nil {
			logError(r, err)
			return
		}

		amount = dbOrder.Amount
		productId = dbOrder.ProductID
	}

	label := "Quantità"
	if productId != 0 {
		product, err := appContext.Database(r).FindProduct(productId)
		if err != nil {
			logError(r, err)
			return
		}

		label += " (" + product.UnitOfMeasure.Symbol + ")"
	}

	appContext.ExecuteTemplate(w, r, "amountInput", components.Input{
		Label:        label,
		Name:         keyOrderAmount,
		DefaultValue: strconv.Itoa(amount),
	})
}
