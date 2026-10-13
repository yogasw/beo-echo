package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseProjectAdvanceConfig_Cors(t *testing.T) {
	cfg, err := ParseProjectAdvanceConfig(`{"delayMs":100,"cors":{"mode":"custom","allowHeaders":["X-Api-Key"],"maxAge":600}}`)
	assert.NoError(t, err)
	assert.Equal(t, 100, cfg.DelayMs)
	assert.Equal(t, CorsModeCustom, cfg.Cors.Mode)
	assert.Equal(t, []string{"X-Api-Key"}, cfg.Cors.AllowHeaders)
	assert.Equal(t, 600, cfg.Cors.MaxAge)
}

func TestParseProjectAdvanceConfig_EmptyHasNoCors(t *testing.T) {
	cfg, err := ParseProjectAdvanceConfig("")
	assert.NoError(t, err)
	assert.Nil(t, cfg.Cors)
}

func TestParseProjectAdvanceConfig_InvalidCors(t *testing.T) {
	_, err := ParseProjectAdvanceConfig(`{"cors":{"mode":"wildcard"}}`)
	assert.Error(t, err)

	_, err = ParseProjectAdvanceConfig(`{"cors":{"mode":"reflect","maxAge":-1}}`)
	assert.Error(t, err)

	_, err = ParseProjectAdvanceConfig(`{"cors":{"mode":"reflect","maxAge":86401}}`)
	assert.Error(t, err)
}

func TestParseProjectAdvanceConfig_LegacyUnknownKeysAccepted(t *testing.T) {
	_, err := ParseProjectAdvanceConfig(`{"global_timeout":5000}`)
	assert.NoError(t, err)
}

func TestDefaultProjectAdvanceConfig_IsReflect(t *testing.T) {
	cfg, err := ParseProjectAdvanceConfig(DefaultProjectAdvanceConfig())
	assert.NoError(t, err)
	assert.Equal(t, CorsModeReflect, cfg.Cors.Mode)
}

func TestParseEndpointAdvanceConfig_Sequence(t *testing.T) {
	cfg, err := ParseEndpointAdvanceConfig(`{"sequence":{"onEnd":"loop"}}`)
	assert.NoError(t, err)
	assert.Equal(t, SequenceLoop, cfg.Sequence.OnEnd)

	_, err = ParseEndpointAdvanceConfig(`{"sequence":{"onEnd":"never"}}`)
	assert.Error(t, err)

	_, err = ParseEndpointAdvanceConfig(`{"timeout":60000,"retries":3}`)
	assert.NoError(t, err, "legacy keys stay accepted")
}

func TestAdvanceConfigToJSON_KeepsCorsAndSequence(t *testing.T) {
	p := &AdvanceConfigProject{Cors: &CorsConfig{Mode: CorsModeReflect}}
	s, err := p.ToJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `{"cors":{"mode":"reflect"}}`, s)

	e := &AdvanceConfigEndpoint{Sequence: &SequenceConfig{OnEnd: SequenceLoop}}
	s, err = e.ToJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `{"sequence":{"onEnd":"loop"}}`, s)

	s, err = (&AdvanceConfigProject{}).ToJSON()
	assert.NoError(t, err)
	assert.Equal(t, "", s)
}
