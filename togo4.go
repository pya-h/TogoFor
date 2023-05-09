package main

import (
	"Togo"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// 2nd project to be done
	// while walking the streets
	// Be'sme BigBang =))))))
	defer func() {
		err := recover()
		if err != nil {
			log.Fatal("Something fucked up: ", err)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	togos, err := Togo.Load(true) // load today's togos,  make(Togo.TogoList, 0)
	if err != nil {
		fmt.Println("Loading failed: ", err)
	}
	for {
		fmt.Print("> ")

		if input, err := reader.ReadString('\n'); err == nil {
			input = input[:len(input)-1] // remove '\n' char from the end of string
			terms := strings.Split(input, "\t")
			num_of_terms := len(terms)
			var now Togo.Date = Togo.Today()
			for i := 0; i < num_of_terms; i++ {
				switch terms[i] {
				case "+":
					if num_of_terms > 1 {

						togo := Togo.Extract(terms[i+1:], togos.NextID())
						if togo.Date.Short() == now.Short() {
							togos = togos.Add(&togo)
							if togo.Date.After(now.Time) {
								togo.Schedule()
							}
						}

						togo.Save()
					} else {
						panic("You must provide some values!")
					}
				case "#":
					if i+1 < num_of_terms && terms[i+1] == "-a" {
						all_togos, err := Togo.Load(false)
						if err != nil {
							panic(err)
						}
						all_togos.Show()
					} else {
						togos.Show()
					}
				case "%":
					var target *Togo.TogoList = &togos
					scope := "Today's"
					if i+1 < num_of_terms && terms[i+1] == "-a" {
						all_togos, err := Togo.Load(false)
						if err != nil {
							panic(err)
						}
						target = &all_togos
						scope = "Total"
					}
					progress, completedInPercent, completed, extra, total := (*target).ProgressMade()
					fmt.Printf("%s Progress: %3.2f%% (%3.2f%% Completed),\nStatistics: %d / %d",
						scope, progress, completedInPercent, completed, total)
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
