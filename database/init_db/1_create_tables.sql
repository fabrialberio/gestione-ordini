USE gestioneordini;

DROP TABLE IF EXISTS ordini;
DROP TABLE IF EXISTS prodotti;
DROP TABLE IF EXISTS unita_di_misura;
DROP TABLE IF EXISTS fornitori;
DROP TABLE IF EXISTS tipologie_prodotto;
DROP TABLE IF EXISTS utenti;
DROP TABLE IF EXISTS ruolo_permesso;
DROP TABLE IF EXISTS permessi;
DROP TABLE IF EXISTS ruoli;

CREATE TABLE ruoli (
    id BIGINT PRIMARY KEY,
    nome VARCHAR(255) UNIQUE
);

CREATE TABLE permessi (
    id BIGINT PRIMARY KEY,
    nome VARCHAR(255) UNIQUE
);

CREATE TABLE ruolo_permesso (
    id_ruolo BIGINT,
    id_permesso BIGINT,
    PRIMARY KEY (id_ruolo, id_permesso)
);
ALTER TABLE ruolo_permesso ADD CONSTRAINT ruolo_permesso_id_ruolo_foreign FOREIGN KEY(id_ruolo) REFERENCES ruoli(id);
ALTER TABLE ruolo_permesso ADD CONSTRAINT ruolo_permesso_id_permesso_foreign FOREIGN KEY(id_permesso) REFERENCES permessi(id);

CREATE TABLE utenti (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    id_ruolo BIGINT,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    nome VARCHAR(255),
    cognome VARCHAR(255),
    creato_il DATETIME
);
ALTER TABLE utenti ADD CONSTRAINT utente_id_ruolo_foreign FOREIGN KEY(id_ruolo) REFERENCES ruoli(id);

CREATE TABLE tipologie_prodotto (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    nome VARCHAR(255) UNIQUE
);

CREATE TABLE fornitori (
    id BIGINT PRIMARY KEY AUTO_INCREMENT
);

CREATE TABLE unita_di_misura (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    simbolo VARCHAR(255) UNIQUE
);

CREATE TABLE prodotti (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    id_tipologia BIGINT,
    id_fornitore BIGINT,
    id_unita_di_misura BIGINT,
    nome VARCHAR(255)
);
ALTER TABLE prodotti ADD CONSTRAINT prodotti_id_fornitore_foreign FOREIGN KEY(id_fornitore) REFERENCES fornitori(id);
ALTER TABLE prodotti ADD CONSTRAINT prodotti_id_tipologia_foreign FOREIGN KEY(id_tipologia) REFERENCES tipologie_prodotto(id);
ALTER TABLE prodotti ADD CONSTRAINT prodotti_id_unita_di_misura_foreign FOREIGN KEY(id_unita_di_misura) REFERENCES unita_di_misura(id);

CREATE TABLE ordini (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    id_prodotto BIGINT,
    id_utente BIGINT,
    quantita BIGINT,
    richiesto_il DATETIME
);
ALTER TABLE ordini ADD CONSTRAINT ordini_id_prodotto_foreign FOREIGN KEY(id_prodotto) REFERENCES prodotti(id);
ALTER TABLE ordini ADD CONSTRAINT ordini_id_utente_foreign FOREIGN KEY(id_utente) REFERENCES utenti(id);
