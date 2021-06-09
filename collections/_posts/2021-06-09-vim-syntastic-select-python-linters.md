---
layout: post
title:  "Choose what Python linters are used for Syntastic plugin in Vim"
date:   2021-06-09 10:27:50 +0200
categories: vim python
excerpt: Some of my files had myriads of errors I didn't care about, even if my local flake8 config was set to ignore them.
#proccessors: pymd
---

## Preface

For a long time I've been annoyed with too many lint errors in my Python files coming from
[syntastic][]. Time to get this cleaned up.

## Selecting active linters

Synastic has a plugin loader system that tries to enable all registered linters that are
present on your system. For me, this means that when I have the binaries flake8, pylint
and pep8 available on my `PATH`, Syntastic tries to run all of them.

You can find a list of all the linter it looks for under [syntastic/syntax_checkers/python](https://github.com/vim-syntastic/syntastic/tree/master/syntax_checkers/python).

Now, you can open a Python file containing various errors and try out all the linters to
see what fits you best. Toggle selected linters at a time by using the below line:

```vim
let g:syntastic_python_checkers=['<plugin name>']
```

Now you can trigger a relint of the document

```vim
:SyntasticCheck
```

I found that for myself, I prefer to only run flake8 for now. It makes it simple for me to
configure it on a per-project basis using *.flake8* files.
That leaves me with this config:

```vim
let g:syntastic_python_checkers=['flake8']
```


## References
- <https://github.com/vim-syntastic/syntastic>
- <https://stackoverflow.com/a/23105873/90674>

[syntastic]: https://github.com/vim-syntastic/syntastic
