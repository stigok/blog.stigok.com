---
layout: post
title:  "Dynamic native browser autocompletion with JavaScript and Vue 3"
date:   2022-06-09 00:26:13 +0200
categories: html, javascript, vue
excerpt: Dynamically update native browser autocomplete list with custom data in Vue 3
#proccessors: pymd
---

## Preface

Some time ago, I noticed that it was possible to load a predefined set of autocomplete
options for HTML `<input type="search">` using the `list` attribute.
By referencing a `<datalist>` element by its `id` we can choose what the
browser will suggest for autocompletion while typing.

```html
<input type="text" list="myDataList">
<datalist id="myDataList">
  <option value="suggestion 1">
  <option value="suggestion 2">
</datalist>
```

I wanted to see if this list could be updated dynamically while typing.
This can definitely be done with vanilla JavaScript, but since I'm interested
in getting to know Vue 3 better, I'm using that. The concepts, however, will
still be useful for those who wants to go vanilla or with other frameworks.

## Setup
### Input field

For Vue we need a container for our app, so that we can control what's inside
of it. I'm going to use a `<form>`. This container is given an `id` so that we
can reference it when we configure the Vue app later on.

For this example I will be simulating a search field. An `<input>` element
with `type="search"` is a very semantically correct choice.

```html
<form id="app">
  <input type="search" list="searchSuggestions" autocomplete="off" @input="onchange">
</form>
```

- `list="searchSuggestions"` references a `<datalist>` element in the same page
  by its `id` (we'll get to that).
- `autocomplete="off"` is used to avoid giving suggestions to the client's previous
  form inputs, that may be unrelated to our demo page -- basically preventing all
  other suggestions but our own to be shown.
- `@input="onchange"` will invoke `onchange` on every input event for the field

### Datalist

The `<datalist>` will hold the autocomplete suggestions. Using `v-for` here
will iterate over all the items in the `suggestions` array in our Vue app
instance so that an `<option>` element will be created for each item. A single
`<option>` represents a single suggestion.

```html
<datalist id="searchSuggestions">
  <!-- Render an option element for each suggested word -->
  <option v-for="s in suggestions" :value="s">
</datalist>
```

### Vue.js app

I'll let this code and the included comments speak for themselves.

```javascript
Vue.createApp({
  // Sets up the intial state
  data: function() {
    return {
      suggestions: []
    }
  },
  // App instance methods
  methods: {
    // Returns matching words for the provided query string.
    // This method should probably go fetch something from the database
    // when used in the real world.
    getMatchingWords(s) {
      const database = ["foo", "bar", "baz", "fourty two", "anaheim"]
      const matches = database.filter(words => words.includes(s))
      return Promise.resolve(matches)
    },
    // Updates the variable that the datalist is rendered from.
    updateSuggestions(query) {
      this.getMatchingWords(query).then(words => {
        // Must splice the existing array to avoid overwriting the
        // array Vue is observing (i.e. the one returned from data()).
        this.suggestions.splice(0, this.suggestions.length, ...words)
      })
    },
    // Fired for every input event because of the `@input` attribute
    // on the form element.
    onchange(event) {
      // Get the value of the search input
      const query = event.target.value

      // Avoid firing the update function excessively while
      // typing fast by "debouncing" the function.
      _.debounce(() => this.updateSuggestions(query), 300)()
    }
  }
}).mount('#app')
```

### Complete solution

```html
<!DOCTYPE html>

<script src="https://unpkg.com/lodash@4.17.21/lodash.min.js"></script>
<script src="https://unpkg.com/vue@3"></script>

<form id="app">
  <datalist id="searchSuggestions">
    <!-- Render an option element for each suggested word -->
    <option v-for="s in suggestions" :value="s">
  </datalist>

  <input type="search" list="searchSuggestions" autocomplete="off" @input="onchange">
</form>

<script>
  'use strict'

  Vue.createApp({
    // Sets up the intial state
    data: function() {
      return {
        suggestions: []
      }
    },
    // App instance methods
    methods: {
      // Returns matching words for the provided query string.
      // This method should probably go fetch something from the database
      // when used in the real world.
      getMatchingWords(s) {
        const database = ["foo", "bar", "baz", "fourty two", "anaheim"]
        const matches = database.filter(words => words.includes(s))
        return Promise.resolve(matches)
      },
      // Updates the variable that the datalist is rendered from.
      updateSuggestions(query) {
        this.getMatchingWords(query).then(words => {
          // Must splice the existing array to avoid overwriting the
          // array Vue is observing (i.e. the one returned from data()).
          this.suggestions.splice(0, this.suggestions.length, ...words)
        })
      },
      // Fired for every input event because of the `@input` attribute
      // on the form element.
      onchange(event) {
        // Get the value of the search input
        const query = event.target.value

        // Avoid firing the update function excessively while
        // typing fast by "debouncing" the function.
        _.debounce(() => this.updateSuggestions(query), 300)()
      }
    }
  }).mount('#app')
</script>
```


## References
- <https://vuejs.org/guide/essentials/list.html>
- <https://developer.mozilla.org/en-US/docs/Web/HTML/Element/datalist>
- <https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/search>
