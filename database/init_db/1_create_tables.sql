USE gestioneordini;

DROP TABLE IF EXISTS Ordini;
DROP TABLE IF EXISTS Prodotti;
DROP TABLE IF EXISTS UnitaDiMisura;
DROP TABLE IF EXISTS Fornitori;
DROP TABLE IF EXISTS TipologiaProdotto;
DROP TABLE IF EXISTS Utenti;
DROP TABLE IF EXISTS RuoloPermesso;
DROP TABLE IF EXISTS Permessi;
DROP TABLE IF EXISTS Ruoli;

CREATE TABLE Ruoli (
    id INT PRIMARY KEY,
    nome VARCHAR(255) UNIQUE
);

CREATE TABLE Permessi (
    id INT PRIMARY KEY,
    nome VARCHAR(255) UNIQUE
);

CREATE TABLE RuoloPermesso (
    id_ruolo INT,
    id_permesso INT,
    PRIMARY KEY (id_ruolo, id_permesso)
);
ALTER TABLE RuoloPermesso ADD CONSTRAINT ruolopermesso_id_ruolo_foreign FOREIGN KEY(id_ruolo) REFERENCES Ruoli(id);
ALTER TABLE RuoloPermesso ADD CONSTRAINT ruolopermesso_id_permesso_foreign FOREIGN KEY(id_permesso) REFERENCES Permessi(id);

CREATE TABLE Utenti (
    id INT PRIMARY KEY AUTO_INCREMENT,
    id_ruolo INT,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    nome VARCHAR(255),
    cognome VARCHAR(255)
);
ALTER TABLE Utenti ADD CONSTRAINT utente_id_ruolo_foreign FOREIGN KEY(id_ruolo) REFERENCES Ruoli(id);

CREATE TABLE TipologiaProdotto (
    id INT PRIMARY KEY AUTO_INCREMENT,
    nome VARCHAR(255) UNIQUE
);

CREATE TABLE Fornitori (
    id INT PRIMARY KEY AUTO_INCREMENT
);

CREATE TABLE UnitaDiMisura (
    id INT PRIMARY KEY AUTO_INCREMENT,
    simbolo VARCHAR(255) UNIQUE
);

CREATE TABLE Prodotti (
    id INT PRIMARY KEY AUTO_INCREMENT,
    id_tipologia INT,
    id_fornitore INT,
    id_unita_di_misura INT,
    nome VARCHAR(255)
);
ALTER TABLE Prodotti ADD CONSTRAINT prodotti_id_fornitore_foreign FOREIGN KEY(id_fornitore) REFERENCES Fornitori(id);
ALTER TABLE Prodotti ADD CONSTRAINT prodotti_id_tipologia_foreign FOREIGN KEY(id_tipologia) REFERENCES TipologiaProdotto(id);
ALTER TABLE Prodotti ADD CONSTRAINT prodotti_id_unita_di_misura_foreign FOREIGN KEY(id_unita_di_misura) REFERENCES UnitaDiMisura(id);

CREATE TABLE Ordini (
    id INT PRIMARY KEY AUTO_INCREMENT,
    id_prodotto INT,
    id_utente INT,
    quantita BIGINT,
    data_richiesta DATE
);
ALTER TABLE Ordini ADD CONSTRAINT ordini_id_prodotto_foreign FOREIGN KEY(id_prodotto) REFERENCES Prodotti(id);
ALTER TABLE Ordini ADD CONSTRAINT ordini_id_utente_foreign FOREIGN KEY(id_utente) REFERENCES Utenti(id);
