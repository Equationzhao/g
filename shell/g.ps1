Remove-Alias -Name ls
function ls {
    param(
        [Switch] $path
    )
    g $args
}

function ll {
    g -1 --perm --icons --time --group --owner --size --title
}

function la {
    g --show-hidden
}

function l {
    g --perm --icons --time --group --owner --size --title --show-hidden
}

# `echo $profile` in PowerShell to find your profile path
# add the following line to your profile
# Invoke-Expression (& { (g --init powershell | Out-String) })
# if you've already remove alias `ls`, you can comment out the first line
# and paste the content to your profile manually
