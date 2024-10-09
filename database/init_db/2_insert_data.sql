USE gestioneordini;

INSERT INTO ruoli (id, nome) VALUES
    (1, "Cuoco"),
    (2, "Magazziniere"),
    (3, "Amministratore");

INSERT INTO permessi (id, nome) VALUES
    (1, "vedi_prodotti"),
    (2, "vedi_tutti_ordini"),
    (3, "modifica_prodotti"),
    (4, "modifica_proprio_ordine"),
    (5, "modifica_tutti_ordini"),
    (6, "modifica_utenti");

INSERT INTO ruolo_permesso (id_ruolo, id_permesso) VALUES
    (1, 1),
    (1, 4),
    (2, 2),
    (2, 3),
    (3, 3),
    (3, 5),
    (3, 6);

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
