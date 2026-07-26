#!/usr/bin/env bash
# Check every http(s) link in the repo's markdown is live. External links are
# flaky and some hosts block bots, so this runs as a dedicated (non-push) check:
# a link is DEAD only if the host is unreachable (000) or returns 404/410. Codes
# like 403/405/429 mean the page exists but rejected the bot, and are reported
# but not failed. Written for bash 3.2 (macOS) as well as 4+ (CI): no mapfile.
set -uo pipefail

urls=()
while IFS= read -r u; do
	urls+=("$u")
done < <(
	grep -rhoE 'https?://[^ )>"`]+' --include='*.md' . |
		sed -E 's/[.,]+$//' | sort -u
)

echo "Checking ${#urls[@]} unique links from markdown..."
fail=0
for url in ${urls[@]+"${urls[@]}"}; do
	code=$(curl -sS -o /dev/null -L --max-time 25 \
		-A 'Mozilla/5.0 link-check' -w '%{http_code}' "$url" 2>/dev/null || echo 000)
	case "$code" in
	000 | 404 | 410)
		printf 'DEAD   (%s)  %s\n' "$code" "$url"
		fail=1
		;;
	2* | 3*) : ;; # live
	*) printf 'note   (%s)  %s\n' "$code" "$url" ;;
	esac
done

if [ "$fail" -eq 0 ]; then
	echo "All links reachable (no 000/404/410)."
else
	echo "Some links are dead."
fi
exit "$fail"
