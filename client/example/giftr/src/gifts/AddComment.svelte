<script>
	import {createEventDispatcher} from 'svelte';
	import client from '../client.js';
	import user from '../user.js';
	let value;
	export let gift;
	function init() {
		value = {
			text: "",
			claim: false,
		}
	}
	async function addComment() {
		await client.api.comment.create({data: {
			author: {connect: {id: $user.id}},
			gift: {connect: {id: gift.id}},
			text: value.text,
			claim: value.claim,
		}})

		init();
		dispatch("save");
	}
	const dispatch = createEventDispatcher();
	init();
</script>
<style>
	.box {
		padding: .75em .75em;
	}
	.comment-box {
		background: var(--color-1-darker);
		color: white;
		border-radius: 0 0 5px 5px;
	}
	.hform {
		display: flex;
		flex-direction: row;
		align-items: baseline;
	}
	.spacer {
		width: 1em;
		height: 1em;
	}
	.fill {
		flex-grow: 1;
	}

	input[type="text"] {
		background: var(--color-1);
		border-bottom: 2px solid var(--color-3);
		border-top: none;
		border-left: none;
		border-right: none;
		margin: 0;
		color: white;
	}
	input[type="text"]::placeholder {
		color: white;
	}
	button {
		background: var(--color-3);
		border-radius: 5px;
		border: none;
		padding: calc(.35em + 2px);
		color: var(--color-1-darker);
		margin: 0;
	}
</style>
<div class="box comment-box hform">
	<input placeholder="Add comment.." class="fill" type=text bind:value={value.text}/>
	<div class="spacer"></div>
	<label>
		<input type=checkbox bind:checked={value.claim}/>
		Claim
	</label>
	<div class="spacer"></div>
	<button on:click={addComment}>Add</button>
</div>
