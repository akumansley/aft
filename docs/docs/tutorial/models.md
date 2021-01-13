---
id: models
title: Models
---

Let's start by adding another store in `store.js` to hold our Todo data, and add some test data to it

```js title="store.js"
export const userstore = Vue.reactive({value: null})


const testData = [{
	text: "Todo",
	done: false,
}, {
	text: "Hello",
	done: true,
}]

export const todoStore = Vue.reactive({todos: testData})
```

And then a new component to display our todo data, `todos.js`.

```js title="todos.js"

export const Todos = {
	data() {
		return {
			todoStore: todoStore,
		}
	},
	template: `
	<div class="list-group" style="width: 20rem;">
		<div v-for="todo in todoStore.todos">{{todo.text}} - {{todo.done}}</div>
	</div>
	`,
}
```

And finally, let's connect it in to our app.


```js title="app.js"
import {api} from './api.js'
import {Login} from './login.js'
import {Todos} from './todos.js'
import {userStore} from './store.js'

const App = {
	components: {
		'login-view': Login,
		'todos': Todos,
	},
	template: `
	<div v-if="!loaded"></div>
	<login-view v-else-if="userStore.value === null" />
	<todos v-else />
	`,
	created() {
		api.me().then((user) => {
			userStore.value = user;
			this.loaded = true
		}, () => {
			this.loaded = true
		});
	},
	data() {
		return {
			userStore: userStore,
			loaded: false,
		}
	}
}

export const app = Vue.createApp(App)
```

## Adding Models

Now let's create a basic Todo object on our backend. We'll have two attributes, `text`, a string which is the text of the Todo, and `done` a boolean indicating whether it's done or not. We'll also add a relationship to a `user` object—the owner of the Todo.

Switch back over to Aft, and navigate to the **Schema** section.

Click the **Add Model** button at the top of the screen. 

Fill in the text field that says *Model name..* with `todo`,

Then click the **add** button under Attributes and fill in our first attribute.

For *Attribute name..* type `text`, and select `String` from the dropdown on the right, indicating its type.

Click **add** again, and this time name the attribute `done`, and select `Bool` from the type dropdown.

Now let's add our relationship to `user` by clicking the **add** button under Relationships.

Set the *Relationship name..* to `user` and select `User` from the dropdown. That tells aft that this property points to a User object.

You can leave the multiple box unchecked, since this relationship will only be to a single user, rather than a list of users.

Once you've done that, click **Save** next to the model name at the top of the page, and you're all done!


## Reading data

Aft automatically adds our new Todo model to the API, so lets try and read some data from it.

First, navigate over to **Terminal** in Aft and run the following function to add a Todo.

```python
def main(aft):
    return aft.api.create("todo", {"data": {
    		"text":"connect the backend", 
	    	"done": False, 
	    	"user": {"connect": {"email": "user@example.come"}}
    	}})
```

Now let's extend our client API slighly, by adding a loadTodos method.

```js title="api.js"
export const api = {
	...
	async loadTodos() {
		return call('api', 'todo.findMany', {});
	},
	...
}
```

And finally, let's update our Todos component to call it!


```js title="todos.js"

export const Todos = {
	data() {
		return {
			loaded: false,
			todoStore: todoStore,
		}
	},
	created() {
		api.loadTodos().then((todos) => {
			todoStore.todos = todos;
			this.loaded = true;
		}, (err) => {
			this.error = err;
			this.loaded = true;
		});
	},
	template: `
	<div v-if="loaded" class="list-group" style="width: 20rem;">
		<div v-for="todo in todoStore.todos">{{todo.text}} - {{todo.done}}</div>
	</div>
	`,
}

```

We can also clean up the test data from `store.js`, now that we've got real server-side data.


```js title="store.js"
export const userstore = Vue.reactive({value: null})
export const todoStore = Vue.reactive({todos: []})
```

In the next section, we'll look at adding some data!
