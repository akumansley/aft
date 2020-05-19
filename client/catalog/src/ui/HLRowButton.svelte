<style>
	.hl-row {
		padding-left: calc(.5em + 1px);
		padding-top: calc(.25em + 1px);
		padding-bottom: calc(.25em + 1px);
		color: var(--text-color-darker);
		cursor: pointer;
	}
	.hl-row:hover {
		padding-left: .5em;
		padding-top: .25em;
		padding-bottom: .25em;
		border: 1px solid var(--border-color);
		color: var(--text-color);
		background: var(--background-highlight);
	}
</style>

<div on:click class="hl-row">
	<slot />
</div>

