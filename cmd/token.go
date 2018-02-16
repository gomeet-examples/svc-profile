package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/gomeet/gomeet/utils/jwt"
)

const (
	// issuer (the principal that issues tokens)
	tokenIssuer = "https://gomeet.com/"

	// default token lifetime expressed in hours
	defaultTokenLifetimeHours = 24
)

var (
	secretSigningKey   string
	tokenLifetimeHours int
	subjectID          string
	customClaims       string

	tokenCmd = &cobra.Command{
		Use:   "token",
		Short: "Creates a JSON Web Token",
		Run:   createToken,
	}
)

func init() {
	RootCmd.AddCommand(tokenCmd)

	// JSON Web Token secret signing key
	tokenCmd.PersistentFlags().StringVarP(&secretSigningKey, "secret-key", "k", "", "JSON Web Token secret signing key")

	// JSON Web Token lifetime
	tokenCmd.PersistentFlags().IntVarP(&tokenLifetimeHours, "lifetime", "l", defaultTokenLifetimeHours, "JSON Web Token lifetime expressed in hours")

	// JSON Web Token subject identifier
	tokenCmd.PersistentFlags().StringVarP(&subjectID, "subject", "s", "", "JSON Web Token subject identifier")

	// JSON Web Token claims
	tokenCmd.PersistentFlags().StringVarP(&customClaims, "claims", "c", "", "JSON Web Token custom claims encoded as a JSON object")
}

func createToken(cmd *cobra.Command, args []string) {
	// ensure the user provided the secret signing key
	if secretSigningKey == "" {
		fmt.Printf("missing the secret signing key\n")
		os.Exit(1)
	}

	// set custom token claims
	var claimsMap jwt.Claims

	if customClaims != "" {
		var (
			claimsData interface{}
			ok         bool
		)
		err := json.Unmarshal([]byte(customClaims), &claimsData)
		if err != nil {
			fmt.Printf("JSON parsing error: %v\n", err)
			os.Exit(1)
		}

		claimsMap, ok = claimsData.(map[string]interface{})
		if !ok {
			fmt.Printf("JSON parsing error: failed type assertion on JSON object - %v\n", err)
			os.Exit(1)
		}
	}

	token, err := jwt.Create(
		"github.com/gomeet/gomeet",
		secretSigningKey,
		tokenLifetimeHours,
		subjectID,
		claimsMap,
	)

	if err != nil {
		fmt.Printf("failed to create JWT : %v\n", err)
		os.Exit(1)
	}

	// display the token
	fmt.Printf("%s\n", token)
}
