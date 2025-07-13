/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package awsauth

import (
	"fmt"
	"regexp"
	"strings"
)

var AllowedAWSAccounts = []string{
	"123456789012",
	"130488995195",
}

// VerifyAWSAccount validates if the provided string is a valid AWS account ID (12 digits)
func VerifyAWSAccount(accountID string) error {
	if accountID == "" {
		return fmt.Errorf("AWS account ID cannot be empty")
	}

	if !regexp.MustCompile(`^\d{12}$`).MatchString(accountID) {
		return fmt.Errorf("invalid AWS account ID format: %s. Must be exactly 12 digits", accountID)
	}

	return nil
}

// VerifyAWSAccountIDArn VerifyAWSAccountID verifies if the AWS account ID is allowed
func VerifyAWSAccountIDArn(arn string) error {
	// Extract account ID from ARN
	arnParts := strings.Split(arn, ":")
	if len(arnParts) < 5 {
		return fmt.Errorf("invalid ARN format: cannot extract account ID from %s", arn)
	}

	accountID := arnParts[4]

	// Validate account ID format
	if err := VerifyAWSAccount(accountID); err != nil {
		return err
	}

	// Check if account ID is in allowed list
	for _, allowed := range AllowedAWSAccounts {
		if accountID == allowed {
			return nil
		}
	}

	return fmt.Errorf("AWS account ID %s is not allowed. Allowed accounts: %s",
		accountID, strings.Join(AllowedAWSAccounts, ", "))
}

// VerifyUserARN validates AWS user ARN format
func VerifyUserARN(arn string) error {
	if arn == "" {
		return fmt.Errorf("userARN cannot be empty")
	}

	// Regular expression to verify ARN format
	arnPattern := `^arn:aws:iam::\d{12}:user/[a-zA-Z0-9+=,.@\-_/]+$`
	matched, _ := regexp.MatchString(arnPattern, arn)
	if !matched {
		return fmt.Errorf("invalid userARN format: %s. Expected format: arn:aws:iam::<account-id>:user/<username>", arn)
	}

	// Verify AWS account ID
	if err := VerifyAWSAccountIDArn(arn); err != nil {
		return err
	}

	return nil
}

func VerifyRoleARN(arn string) error {
	if arn == "" {
		return fmt.Errorf("roleARN cannot be empty")
	}

	// Regular expression to verify ARN format
	arnPattern := `^arn:aws:iam::\d{12}:role/[a-zA-Z0-9+=,.@\-_/]+$`
	matched, _ := regexp.MatchString(arnPattern, arn)
	if !matched {
		return fmt.Errorf("invalid roleARN format: %s. Expected format: arn:aws:iam::<account-id>:role/<rolename>", arn)
	}

	return nil
}

// VerifyUsername validates username format
func VerifyUsername(username string, enforceEmail bool) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	// Check length
	if len(username) < 3 || len(username) > 254 {
		return fmt.Errorf("username length must be between 3 and 254 characters")
	}

	// Validate email format if enforceEmail is true or if username contains @
	if enforceEmail || strings.Contains(username, "@") {
		emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(emailPattern, username)
		if !matched {
			return fmt.Errorf("invalid email format in username: %s", username)
		}
	} else {
		// Validate regular username format
		usernamePattern := `^[a-zA-Z0-9][a-zA-Z0-9._-]*[a-zA-Z0-9]$`
		matched, _ := regexp.MatchString(usernamePattern, username)
		if !matched {
			return fmt.Errorf("username can only contain letters, numbers, dots, hyphens, and underscores: %s", username)
		}
	}

	return nil
}

// VerifyGroups validates group names
func VerifyGroups(groups []string) error {
	if len(groups) == 0 {
		return fmt.Errorf("at least one group must be specified")
	}

	for _, group := range groups {
		if group == "" {
			return fmt.Errorf("group name cannot be empty")
		}

		if len(group) > 253 {
			return fmt.Errorf("group name cannot be longer than 253 characters: %s", group)
		}
	}

	return nil
}

// SetAllowedAWSAccounts updates the list of allowed AWS account IDs
func SetAllowedAWSAccounts(accounts []string) error {
	// Validate format of all accounts
	for _, acc := range accounts {
		// Validate account ID format
		if err := VerifyAWSAccount(acc); err != nil {
			return err
		}
	}

	AllowedAWSAccounts = accounts
	return nil
}

// ConvertMapAccountToString конвертує MapAccountSpec в строку accountid
func ConvertMapAccountToString(accounts []*MapAccount) []string {
	result := make([]string, len(accounts))
	for i, account := range accounts {
		result[i] = account.AccountID
	}
	return result
}

// ConvertStringToMapAccount конвертує строки accountid в MapAccountSpec
func ConvertStringToMapAccount(accountIDs []string) []*MapAccount {
	result := make([]*MapAccount, len(accountIDs))
	for i, accountID := range accountIDs {
		result[i] = &MapAccount{
			AccountID: accountID,
		}
	}
	return result
}

/*
	// Конвертація з MapAccountSpec в строки
	specs := []MapAccountSpec{
	{AccountID: "111111111111"},
	{AccountID: "222222222222"},
	}
	strings := ConvertMapAccountToString(specs)
	// результат: []string{"111111111111", "222222222222"}

	// Конвертація зі строк в MapAccountSpec
	accountIDs := []string{"111111111111", "222222222222"}
	specs = ConvertStringToMapAccount(accountIDs)
	// результат: []MapAccountSpec{{AccountID: "111111111111"}, {AccountID: "222222222222"}}
*/
