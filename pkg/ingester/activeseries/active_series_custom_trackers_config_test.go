// SPDX-License-Identifier: AGPL-3.0-only

package activeseries

import (
	"flag"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func mustNewActiveSeriesCustomTrackersConfigFromMap(t *testing.T, source map[string]string) *ActiveSeriesCustomTrackersConfig {
	m, err := NewActiveSeriesCustomTrackersConfig(source)
	require.NoError(t, err)
	return &m
}

func mustNewActiveSeriesCustomTrackersConfigFromString(t *testing.T, source string) *ActiveSeriesCustomTrackersConfig {
	m := ActiveSeriesCustomTrackersConfig{}
	err := m.Set(source)
	require.NoError(t, err)
	return &m
}

func mustNewActiveSeriesCustomTrackersConfigDeserializedFromYaml(t *testing.T, yamlString string) *ActiveSeriesCustomTrackersConfig {
	m := ActiveSeriesCustomTrackersConfig{}
	err := yaml.Unmarshal([]byte(yamlString), &m)
	require.NoError(t, err)
	return &m
}

func TestActiveSeriesCustomTrackersConfigs(t *testing.T) {
	for _, tc := range []struct {
		name     string
		flags    []string
		expected *ActiveSeriesCustomTrackersConfig
		error    error
	}{
		{
			name:     "empty flag value produces empty config",
			flags:    []string{`-ingester.active-series-custom-trackers=`},
			expected: &ActiveSeriesCustomTrackersConfig{},
		},
		{
			name:  "empty matcher fails",
			flags: []string{`-ingester.active-series-custom-trackers=foo:`},
			error: errors.New(`invalid value "foo:" for flag -ingester.active-series-custom-trackers: semicolon-separated values should be <name>:<matcher>, but one of the sides was empty in the value 0: "foo:"`),
		},
		{
			name:  "empty whitespace-only matcher fails",
			flags: []string{`-ingester.active-series-custom-trackers=foo: `},
			error: errors.New(`invalid value "foo: " for flag -ingester.active-series-custom-trackers: semicolon-separated values should be <name>:<matcher>, but one of the sides was empty in the value 0: "foo: "`),
		},
		{
			name:  "second empty whitespace-only matcher fails",
			flags: []string{`-ingester.active-series-custom-trackers=foo: ;bar:{}`},
			error: errors.New(`invalid value "foo: ;bar:{}" for flag -ingester.active-series-custom-trackers: semicolon-separated values should be <name>:<matcher>, but one of the sides was empty in the value 0: "foo: "`),
		},
		{
			name:  "empty name fails",
			flags: []string{`-ingester.active-series-custom-trackers=:{}`},
			error: errors.New(`invalid value ":{}" for flag -ingester.active-series-custom-trackers: semicolon-separated values should be <name>:<matcher>, but one of the sides was empty in the value 0: ":{}"`),
		},
		{
			name:  "empty whitespace-only name fails",
			flags: []string{`-ingester.active-series-custom-trackers= :{}`},
			error: errors.New(`invalid value " :{}" for flag -ingester.active-series-custom-trackers: semicolon-separated values should be <name>:<matcher>, but one of the sides was empty in the value 0: " :{}"`),
		},
		{
			name:     "one matcher",
			flags:    []string{`-ingester.active-series-custom-trackers=foo:{foo="bar"}`},
			expected: mustNewActiveSeriesCustomTrackersConfigFromMap(t, map[string]string{`foo`: `{foo="bar"}`}),
		},
		{
			name:     "whitespaces are trimmed from name and matcher",
			flags:    []string{`-ingester.active-series-custom-trackers= foo :	{foo="bar"}` + "\n "},
			expected: mustNewActiveSeriesCustomTrackersConfigFromMap(t, map[string]string{`foo`: `{foo="bar"}`}),
		},
		{
			name:     "two matchers in one flag value",
			flags:    []string{`-ingester.active-series-custom-trackers=foo:{foo="bar"};baz:{baz="bar"}`},
			expected: mustNewActiveSeriesCustomTrackersConfigFromMap(t, map[string]string{`foo`: `{foo="bar"}`, `baz`: `{baz="bar"}`}),
		},
		{
			name:     "two matchers in two flag values",
			flags:    []string{`-ingester.active-series-custom-trackers=foo:{foo="bar"}`, `-ingester.active-series-custom-trackers=baz:{baz="bar"}`},
			expected: mustNewActiveSeriesCustomTrackersConfigFromMap(t, map[string]string{`foo`: `{foo="bar"}`, `baz`: `{baz="bar"}`}),
		},
		{
			name:  "two matchers with same name in same flag",
			flags: []string{`-ingester.active-series-custom-trackers=foo:{foo="bar"};foo:{boo="bam"}`},
			error: errors.New(`invalid value "foo:{foo=\"bar\"};foo:{boo=\"bam\"}" for flag -ingester.active-series-custom-trackers: matcher "foo" for active series custom trackers is provided twice`),
		},
		{
			name:  "two matchers with same name in separate flags",
			flags: []string{`-ingester.active-series-custom-trackers=foo:{foo="bar"}`, `-ingester.active-series-custom-trackers=foo:{boo="bam"}`},
			error: errors.New(`invalid value "foo:{boo=\"bam\"}" for flag -ingester.active-series-custom-trackers: matcher "foo" for active series custom trackers is provided more than once`),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			flagSet := flag.NewFlagSet("test", flag.ContinueOnError)

			var config ActiveSeriesCustomTrackersConfig
			flagSet.Var(&config, "ingester.active-series-custom-trackers", "...usage docs...")
			err := flagSet.Parse(tc.flags)

			if tc.error != nil {
				assert.EqualError(t, err, tc.error.Error())
				return
			}

			require.Equal(t, tc.expected, &config)

			// Check that ActiveSeriesCustomTrackersConfig.String() value is a valid flag value.
			flagSetAgain := flag.NewFlagSet("test-string", flag.ContinueOnError)
			var configAgain ActiveSeriesCustomTrackersConfig
			flagSetAgain.Var(&configAgain, "ingester.active-series-custom-trackers", "...usage docs...")
			require.NoError(t, flagSetAgain.Parse([]string{"-ingester.active-series-custom-trackers=" + config.String()}))

			require.Equal(t, tc.expected, &configAgain)
		})
	}
}

