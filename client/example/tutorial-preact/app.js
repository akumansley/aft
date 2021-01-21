import {html, Component, useState, useEffect, useCallback} from 'https://unpkg.com/htm/preact/standalone.module.js'
import aft from './aft.js'


export function App (props) {
	const [user, setUser] = useState(null);
	const [loaded, setLoaded] = useState(false);

	useEffect(async () => {
		try {
			setUser(await aft.rpc.me());
		} catch {

		} finally {
			setLoaded(true);
		}
	}, []);

	if (!loaded) {
		return html``
	} else if (user === null) {
		return html`<${Login} setUser=${setUser} />`
	} else {
		return html`<${Todos} user=${user}/>`
	}
}

function Login(props) {
	const setUser = props.setUser;
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

	let form = html`
	<div class="login-form">
		<input type=email placeholder="Email" value=${email} onInput=${(e) => setEmail(e.target.value)}/>
		<input type=password placeholder="Password" value=${password} onInput=${(e) => setPassword(e.target.value)}/>
		${errorMessage && html`<div class=error>${errorMessage}</div>`}
		<button onClick=${submit} >Sign in</button>
	</div>`
	return form
}

function Todos({user}) {
	const [todos, setTodos] = useState([]);
	const [loaded, setLoaded] = useState(false);
	const [addText, setAddText] = useState("");
	const [errorMessage, setErrorMessage] = useState(null);

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

	const addTodo = useCallback(async () => {
		setErrorMessage(null)
		try {
			const todo = await aft.api.todo.create({
				data: {
					text: addText,
					user: {connect: {id: user.id}},
				}
			});
			setAddText("");
			setTodos((todos) => [...todos, todo]);
		} catch (e) {
			setErrorMessage(e.message);
		}
	}, [todos, user, addText])

	if (!loaded) {
		return html``
	}

	const signout = () => aft.rpc.logout()

	return html`<div class="todo-list">
				<div class="row">${user.email} <a onClick=${signout}>sign out</a></div>
		${todos.map(todo => {

			const toggle = async () => {
				const updated = await aft.api.todo.update({
					where: {id: todo.id},
					data: {done: !todo.done},
				})

				setTodos(todos => {
					return todos.map(t => t.id === todo.id? updated: t)
				})
			}

			return html`<${Todo} key=${todo.id} toggle=${toggle} todo=${todo} />`
		})}
		${errorMessage && html`<div class=${"row error"}>${errorMessage}</div>`}
		<div class="add-todo">
			<input type=text value=${addText} onInput=${e => setAddText(e.target.value)} />
			<button onClick=${addTodo}>Add</button>
		</div>
	</div>`
}

function Todo({todo, toggle}) {
	return html`
	<div class="row" onClick=${toggle}>
		<div>${todo.text}</div>
		<input type="checkbox" checked=${todo.done} />
	</div>`

}

