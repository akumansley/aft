import {html, useState, useEffect, useCallback} from 'https://unpkg.com/htm/preact/standalone.module.js'
import aft from './aft.js'
import {Box, Row, Fill, Button, List} from './ui.js'


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

	return html`
	<${Box}>
		<input type=email placeholder="Email" 
			value=${email} 
			onInput=${(e) => setEmail(e.target.value)}/>
		<input type=password placeholder="Password" 
			value=${password} 
			onInput=${(e) => setPassword(e.target.value)}/>
		${errorMessage && html`<div class=error>${errorMessage}</div>`}
		<${Button} onClick=${submit}>Sign in<//>
	</${Box}>`
}

function Todos({user, setUser}) {
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

	const signout = async () => {
		await aft.rpc.logout();
		setUser(null);
	}

	return html`<${List}>
	<${Row}>
		<${Fill}><b>Todos</b><//><a onClick=${signout}>Sign out</a>
	<//>
	${todos.map(todo => {
		const toggle = useToggle(setTodos, todo);
		return html`<${Todo} key=${todo.id} toggle=${toggle} todo=${todo} />`
	})}
	${errorMessage && html`<div class=${"row error"}>${errorMessage}</div>`}
	<${AddTodo} addText=${addText} setAddText=${setAddText} addTodo=${addTodo}/>
	<//>`
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

function AddTodo({addText, setAddText, addTodo}) {
	return html`<${Row} padding=${".5em"}>
		<input style="flex-grow: 1; margin-right:.5em" 
			type=text 
			value=${addText} 
			onInput=${e => setAddText(e.target.value)} />
		<${Button} onClick=${addTodo}>Add<//>
	<//>`
}

function Todo({todo, toggle}) {
	return html`<${Row}>
		<${Fill}>${todo.text}<//>
		<input type="checkbox" onClick=${toggle} checked=${todo.done} />
	<//>`
}

