---
id: running
title: Running Aft
---

If you haven't installed Aft yet, read the [Getting Started](../getting-started) page.

Make a new directory somewhere on your computer for the tutorial:

```bash
mkdir aft-tutorial
cd aft-tutorial
```

If you have Aft installed, run it:

```bash
aft -db tutorial.dbl -authed=false
```

You should see it print out:

```bash
Serving dev on http://localhost:8081
```

Open that link in your browser. You should see a page labeled `Schema`, and a listing of system models.

Before we set up our backend, let's get a client started.