package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/pflag"
)

var (
	cliAlgorithm = pflag.StringP("algorithm", "", "HS256", "Signing algorithm - possible values are HS256, HS384, HS512")
	cliTimeout   = pflag.DurationP("timeout", "", 2*time.Hour, "JWT token expires time")
	help         = pflag.BoolP("help", "h", false, "Print this help message")
)

func main() {
	pflag.Usage = func() {
		fmt.Println(`Usage: gentoken [OPTIONS] SECRETID SECRETKEY`)
		pflag.PrintDefaults()
	}
	pflag.Parse()

	if *help {
		pflag.Usage()
		return
	}

	if pflag.NArg() != 2 {
		pflag.Usage()
		os.Exit(1)
	}

	token, err := createJWTToken(*cliAlgorithm, *cliTimeout, os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Println(token)
}

func createJWTToken(algorithm string, timeout time.Duration, secretID, secretKey string) (string, error) {
	expire := time.Now().Add(timeout)

	token := jwt.NewWithClaims(jwt.GetSigningMethod(algorithm), jwt.MapClaims{
		"jti": secretID,
		"exp": expire.Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString([]byte(secretKey))
}
