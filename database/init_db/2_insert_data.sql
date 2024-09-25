USE gestioneordini;

INSERT INTO ruoli (id, nome) VALUES
    (0, "Cuoco"),
    (1, "Magazziniere"),
    (2, "Ammministratore");

INSERT INTO permessi (id, nome) VALUES
    (0, "vedi_prodotti"),
    (1, "vedi_tutti_ordini"),
    (2, "modifica_prodotti"),
    (3, "modifica_proprio_ordine"),
    (4, "modifica_tutti_ordini"),
    (5, "modifica_utenti");

INSERT INTO ruolo_permesso (id_ruolo, id_permesso) VALUES
    (0, 0),
    (0, 3),
    (1, 1),
    (2, 2),
    (2, 4),
    (2, 5);

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
