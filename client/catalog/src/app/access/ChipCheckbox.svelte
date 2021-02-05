<script>
	export let checked = false;
</script>
<style>
	.hl-check {
		display: none;
	}
	.hl-label.checked {
		background: var(--background-highlight);
		color: var(--text-color);
	}
	.hl-label:hover {
		background: var(--background-highlight);
		color: var(--text-color);
	}
	.hl-label {
		background: var(--background);
		cursor: pointer;
		border: 1px solid var(--border-color);
		display: flex;
		flex-direction: row;
		align-items: center;
		padding: .5em 1em;
		font-size: var(--scale--1);
		color: var(--text-color-darker);
	}
	.check {
		margin-right: .5em;
	}
</style>

<label class:checked="{checked}"  class="hl-label">
	{#if checked}
	<svg class="check" xmlns="http://www.w3.org/2000/svg" width="1em" height="13px" viewBox="0 0 24 24"><path style="fill: var(--text-color)" d="M9 22l-10-10.598 2.798-2.859 7.149 7.473 13.144-14.016 2.909 2.806z"/></svg>
	{/if}
	<input class="hl-check" on:change type=checkbox bind:checked={checked} />
	<slot/>
</label>
