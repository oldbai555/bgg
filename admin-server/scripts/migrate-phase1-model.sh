#!/usr/bin/env bash
# Phase 1: migrate internal/model/*.go to internal/model/<domain>/
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
MODEL_DIR="$ROOT/internal/model"

declare -A DOMAIN_MAP=(
  # iam
  [adminusermodel]=iam
  [adminrolemodel]=iam
  [adminpermissionmodel]=iam
  [adminmenumodel]=iam
  [admindepartmentmodel]=iam
  [adminuserrolemodel]=iam
  [adminrolepermissionmodel]=iam
  [adminapimodel]=iam
  [adminpermissionmenumodel]=iam
  [adminpermissionapimodel]=iam
  # blog
  [blogtagmodel]=blog
  [blogarticlemodel]=blog
  [blogarticletagmodel]=blog
  [blogarticleauditmodel]=blog
  [blogfriendlinkmodel]=blog
  [blogsocialinfomodel]=blog
  # video
  [videomodel]=video
  # chat
  [chatmodel]=chat
  [chatusermodel]=chat
  [chatmessagemodel]=chat
  # sdk
  [sdkkeymodel]=sdk
  [sdkinterfacemodel]=sdk
  [sdkkeyapimodel]=sdk
  [sdkcalllogmodel]=sdk
  # task
  [admintaskmodel]=task
  # monitoring
  [adminoperationlogmodel]=monitoring
  [adminloginlogmodel]=monitoring
  [auditlogmodel]=monitoring
  [adminperformancelogmodel]=monitoring
  # system
  [adminconfigmodel]=system
  [admindicttypemodel]=system
  [admindictitemmodel]=system
  [adminfilemodel]=system
  [adminnoticemodel]=system
  [adminnotificationmodel]=system
  [filemodel]=system
  # misc
  [demomodel]=misc
  [dailyshortsentencemodel]=misc
)

for domain in iam blog video chat sdk task monitoring system misc; do
  mkdir -p "$MODEL_DIR/$domain"
done

# vars.go stays in internal/model as shared package
if [[ -f "$MODEL_DIR/vars.go" ]]; then
  echo "Keeping vars.go in internal/model (shared ErrNotFound)"
fi

for f in "$MODEL_DIR"/*.go; do
  [[ -f "$f" ]] || continue
  base=$(basename "$f" .go)
  # strip _gen suffix for lookup
  lookup=$base
  lookup=${lookup/_gen/}

  if [[ "$base" == "vars" ]]; then
    continue
  fi

  domain=${DOMAIN_MAP[$lookup]:-}
  if [[ -z "$domain" ]]; then
    echo "ERROR: no domain mapping for $base" >&2
    exit 1
  fi

  dest="$MODEL_DIR/$domain/$base.go"
  if [[ -f "$dest" ]]; then
    echo "SKIP exists: $dest"
    continue
  fi

  sed 's/^package model$/package '"$domain"'/' "$f" > "$dest"
  rm "$f"
  echo "MOVED $base -> $domain/"
done

echo "Phase 1 file migration done."
