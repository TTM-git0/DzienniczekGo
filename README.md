# API Dzienniczka Ocen (Go & PostgreSQL)

Prosty backend REST API napisany w języku Go, służący do zarządzania ocenami uczniów.

# Technologie
- Język: Go (Golang)
- Framework WWW: Gin
- Baza Danych: PostgreSQL
- Sterownik bazy: github.com/lib/pq

# Jak uruchomić projekt?

1. Konfiguracja bazy danych:
   Upewnij się, że masz zainstalowany PostgreSQL. Utwórz bazę o nazwie "szkola" i wykonaj poniższy skrypt SQL, aby utworzyć wymaganą tabelę:

   ```sql
   CREATE TABLE oceny (
       id SERIAL PRIMARY KEY,
       wartosc FLOAT8 NOT NULL
   );
   ```

2. Uruchomienie serwera:
Pobierz wymagane zależności i uruchom kod za pomocą terminala:


go mod tidy
go run main.go

Serwer startuje na porcie 9000 (adres: http://localhost:9000)

# Endpointy API
GET /wyswietl - Pobiera listę wszystkich ocen

GET /srednia - Zwraca wyliczoną średnią ocen bezpośrednio z bazy danych

GET /podsumowanie - Zwraca listę wszystkich ocen oraz ich średnią w jednej odpowiedzi JSON

POST /dodaj - Dodaje nową ocenę do bazy. Wymaga podania w ciele żądania obiektu JSON, np.: {"wartosc": 4.5}