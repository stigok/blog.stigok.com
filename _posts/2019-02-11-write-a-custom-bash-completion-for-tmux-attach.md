---
layout: post
title:  "Write a custom bash completion script for tmux attach"
date:   2019-02-11 02:22:16 +0100
categories: tmux bash
---

I have a script I use to attach to existing, or create a new, tmux session.
`~/bin/tat` looks like the following

```bash
#!/bin/bash
# List existing tmux sessions, attach to existing, or create a new session
# stigok, feb 2019

if [ -z $1 ]; then
  tmux list-sessions
else
  tmux new -As $1
fi
```

Then I wanted to be able to tab-complete existing session names, instead of
first having to write `tmux ls` (or `tat`) to get list of open sessions.
This led to `/etc/bash_completion.d/tat`:

```bash
# stigok, feb 2019
_tat()
{
    local cur opts arglen
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"

    # All existing tmux session names
    opts="$(tmux ls -F '#{session_name}')"

    # The `tat` binary should only accept a single argument.
    # Only trigger completion on the first arg (after the binary name itself)
    arglen=${#COMP_WORDS[@]}
    if [ $arglen -eq 2 ]; then
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        return 0
    fi
}
complete -F _tat tat
```

Works beautifully. Now I wonder how I should package this in order to keep the
completion code closer to the shell script itself...


## References
- `man tmux`
- https://unix.stackexchange.com/questions/1800/how-to-specify-a-custom-autocomplete-for-specific-commands
- https://debian-administration.org/article/317/An_introduction_to_bash_completion_part_2
- https://www.cyberciti.biz/faq/finding-bash-shell-array-length-elements/
