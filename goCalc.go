package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"regexp"
	"strconv"
	"strings"
)
import fyne "fyne.io/fyne/v2"
import "fyne.io/fyne/v2/layout"

type Button struct {
	Handler func(label string, entryText string) string
	Label   string
	Color   string
}

func main() {
	a := app.New()
	w := a.NewWindow("Calculator by Yesayi")
	w.Resize(fyne.NewSize(float32(300), float32(30)))

	buttonMatrix := [][]Button{
		{
			{Handler: func(label string, entryText string) string {
				return ""
			}, Label: "C", Color: "red"},
			{Label: "(", Color: "gray"},
			{Label: ")", Color: "gray"},
			{Handler: actionHandler, Label: "/", Color: "gray"},
		},
		{
			{Label: "7", Color: "gray"},
			{Label: "8", Color: "gray"},
			{Label: "9", Color: "gray"},
			{Handler: actionHandler, Label: "*", Color: "gray"},
		},
		{
			{Label: "4", Color: "gray"},
			{Label: "5", Color: "gray"},
			{Label: "6", Color: "gray"},
			{Handler: actionHandler, Label: "+", Color: "gray"},
		},
		{
			{Label: "1", Color: "gray"},
			{Label: "2", Color: "gray"},
			{Label: "3", Color: "gray"},
			{Handler: actionHandler, Label: "-", Color: "gray"},
		},
		{
			{Label: "0", Color: "gray"},
			{Handler: func(_ string, entryText string) string {
				return findAnswer(entryText)
			}, Label: "=", Color: "lightblue"},
		},
	}
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Calculate...")
	entry.OnChanged = func(text string) {
		// Validate and filter the input
		filteredText := filterInput(text)

		// Set the filtered text back to the entry
		entry.SetText(filteredText)
	}

	mainContainer := container.New(layout.NewGridLayout(1))

	for _, row := range buttonMatrix {
		buttonsContainer := container.New(layout.NewGridLayout(len(row)))
		for _, label := range row {
			label := label // Capture the variable for the closure
			button := widget.NewButton(label.Label, func() {
				if label.Handler == nil {
					entry.SetText(defaultHandler(label.Label, entry.Text))
				} else {
					entry.SetText(label.Handler(label.Label, entry.Text))
				}
				entry.Refresh()
			})
			buttonsContainer.Add(button)
		}
		mainContainer.Add(buttonsContainer)
	}

	// Create a container for the entry and buttons
	content := container.NewVBox(
		entry,
		mainContainer,
	)

	// Set the content of the window
	w.SetContent(content)

	w.ShowAndRun()
}

func defaultHandler(label string, entryValue string) string {
	if entryValue == "0" {
		return label
	}
	return entryValue + label
}

func actionHandler(label string, entryText string) string {
	if len(entryText) == 0 {
		return ""
	}

	lastSymbol := entryText[len(entryText)-1]

	actionSymbols := []string{"+", "-", "/", "*"}
	for _, symbol := range actionSymbols {
		if symbol == string(lastSymbol) {
			return entryText
		}
	}

	return defaultHandler(label, entryText)
}

func findAnswer(expression string) string {
	for strings.Contains(expression, "(") {
		subExpression, deepStart, deepEnd, err := findDeepestScopeContent(expression)
		if err != nil {
			fmt.Println(err)
		}
		subExpressionAnswer := findAnswer(subExpression)

		expression = expression[0:deepStart] + subExpressionAnswer + expression[deepEnd:(len(expression)-1)]
	}
	possibleActions := []string{"*", "/", "+", "-"}
	for len(possibleActions) > 0 {
		action := possibleActions[0]
		pattern := `[0-9]+|[+\-*/]`
		re := regexp.MustCompile(pattern)
		operands := re.FindAllString(expression, -1)
		if len(operands) >= 0 && len(operands) < 3 {
			break
		}
		for key, item := range operands {
			if item == action {
				subResult := doOperation(floatVal(operands[key-1]), item, floatVal(operands[key+1]))
				expression = strings.Join(operands[0:(key-1)], "") + strconv.FormatFloat(subResult, 'f', -1, 64) + strings.Join(operands[(key+2):], "")
				break
			}
		}
		if !strings.Contains(expression, action) {
			possibleActions = possibleActions[1:]
		}
	}

	return expression
}

func findDeepestScopeContent(expression string) (string, int, int, error) {
	// Find the deepest scope by identifying the innermost parentheses
	deepestStart := strings.LastIndex(expression, "(")
	if deepestStart == -1 {
		return "", -1, -1, nil
	}

	deepestEnd := strings.Index(expression[deepestStart:], ")")
	if deepestEnd == -1 {
		return "", -1, -1, fmt.Errorf("no matching ')' found")
	}

	deepestEnd += deepestStart

	// Extract the content of the deepest scope
	deepestScopeContent := expression[deepestStart+1 : deepestEnd]
	return deepestScopeContent, deepestStart, deepestEnd, nil
}

func doOperation(left float64, operator string, right float64) float64 {
	switch operator {
	case "*":
		return left * right
	case "/":
		return left / right
	case "+":
		return left + right
	case "-":
		return left - right
	default:
		return float64(0)
	}
}

func floatVal(numeric string) float64 {
	floatValue, err := strconv.ParseFloat(numeric, 32)

	if err != nil {
		fmt.Println(err)
	}

	return float64(floatValue)
}

func filterInput(input string) string {
	// Define a regular expression that matches numeric values
	numericRegex := regexp.MustCompile("[0-9]+|[+\\-*/]|[()]")

	// Find all matches in the input string
	matches := numericRegex.FindAllString(input, -1)

	// Concatenate the matches to get the valid numeric text
	validText := ""
	for _, match := range matches {
		validText += match
	}

	return validText
}
