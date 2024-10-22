USE gestioneordini;

INSERT INTO ruoli (id, nome) VALUES
    (1, "Cuoco"),
    (2, "Magazziniere"),
    (3, "Amministratore");

INSERT INTO permessi (id, nome) VALUES
    (1, "ordini"),
    (2, "proprio_ordine"),
    (3, "tipologie_prodotto"),
    (4, "fornitori"),
    (5, "unita_di_misura"),
    (6, "prodotti"),
    (7, "utenti"),
    (8, "proprio_utente");

INSERT INTO ruolo_permesso (id_ruolo, id_permesso) VALUES
    (1, 2),
    (1, 6),
    (1, 8),
    (2, 1),
    (2, 3),
    (2, 4),
    (2, 5),
    (2, 6),
    (2, 8),
    (3, 1),
    (3, 3),
    (3, 4),
    (3, 5),
    (3, 6),
    (3, 7);

INSERT INTO tipologie_prodotto (nome) VALUES
    ("Carne"),
    ("Pesce"),
    ("Frutta"),
    ("Verdura"),
    ("Cereali"),
    ("Latticini");

INSERT INTO unita_di_misura (simbolo) VALUES
    ("pz"),
    ("Kg"),
    ("Lt");
