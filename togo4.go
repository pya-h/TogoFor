package main

import (
	"fmt"
	"time"
	"strings"
	"bufio"
	"log"
	"os"
)

type Date struct {time.Time}
	func (this Date) Get() string {
		return fmt.Sprintf("%d-%d-%d", this.Year(), this.Month(), this.Day())
	}


type Togo struct {
		id uint64
		title string
		description string
		weight uint16
		extra bool
		progress uint8
		done bool
		date Date
		duration time.Duration
	}
	func (this Togo) Info () string {
		return fmt.Sprintf("Togo #%d) %s:\t%s\nWeight: %d\nExtra: %t\nProgress: %d\nAt: %s, about %.1f minutes\nCompleted: %t",
		this.id, this.title, this.description, this.weight, this.extra, this.progress, this.date.Get(), this.duration.Minutes(), this.done)
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


func isCommand(term string)  bool {
	return term == "+" || term == "$"
}

func NewTogo(terms []string, nextID uint64) (togo Togo) {
	togo.title = terms[0]
	togo.id = nextID
	togo.done = false
	num_of_terms := len(terms)
	for i := 1; i < num_of_terms && !isCommand(terms[i]); i++ {
		switch terms[i] {

			case "=":
				i++
				togo.description = terms[i]
			case "+x":
				togo.extra = true
			case "-x":
				togo.extra = false
			case "+p":
				i++
				
				if _, err := fmt.Sscan(terms[i], &togo.progress); err != nil {
					panic(err)
				} else if togo.progress < 0 {
					togo.progress = 0
				} else if togo.progress > 100 {
					togo.progress = 100
				}
			case "@":
				// im++
				togo.date = Date{time.Now()}
				// get the actual date here
			case "->":
				i++
				if _, err := fmt.Sscan(terms[i], &togo.duration); err != nil {
					panic(err)
				} else if togo.duration > 0 {
					togo.duration *= time.Minute
				} else {
					panic("Duration must be positive intsger!")
				}
			case "+w":
				i++
				
				if _, err := fmt.Sscan(terms[i], &togo.weight); err != nil {
					panic(err)
				} else if togo.weight < 0 {
					togo.weight = 0
				}
			
		}
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
	togos := make(TogoList, 0, 0)
	for {
		fmt.Print("> ")

		if input, err := reader.ReadString('\n'); err == nil {
			input = input[:len(input) - 1] // remove '\n' char from the end of string
			terms := strings.Split(input, "\t")
			num_of_terms := len(terms)
			for i := 0; i < num_of_terms; i++{
				switch(terms[i]) {
					case "+":
						togo := NewTogo(terms[i+1:], togos.NextID())
						//fmt.Println(togo.Info())

						togos = togos.Add(&togo)
					case "$":
						togos.Show()
				}
			}
			
			// process 
		} else {
			panic(err)	
		}

	}

}
