package utils 
import (
    "fmt"
)


func ProceedTask(message string) bool {
	var proceed string

	fmt.Print(message)
	_, err := fmt.Scanln(&proceed)

	if err != nil {
		return false
	}

	switch proceed {
	case "y", "Y", "yes", "Yes":
		return true
	default:
		return false
	}
}
