package utils

import (
	"bufio"
	"fmt"
	"math"
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
		operator["+"] = Operator{precedence: 1, associativity: "left"}
		operator["-"] = Operator{precedence: 1, associativity: "left"}
		operator["*"] = Operator{precedence: 1, associativity: "left"}
		operator["/"] = Operator{precedence: 1, associativity: "left"}
		operator["**"] = Operator{precedence: 1, associativity: "right"}
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

func check_type(s string) string {
	var output_type string = ""

	if s == "(" {
		output_type = "lp"
	}

	if s == ")" {
		output_type = "rp"
	}

	if output_type == "" {
		for _, op := range valid_operator {
			if op == s {
				output_type = "op"
				break
			}
		}
	}

	return output_type

}

func clean_operators(user_input string) string {

	// Clean + and 0 operators
	re := regexp.MustCompile(`[+-]{2,}`)

	var output string = re.ReplaceAllStringFunc(user_input, sub_clean_operators)

	// Opening Parentheses
	re = regexp.MustCompile(`\({2,}`)
	output = re.ReplaceAllString(output, "(")

	// Closing Parentheses
	re = regexp.MustCompile(`\){2,}`)
	output = re.ReplaceAllString(output, ")")

	return output
}

func postfix_evaluator(tokenize_input []interface{}) float32 {
	var output_stack = []float32{}
	var output float32 = 0

	for _, token := range tokenize_input {
		// Check whether this is operator or number
		switch token.(type) {

		// This is operators
		case string:
			var output float32 = 0
			right_operand := output_stack[len(output_stack)-1]
			left_operand := output_stack[len(output_stack)-2]
			output_stack = output_stack[0 : len(output_stack)-2]

			// Apply logic
			switch token {
			case "+":
				output = left_operand + right_operand
			case "-":
				output = left_operand - right_operand
			case "*":
				output = left_operand * right_operand
			case "/":
				output = left_operand / right_operand
			case "**":
				output = float32(math.Pow(float64(left_operand), float64(right_operand)))
			}

			output_stack = append(output_stack, output)

		// This is number
		case int, int32, int64, float32, float64:
			if num, ok := token.(float32); ok {
				output_stack = append(output_stack, num)
			} else if num, ok := token.(float64); ok {
				output_stack = append(output_stack, float32(num)) // Convert to float32
			}
		}

	}

	output = output_stack[0]

	// Validate output
	for _, val := range output_stack {
		val := float64(val)
		if math.IsInf(val, -1) || math.IsInf(val, 1) || math.IsNaN(val) {
			output = float32(val)
			break
		}

	}

	return output
}

func shunting_yard(input_string string, user_select_mode int) []interface{} {

	// Initialize
	operator := get_operator_precedence(user_select_mode) // Get operator along with its precedence and associativity

	var output_queue = []interface{}{} // Create queue for output
	var output_stack = []string{}      // Create stack for operators

	// Tokenize all input_string
	re := regexp.MustCompile(`\d+(\.\d+)?|\+|\-|\*{1,2}|/|\(|\)`)
	tokenize_input_string := re.FindAllString(input_string, -1)

	// Loop through each token
	for _, token := range tokenize_input_string {
		ctoken, err := strconv.ParseFloat(token, 32) // Try to convert to float

		// If current token is number
		if err == nil {
			output_queue = append(output_queue, ctoken) // Append to queue
		} else {
			// If current token is not number
			// Check type
			current_token_type := check_type(token)

			if current_token_type == "op" {

				if len(output_stack) == 0 {
					output_stack = append(output_stack, token)
				} else {
					for len(output_stack) > 0 {
						top := output_stack[len(output_stack)-1]

						current_top_type := check_type(top)

						if current_top_type == "lp" {
							output_stack = append(output_stack, token)
							break // Exit loop
						}

						// Check precedence of current operators
						current_top_precedence := operator[top].precedence
						current_top_associativity := operator[top].associativity
						current_token_precedence := operator[token].precedence

						if current_top_precedence > current_token_precedence {
							output_queue = append(output_queue, top)
							output_stack = output_stack[0 : len(output_stack)-1]

						} else if (current_top_precedence == current_token_precedence) && (current_top_associativity == "left") {
							output_queue = append(output_queue, top)
							output_stack = output_stack[0 : len(output_stack)-1]

						} else {
							output_stack = append(output_stack, token)
							break // Exit loop
						}
					}
				}

			} else if current_token_type == "lp" {
				output_stack = append(output_stack, token)
			} else if current_token_type == "rp" {
				for len(output_stack) > 0 {
					top := output_stack[len(output_stack)-1]
					output_stack = output_stack[0 : len(output_stack)-1]
					current_top_type := check_type(top)

					if current_top_type == "lp" {
						// Exit loop
						break
					}
					output_queue = append(output_queue, top)

				}
			}

		}

	}

	// After all tokens are read
	for len(output_stack) > 0 {
		top := output_stack[len(output_stack)-1]
		output_stack = output_stack[0 : len(output_stack)-1]
		output_queue = append(output_queue, top)

	}

	return output_queue

}

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
		postfix_input := shunting_yard(clean_user_input, user_select_mode)

		// Postfix evaluator
		output := postfix_evaluator(postfix_input)
		fmt.Println("Calculate results:", output)

	}
}
