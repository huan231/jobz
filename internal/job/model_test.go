package job

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalJSONStatusRunning(t *testing.T) {
	v, _ := json.Marshal(Running)

	assert.Equal(t, []byte(`"running"`), v)
}

func TestMarshalJSONStatusSucceeded(t *testing.T) {
	v, _ := json.Marshal(Succeeded)

	assert.Equal(t, []byte(`"succeeded"`), v)
}

func TestMarshalJSONStatusFailed(t *testing.T) {
	v, _ := json.Marshal(Failed)

	assert.Equal(t, []byte(`"failed"`), v)
}

func TestMarshalJSONInvalidStatus(t *testing.T) {
	_, err := json.Marshal(Status(4))

	assert.NotNil(t, err)
}
