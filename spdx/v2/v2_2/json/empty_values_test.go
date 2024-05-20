package json

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spdx/tools-golang/spdx/v2/common"
	spdx "github.com/spdx/tools-golang/spdx/v2/v2_2"
)

func Test_omitsAppropriateProperties(t *testing.T) {
	tests := []struct {
		name     string
		pkg      spdx.Package
		validate func(t *testing.T, got map[string]interface{})
	}{
		{
			name: "include packageVerificationCode exclusions",
			pkg: spdx.Package{
				PackageVerificationCode: common.PackageVerificationCode{
					ExcludedFiles: []string{},
				},
			},
			validate: func(t *testing.T, got map[string]interface{}) {
				require.Contains(t, got, "packageVerificationCode")
			},
		},
		{
			name: "include packageVerificationCode value",
			pkg: spdx.Package{
				PackageVerificationCode: common.PackageVerificationCode{
					Value: "1234",
				},
			},
			validate: func(t *testing.T, got map[string]interface{}) {
				require.Contains(t, got, "packageVerificationCode")
			},
		},
		{
			name: "omit empty packageVerificationCode",
			pkg: spdx.Package{
				PackageVerificationCode: common.PackageVerificationCode{},
			},
			validate: func(t *testing.T, got map[string]interface{}) {
				require.NotContains(t, got, "packageVerificationCode")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := json.Marshal(test.pkg)
			require.NoError(t, err)
			var unmarshalled map[string]interface{}
			err = json.Unmarshal(got, &unmarshalled)
			require.NoError(t, err)
			test.validate(t, unmarshalled)
		})
	}
}
