package database

type TipologiaProdotto struct {
	ID   int
	Nome string
}

type Fornitore struct {
	ID int
}

type UnitaDiMisura struct {
	ID      int
	Simbolo string
}

type Prodotto struct {
	ID              int
	IDTipologia     int
	IDFornitore     int
	IDUnitaDiMisura int
	Nome            string
}
