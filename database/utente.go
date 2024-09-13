package database

import (
	"database/sql"
)

type Utente struct {
	ID       int
	IDRuolo  int
	Username string
	Nome     string
	Cognome  string
}

const (
	IDRuoloCuoco int = iota
	IDRuoloMagazziniere
	IDRuoloAmministratore
)

const (
	IDPermessoVediProdotti int = iota
	IDPermessoVediTuttiOrdini
	IDPermessoModificaProdotti
	IDPermessoModificaProprioOrdine
	IDPermessoModificaTuttiOrdini
	IDPermessoModificaUtenti
)

func (db *Database) GetUtente(id int) (*Utente, error) {
	query := "SELECT id, id_ruolo, username, nome, cognome FROM Utenti WHERE id = ?"
	var utente Utente

	err := db.conn.QueryRow(query, id).Scan(
		&utente.ID,
		&utente.IDRuolo,
		&utente.Username,
		&utente.Nome,
		&utente.Cognome,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &utente, nil
}

func (db *Database) UtenteHasPermesso(id_utente int, id_permesso int) (bool, error) {
	query := "SELECT u.username FROM Utenti u JOIN Ruoli r ON u.id_ruolo = r.id JOIN RuoloPermesso rp ON r.id = rp.id_ruolo JOIN Permessi p ON rp.id_permesso = p.id WHERE u.id = ? AND p.id = ?"
	row := db.conn.QueryRow(query, id_utente, id_permesso)
	var username string

	err := row.Scan(&username)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
