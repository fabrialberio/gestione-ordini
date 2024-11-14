USE gestioneordini;

INSERT INTO ruoli (id, nome) VALUES
    (1, "Cuoco"),
    (2, "Magazziniere"),
    (3, "Amministratore");

INSERT INTO tipologie_prodotto (nome) VALUES
    ("Insaccati"),
    ("Bibite"),
    ("Carne"),
    ("Verdure"),
    ("Dolce"),
    ("Formaggi"),
    ("Frutta"),
    ("Pasta"),
    ("Pesce"),
    ("Spezie"),
    ("Farina"),
    ("Varie");

INSERT INTO unita_di_misura (simbolo) VALUES
    ("pz"),
    ("Kg"),
    ("Lt"),
    ("Scatole");
