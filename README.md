<div align="center">
  <img src="https://github.com/fabrialberio/gestione-ordini/blob/75d0803e0e40ad218dd403b64b22f1f49f7e3335/screenshot.png">
</div>

## Gestione ordini
Un'interfaccia per gestire gli ordini e organizzarli in base ai fornitori.

### Installazione
1. [Installa Docker](https://docs.docker.com/engine/install/)
2. [scarica i file](https://github.com/fabrialberio/gestione-ordini/archive/refs/heads/main.zip) oppure clona il repository con Git:
```
git clone https://github.com/fabrialberio/gestione-ordini
```
3. crea una chiave privata:
```
openssl genrsa -out private.key 4096
```
4. crea una chiave pubblica corrispondente:
```
openssl rsa -in private.key -out public.key -outform PEM -pubout
```
5. crea un file `.env` con le credenziali necessarie per l'avvio, ad esempio:
```env
MYSQL_USER: gestioneordini
MYSQL_PASSWORD: example
MYSQL_ROOT_PASSWORD: example
MYSQL_DATABASE: gestioneordini
MYSQL_CONTAINER_NAME: gestione-ordini-mysql
ADMIN_PASSWORD: example
```
6. avvia il server di sviluppo:
```
docker compose -f docker-compose.dev.yml up
```
7. puoi accedere al server su `http://localhost:80`, con il nome utente `amdin` e la password inserita nel file `.env`.

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
