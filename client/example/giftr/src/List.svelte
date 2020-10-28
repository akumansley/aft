<script>
	import client from './client.js';
	import user from './user.js';
	export let listUserId;
	import { Link } from "svelte-routing";
	import Gift from './gifts/Gift.svelte';
	import Comment from './gifts/Comment.svelte';
	import AddComment from './gifts/AddComment.svelte';

	const isOwnList = $user.id === listUserId;

	let newGift;
	async function addGift() {
		await client.api.gift.create({data: {
			description: newGift,
			user: {connect: {id: $user.id}},
		}})
		newGift = "";
		load();
	}

	let gifts = [];
	async function load() {
		gifts = await client.api.gift.findMany({
			where: {
				user: {id: listUserId},
			},
			include: {
				comments: true,
			}
		});
	}
	let loaded = load();
</script>

<style>
	.box {
		padding: .5em .75em;
	}
	.double-box {
		padding: 1em .75em;
	}
	.spacer {
		height: 1em;
		width: 1em;
	}
	.indent {
		padding-left: 3em;
	}

	.gift-box {
		background: var(--color-1-darker);
		border-radius: 5px;
		color: white;
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

<Link to="/">&larr; Back</Link>

<div class="spacer"></div>
{#await loaded then _}
<div class="page">

	{#if isOwnList}
	<div class="double-box gift-box">
		<input placeholder="Add a gift.." type=text bind:value={newGift}/> 
		<button on:click={addGift}>Add</button>
	</div>

	<div class="spacer"></div>
	{/if}

	{#if gifts.length === 0 }
	<div class="box">
		No gifts yet!
	</div>
	{:else}

	{#each gifts as gift}
	<Gift value={gift} />

	<div class="indent">
	{#each gift.comments as comment}
	<Comment value={comment} />
	{/each}
	<AddComment gift={gift} on:save={load} />
	</div>

	<div class="spacer"></div>

	{/each}
	{/if}

</div>
{/await}