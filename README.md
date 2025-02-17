# Discord Token Checker

This is a Discord Token Checker written in Go. It takes a list of tokens from a `tokens.txt` file, validates them by checking the associated user data via Discord's API, categorizes the tokens based on their subscription status, and stores them in separate output folders.

## Features:
- Checks whether Discord tokens are valid and retrieves user data (ID, username, email, phone, and subscriptions).
- Categorizes tokens based on subscription status (e.g., "1 Month", "3 Month", "UsedTokens", and invalid tokens).
- Saves valid and invalid tokens in separate files, categorized by subscription and validity.
- Provides a color-coded output for the status of each token (success or failure).
- Generates output folders with a timestamp for each run and stores categorized tokens in subfolders.

## Requirements:
- Go 1.16 or higher.
- A `tokens.txt` file containing Discord tokens (one token per line).

## Installation:

1. Clone this repository or download the Go code.
2. Install the required Go dependencies:

   ```bash
   go get github.com/fatih/color
   ```

3. Create a `tokens.txt` file in the root directory of your project. Add one Discord token per line.

## Usage:

1. Run the program by executing:

   ```bash
   go run main.go
   ```

2. The program will check the tokens, categorize them, and save the results in an output folder with a timestamp in the format `YYYY-MM-DD_HH-MM-SS`.

3. The output directories and files:
   - `output/YYYY-MM-DD_HH-MM-SS/`: Main directory for each run.
   - `1Month/`: Tokens with a "1 Month" subscription.
   - `3Month/`: Tokens with a "3 Month" subscription.
   - `UsedTokens/`: Tokens used but without valid subscriptions.
   - `InvalidTokens/`: Invalid tokens that couldn't be verified.

4. The status of each token is displayed on the console, color-coded as follows:
   - `[SUCCESS]`: The token is valid and successfully verified.
   - `[FAILED]`: The token is invalid or couldn't be verified.

## Code Overview:

### Key Functions:
- `checkToken`: Sends an HTTP request to Discord's API to verify the token and retrieve user information.
- `parseToken`: Extracts the token from each line of the `tokens.txt` file.
- `formatOutput`: Prints the status of each token in the console, formatted with colors.
- `getFlags`: Generates a list of flags based on user data (e.g., "VERIFIED", "NITRO").
- `categorizeTokens`: Categorizes tokens into the appropriate directory based on their subscription and validity status.
- `readTokens`: Reads tokens from the `tokens.txt` file.

### Concurrent Execution:
The program processes the tokens concurrently, using Goroutines and WaitGroups to handle multiple tokens at once.

## Example Output:

```
Starting Discord Token Checker...

15:04:05 | [SUCCESS] | Successfully Checked Token [Token: abc123xyz...] | Flags: [VERIFIED][1 Month]
15:04:05 | [FAILED] | Successfully Checked Token [Token: def456xyz...] | Flags: [INVALID]

COMPLETED | Finished Checking Tokens.
```

## License:
This project is licensed under the MIT License - see the LICENSE file for details.
