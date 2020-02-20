package main

import (  
	"fmt"
	"io/ioutil"
	"bufio"
	"os"
	"encoding/json"
)

type Pda struct {
	Name string
	States [] string
	Input_alphabet [] string
	Stack_alphabet [] string
	Accepting_states [] string
	Start_state string
	Transitions [][]string
}

func print(p Pda) {
	fmt.Println(p.Name)
	fmt.Println(p.States)
	fmt.Println(p.Input_alphabet)
	fmt.Println(p.Transitions)
}

func check(e error) {
    if e != nil {
        fmt.Print(e)
    }
}

func is_accepted(p Pda, cs string) {
	flag := 0
	accepting_states := p.Accepting_states
	for i:= 0; i<len(accepting_states); i++ {
		fmt.Println("as", accepting_states[i])
		if cs == accepting_states[i] {
			fmt.Println("Accepted")
			flag = 1
			break
		}
	}

	if flag == 0 {
		fmt.Println("Rejected")
	}
}


func transition(p Pda, s string) {
	inp_len := len(s)
	transitions := p.Transitions
	tran_len := len(transitions)

	current_state := p.Start_state
	currentStackSymbol := "null"
	stack := []string {"null"}

	for i := 0; i < inp_len; i++ {
		//fmt.Println(string(s[i]))
		char := string(s[i])
		for j := 0; j < tran_len; j++ {
			fmt.Println("Current state", current_state)
			t := transitions[j]
			if t[0] == current_state && t[1] == char && t[2] == currentStackSymbol {
				fmt.Println("found matching transition")
				current_state = t[3]

				if t[4] == "0" {
					stack = append(stack, "0")
					fmt.Println("Stack", stack)
				} else if t[4] == "null" {
					stack = stack[:len(stack) -1]
					fmt.Println("Stack", stack)
				}

				currentStackSymbol = stack[len(stack)-1]
				//fmt.Println("Current state", current_state)
				break
				
			}			
		}

		fmt.Println("Outside inside for")
		fmt.Println("Current state", current_state)
		fmt.Println("len", len(stack))

		if current_state == "q3" && len(stack) == 1 {
			current_state = "q4"
			break
		}
		
	}

	
	is_accepted(p, current_state)
	
}

func main(){

	dat, err := ioutil.ReadFile("./input.json")
	check(err)

	var p Pda
	err = json.Unmarshal(dat, &p)
	check(err)

	print(p)

	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Enter input string: ")
	// text, _ := reader.ReadString()

	

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	//fmt.Println(scanner.Text())
	transition(p, scanner.Text())
}

