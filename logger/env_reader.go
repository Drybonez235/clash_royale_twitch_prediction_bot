package logger

import (
	"bufio"
	"fmt"
	"os"
	"strings"
) 

//Subscripion info
type Env_variables struct{
	APP_ID string
	APP_SECRET string
	CLASH_API_SECRET string
	ENCRYPTION_SECRET string
	SUBSCRIPTION_INFO_URI string
	OAUTH_AUTHORIZE_URI string
	OAUTH_REFRESH_TOKEN_URI string
	OAUTH_USERINFO_URI string
	USER_INFO_URI string
	OAUTH_CLAIMS_INFO_URI string
	PREDICTION_URI string
	OAUTH_VALIDATE_TOKEN_URI string
	ROYALE_BETS_URL string
	ROYALE_BETS_REDIRECT_URL string
}

func Get_env_variables(file_path string) (Env_variables, error) {
	file, err := os.Open(file_path)
	if err != nil {
		return Env_variables{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var Env_struct Env_variables

	for scanner.Scan() {
		line := scanner.Text()
		token := strings.SplitN(line, "=", 2) // Use SplitN to avoid index issues

		if len(token) != 2 {
			return Env_variables{}, fmt.Errorf("invalid line format: %s", line)
		}

		key, value := strings.TrimSpace(token[0]), strings.TrimSpace(token[1])

		switch key {
		case "APP_ID":
			Env_struct.APP_ID = value
		case "APP_SECRET":
			Env_struct.APP_SECRET = value
		case "CLASH_API_SECRET":
			Env_struct.CLASH_API_SECRET = value
		case "ENCRYPTION_SECRET":
			Env_struct.ENCRYPTION_SECRET = value
		case "SUBSCRIPTION_INFO_URI":
			Env_struct.SUBSCRIPTION_INFO_URI = value
		case "ROYALE_BETS_URL":
			Env_struct.ROYALE_BETS_URL = value
		case "OAUTH_REFRESH_TOKEN_URI":
			Env_struct.OAUTH_REFRESH_TOKEN_URI = value
		case "OAUTH_AUTHORIZE_URI":
			Env_struct.OAUTH_AUTHORIZE_URI = value
		case "OAUTH_USERINFO_URI":
				Env_struct.OAUTH_USERINFO_URI = value
		case "USER_INFO_URI":
				Env_struct.USER_INFO_URI = value
		case "PREDICTION_URI":
				Env_struct.PREDICTION_URI = value
		case "OAUTH_CLAIMS_INFO_URI":
				Env_struct.OAUTH_CLAIMS_INFO_URI = value
		case "OAUTH_VALIDATE_TOKEN_URI":
				Env_struct.OAUTH_VALIDATE_TOKEN_URI = value
		case "ROYALE_BETS_REDIRECT_URL":
				Env_struct.ROYALE_BETS_REDIRECT_URL = value
		default:
			return Env_variables{}, fmt.Errorf("unexpected key: %s", key)
		}
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		return Env_variables{}, fmt.Errorf("error reading file: %w", err)
	}

	return Env_struct, nil
}