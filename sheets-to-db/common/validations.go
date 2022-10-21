package common

import (
	"fmt"
	"os"

	multierror "github.com/hashicorp/go-multierror"
)

// HasEnvironmentVars checks if all the listed environment variables are set.
func HasEnvironmentVars(keys []string) error {
	var result error
	for _, k := range keys {
		_, found := os.LookupEnv(k)
		if !found {
			result = multierror.Append(result, fmt.Errorf("expected environment var %q missing", k))
		}
	}
	return result
}
