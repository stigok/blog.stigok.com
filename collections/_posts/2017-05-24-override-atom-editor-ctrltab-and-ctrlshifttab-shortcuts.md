---
layout: post
title: "Override Atom editor ctrl+tab and ctrl+shift+tab shortcuts"
date: 2017-05-24 19:58:28 +0200
categories: atom
redirect_from:
  - /post/override-atom-editor-ctrltab-and-ctrlshifttab-shortcuts
---

I want Ctrl+Tab and Ctrl+Shift+Tab to quite simply go to Next and Previous pane. I don't want no smart *last pane used* functionality. I got that in my brain. Override it like this in your *keymap.cson*:

    'body':
      'ctrl-tab': 'pane:show-next-item'
      'ctrl-shift-tab': 'pane:show-previous-item'
      'ctrl-tab ^ctrl': 'unset!'
      'ctrl-shift-tab ^ctrl': 'unset!'

## References
- https://github.com/atom/settings-view/issues/124#issuecomment-261105431