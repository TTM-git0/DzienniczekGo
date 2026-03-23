package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Ocenki struct {
	db *sql.DB
}

type NewGrade struct {
	Grade float64 `json:"wartosc"`
}

func main() {
	connStr := "host=localhost port=5432 user=postgres password=Twoje_haslo_do_bazy dbname=szkola sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println("BŁĄD: Nie można połączyć z bazą! Sprawdź hasło.")
		panic(err)
	}

	app := Ocenki{db: db}
	r := gin.Default()

	r.GET("/srednia", app.pobierzSrednia)
	r.POST("/dodaj", app.dodajOcene)
	r.GET("/wyswietl", app.wyswietlDane)
	r.GET("/podsumowanie", app.podsumowanie)

	fmt.Println("Serwer ruszył na porcie :9000")
	r.Run(":9000")
}

func (a *Ocenki) podsumowanie(c *gin.Context) {
	var srednia float64
	errSrednia := a.db.QueryRow("SELECT COALESCE(AVG(wartosc), 0) FROM oceny").Scan(&srednia)
	if errSrednia != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd średniej"})
		return
	}

	var wszystkieoceny []float64
	rows, errLista := a.db.Query("SELECT wartosc FROM oceny")
	if errLista != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd wyswietlania danych"})
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var Jednaocena float64
		if err := rows.Scan(&Jednaocena); err != nil {
			continue
		}
		wszystkieoceny = append(wszystkieoceny, Jednaocena)
	}
	c.JSON(http.StatusOK, gin.H{
			"srednia": srednia,
			"oceny":   wszystkieoceny,
		})

}
func (a *Ocenki) wyswietlDane(c *gin.Context) {
	var wszystkieOceny []float64
	rows, err := a.db.Query("SELECT wartosc FROM oceny")
	if err != nil {
		c.JSON(500, gin.H{"error": "Błąd bazy"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var jednaOcena float64
		if err := rows.Scan(&jednaOcena); err != nil {
			continue
		}
		wszystkieOceny = append(wszystkieOceny, jednaOcena)
	}

	c.JSON(http.StatusOK, gin.H{
		"oceny": wszystkieOceny,
	})
}
func (a *Ocenki) pobierzSrednia(c *gin.Context) {
	var srednia float64
	err := a.db.QueryRow("SELECT COALESCE(AVG(wartosc), 0) FROM oceny").Scan(&srednia)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd bazy"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"srednia":   srednia,
		"wiadomosc": "Dane pobrane z PostgreSQL!",
	})
}

func (a *Ocenki) dodajOcene(c *gin.Context) {
	var dane NewGrade
	if err := c.ShouldBindJSON(&dane); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Zły JSON"})
		return
	}

	_, err := a.db.Exec("INSERT INTO oceny (wartosc) VALUES ($1)", dane.Grade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Nie udało się zapisać"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Ocena dodana do bazy!"})
}
