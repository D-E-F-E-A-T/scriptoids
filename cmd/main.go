package main

import (
	"fmt"
	"github.com/dhsavell/scriptoids/pkg/scriptoids"
	"log"
)

func main() {
	env, err := scriptoids.NewEnvironmentFromEnvVars()

	if err != nil {
		log.Fatal(err)
	}

	_, e := env.GetInstalledPackageByName("a")
	fmt.Printf("installed: %v", e)
}