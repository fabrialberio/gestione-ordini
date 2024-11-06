package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/components"
	"net/http"
	"time"
)

func GetChef(w http.ResponseWriter, r *http.Request) {
	productTypes, err := appContext.Database(r).FindAllProductTypes()
	if err != nil {
		ShowError(w, r, err)
		return
	}

	data := struct {
		ProductAmountInput components.ProductAmountInput
		ExpiresAtInput     components.Input
	}{
		ProductAmountInput: components.ProductAmountInput{
			ProductSelectName: keyOrderProductID,
			SearchDialog: components.ProductSearchDialog{
				ProductSearchURL: DestProductSearch,
				SearchInputName:  keyProductSearchQuery,
				ProductTypesName: keyProductSearchProductTypes,
				ProductTypes:     productTypes,
			},
			SelectedProduct:         0,
			AmountInputName:         keyOrderAmount,
			AmountInputDefaultValue: 1,
		},
		ExpiresAtInput: components.Input{
			Label:        "Scadenza",
			Name:         keyOrderRequestedAt,
			Type:         "date",
			DefaultValue: time.Now().Format(dateFormat),
		},
	}

	appContext.ExecuteTemplate(w, r, "chef.html", data)
}
