-- links-to-html.lua rewrites links to the site's own Markdown pages so the HTML
-- version links to the HTML version: "user_manual.md#keys" becomes
-- "user_manual.html#keys". The Markdown source keeps .md links, because the
-- Markdown is what a BonBon browser and cmarkbrowser read.
--
-- Only a relative link (no scheme, no leading "//" or "/") whose path ends in
-- .md or .markdown is changed; an address such as https://example.org/a.md
-- points at somebody else's page and is left alone.
function Link(el)
  local target = el.target
  if target:match("^[%a][%w+.-]*:") or target:match("^/") or target:match("^#") then
    return nil
  end
  local path, rest = target:match("^([^#?]*)(.*)$")
  local base = path:match("^(.*)%.md$") or path:match("^(.*)%.markdown$")
  if base == nil then
    return nil
  end
  if base == "README" then
    base = "index"
  end
  el.target = base .. ".html" .. rest
  return el
end
