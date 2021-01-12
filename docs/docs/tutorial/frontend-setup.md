---
id: frontend-setup
title: Frontend Setup
---

Create a new directory called "client":

```bash
mkdir client
cd client
```

Open a new file in that directory called "index.html", and start it off with the following contents

```html title="index.html"
<!doctype html>
<head>
	<title>Aft Tutorial</title>
</head>

<div class="container">
	<h1>Aft Tutorial</h1>
</div>
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

Open up the URL for the client in your browser, and you should see your HTML—off to a good start! 

## Adding Frameworks

For this tutorial, we're going to use [Bootstrap](https://getbootstrap.com/) to give us some sensible default styles, and Vue to give us a nice component-based UI system.

Aft doesn't come with its own frontend, so you can just as easily make an app using React, Flutter, native iOS or Android frameworks, or low-code frontend builders like [Bubble](https://bubble.io).

We're using Vue and Bootstrap mainly becuase they should be familiar to developers who have some web development experience, and becuase they can be used without a build step, which helps keep the tutorial simple.

So, with that said, let's add Bootstrap inside of our `<head>` tag:

```html title="index.html"
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" crossorigin="anonymous"> 
```

And then do the same for Vue:

```html title="index.html"
<script src="https://unpkg.com/vue@3.0.5"></script>
```

Then with those two added, we can make our UI slightly more styled. Update the contents of the "container" `div`:

```html title="index.html"
<div class="container">
	<div class="row justify-content-center align-items-center mt-4">
		<h1>Aft Tutorial</h1>
	</div>
	<div class="row justify-content-center align-items-center mt-3">
		<div id="app"></div>
	</div>
</div>
```

Refresh the client in your browser to see your updated styles!

Next, we'll get started on our app and implement our first RPC, login.
