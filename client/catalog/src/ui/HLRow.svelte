<style>
	.hl-row {
		border-bottom: 1px solid var(--border-color);
		padding-left: .5em;
		padding-top: .25em;
		padding-bottom: .25em;
	}
	.hl-row:hover {
		background: var(--background-highlight);
	}
	.hl-row:last-child {
		border-bottom: none;
	}
</style>

<div class="hl-row">
	<slot />
</div>

