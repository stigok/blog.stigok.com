---
layout: post
title:  "Verify a string consists strictly of a substring of itself in javascript"
date:   2019-04-30 23:15:00 +0200
categories: javascript
---

I stumbled over a Stack Overflow question titled *"[How do I check if a string is entirely made of the same substring?][so-question]"*.
Before reading through the OP's current solution, I thought I'd give it a shot myself. I felt very clever, possibly even unique at thought, but I ended up with a solution almost identical to his.

```javascript
function check (str) {
  // Length of the substring cannot be larger than half of the str length
  const maxlen = Math.floor(str.length / 2) || 1;

  for (let i = 1; i <= maxlen; i++) {
    // Skip indices that doesn't divide with the string length
    if (str.length % i !== 0) continue;

    // Remove all substr occurences and check the returned string's length
    const substr = str.substring(0, i);
    if (str.replace(new RegExp(substr, 'g'), '').length === 0)
      return true;
  }

  return false;
}
```

Then I continued to read on in the answers section to find out that there exist theorems for this very problem. Below quotes from SO user [templatetypedef][]

> A string consists of the same pattern repeated multiple times if and only if the string is a nontrivial rotation of itself.

Then follows up with another one

>  If x and y are strings of the same length, then x is a rotation of y if and only if x is a substring of yy.

He then followed up with the following solution for the `check` function

```javascript
function check(str) {
    return (str + str).indexOf(str, 1) !== str.length;
}
```

I was humbled and amazed.
I urge you to go read the [whole question][so-question] and the different answers.

[so-question]: https://stackoverflow.com/questions/55823298/how-do-i-check-if-a-string-is-entirely-made-of-the-same-substring
[templatetypedef]: https://stackoverflow.com/users/501557/templatetypedef
