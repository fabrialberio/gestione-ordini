package database

import (
	"database/sql"
)

type Utente struct {
	ID           int
	IDRuolo      int
	Username     string
	PasswordHash string
	Nome         string
	Cognome      string
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
	query := "SELECT id, id_ruolo, username, password_hash, nome, cognome FROM Utenti WHERE id = ?"
	var utente Utente

	err := db.conn.QueryRow(query, id).Scan(
		&utente.ID,
		&utente.IDRuolo,
		&utente.Username,
		&utente.PasswordHash,
		&utente.Nome,
		&utente.Cognome,
	)
	if err != nil {
		return nil, err
	}

	return &utente, nil
}

func (db *Database) GetUtenteByUsername(username string) (*Utente, error) {
	query := "SELECT id FROM Utenti WHERE username = ?"
	var id int

	err := db.conn.QueryRow(query, username).Scan(&id)
	if err != nil {
		return nil, err
	}

	return db.GetUtente(id)
}

func (db *Database) AddUtente(
	idRuolo int,
	username string,
	passwordHash string,
	nome string,
	cognome string,
) error {
	query := "INSERT INTO Utenti (id_ruolo, username, password_hash, nome, cognome) VALUES (?, ?, ?, ?, ?)"
	_, err := db.conn.Exec(query, idRuolo, username, passwordHash, nome, cognome)

	return err
}

func (db *Database) UtenteHasPermesso(idUtente int, idPermesso int) (bool, error) {
	query := "SELECT u.username FROM Utenti u JOIN Ruoli r ON u.id_ruolo = r.id JOIN RuoloPermesso rp ON r.id = rp.id_ruolo JOIN Permessi p ON rp.id_permesso = p.id WHERE u.id = ? AND p.id = ?"
	row := db.conn.QueryRow(query, idUtente, idPermesso)
	var username string

	err := row.Scan(&username)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