func TestActiveSeriesCustomTrackerConfig_Equality(t *testing.T) {
	configSets := [][]ActiveSeriesCustomTrackersConfig{
		{
			*mustNewActiveSeriesCustomTrackersConfigFromString(t, `foo:{foo='bar'};baz:{baz='bar'}`),
			*mustNewActiveSeriesCustomTrackersConfigFromMap(t, map[string]string{
				"baz": `{baz='bar'}`,
				"foo": `{foo='bar'}`,
			}),
			*mustNewActiveSeriesCustomTrackersConfigDeserializedFromYaml(t,
				`
                baz: "{baz='bar'}"
                foo: "{foo='bar'}"`),
		},
		{
			*mustNewActiveSeriesCustomTrackersConfigFromString(t, `test:{test='true'}`),
			*mustNewActiveSeriesCustomTrackersConfigFromMap(t, map[string]string{"test": `{test='true'}`}),
			*mustNewActiveSeriesCustomTrackersConfigDeserializedFromYaml(t, `test: "{test='true'}"`),
		},
		{
			*mustNewActiveSeriesCustomTrackersConfigDeserializedFromYaml(t,
				`
        baz: "{baz='bar'}"
        foo: "{foo='bar'}"
        extra: "{extra='extra'}"`),
		},
	}

	for _, configSet := range configSets {
		t.Run("EqualityBetweenSet", func(t *testing.T) {
			for i := 0; i < len(configSet); i++ {
				for j := i + 1; j < len(configSet); j++ {
					assert.Equal(t, configSet[i].String(), configSet[j].String(), "matcher configs should be equal")
				}
			}
		})
	}

	t.Run("NotEqualsAcrossSets", func(t *testing.T) {
		var activeSeriesMatchers []*ActiveSeriesCustomTrackersConfig
		for _, matcherConfigs := range configSets {
			activeSeriesMatchers = append(activeSeriesMatchers, &matcherConfigs[0])
		}

		for i := 0; i < len(activeSeriesMatchers); i++ {
			for j := i + 1; j < len(activeSeriesMatchers); j++ {
				assert.NotEqual(t, activeSeriesMatchers[i].String(), activeSeriesMatchers[j].String(), "matcher configs should NOT be equal")
			}
		}
	})

}

func TestActiveSeriesCustomTrackersConfigs_Deserialization(t *testing.T) {
	correctInput := `
        baz: "{baz='bar'}"
        foo: "{foo='bar'}"
    `
	malformedInput :=
		`
        baz: "123"
        foo: "{foo='bar'}"
    `
	t.Run("ShouldDeserializeCorrectInput", func(t *testing.T) {
		config := ActiveSeriesCustomTrackersConfig{}
		err := yaml.Unmarshal([]byte(correctInput), &config)
		assert.NoError(t, err, "failed do deserialize ActiveSeriesMatchers")
		expectedConfig, err := NewActiveSeriesCustomTrackersConfig(map[string]string{
			"baz": "{baz='bar'}",
			"foo": "{foo='bar'}",
		})
		require.NoError(t, err)
		assert.Equal(t, expectedConfig.String(), config.String())
	})

	t.Run("ShouldErrorOnMalformedInput", func(t *testing.T) {
		config := ActiveSeriesCustomTrackersConfig{}
		err := yaml.Unmarshal([]byte(malformedInput), &config)
		assert.Error(t, err, "should not deserialize malformed input")
	})
}
