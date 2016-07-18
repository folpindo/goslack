package rtmtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalEvent(t *testing.T) {
	raw := []byte(`{ "type": "message" }`)
	event, err := Unmarshal(raw)
	require.Nil(t, err)
	assert.Equal(t, "message", event.Type)
	assert.Equal(t, raw, event.raw)
}

func TestUnmarshalInvalidEvent(t *testing.T) {
	raw := []byte(`{}`)
	_, err := Unmarshal(raw)
	assert.Error(t, err)
}

func TestCastToRtmMessage(t *testing.T) {
	raw := []byte(`{ "type": "message", "id": 123456, "channel": "CA1B2C3", "user": "UD1E2F3", "text": "Message content" }`)
	event, err := Unmarshal(raw)
	require.Nil(t, err)

	message, err := event.RtmMessage()
	require.Nil(t, err)

	assert.Equal(t, "message", message.Type)
	assert.Equal(t, uint64(123456), message.ID)
	assert.Equal(t, "CA1B2C3", message.Channel)
	assert.Equal(t, "UD1E2F3", message.User)
	assert.Equal(t, "Message content", message.Text)
}

func TestCastToInvalidRtmMessage(t *testing.T) {
	raw := []byte(`{ "type": "banana" }`)
	event, err := Unmarshal(raw)
	require.Nil(t, err)

	_, err = event.RtmMessage()
	assert.Error(t, err)
}

func TestCastToRtmUserChange(t *testing.T) {
	raw := []byte(`{ "type": "user_change", "user": { "id": "U123ABC", "name": "Bananaman" } }`)
	event, err := Unmarshal(raw)
	require.Nil(t, err)

	userChange, err := event.RtmUserChange()
	require.Nil(t, err)

	assert.Equal(t, "user_change", userChange.Type)
	assert.Equal(t, "U123ABC", userChange.User.ID)
	assert.Equal(t, "Bananaman", userChange.User.Name)
}

func TestCastToInvalidRtmUserChange(t *testing.T) {
	raw := []byte(`{ "type": "banana" }`)
	event, err := Unmarshal(raw)
	require.Nil(t, err)

	_, err = event.RtmUserChange()
	assert.Error(t, err)
}
