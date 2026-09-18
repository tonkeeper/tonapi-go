#!/usr/bin/env bash
#
# Re-applies the local patches the upstream TonAPI spec is missing: the bearer
# security scheme, without which ogen drops oas_security_gen.go and the package
# stops compiling, and the Action.type value that must match its payload field.
#
# Idempotent. Exits non-zero when a patch can neither be applied nor be found
# already applied, which means upstream moved something.

set -euo pipefail

spec="${1:-api/openapi.yml}"
tmp="$spec.patch-openapi.tmp"
trap 'rm -f "$tmp"' EXIT

note() { printf 'patch-openapi: %s\n' "$1"; }
fail() { printf 'patch-openapi: %s\n' "$1" >&2; exit 1; }

[ -f "$spec" ] || fail "spec not found: $spec"

# Top-level security: bearer token, or anonymous.
if grep -q '^security:' "$spec"; then
	note 'security: already present'
else
	grep -q '^paths:$' "$spec" || fail "no top-level 'paths:' key to anchor the security block to"
	awk '
		/^paths:$/ && !inserted {
			print "security:"
			print "  - bearerAuth: [ ]"
			print "  - { }"
			print ""
			inserted = 1
		}
		{ print }
	' "$spec" >"$tmp" && mv "$tmp" "$spec"
	note 'security: restored'
fi

# The bearerAuth scheme the block above refers to.
if grep -q '^  securitySchemes:$' "$spec"; then
	note 'components.securitySchemes: already present'
else
	grep -q '^components:$' "$spec" || fail "no top-level 'components:' key to anchor securitySchemes to"
	awk '
		{ print }
		/^components:$/ && !inserted {
			print "  securitySchemes:"
			print "    bearerAuth:"
			print "      type: http"
			print "      scheme: bearer"
			inserted = 1
		}
	' "$spec" >"$tmp" && mv "$tmp" "$spec"
	note 'components.securitySchemes: restored'
fi

# Only the enum entry and the property key are renamed, the schema they point at
# really is called SetSignatureAllowedAction.
if grep -q '^            - SetSignatureAllowedAction$' "$spec" ||
	grep -q '^        SetSignatureAllowedAction:$' "$spec"; then
	sed -e 's/^            - SetSignatureAllowedAction$/            - SetSignatureAllowed/' \
		-e 's/^        SetSignatureAllowedAction:$/        SetSignatureAllowed:/' \
		"$spec" >"$tmp" && mv "$tmp" "$spec"
	note 'Action.SetSignatureAllowed: restored'
else
	note 'Action.SetSignatureAllowed: already correct'
fi

# Every patch must be observable in the result, whichever branch produced it.
while IFS= read -r probe; do
	grep -q "$probe" "$spec" || fail "verification failed: /$probe/ not found in $spec after patching"
done <<'PROBES'
^security:$
^  - bearerAuth: \[ \]$
^    bearerAuth:$
^            - SetSignatureAllowed$
^        SetSignatureAllowed:$
PROBES

note "$spec is patched"
