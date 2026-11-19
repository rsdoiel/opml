<#
.SYNOPSIS
PowerShell script for building the opml website pages with pandoc. The
counterpart of website.mak: keep the page list the same in both.
#>
$project = "opml"
Write-Output "Building website for ${project}"
$pandoc = Get-Command pandoc -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Source
if (-not $pandoc) {
    Write-Error "website: pandoc not found"
    exit 66
}

$programs = @("opml2json", "opml2urls", "opmlcat", "opmlsort", "urls2opml")
$sitePages = @("README", "about", "INSTALL", "INSTALL_NOTES_macOS", "INSTALL_NOTES_Windows", "NOTES", "user-manual") +
    ($programs | ForEach-Object { "$_.1" })

foreach ($page in $sitePages) {
    $out = if ($page -eq "README") { "index.html" } else { "$page.html" }
    & $pandoc "--metadata" "title=$page" "-s" "--to" "html5" "$page.md" "-o" $out `
        "--lua-filter=links-to-html.lua" `
        "--lua-filter=add-col-scope.lua" `
        "--template=page.tmpl"
    if ($LASTEXITCODE -ne 0) {
        Write-Error "website: pandoc failed on $page.md"
        exit $LASTEXITCODE
    }
}
