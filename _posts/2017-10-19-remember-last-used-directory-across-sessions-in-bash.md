---
layout: post
title: "Remember last used directory across sessions in bash"
date: 2017-10-19 13:28:14 +0200
categories: bash
redirect_from:
  - /post/remember-last-used-directory-across-sessions-in-bash
---

Stumbled around to replace a function I used to have in `zsh` which remember my last used directory. I recently went back to using `bash` in yet another step towards a more simpler stack (and life).

Anyway, I ended up appending the below to my `.bashrc`:

    # Rembember last used directory
    LAST_DIR_FILE="/tmp/.lastdir"
    function cd() {
      builtin cd "$@"
      pwd > $LAST_DIR_FILE
    }
    if [ -f "$LAST_DIR_FILE" ]; then
      cd `< $LAST_DIR_FILE` 
    fi

The `builtin` command prevents endless recursion when calling `cd` by calling on the *actual* builtin function

## References
- https://unix.stackexchange.com/questions/102746/how-to-invoke-a-shell-built-in-explicitly