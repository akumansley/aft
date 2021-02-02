import {html} from 'https://unpkg.com/htm/preact/standalone.module.js'

export function Box({children}) {
	const styles = {
		width: "var(--page-width)",
		display: "grid",
		gap: ".5em",
		padding: "1.5em 1em",
		background: "var(--light)",
		border: "1px solid var(--border)",
	}
	return html`<div style=${styles}>${children}</div>`
}

export function Button({onClick, children}) {
	const styles = {
		background: "var(--primary)",
		color: "white",
	}
	return html`
	<button onClick=${onClick} style=${styles}>${children}</button>
	`
}

export function List({children}) {
	const styles = {
		width: "var(--page-width)",
		borderBottom: "1px solid var(--border)",
		background: "var(--light)",
	}
	return html`<div style=${styles}>${children}</div>`
}

export function Row({children, padding}) {
	const styles = {
		"padding": padding || ".5em 1em",
		border: "solid var(--border)",
		borderWidth: "1px 1px 0 1px",
		display: "flex",
		alignItems: "baseline",
		flexDirection: "row"
	}
	return html`<div style=${styles}>${children}</div>`
}

export const Fill = ({children}) => html`<div style="flex-grow: 1">${children}</div>`
