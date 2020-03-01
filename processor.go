package main

// Libraries needed:
// fmt: to print output / statements.
// ioutil: file operations like open file, read file and close file.
// bufio: read input from standard input using the Scanner class.
// os: used to read command line arguments.
// encoding/json: used to marshal / unmarshal json input as needed.


import (  
	"fmt"
	"io/ioutil"
	"bufio"
	"os"
	"encoding/json"
)

// A structure that defines the a Push Down Automata
// This class contains the attributes of a PDA
type Pda struct {
	Name string
	States [] string
	Input_alphabet [] string
	Stack_alphabet [] string
	Accepting_states [] string
	Start_state string
	Transitions [][]string
	Eos string
}

// This class simulates a PDA processor that is it runs the PDA for teh provided input.
type PDAProcessor struct{
	Stack [] string
	Pda Pda
	Current_State string
}

// Function to push data on to the stack when executing the PDA. This is a function of the PDA processor. It modifies the stack.
func push (p *PDAProcessor, val string) {
	p.Stack = append(p.Stack, val)
}

// Function to pop data from the stack when executing the PDA. This is a function of the PDA processor. It modifies the stack.
func pop (p *PDAProcessor) {
	p.Stack = p.Stack[:len(p.Stack) -1]
}

// Function to obtain the top n elements of the stack. This function does not modify the stack.
func peek (p *PDAProcessor, k int) []string {
	top := [] string{}
	l := len(p.Stack)
	if ( k == 1) {
		top = append(top, p.Stack[l-1])
	} else {
		top = p.Stack[l-k:l-1]
	}
	return top
}

// Function reset the PDA and the stack. This deletes everything from the stack so that we can start anew.
func reset (p *PDAProcessor) {
	p.Stack = []string {"null"}
	p.Current_State = p.Pda.Start_state
}

// Function to open the grammar file. This function opens the json file, reads it and unmarshal's its input into the PDA structure.
func open(fn string, p *Pda) bool {
	dat, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Print(err)
		return false
    }

	err = json.Unmarshal(dat, &p)
	if err != nil {
		fmt.Print(err)
		return false
	}
	
	return true
}

// Function to check if the Automata is in accepting after the input has been processed.
func is_accepted(p Pda, cs string) {
	flag := 0
	accepting_states := p.Accepting_states
	
	for i:= 0; i < len(accepting_states); i++ {
		if cs == accepting_states[i] {
			flag = 1
			fmt.Println("Accepted")
			break
		}
	}

	if flag == 0 {
		fmt.Println("Rejected")
	}
}

func put(proc PDAProcessor, p Pda, s string) int {
	inp_len := len(s)
	transitions := p.Transitions
	tran_len := len(transitions)
	transition_count := 0

	current_state := p.Start_state
	currentStackSymbol := "null"

	for i := 0; i < inp_len; i++ {
		
		char := string(s[i])
		matching_transition := false
		for j := 0; j < tran_len; j++ {
			t := transitions[j]
			if t[0] == current_state && t[1] == char && t[2] == currentStackSymbol {
				matching_transition = true
				current_state = t[3]

				if t[4] == "0" {
					push(&proc, "0")
				} else if t[4] == "null" {
					pop(&proc)
				}

				top := peek(&proc, 1)[0]
				currentStackSymbol = top
				transition_count = transition_count + 1
				break
			}	       
		}

		if (!matching_transition) {
			break
		}

		top := peek(&proc, 1)[0]
		if current_state == "q3" && top == "null" && i == inp_len-1 {
			current_state = "q4"
			transition_count = transition_count + 1
			break
		}
	}

	is_accepted(p, current_state)

	return transition_count
}

func main(){

	fn := os.Args[1]

	var p Pda
	open(fn, &p)

	proc := PDAProcessor {
		Pda: p,
	}
	reset(&proc)

	transition_count := 0

	if len(os.Args) < 3 {

		fmt.Print("Enter input string: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		transition_count = put(proc, p, scanner.Text())

	} else {
		inp := os.Args[2]
		file, err := os.Open(inp)

		input := make([]byte, 32)
		n, err := file.Read(input)
		if err != nil {
			fmt.Print(err)
		}

		transition_count = put(proc, p, string(input[:n]))
	}

	fmt.Println("Number of transitions = ", transition_count)
}