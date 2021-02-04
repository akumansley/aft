---
id: frontend-setup
title: Frontend Setup
---

Create a new directory called "client":

```bash
mkdir client
cd client
```

Open a new file in that directory called "index.html", and start it off with a nearly empty page and a link to a stylesheet.

```html title="index.html"
<head>
	<link rel=stylesheet href="./styles.css" />
	<title>Aft Tutorial</title>
</head>
<body>
	<h1>Hello Aft</h1>
</body>
```

Then let's add our stylesheet. We'll just paste in all the styles we'll need for the tutorial at once for convenience.

```css title="styles.css"
:root {
	--page-width: 24em;
	--primary: #1266F1;
	--danger: #F93154;
	--light: #FBFBFB;
	--border: 1px solid rgba(0,0,0,.125);
}

* {
	box-sizing: border-box; 
	margin: 0;
	padding: 0;
}

a {
	color: var(--primary);
	cursor: pointer;
}

.error {
	color: var(--danger);
}

body {
	font-family: system-ui, sans-serif;
	display: flex;
	justify-content: center;
	align-items: center;
	font-size: 16px;
    font-weight: 400;
    line-height: 1.7;
}

input, button {
	border: var(--border); 
	padding: 0.25em 0.375em; 
	font: inherit;
}

button {
	background: var(--primary);
	color: white;
}

.box {
	width: var(--page-width);
	background: var(--light);
	border: var(--border);
}

.stack {
	display: grid;
	gap: .5em;
	padding: 1em;
	grid-template-columns: 1fr;
}

.row {
	padding: .5em 1em;
	border-top: var(--border);
	display: flex;
	align-items: baseline;
	flex-direction: row;
}

.row:first-child {
	border-top: none;
}

.row > *:first-child {
	flex-grow: 1;
}
```

Then restart `aft`, pointing at our client directory:

```bash
aft -db ./tutorial.dbl -authed=false -serve_dir=client
```

This time you should see aft print out two messages:

```bash
Serving client on http://localhost:8080
Serving dev on http://localhost:8081
```

Open up the URL for the client in your browser, and you should see the page title from your HTML—off to a good start! 

## Adding Frameworks

For a sophisticated frontend, we'd probably want to use a full-featured framework with some kind of build step, but to keep things easy for this tutorial, we're going to rely on a single frontend library, Preact, and not worry about supporting older browsers.

Aft doesn't come with its own frontend, so you can just as easily make an app using Vue, React, Flutter, native iOS or Android frameworks, or low-code tools like [Bubble](https://bubble.io).

Let's add our first component in a module script tag, and drop the body tag—we don't need it anymore!

```html title="index.html"
<head>
	<link rel=stylesheet href="./styles.css" />
	<script type=module>
		import {html, render} from 'https://unpkg.com/htm/preact/standalone.module.js'
		function App() {
			return html`<h1>Hello Aft</h1>`
		}
		render(html`<${App} />`, document.body);
	</script>
</head>
```

For those who've done some React development before, the use of template literals fills a similar role to JSX, but with no build step required.

Next, we'll get started on our app and implement our first RPC, login.
