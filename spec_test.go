package tonapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const (
	specPath    = "api/openapi.yml"
	liveSpecURL = "https://tonapi.io/v2/openapi.json"
)

// knownSpecDivergence lists paths where specPath deliberately differs from the served spec.
// Each is applied by scripts/patch-openapi.sh.
var knownSpecDivergence = map[string]string{
	"components.schemas.Action.properties.type": "upstream sends SetSignatureAllowedAction, the payload field is SetSignatureAllowed",
}

func loadLocalSpec(t *testing.T) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(specPath)
	require.NoError(t, err)
	var spec map[string]any
	require.NoError(t, yaml.Unmarshal(raw, &spec))
	return spec
}

// TestSpecDeclaresBearerAuth tests that the security section survived the last spec refresh.
// Without it ogen emits no oas_security_gen.go and the package stops compiling.
func TestSpecDeclaresBearerAuth(t *testing.T) {
	var spec struct {
		Security   []map[string]any `yaml:"security"`
		Components struct {
			SecuritySchemes map[string]struct {
				Type   string `yaml:"type"`
				Scheme string `yaml:"scheme"`
			} `yaml:"securitySchemes"`
		} `yaml:"components"`
	}
	raw, err := os.ReadFile(specPath)
	require.NoError(t, err)
	require.NoError(t, yaml.Unmarshal(raw, &spec))

	scheme, ok := spec.Components.SecuritySchemes["bearerAuth"]
	require.Truef(t, ok, "components.securitySchemes.bearerAuth missing from %s, run `make patch-openapi`", specPath)
	require.Equal(t, "http", scheme.Type)
	require.Equal(t, "bearer", scheme.Scheme)

	var bearer, anonymous bool
	for _, requirement := range spec.Security {
		if _, ok := requirement["bearerAuth"]; ok {
			bearer = true
		}
		if len(requirement) == 0 {
			anonymous = true
		}
	}
	require.True(t, bearer, "top-level security does not offer bearerAuth")
	require.True(t, anonymous, "top-level security does not offer the anonymous alternative")
}

// TestSpecActionUnionIsConsistent tests that every Action.type value names the property holding
// that action's payload. The generated code compiles either way, so nothing else notices.
func TestSpecActionUnionIsConsistent(t *testing.T) {
	spec := loadLocalSpec(t)

	schemas, ok := dig(spec, "components", "schemas")
	require.True(t, ok, "components.schemas missing")
	action, ok := schemas.(map[string]any)["Action"].(map[string]any)
	require.True(t, ok, "components.schemas.Action missing")
	properties, ok := action["properties"].(map[string]any)
	require.True(t, ok, "Action has no properties")

	// Carried by every action whatever its type.
	common := map[string]bool{"type": true, "status": true, "simple_preview": true, "base_transactions": true}
	// The one type value with no payload of its own.
	const noPayload = "Unknown"

	payloads := map[string]bool{}
	for name := range properties {
		if !common[name] {
			payloads[name] = true
		}
	}

	typeProperty, ok := properties["type"].(map[string]any)
	require.True(t, ok, "Action has no type property")
	values, ok := typeProperty["enum"].([]any)
	require.True(t, ok, "Action.type has no enum")
	require.NotEmpty(t, values)

	for _, raw := range values {
		value := fmt.Sprint(raw)
		if value == noPayload {
			continue
		}
		require.Truef(t, payloads[value], "Action.type value %q has no matching payload property on Action", value)
		delete(payloads, value)
	}
	require.Empty(t, payloads, "Action carries payload properties that no Action.type value selects")
}

// TestSpecEnumsMatchLiveAPI tests that specPath lists every enum value TonAPI currently returns,
// since an unlisted one fails response validation. Skips when the API is unreachable.
func TestSpecEnumsMatchLiveAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live spec comparison in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, liveSpecURL, nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Skipf("cannot reach %s: %v", liveSpecURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Skipf("%s returned %s", liveSpecURL, resp.Status)
	}

	var served map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&served))

	live, local := collectEnums(served), collectEnums(loadLocalSpec(t))
	require.NotEmpty(t, live, "served spec has no enums")

	for _, path := range sortedKeys(live) {
		if reason, ok := knownSpecDivergence[path]; ok {
			t.Logf("%s: known local divergence (%s)", path, reason)
			continue
		}
		if _, ok := local[path]; !ok {
			t.Errorf("%s is served but missing from %s, refresh the spec and run `make generate-client`", path, specPath)
			continue
		}
		for _, value := range live[path] {
			require.Containsf(t, local[path], value,
				"%s: TonAPI returns %q, which %s does not list", path, value, specPath)
		}
	}
}

// collectEnums returns every enum under components.schemas, keyed by its dotted path.
func collectEnums(doc map[string]any) map[string][]string {
	found := map[string][]string{}

	var walk func(node any, path string)
	walk = func(node any, path string) {
		switch n := node.(type) {
		case map[string]any:
			if raw, ok := n["enum"].([]any); ok {
				values := make([]string, 0, len(raw))
				for _, value := range raw {
					values = append(values, fmt.Sprint(value))
				}
				found[path] = values
			}
			for key, value := range n {
				if key == "enum" {
					continue
				}
				walk(value, path+"."+key)
			}
		case []any:
			for i, value := range n {
				walk(value, fmt.Sprintf("%s[%d]", path, i))
			}
		}
	}

	if schemas, ok := dig(doc, "components", "schemas"); ok {
		walk(schemas, "components.schemas")
	}
	return found
}

// dig follows a chain of mapping keys, reporting whether the whole chain resolved.
func dig(node any, keys ...string) (any, bool) {
	for _, key := range keys {
		mapping, ok := node.(map[string]any)
		if !ok {
			return nil, false
		}
		node, ok = mapping[key]
		if !ok {
			return nil, false
		}
	}
	return node, true
}

func sortedKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
