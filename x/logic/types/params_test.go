package types_test

import (
	"testing"

	"github.com/okp4/okp4d/x/logic/types"
	"github.com/stretchr/testify/require"
)

func Test_validateParams(t *testing.T) {
	params := types.DefaultParams()

	// default params have no error
	require.NoError(t, params.Validate())
}
