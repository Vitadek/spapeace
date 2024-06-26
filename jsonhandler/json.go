package jsonhandler

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"spock/jsonparser"
	"spock/term"
)

var jsonFilePath string
var jsonData *jsonparser.JSONData
var searchQuery string

func ParseJSON() error {
	for {
		term.DisableRawMode()      // Temporarily disable raw mode for getting the file path
		defer term.EnableRawMode() // Re-enable raw mode after getting the input

		fmt.Print("Enter JSON file path: ")
		fmt.Scanln(&jsonFilePath)

		var err error
		jsonData, err = jsonparser.ParseJSONFile(jsonFilePath)
		if err != nil {
			fmt.Println("File not found, please try again.")
			continue
		}

		break
	}

	term.EnableRawMode() // Ensure raw mode is enabled before displaying the results

	displayGroupsWithStatus(jsonData)
	return nil
}

func SaveJSON() error {
	return jsonparser.SaveJSONFile(jsonFilePath, jsonData)
}

func updateGroupStatus(group *jsonparser.Rule, newStatus string) {
	group.Status = newStatus
}

func displayGroupsWithStatus(data *jsonparser.JSONData) {
	var selected int
	sorting := false
	filterStatus := ""

	for {
		groups := filterGroups(data, filterStatus)

		// Apply sorting
		if sorting {
			sort.Slice(groups, func(i, j int) bool {
				return strings.Compare(groups[i].GroupID, groups[j].GroupID) < 0
			})
		}

		// Apply search
		if searchQuery != "" {
			groups = searchGroups(groups, searchQuery)
		}

		// Clear screen and move cursor to top-left
		fmt.Print("\033[2J\033[H")

		// Display the active filter at the top and the search query if any
		fmt.Printf("Active Filter: %s\n", getDisplayFilter(filterStatus))
		if searchQuery != "" {
			fmt.Printf("Search Query: %s\n", searchQuery)
		}
		fmt.Println()

		// Display group options with status
		fmt.Println("Groups:")
		for i, group := range groups {
			if i == selected {
				fmt.Printf("> %s [%s]\n", group.GroupID, group.Status)
			} else {
				fmt.Printf("%s [%s]\n", group.GroupID, group.Status)
			}
		}

		// Display instructions on the right side
		fmt.Println("\nKeys:")
		fmt.Println("j: Move down")
		fmt.Println("k: Move up")
		fmt.Println("h: Move left (previous filter)")
		fmt.Println("l: Move right (next filter)")
		fmt.Println("e: Select option")
		fmt.Println("n: Change status to not_a_finding")
		fmt.Println("x: Change status to not_applicable")
		fmt.Println("r: Change status to not_reviewed")
		fmt.Println("o: Change status to open")
		fmt.Println("s: Toggle sorting by GroupID")
		fmt.Println("/: Search")
		fmt.Println("Esc: Back to main menu")

		// Move the cursor to the bottom of the screen
		fmt.Print("\033[999B")

		// Read keyboard input
		key := term.GetKey()

		switch key {
		case "k":
			selected = (selected - 1 + len(groups)) % len(groups)
		case "j":
			selected = (selected + 1) % len(groups)
		case "h":
			filterStatus = prevFilterStatus(filterStatus)
		case "l":
			filterStatus = nextFilterStatus(filterStatus)
		case "e":
			if selected < len(groups) {
				runGroupSubmenu(groups[selected])
			}
		case "n":
			if selected < len(groups) {
				updateGroupStatus(groups[selected], "not_a_finding")
				SaveJSON()
			}
		case "x":
			if selected < len(groups) {
				updateGroupStatus(groups[selected], "not_applicable")
				SaveJSON()
			}
		case "r":
			if selected < len(groups) {
				updateGroupStatus(groups[selected], "not_reviewed")
				SaveJSON()
			}
		case "o":
			if selected < len(groups) {
				updateGroupStatus(groups[selected], "open")
				SaveJSON()
			}
		case "s":
			sorting = !sorting
		case "/":
			term.DisableRawMode()
			fmt.Print("Enter search query: ")
			reader := bufio.NewReader(os.Stdin)
			searchQuery, _ = reader.ReadString('\n')
			searchQuery = strings.TrimSpace(searchQuery)
			term.EnableRawMode()
		case "esc":
			return
		}
	}
}

