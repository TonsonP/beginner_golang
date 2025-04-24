package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var computation_mode_mapping = map[int]string{
	1: "Scientific",
	2: "Accounting",
}

type Operator struct {
	precedence    int
	associativity string // left or right
}

var valid_operator = []string{"+", "-", "*", "/", "**"}

func get_operator_precedence(user_select_mode int) map[string]Operator {
	operator := map[string]Operator{}

	if user_select_mode == 1 {
		operator["+"] = Operator{precedence: 2, associativity: "left"}
		operator["-"] = Operator{precedence: 2, associativity: "left"}
		operator["*"] = Operator{precedence: 3, associativity: "left"}
		operator["/"] = Operator{precedence: 3, associativity: "left"}
		operator["**"] = Operator{precedence: 4, associativity: "right"}
	} else if user_select_mode == 2 {
		operator["+"] = Operator{precedence: 2, associativity: "left"}
		operator["-"] = Operator{precedence: 2, associativity: "left"}
		operator["*"] = Operator{precedence: 3, associativity: "left"}
		operator["/"] = Operator{precedence: 3, associativity: "left"}
		operator["**"] = Operator{precedence: 4, associativity: "right"}
	} else {
		operator["+"] = Operator{precedence: 2, associativity: "left"}
		operator["-"] = Operator{precedence: 2, associativity: "left"}
		operator["*"] = Operator{precedence: 3, associativity: "left"}
		operator["/"] = Operator{precedence: 3, associativity: "left"}
		operator["**"] = Operator{precedence: 4, associativity: "right"}
	}
	return operator
}

func sub_clean_operators(s string) string {
	// Count the number of minus signs
	count_minus := strings.Count(s, "-")
	if count_minus%2 == 0 {
		return "+"
	}
	return "-"
}

func clean_operators(user_input string) string {

	// Clean + and 0 operators
	re := regexp.MustCompile(`[+-]{2,}`)

	var output string = re.ReplaceAllStringFunc(user_input, sub_clean_operators)

	// Opening Parentheses
	re = regexp.MustCompile(`\({2,}`)
	output = re.ReplaceAllString(output, "(")

	// Closing Parentheses
	re = regexp.MustCompile(`\({2,}`)
	output = re.ReplaceAllString(output, ")")

	return output
}

// Initialize:
//     output_queue ← empty list
//     operator_stack ← empty stack

// For each token in the input:
//     If token is a number:
//         Append token to output_queue

//     Else if token is an operator (e.g., +, -, *, /, **):
//         While there is an operator at the top of the operator_stack
//               and it is not a left parenthesis
//               and (it has greater precedence OR
//                    it has equal precedence and is left-associative):
//             Pop operator from operator_stack and append to output_queue
//         Push current operator token to operator_stack

//     Else if token is a left parenthesis "(":
//         Push it to operator_stack

//     Else if token is a right parenthesis ")":
//         While the operator at the top of the stack is not a left parenthesis:
//             Pop operator from stack and append to output_queue
//         Pop the left parenthesis from the stack and discard it

// After all tokens are read:
//     While there are still operators on the stack:
//         If the operator is a parenthesis:
//             Error: Mismatched parenthesis
//         Pop operator and append to output_queue

// Return output_queue

