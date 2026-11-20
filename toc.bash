#!/bin/bash
#
# toc.bash writes the site's table of contents, /.well-known/toc.opml, for a
# BonBon browser or cmarkbrowser, from the Markdown pages named on the command
# line.
#
#   toc.bash BASE_URL SITE_TITLE PAGE.md [PAGE.md ...] > .well-known/toc.opml
#
# BASE_URL is where the site is published, ending in a slash. Pages are listed
# in the order given, in one "Pages" group, each with its title (the title: of
# its YAML front matter, else its first "# " heading) and the address of its
# Markdown file. README.md is listed as
# index.md, the page a BonBon browser asks for at the site's address. The output
# is the same for the same input (the date is the last commit's, not today's),
# so a rebuild changes nothing in git.
#
set -eu

if [ "$#" -lt 3 ]; then
	echo "usage: toc.bash BASE_URL SITE_TITLE PAGE.md [PAGE.md ...]" >&2
	exit 2
fi
BASE_URL="$1"
SITE_TITLE="$2"
shift 2

case "${BASE_URL}" in
http://*/ | https://*/) ;;
*)
	echo "toc.bash: BASE_URL must be an http or https address ending in /: ${BASE_URL}" >&2
	exit 2
	;;
esac

# xml_escape quotes the characters that matter in an XML attribute or text. It
# uses sed: in bash 5.2 an "&" in the replacement of ${s//x/y} means "the
# matched text", which turned "<" into "<lt;".
xml_escape() {
	printf '%s' "$1" | sed -e 's/&/\&amp;/g' -e 's/</\&lt;/g' -e 's/>/\&gt;/g' -e 's/"/\&quot;/g'
}

# page_title is the title a page gives itself: the title: of its YAML front
# matter, else the text of its first "# " heading. (A man page's first heading
# is NAME, so its front matter is what names it.)
page_title() {
	local t
	t=$(awk 'NR==1 && $0!="---" {exit} NR>1 && $0=="---" {exit} NR>1 && /^title:/ {sub(/^title: */, ""); gsub(/^["'\'']|["'\'']$/, ""); print; exit}' "$1")
	if [ -z "${t}" ]; then
		t=$(sed -n 's/^# \{1,\}\(.*[^ ]\) *$/\1/p' "$1" | head -1)
	fi
	printf '%s' "${t}"
}

# The date of the last commit that touched a page, or the newest file time
# when there is no git history; never the current time.
modified() {
	local d
	d=$(git log -1 --format=%cD -- "$@" 2>/dev/null || true)
	if [ -z "${d}" ]; then
		d=$(LC_ALL=C date -u -r "$(ls -t "$@" | head -1)" '+%a, %d %b %Y %H:%M:%S GMT')
	fi
	printf '%s' "${d}"
}

printf '<?xml version="1.0" encoding="UTF-8"?>\n'
printf '<opml version="2.0">\n'
printf '\t<head>\n\t\t<title>%s</title>\n\t\t<dateModified>%s</dateModified>\n\t</head>\n' \
	"$(xml_escape "${SITE_TITLE}")" "$(modified "$@")"
printf '\t<body>\n\t\t<outline text="Pages">\n'
for page in "$@"; do
	if [ ! -f "${page}" ]; then
		echo "toc.bash: ${page}: no such page" >&2
		exit 66
	fi
	title=$(page_title "${page}")
	if [ -z "${title}" ]; then
		title="${page%.md}"
	fi
	published="${page}"
	if [ "${page}" = "README.md" ]; then
		published="index.md"
	fi
	printf '\t\t\t<outline text="%s" url="%s%s"></outline>\n' \
		"$(xml_escape "${title}")" "${BASE_URL}" "${published}"
done
printf '\t\t</outline>\n\t</body>\n</opml>\n'
