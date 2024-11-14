## Gestione ordini
Un'interfaccia per gestire gli ordini e organizzarli in base ai fornitori.

### Installazione
1. [Installare Docker](https://docs.docker.com/engine/install/)
2. Scaricare i file o clonare il repository con `git clone https://github.com/fabrialberio/gestione-ordini`
3. Aprire un terminale nella cartella ed seguire `docker compose up --build`

### Funzionalità
Pagine del sito:
- pagina per i cuochi per aggiungere e modificare gli ordini
- console per magazzinieri e amministratori con:
  - pagina per aggiungere ordini
  - pagina per visualizzare modificare tutti gli ordini e scaricarli in diversi formati
  - pagina per visualizzare e modificare i prodotti
  - pagina per visualizzare e modificare i fornitori
- pagine esclusive per gli amministratori:
  - pagina per visualizzare e modificare gli utenti
  - pagina per importare prodotti e utenti da file CSV

Ottimizzazione:
- interfaccia responsiva utilizzabile da telefono
- interfaccia utilizzabile da tastiera
- sito compatibile e testato con Safari, Chrome e Firefox
- tabella dei prodotti con caricamento dinamico
- ricerca prodotti con caricamento in background

Sicurezza:
- reverse proxy
- firewall che permette l'accesso al server solo per https, riducendo la superficie d'attacco
- protocollo https con certificato TLS, per cifrare i dati durante il transito
- autenticazione con JWT e cifratura sicura delle password
- controllo dell'autenticazione per ogni richiesta
