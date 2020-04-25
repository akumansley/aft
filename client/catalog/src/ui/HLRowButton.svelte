<style>
	.hl-row {
		border-bottom: 1px solid var(--border-color);
		padding-left: .5em;
		padding-top: .25em;
		padding-bottom: .25em;
		color: var(--text-color-darker);
		cursor: pointer;
	}
	.hl-row:hover {
		color: var(--text-color);
		background: var(--background-highlight);
	}
	.hl-row:last-child {
		border-bottom: none;
	}
</style>

<div on:click class="hl-row">
	<slot />
</div>

