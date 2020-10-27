<script>
	import client from './client.js';
	import user from './user.js';
	export let listUserId;
	import { Link } from "svelte-routing";

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
			}
		});
	}
	let loaded = load();
</script>
<style>
	.box {
		padding: .5em .75em;
	}
	.spacer {
		height: 1em;
		width: 1em;
	}
	.page {
		background: var(--color-1-darker);
		border: 2px solid var(--color-3);
		border-radius: 5px;
		color: white;
	}
	.gift-box {
		border-bottom: 2px solid var(--color-3);
	}
	.gift-box:last-child {
		border-bottom: none;
	}
	input[type="text"] {
		background: var(--color-1);
		border-bottom: 2px solid var(--color-3);
		border-top: none;
		border-left: none;
		border-right: none;
	}
	input[type="text"]::placeholder {
		color: white;
	}
	button {
		background: var(--color-1-darker);
		border: 2px solid var(--color-3);
		border-radius: 5px;
		color: white;
	}
</style>

<Link to="/">&larr; Back</Link>

<div class="spacer"></div>
{#await loaded then _}
<div class="page">
	{#if isOwnList}
	<div class="spacer"></div>
	<div class="box gift-box">
		<input placeholder="Add a gift.." type=text bind:value={newGift}/>
		<button on:click={addGift}>Add</button>
	</div>
	{/if}

	{#if gifts.length === 0 }
	<div class="box">
		No gifts yet!
	</div>
	{:else}
	{#each gifts as gift}
	<div class="box gift-box">
		<div class="gift">
			{gift.description}
		</div>
	</div>
	{/each}
	{/if}
</div>
{/await}