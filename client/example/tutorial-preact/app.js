import {html, useState, useEffect, useCallback} from 'https://unpkg.com/htm/preact/standalone.module.js'
import aft from './aft.js'


export function App (props) {
	const [user, setUser] = useState(null);
	const [loaded, setLoaded] = useState(false);

	useEffect(async () => {
		try {
			setUser(await aft.rpc.me());
		} catch {} finally {
			setLoaded(true);
		}
	}, []);

	if (!loaded) {
		return html``
	} else if (user === null) {
		return html`<${Login} setUser=${setUser} />`
	} else {
		return html`<${Todos} user=${user} setUser=${setUser}/>`
	}
}

function Login({setUser}) {
	const [errorMessage, setErrorMessage] = useState(null);
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");

	const submit = useCallback(async () => {
		setErrorMessage(null);
		try {
			const user = await aft.rpc.login({
				email: email,
				password: password,
			});
			setUser(user);
		} catch (e) {
			setErrorMessage(e.message);
		}
	}, [email, password])

	return html`
	<div class="box stack">
		<input type=email placeholder="Email" 
			value=${email} 
			onInput=${(e) => setEmail(e.target.value)}/>
		<input type=password placeholder="Password" 
			value=${password} 
			onInput=${(e) => setPassword(e.target.value)}/>
		${errorMessage && html`<div class=error>${errorMessage}</div>`}
		<button onClick=${submit}>Sign in</button>
	</div>`
}

function Todos({user, setUser}) {
	const [todos, setTodos] = useState([]);
	const [loaded, setLoaded] = useState(false);

	useEffect(async () => {
		try {
			setTodos(await aft.api.todo.findMany({
				where: {
					done: false,
				}
			}));
		} catch {
		} finally {
			setLoaded(true);
		}
	}, []);


	if (!loaded) {
		return html``
	}

	const signout = async () => {
		await aft.rpc.logout();
		setUser(null);
	}

	return html`
	<div class="box">
		<div class="row">
			<div><b>Todos</b></div><a onClick=${signout}>Sign out</a>
		</div>
		${todos.map(todo => {
			const toggle = useToggle(setTodos, todo);
			return html`<${Todo} key=${todo.id} toggle=${toggle} todo=${todo}/>`
		})}
		<${AddTodo} user=${user} setTodos=${setTodos} todos=${todos}/>
	</div>`
}

function useToggle(setTodos, todo) {
	const toggle = async () => {
		const updated = await aft.api.todo.update({
			where: {id: todo.id},
			data: {done: !todo.done},
		})
		setTodos(todos => {
			return todos.map(t => t.id === todo.id? updated: t)
		})
	}
	return toggle
}

function AddTodo({user, setTodos, todos}) {
	const [text, setText] = useState("");
	const [errorMessage, setErrorMessage] = useState(null);

	const addTodo = useCallback(async () => {
		setErrorMessage(null)
		try {
			const todo = await aft.api.todo.create({
				data: {
					text: text,
					user: {connect: {id: user.id}},
				}
			});
			setText("");
			setTodos((todos) => [...todos, todo]);
		} catch (e) {
			setErrorMessage(e.message);
		}
	}, [todos, user, text])

	return html`
	${errorMessage && html`<div class="row error">${errorMessage}</div>`}
	<div class="row">
		<input type=text
			value=${text} 
			onInput=${e => setText(e.target.value)} />
		<button onClick=${addTodo}>Add</button>
	</div>`
}

function Todo({todo, toggle}) {
	return html`
	<div class="row">
		<div>${todo.text}</div>
		<input type=checkbox onClick=${toggle} checked=${todo.done} />
	</div>`
}