func nextFilterStatus(current string) string {
	statuses := []string{"", "not_reviewed", "not_a_finding", "open", "not_applicable"}
	for i, status := range statuses {
		if status == current {
			return statuses[(i+1)%len(statuses)]
		}
	}
	return statuses[1]
}

func prevFilterStatus(current string) string {
	statuses := []string{"", "not_reviewed", "not_a_finding", "open", "not_applicable"}
	for i, status := range statuses {
		if status == current {
			return statuses[(i-1+len(statuses))%len(statuses)]
		}
	}
	return statuses[len(statuses)-1]
}

func getDisplayFilter(filter string) string {
	switch filter {
	case "":
		return "None"
	case "not_reviewed":
		return "Not Reviewed"
	case "not_a_finding":
		return "Not A Finding"
	case "not_applicable":
		return "Not Applicable"
	case "open":
		return "Open"
	}
	return filter
}

func runGroupSubmenu(group *jsonparser.Rule) {
	options := []string{"Status", "Fix Text", "Check Text"}
	selected := 0

	for {
		// Clear screen and move cursor to top-left
		fmt.Print("\033[2J\033[H")

		// Display details and submenu options
		fmt.Printf("Group ID: %s\n", group.GroupID)
		fmt.Printf("Status: %s\n", group.Status)
		fmt.Printf("Fix Text: %s\n", group.FixText)
		fmt.Printf("Check Text: %s\n", group.CheckText)
		fmt.Println("\nOptions:")
		for i, option := range options {
			if i == selected {
				fmt.Printf("> %s\n", option)
			} else {
				fmt.Println(option)
			}
		}

		// Display instructions on the right side
		fmt.Println("\nKeys:")
		fmt.Println("j: Move down")
		fmt.Println("k: Move up")
		fmt.Println("e: Edit option")
		fmt.Println("Esc: Back to group menu")

		// Move the cursor to the bottom of the screen
		fmt.Print("\033[999B")

		// Read keyboard input
		key := term.GetKey()

		switch key {
		case "k":
			selected = (selected - 1 + len(options)) % len(options)
		case "j":
			selected = (selected + 1) % len(options)
		case "e":
			editGroupDetail(options[selected], group)
			SaveJSON()
		case "esc":
			return
		}
	}
}

func editGroupDetail(option string, group *jsonparser.Rule) {
	term.DisableRawMode()      // Temporarily disable raw mode for input
	defer term.EnableRawMode() // Re-enable raw mode after getting the input

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Current %s: ", option)
	switch option {
	case "Status":
		fmt.Println(group.Status)
	case "Fix Text":
		fmt.Println(group.FixText)
	case "Check Text":
		fmt.Println(group.CheckText)
	}

	fmt.Printf("Enter new %s: ", option)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input) // Remove surrounding whitespace, including newline character

	switch option {
	case "Status":
		group.Status = input
	case "Fix Text":
		group.FixText = input
	case "Check Text":
		group.CheckText = input
	}
	fmt.Println("Value changed, press any key to continue...")
	fmt.Scanln() // Wait for input to resume raw mode
}

func extractGroups(data *jsonparser.JSONData) []*jsonparser.Rule {
	var groups []*jsonparser.Rule
	for i := range data.Stig {
		for j := range data.Stig[i].Rule {
			groups = append(groups, &data.Stig[i].Rule[j])
		}
	}
	return groups
}

func filterGroups(data *jsonparser.JSONData, filterStatus string) []*jsonparser.Rule {
	groups := extractGroups(data)

	// Apply filtering
	if filterStatus != "" {
		var filteredGroups []*jsonparser.Rule
		for _, group := range groups {
			if group.Status == filterStatus {
				filteredGroups = append(filteredGroups, group)
			}
		}
		return filteredGroups
	}
	return groups
}

func searchGroups(groups []*jsonparser.Rule, query string) []*jsonparser.Rule {
	var result []*jsonparser.Rule
	for _, group := range groups {
		if strings.Contains(group.GroupID, query) {
			result = append(result, group)
		}
	}
	return result
}
