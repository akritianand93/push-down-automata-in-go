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

// Function to push data on to the stack when executing the PDA. It modifies the stack.
func push(p *PDAProcessor, val string) {
	p.Stack = append(p.Stack, val)
}

// Function to pop data from the stack when executing the PDA. It modifies the stack.
func pop(p *PDAProcessor) {
	p.Stack = p.Stack[:len(p.Stack) -1]
}

// Function to obtain the top n elements of the stack. This function does not modify the stack.
func peek(p PDAProcessor, k int) []string {
	top := [] string{}
	l := len(p.Stack)
	if (l <= k) {
		top = p.Stack
	} else if ( k == 1) {
		top = append(top, p.Stack[l-1])
	} else {
		top = p.Stack[l-k:l-1]
	}
	return top
}

// Function to reset the PDA and the stack. This deletes everything from the stack 
// and sets the current state to the start state so that we can start anew.
func reset(p *PDAProcessor) {
	p.Stack = make([]string, 0)
	p.Current_State = p.Pda.Start_state
}

// Function to open the grammar file. This function opens the json file, reads it and 
// unmarshal's its input into the PDA structure.
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

// Function to check if the input string has been accepted by the pda 
func is_accepted(proc PDAProcessor) {
	flag := 0
	accepting_states := proc.Pda.Accepting_states
	cs := proc.Current_State

	if len(proc.Stack) == 0 {
		for i:= 0; i < len(accepting_states); i++ {
			if cs == accepting_states[i] {
				flag = 1
				fmt.Println("Input token Accepted")
				done(proc, true)
				break
			}
		}
	}

	if flag == 0 {
		fmt.Println("Input token Rejected")
		done(proc, false)
	}
}

// The done returns the final status of the current state and the stack after the input string is processed.
func done(proc PDAProcessor, is_accepted bool){
	fmt.Println("pda_name: ", proc.Pda.Name)
	fmt.Println("is_method_accepted: ", is_accepted)
	fmt.Println("current_state: ", proc.Current_State)
	fmt.Println("stack_symbols", peek(proc, 5))
}


// This function accepts the input string and performs the necessary transitions and 
// stack operations for every token,
func put(proc PDAProcessor, p Pda, s string) int {
	inp_len := len(s)
	transitions := p.Transitions
	tran_len := len(transitions)
	transition_count := 0

	proc.Current_State = p.Start_state
	currentStackSymbol := "null"
	i := 0
	for ; i < inp_len; i++ {
		
		char := string(s[i])
		matching_transition := false
		for j := 0; j < tran_len; j++ {
			t := transitions[j]
			transition_count = check_for_dead_moves(t,&proc,transition_count) 
			if t[0] == proc.Current_State && t[1] == char && t[2] == currentStackSymbol {
				matching_transition = true
				proc.Current_State = t[3]

				if t[4] != "null" {
					push(&proc, t[4])
				} else {
					pop(&proc)
				}

				top := peek(proc, 1)[0]
				currentStackSymbol = top
				transition_count = transition_count + 1
				break
			}	       
		}

		if (!matching_transition) {
			break
		}
	}
	if i == inp_len {
		transition_count = eos(proc,transition_count)
	} else {
		is_accepted(proc)
	}

	return transition_count
}

// Performs the last transition to move the Automata to accepting state after the input
// string has been successfully parsed. 
func eos(proc PDAProcessor, transition_count int)int {
	length_of_stack := len(proc.Stack)
	allowed_transitions := proc.Pda.Transitions
	target_state := ""

	for j := 0; j < len(allowed_transitions); j++ {
		var allowed_current_state = allowed_transitions[j][0]
		if allowed_current_state == proc.Current_State {
			target_state = allowed_transitions[j][3]
		}
	}

	if peek(proc, 1)[0] == proc.Pda.Eos {
		proc.Current_State = target_state
		if length_of_stack > 0 {
			pop(&proc)
		}
		transition_count = transition_count + 1
	}
	is_accepted(proc)
	return transition_count
}

// Pushes initial EOS token into the stack and moves to the next state indicating
// the start of transitions
func check_for_dead_moves(transition []string, proc *PDAProcessor, transition_count int) int{
	allowed_current_state := transition[0]
	input := transition[1]
	allowed_top_of_stack := transition[2]
	target_state := transition[3]
	action_item := transition[4]
	if allowed_current_state == proc.Current_State && input == "null" && allowed_top_of_stack == "null"{
        proc.Current_State = target_state
        push(proc, action_item)
        transition_count = transition_count + 1
	} 
	return transition_count
}


// main function to start input processing. It reads the grammar json file
// to create an instance of PDA. The input is obtained either from a text file or
// standard input.
func main(){

	fn := os.Args[1]

	var p Pda
	open(fn, &p)

	proc := PDAProcessor {
		Pda: p,
	}
	reset(&proc)
	transition_count := 0
	input_token := ""

	if len(os.Args) < 3 {

		fmt.Print("Enter input string: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input_token = scanner.Text()
	} else {
		inp := os.Args[2]
		file, err := os.Open(inp)

		input := make([]byte, 32)
		n, err := file.Read(input)
		if err != nil {
			fmt.Print(err)
		}
		input_token = string(input[:n])
	}

	if input_token == "" {
		is_accepted(proc)
	} else{
		transition_count = put(proc, p, input_token)
	}

	fmt.Println("Number of transitions = ", transition_count)
}