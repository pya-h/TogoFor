package Togo

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const DATABASE_NAME string = "./togos.db"

type Date struct{ time.Time }

func (d Date) Get() string {

	return fmt.Sprintf("%d-%d-%d\t%d:%d", d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute())
}

// Struct Togo start
type Togo struct {
	Id          uint64
	Title       string
	Description string
	Weight      uint16
	Progress    uint8
	Extra       bool
	Date        Date
	Duration    time.Duration
}

func (togo Togo) Info() string {
	return fmt.Sprintf("Togo #%d) %s:\t%s\nWeight: %d\nExtra: %t\nProgress: %d\nAt: %s, about %.1f minutes",
		togo.Id, togo.Title, togo.Description, togo.Weight, togo.Extra, togo.Progress, togo.Date.Get(), togo.Duration.Minutes())
}
func (togo Togo) Save() {
	const CREATE_TABLE_QUERY string = `CREATE TABLE IF NOT EXISTS togos (id INTEGER NOT NULL PRIMARY KEY,
			title TEXT NOT NULL, description TEXT, weight INTEGER, extra INTEGER,
			progress INTEGER, date DATETIME, duration INTEGER)`

	db, err := sql.Open("sqlite3", DATABASE_NAME)

	if err != nil {
		panic(err)
	}
	defer db.Close()
	if _, err := db.Exec(CREATE_TABLE_QUERY); err != nil {
		panic(err)
	}
	extra := 0
	if togo.Extra {
		extra = 1
	}
	if _, err := db.Exec("INSERT INTO togos VALUES (?,?,?,?,?,?,?,?)", togo.Id,
		togo.Title, togo.Description, togo.Weight, extra, togo.Progress,
		togo.Date.Time, togo.Duration.Minutes()); err != nil {
		panic(err)
	}
}

// Struct Togo end

// TogoList start
type TogoList []Togo

func (these TogoList) Show() {
	fmt.Println("------------------------------------------------")
	for _, el := range these {
		fmt.Println(el.Info())
		fmt.Println("-----------------------------------------------")
	}
}
func (these TogoList) Add(new_togo *Togo) TogoList {
	return append(these, *new_togo)
}
func (togos TogoList) NextID() (id uint64) {
	id = uint64(len(togos)) // temporary
	return
}
func (togos TogoList) ProgressMade() (progress float64, completedInPercent float64, completed uint64, extra uint64, total uint64) {
	totalInPercent := uint64(0)
	for _, togo := range togos {
		progress += float64(togo.Progress) * float64(togo.Weight)
		if togo.Progress == 100 {
			completed++
			completedInPercent += float64(togo.Progress) * float64(togo.Weight)
		}
		if !togo.Extra {
			totalInPercent += uint64(100 * togo.Weight)
			total++
		} else {
			extra++
		}
	}
	progress *= 100 / float64(totalInPercent) // CHECK IF IT CALCULAFES DECIMAL PART OR NOT
	completedInPercent *= 100 / float64(totalInPercent)
	return
}

// TogoList end

func Load() (togos TogoList, err error) {

	togos = make(TogoList, 0)
	err = nil
	if db, e := sql.Open("sqlite3", DATABASE_NAME); e == nil {
		defer db.Close()
		const SELECT_QUERY string = "SELECT * FROM togos"
		rows, e := db.Query(SELECT_QUERY)
		if e != nil {
			err = e
			return
		}
		for rows.Next() {
			var togo Togo
			var date time.Time

			err = rows.Scan(&togo.Id, &togo.Title, &togo.Description, &togo.Weight, &togo.Extra, &togo.Progress, &date, &togo.Duration)
			togo.Date = Date{date}
			togo.Duration *= time.Minute
			if err != nil {
				panic(err)
			}
			togos = togos.Add(&togo)
		}

	} else {
		err = e
	}
	return
}

func isCommand(term string) bool {
	return term == "+" || term == "%" || term == "#"
}

func NewTogo(terms []string, nextID uint64) (togo Togo) {
	// setting default values
	if togo.Title = terms[0]; togo.Title == "" {
		togo.Title = "Untitled"
	}
	togo.Id = nextID
	togo.Weight = 1
	togo.Date = Date{time.Now()}

	num_of_terms := len(terms)
	for i := 1; i < num_of_terms && !isCommand(terms[i]); i++ {
		switch terms[i] {
		case "=", "+w":
			i++

			if _, err := fmt.Sscan(terms[i], &togo.Weight); err != nil {
				panic(err)
			}

		case ":", "+d":
			i++
			togo.Description = terms[i]
		case "+x":
			togo.Extra = true
		case "-x":
			togo.Extra = false
		case "+p":
			i++

			if _, err := fmt.Sscan(terms[i], &togo.Progress); err != nil {
				panic(err)
			} else if togo.Progress > 100 {
				togo.Progress = 100
			}
		case "@":
			// im++
			i++
			today := time.Now()
			var delta int
			if _, err := fmt.Sscan(terms[i], &delta); err != nil {
				panic(err)
			}
			today = today.AddDate(0, 0, delta)
			i++
			temp := strings.Split(terms[i], ":")
			var hour, min int
			if _, err := fmt.Sscan(temp[0], &hour); err != nil {
				panic(err)
			} else if hour >= 24 || hour < 0 {
				panic("Hour part must be between 0 and 23!")
			}
			if _, err := fmt.Sscan(temp[1], &min); err != nil {
				panic(err)
			} else if min >= 60 || min < 0 {
				panic("Minute part must be between 0 and 59!")
			}
			togo.Date = Date{time.Date(today.Year(), today.Month(), today.Day(), hour, min, 0, 0, time.Local)}
			// get the actual date here
		case "->":
			i++
			if _, err := fmt.Sscan(terms[i], &togo.Duration); err != nil {
				panic(err)
			} else if togo.Duration > 0 {
				togo.Duration *= time.Minute
			} else {
				panic("Duration must be positive integer!")
			}
		}

	}
	return
}
