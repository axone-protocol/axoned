#!/usr/bin/env bash

set -eo pipefail

major_version=$(cut -d. -f1 < version)

if [ "${major_version}" -gt 1 ]; then
  module_name=$(go mod edit -json | jq -r '.Module.Path')
  module_name_unversioned=$(echo "${module_name}" | sed -E 's|/v[0-9]+$||')
  module_name_versioned="${module_name_unversioned}/v${major_version}"
  echo "🔬 major version detected, updating module path to ${module_name_versioned}"

  go mod edit -module "${module_name_versioned}"
  echo "✅ module name updated to ${module_name_versioned} in go.mod"

  sed_i_flag=""
  if [ "$(uname)" = "Darwin" ]; then
    sed_i_flag=(-i '')
  else
    sed_i_flag=(-i)
  fi
  echo "⬆️ updating ${module_name} to ${module_name_versioned}..."
  find . -type f \( -name "*.go" \) \
      -exec echo "  - processing {}" \; \
      -exec sed "${sed_i_flag[@]}" "s|\"${module_name}|\"${module_name_versioned}|g" {} \;

  find . -type f \( -name "README.md" \) \
      -exec echo "  - processing {}" \; \
      -exec sed "${sed_i_flag[@]}" "s|${module_name}|${module_name_versioned}|g" {} \;

  echo "✅ packages updated to ${module_name_versioned} in source files"
  echo "🧹 cleaning up go.sum"
  go mod tidy
else
  echo "🙅version is not greater than 1, no need to update module path"
fi
