---
layout: post
title:  "Helpful bash aliases for git"
date:   2021-12-31 11:42:29 +0100
categories: git bash
excerpt: I use some aliases in bash to operate git faster in my projects.
#proccessors: pymd
---

I have this in my .bashrc

```shell
alias g='git status'
alias gap='git add -p'
alias gc='git commit'
alias gcf='git commit --fixup'
alias gri='git rebase -i'
alias gca='git commit --amend --no-edit'
```

This makes me so much faster in the shell when using git.
Let's go through them one by one.

## `g`

Whenever I enter a repo, I run `g` to see if I'm in a clean workspace.
I try to do this before exiting too, so I don't leave garbage around for
the next time I come around. I hate entering a dirty workspace, because
I usually have no rememberance of what I was doing. And without changes
being comitted, the changes does not have a message (or a context).

```
$ g
On branch main

No commits yet

Untracked files:
  (use "git add <file>..." to include in what will be committed)
  index.html

nothing added to commit but untracked files present (use "git add" to track)
```

## `gap`

This helps keeping my commits clean. Before I use it, I might do a `git diff`
to get a quick overview over all the changes. But even if I want to commit all
of it, I usually use `gap`, out of habit, to make sure I see all the changes
I am comitting.

This is the equivalent of `git add --patch` which lets you stage smaller "hunks"
of the diff interactively.

If I want to pick out unstaged hunks from a single file:

```
$ gap
diff --git a/index.html b/index.html
index 323fae0..e17c29d 100644
--- a/index.html
+++ b/index.html
@@ -1 +1,5 @@
+<h1>a new header</h1>
+
 foobar
+
+<footer>a new footer</footer>
(1/1) Stage this hunk [y,n,q,a,d,s,e,?]? s
Split into 2 hunks.
@@ -1 +1,3 @@
+<h1>a new header</h1>
+
 foobar
(1/2) Stage this hunk [y,n,q,a,d,j,J,g,/,e,?]?
```

## `gc`

As simple as `git commit`. It just helps me write less. If I just write `gc` an
interactive window pops up where I can write a longer commit message.
If I'm just going to write a single line, I append the message argument:

```
$ gc -m 'add header and footer'
[main e1bf653] add header and footer
 1 file changed, 4 insertions(+)
```

## `gcf`

This is when I want to make some changes to a previous commit, or merge two commits
into one, without updating the commit message of the destination commit.
If the destination commit is the previous one, I would rather use `gca`, but if it's
further back in the log, I will use `gcf`.

First I have to find the sha of the destination commit I want to *fixup*. I use `git log`
for this, or even `git commit --oneline` for a more compact output.

```
$ gcf 13c018d
[main 74a7a19] fixup! first
 1 file changed, 2 insertions(+), 2 deletions(-)
```

The git log now has a new commit with a `fixup` instruction. This can be processed
by `git rebase --interactive --autosquash`, which re-orders the fixup-commits
automatically.

## `gri`

This is my favorite git function nowadays. I'm a perfectionist who likes to have as
clean commit history as possible before pushing to the remotes and submitting PRs.

This is almost always preceded by a `gcf` command, as I run this after I've fixed up
a previous commit. However, I need an extra argument to tell git how far back in the
history I want to go. One way to do this is to count how many commits back the target
of the fixup was using `git log`, then running

```
$ gcf @~7
```

But it is cumbersome to have to count every time. A much neater alternative is to
reference *the parent of a sha*. I already know the sha, as it's the same argument as used
in the `gcf` command, and git has an operator `~` we can use to specify the parent.

```
$ gcf 13c018d
$ gri 13c018d~
```

Now I get an interactive window which I can often just close immediately to get the
commits fixed up (merged).

I have `rebase.autosquash` enabled by default in my global git config, so this alias
doesn't need the `--autosquash` argument.

## `gca`

When I'm just editing the previous commit, I can amend the commit with `git commit --amend`.
I usually don't want to update the commit message, so I can just run `gca` and be done
with it.

```
$ gca
[main 3aee9cd] add header and footer
 Date: Fri Dec 31 12:08:40 2021 +0100
 3 files changed, 4 insertions(+)
 create mode 100644 bar
 create mode 100644 foo
```

And if I *wanted* to update the message here, I can run `gca --edit`.
