package main

import (
	"gestione-ordini/database"
	"net/http"
)

func cook(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if err := checkRole(r, database.RoleIDCook); err != nil {
		return nil, err
	}

	return nil, nil
}
