<#
generated with CMTools 0.0.1 7534e43, then fixed by hand (see below)

.SYNOPSIS
Publish script for GitHub pages. It expects the gh-pages branch to already exist.

.DESCRIPTION
Every step is checked, as in publish.bash. The script used to rely on `git pull`
having a default way to reconcile diverged branches, and carried on after a failed
pull to push nothing. Now each pull says --no-rebase --no-edit, a commit is made
only when there is something to commit, the first failure stops the script and
names the step, and the push is confirmed by comparing gh-pages with
origin/gh-pages. A merge conflict leaves you on gh-pages in the middle of the
merge, to resolve by hand.

This port of publish.bash has not been run: there is no pwsh on the machine it
was written on.
#>

function Stop-Publish([string]$Step) {
    [Console]::Error.WriteLine("publish stopped: $Step")
    [Console]::Error.WriteLine("You are on branch $(git branch --show-current).")
    exit 1
}

# Invoke-Git runs git and stops the script if it fails.
function Invoke-Git([string]$Step, [string[]]$GitArgs) {
    & git @GitArgs
    if ($LASTEXITCODE -ne 0) { Stop-Publish $Step }
}

# Save-All commits tracked changes, if there are any.
function Save-All([string]$Message) {
    git diff --quiet HEAD --
    if ($LASTEXITCODE -ne 0) { Invoke-Git "git commit" @("commit", "-am", $Message) }
}

$workingBranch = git branch --show-current
if ($workingBranch -eq "gh-pages") {
    Save-All "publishing to gh-pages branch"
    Invoke-Git "git push origin gh-pages" @("push", "origin", "gh-pages")
} else {
    Write-Output "You're in $workingBranch branch"
    Write-Output "You need to pull in changes to the gh-pages branch to publish"
    $yesNo = Read-Host "Process Y/n"
    if ($yesNo -eq "Y" -or $yesNo -eq "y") {
        Write-Output "Committing and pushing to $workingBranch"
        Save-All "committing to $workingBranch"
        Invoke-Git "git push origin $workingBranch" @("push", "origin", $workingBranch)
        Write-Output "Changing branches from $workingBranch to gh-pages"
        Invoke-Git "git checkout gh-pages" @("checkout", "gh-pages")
        Write-Output "Merging changes from origin gh-pages"
        Invoke-Git "git pull origin gh-pages" @("pull", "--no-rebase", "--no-edit", "origin", "gh-pages")
        Write-Output "Merging changes from $workingBranch into gh-pages"
        Invoke-Git "git pull origin $workingBranch (resolve the conflict, commit, then git push origin gh-pages)" @("pull", "--no-rebase", "--no-edit", "origin", $workingBranch)
        Write-Output "Pushing changes up and publishing"
        Invoke-Git "git push origin gh-pages" @("push", "origin", "gh-pages")
        git fetch origin gh-pages
        if ((git rev-parse gh-pages) -ne (git rev-parse origin/gh-pages)) {
            Stop-Publish "gh-pages and origin/gh-pages differ after the push"
        }
        Write-Output "Published gh-pages at $(git rev-parse --short gh-pages)"
        Write-Output "Changing back to your working branch $workingBranch"
        Invoke-Git "git checkout $workingBranch" @("checkout", $workingBranch)
    }
}
