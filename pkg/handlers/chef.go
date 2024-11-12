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
	amount, err := strconv.Atoi(r.FormValue(keyOrderAmount))
	if err != nil {
		id, err := strconv.Atoi(r.FormValue(keyOrderID))
		if err != nil {
			amount = 1
		} else {
			order, err := appContext.Database(r).FindOrder(id)
			if err != nil {
				logError(r, err)
			}

			amount = order.Amount
		}
	}

	baseLabel := "Quantità"
	label, err := composeLabel(r, baseLabel)
	if err != nil {
		label = baseLabel
	}

	appContext.ExecuteTemplate(w, r, "amountInput", components.Input{
		Label:        label,
		Name:         keyOrderAmount,
		DefaultValue: strconv.Itoa(amount),
	})
}

func composeLabel(r *http.Request, baseLabel string) (string, error) {
	var productId int

	parsedOrder, err := parseOrderForm(r)
	if err != nil {
		return "", err
	}

	if parsedOrder.ProductID != 0 {
		productId = parsedOrder.ProductID
	} else {
		if parsedOrder.ID != 0 {
			order, err := appContext.Database(r).FindOrder(parsedOrder.ID)
			if err != nil {
				return "", err
			}

			productId = order.ProductID
		}
	}

	product, err := appContext.Database(r).FindProduct(productId)
	if err != nil {
		return "", err
	}

	return baseLabel + " (" + product.UnitOfMeasure.Symbol + ")", nil
}
