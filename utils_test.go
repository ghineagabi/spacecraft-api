package main

import (
	"fmt"
	"testing"
)

func TestSHA512(t *testing.T) {
	strings := []string{"admin", "root", "sysadmin"}
	var shaStrings []string
	actual := make(map[string]string)

	for _, s := range strings {
		shaStrings = append(shaStrings, SHA512(s))
	}

	for i, s := range strings {
		fmt.Printf("%s: %s \n", s, shaStrings[i])
		actual[s] = shaStrings[i]
	}

	expected := map[string]string{"admin": "30bb8411dd0cbf96b10a52371f7b3be1690f7afa16c3bd7bc7d02c0e2854768d",
		"root":     "a09b0c9c4bf4bd4c9c5ec45f574531c7d0c8580825d6d8586954eddfbfefb250",
		"sysadmin": "7c8dc8ce5e2d5756eaed391bc6a5cb3e5879c80c049ccb89f5d366456be3ceb9"}

	for k, v := range actual {
		if val, ok := expected[k]; ok {
			if val != v {
				t.Fatalf("SHA value for %s unfound. expected: %s, actual: %s", k, val, v)
			}
		}
	}
}