func shunting_yard(input_string string, user_select_mode int) float32 {

	// Initialize
	operator := get_operator_precedence(user_select_mode) // Get operator along with its precedence and associativity

	var output_queue = []interface{}{} // Create queue for output
	var output_stack = []string{}      // Create stack for operators

	// Tokenize all input_string
	re := regexp.MustCompile(`\d+(\.\d+)?|\+|\-|\*{1,2}|/|\(|\)`)
	tokenize_input_string := re.FindAllString(input_string, -1)

	// Loop through each token
	for _, token := range tokenize_input_string {
		ctoken, err := strconv.ParseFloat(token, 16) // Try to convert to float

		// Check whether current token is number or not
		if err == nil {
			// if token is number -> add to queue
			output_queue = append(output_queue, ctoken)
		} else {
			// In case token is not number

			// If output_stack is empty just append it and continue
			if len(output_stack) == 0 {
				// Append
				output_stack = append(output_stack, token)
				continue
			}

			// Peek at top stack
			// Check whether it is operators
			found := false
			for _, op := range valid_operator {
				if output_stack[len(output_stack)-1] == op {
					found = true
					break
				}
			}

			// If top stack is an operators
			if found {
				// Initialize
				current_operator_precedence := operator[token].precedence
				current_operator_associativity := operator[token].associativity

				for {

					// Guard clauses
					if len(output_stack) == 0 {
						break
					}

					top_stack_precedence := operator[output_stack[len(output_stack)-1]].precedence
					if (top_stack_precedence > current_operator_precedence) ||
						((top_stack_precedence == current_operator_precedence) &&
							(current_operator_associativity == "left")) {
						// pop operator stack to output queue
						output_queue = append(output_queue, output_stack[len(output_stack)-1]) // Push last element in stack to queue
						output_stack = output_stack[0 : len(output_stack)-1]                   // Remove last element

					} else {
						break
					}
				}
				// push current operator to operator stack
				output_stack = append(output_stack, token)

			} else {
				// In case of parentheses
				if token == "(" {

					// If left parenthese push to stack
					output_stack = append(output_stack, token)

				} else {
					// In case of right parentheses
					// Pop operator until "(" if found
					for len(output_stack) > 0 {
						top := output_stack[len(output_stack)-1]
						output_stack = output_stack[:len(output_stack)-1] // pop

						if top == "(" {
							break // discard the "(" and stop popping
						}
						output_queue = append(output_queue, top)
					}
				}
			}

		}
	}

	// After we ran out of token pop any remaining operators from operators stack to output queue.
	// for i := len(output_stack) - 1; i >= 0; i-- {
	// 	current_output := output_stack[i]
	// 	output_queue = append(output_queue, current_output)
	// }

	// if token is operator -> check the top of stack
	// While the operator at the top of the stack has:
	// - greater precedence, OR
	// - equal precedence AND the **current** operator is left-associative
	// Do:
	// 	Pop from the operator stack to the output queue
	// push current operator to operator stack

	// if current token is left parentheses, push to opeatorstack
	// if right prantheses, pop operator from operatorstack to output queue until "(" is found
	// pop/trash "("

	// after all tokens, pop any remaining operators from operatorstack to outputqueue

	fmt.Println(output_queue...)

	return 3.14

}

// func calculate(clean_user_input string, user_select_mode int) string {

// 	// Calculate the one in parentheses first.
// 	parentheses_re := regexp.MustCompile(`\((.*?)\)`)

// 	parentheses_re_matchs := parentheses_re.FindAllStringSubmatch(clean_user_input, -1)

// 	for _, match := range parentheses_re_matchs {
// 		fmt.Println(match)
// 	}

// 	return ""

// }

func validate_input(user_input string) string {
	// Define token patter, matches any: 5, 3.14, +, -, *, **, /
	token_pattern := regexp.MustCompile(`\d+(\.\d+)?|\+|\-|\*{1,2}|/|\(|\)`)
	user_input = strings.ReplaceAll(user_input, " ", "")
	tokens := token_pattern.FindAllString(user_input, -1)
	validate_input := strings.Join(tokens, "")
	return validate_input
}

func Calculator() {
	fmt.Println("Welcome to calculator application, please select computation mode")
	fmt.Println(computation_mode_mapping)
	var user_select_mode int

	for {
		fmt.Scanln(&user_select_mode)

		if val, exists := computation_mode_mapping[user_select_mode]; exists {
			message := fmt.Sprintf("You have selected: %d %s", user_select_mode, val)
			fmt.Println(message)
			break
		}

		fmt.Println("Please select valid computation mode from the mapping list")
	}

	for {
		var user_input string
		var reader *bufio.Reader = bufio.NewReader(os.Stdin)
		fmt.Print("Enter numerical value(s) with operator (s):")
		user_input, _ = reader.ReadString('\n')
		user_input = strings.TrimSpace(user_input)
		user_input = strings.ReplaceAll(user_input, " ", "")

		fmt.Println("User Input ", user_input)
		var validate_input string = validate_input(user_input)
		fmt.Println("Validate User Input ", validate_input)
		clean_user_input := clean_operators(user_input)
		fmt.Println("Clean User Input ", clean_user_input)

		// calculate_results := calculate(clean_user_input, user_select_mode)
		calculate_results := shunting_yard(clean_user_input, user_select_mode)
		fmt.Println("Calculate results:", calculate_results)
		// var prev_num, current_num, num1, num2, num3, num4 float64

	}
}
