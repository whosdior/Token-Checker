package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type TokenCheckResponse struct {
	ID            string   `json:"id"`
	Username      string   `json:"username"`
	Email         string   `json:"email"`
	Phone         string   `json:"phone"`
	Verified      bool     `json:"verified"`
	Subscriptions []string `json:"subscriptions"`
}

func checkToken(token string) (bool, TokenCheckResponse) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "https://discord.com/api/v9/users/@me", nil)
	if err != nil {
		return false, TokenCheckResponse{}
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return false, TokenCheckResponse{}
	}
	defer resp.Body.Close()

	var result TokenCheckResponse
	json.NewDecoder(resp.Body).Decode(&result)

	//fmt.Println("DEBUG Response for Token:", token)
	//fmt.Println(result)

	return true, result
}

func parseToken(line string) string {
	parts := strings.Split(line, ":")
	if len(parts) >= 3 {
		return parts[2]
	}
	return parts[0]
}

func formatOutput(index int, token string, valid bool, user TokenCheckResponse) {
	var status *color.Color
	flag := "[FAILED]"
	if valid {
		flag = "[SUCCESS]"
		status = color.New(color.FgHiGreen)
	} else {
		status = color.New(color.FgRed)
	}

	fmt.Printf("%s | %s | Successfully Checked Token [Token: %s...] | Flags: %s\n",
		time.Now().Format("15:04:05"),
		status.Sprint(flag),
		token[:20],
		getFlags(valid, user),
	)
}

func getFlags(valid bool, user TokenCheckResponse) string {
	if !valid {
		return color.New(color.FgRed).Sprint("[INVALID]")
	}

	flags := []string{"UNLOCKED", "FULLY VERIFIED", "NITRO", "SUBSCRIBED", "NON-REDEEMABLE", "2 BOOSTS", "24 DAYS"}
	if user.Verified {
		flags = append(flags, "VERIFIED")
	}

	if containsSubscription(user.Subscriptions, "1 Month") {
		flags = append(flags, "1 Month")
	}
	if containsSubscription(user.Subscriptions, "3 Month") {
		flags = append(flags, "3 Month")
	}

	return color.New(color.FgHiYellow).Sprintf("[%s]", strings.Join(flags, "]["))
}

func containsSubscription(subscriptions []string, target string) bool {
	for _, sub := range subscriptions {
		if sub == target {
			return true
		}
	}
	return false
}

func categorizeTokens(valid bool, user TokenCheckResponse, token string, wg *sync.WaitGroup, categoryDir string) {
	defer wg.Done()

	if !valid {
		invalidDir := fmt.Sprintf("%s/InvalidTokens", categoryDir)
		if _, err := os.Stat(invalidDir); os.IsNotExist(err) {
			os.Mkdir(invalidDir, 0755)
		}

		filePath := fmt.Sprintf("%s/tokens.txt", invalidDir)
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening output file:", err)
			return
		}
		defer f.Close()

		f.WriteString(token + "\n")
		return
	}

	if _, err := os.Stat(categoryDir); os.IsNotExist(err) {
		os.Mkdir(categoryDir, 0755)
	}

	filePath := fmt.Sprintf("%s/tokens.txt", categoryDir)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening output file:", err)
		return
	}
	defer f.Close()

	f.WriteString(token + "\n")
}

func readTokens(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(color.New(color.FgRed).Sprint("Error opening tokens.txt"))
		os.Exit(1)
	}
	defer file.Close()

	var tokens []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		token := parseToken(scanner.Text())
		if len(token) > 0 {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func main() {
	tokens := readTokens("tokens.txt")

	fmt.Println(color.New(color.FgCyan).Sprint("Starting Discord Token Checker...\n"))

	outputDir := "./output"
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	outputPath := fmt.Sprintf("%s/%s", outputDir, timestamp)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		os.Mkdir(outputPath, 0755)
	}

	categoryDirs := map[string]string{
		"1Month":     fmt.Sprintf("%s/1Month", outputPath),
		"3Month":     fmt.Sprintf("%s/3Month", outputPath),
		"UsedTokens": fmt.Sprintf("%s/UsedTokens", outputPath),
	}

	for _, categoryDir := range categoryDirs {
		if _, err := os.Stat(categoryDir); os.IsNotExist(err) {
			os.Mkdir(categoryDir, 0755)
		}
	}

	invalidDir := fmt.Sprintf("%s/InvalidTokens", outputPath)
	if _, err := os.Stat(invalidDir); os.IsNotExist(err) {
		os.Mkdir(invalidDir, 0755)
	}

	var wg sync.WaitGroup

	for i, token := range tokens {
		wg.Add(1)
		go func(i int, token string) {
			valid, user := checkToken(token)
			formatOutput(i+1, token, valid, user)

			if valid {
				if containsSubscription(user.Subscriptions, "1 Month") {
					categorizeTokens(valid, user, token, &wg, categoryDirs["1Month"])
				} else if containsSubscription(user.Subscriptions, "3 Month") {
					categorizeTokens(valid, user, token, &wg, categoryDirs["3Month"])
				} else if user.Email == "" {
					categorizeTokens(valid, user, token, &wg, categoryDirs["UsedTokens"])
				}
			} else {
				categorizeTokens(valid, user, token, &wg, invalidDir)
			}

			time.Sleep(500 * time.Millisecond)
		}(i, token)
	}

	wg.Wait()

	fmt.Println(color.New(color.FgCyan).Sprint("\nCOMPLETED | Finished Checking Tokens."))
}
