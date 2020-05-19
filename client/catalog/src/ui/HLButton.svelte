<script>
	export let value;
</script>
<style>
button {
		-webkit-appearance: none;
		-moz-appearance: none;
		appearance: none;
		background-color: var(--background);
		border-radius: 0;
		margin: 0;
		color: var(--text-color-darker);
		border-color: var(--border-color);
}
button:active {
		-webkit-appearance: none;
		-moz-appearance: none;
		appearance: none;
		background-color: var(--background-highlight);
		border-radius: 0;
		margin: 0;
		color: inherit;
		border-color: var(--border-color);
}

</style>
<button on:click>
	<slot />
</button>


