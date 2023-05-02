package main

import (
	"fmt"

	"time"
	"strings"
	"bufio"
	"log"
	"os"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const DATABASE_NAME string = "./togos.db"

type Date struct {time.Time}
	func (d Date) Get() string {

		return fmt.Sprintf("%d-%d-%d", d.Year(), d.Month(), d.Day())
	}


type Togo struct {
		Id uint64
		Title string
		Description string
		Weight uint16
		Progress uint8
		Extra bool
		Date Date
		Duration time.Duration
	}
	func (togo Togo) Info () string {
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

// TogoList:
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
		progress *= 100 / float64(totalInPercent)  // CHECK IF IT CALCULAFES DECIMAL PART OR NOT
		completedInPercent *= 100 / float64(totalInPercent)
		return
	}

func isCommand(term string)  bool {
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
			case "=":
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
				togo.Date = Date{time.Now()}
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
			case "+w":
				i++
				
				if _, err := fmt.Sscan(terms[i], &togo.Weight); err != nil {
					panic(err)
				}
			
		}
	}
	return
}


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


func main() {
	// 2nd project to be done
	// while walking the streets
	// Be'sme BigBang =))))))
	defer func() {
		err := recover()
		if err != nil {
			log.Fatal("Something fucked up: ", err)
		}
	} ()

	reader := bufio.NewReader(os.Stdin)
	togos, err := Load()// make(TogoList, 0)
	if err != nil {
		fmt.Println("Loading failed: ", err)
	}
	for {
		fmt.Print("> ")

		if input, err := reader.ReadString('\n'); err == nil {
			input = input[:len(input) - 1] // remove '\n' char from the end of string
			terms := strings.Split(input, "\t")
			num_of_terms := len(terms)
			for i := 0; i < num_of_terms; i++{
				switch(terms[i]) {
				case "+":
					if num_of_terms > 1 {

						togo := NewTogo(terms[i+1:], togos.NextID())
						togos = togos.Add(&togo)
						
						togo.Save()
					} else {
						panic("You must provide some values!")
					}
				case "#":
					togos.Show()
				case "%":
					progress, completedInPercent, completed, extra, total := togos.ProgressMade()
					fmt.Printf("Progress: %3.2f%% (%3.2f%% Completed),\nStatistics: %d / %d",
						progress, completedInPercent, completed, total)
					if extra > 0 {
						fmt.Printf("[+%d]\n", extra)
					}
					fmt.Println()
				case "><":
					fmt.Println("Fuck U & Have a nice day.")
					return
				}
			}
			
			// process 
		} else {
			panic(err)	
		}

	}

}
