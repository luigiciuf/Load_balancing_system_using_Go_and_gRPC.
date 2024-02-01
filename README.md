# Mini-Progetto SDCC

Questo progetto implementa un sistema di bilanciamento del carico per la verifica di numeri primi utilizzando  Go e gRPC.

## Struttura del Progetto

Il progetto è suddiviso tre diversi file `.go`:

- `main.go`: Contiene il codice del client che interagisce con il server di bilanciamento del carico.
- `loadbalancer.go`: Implementa il server di bilanciamento del carico che distribuisce le richieste tra le repliche.
- `replica.go`: Implementa le repliche che generano e verificano numeri primi.

## Configurazione

La configurazione del progetto è definita nel file `config.json` dove è possibile personalizzare gli indirizzi del server di bilanciamento del carico,delle repliche.
E' possibile aggiungere  o rimuovere indirizzi dalla lista `replica_addresses` per gestire dinamicamente il numero di repliche nel sistema.

## Avvio
Per avviare il progetto, esegui i seguenti comandi:

1. Avvia il Load Balancer utilizzando il comando `go run .\Load_Balancer\load_balancer.go`.
2. Avvia le Repliche utilizzando il comando `go run .\Load_Balancer\load_balancer.go`.
3. Avvia il Client utilizzando il comando `go run .\Client\client.go`.

Dopo aver eseguito questi comandi nell'ordine specificato, il Client sarà pronto per accettare un numero da verificare. Il risultato della verifica sarà restituito con il numero replica che ha gestito la richiesta. 

