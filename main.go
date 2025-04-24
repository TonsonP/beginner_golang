package main

import (
	"fmt"
	"beginner_golang/utils"
)

func main() {
	var user_select_mode int
	fmt.Println(`Welcome to my beginner golang applications, select which utility tools you want to use`)
	fmt.Println(utils.Node_mapping)

	for {
		fmt.Scanln(&user_select_mode)

		// _ will return value if that value associate with key.
		if _, exists := utils.Node_mapping[user_select_mode]; exists {
			var message string = fmt.Sprintf("You have select: %d %s", user_select_mode, utils.Node_mapping[user_select_mode])
			fmt.Println(message)
			break
		}

		fmt.Println("Please select utility tools number from the mapping list")

	}

	switch user_select_mode {
	
	// Proceed scenario Calculator
	case 1:
		utils.Calculator()

	}
	fmt.Println("Hello World")
}

