<script>
	export let value = null;
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();
	function dispatchSelect(e) {
		dispatch('change');
	}
</script>
<style>
	select {
		-webkit-appearance: none;
		-moz-appearance: none;
		appearance: none;
		background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='292.4' height='292.4'%3E%3Cpath fill='%23e4e1e8' d='M287 69.4a17.6 17.6 0 0 0-13-5.4H18.4c-5 0-9.3 1.8-12.9 5.4A17.6 17.6 0 0 0 0 82.2c0 5 1.8 9.3 5.4 12.9l128 127.9c3.6 3.6 7.8 5.4 12.8 5.4s9.2-1.8 12.8-5.4L287 95c3.5-3.5 5.4-7.8 5.4-12.8 0-5-1.9-9.2-5.5-12.8z'/%3E%3C/svg%3E");
		background-repeat: no-repeat, repeat;
		background-position: right .7em top 55%, 0 0;
		background-size: .5em auto, 100%;
		padding-right: 1.8em;

		background-color: var(--background-highlight);
		border-radius: 0;
		margin: 0;
		color: inherit;
		border-color: var(--border-color);
	}
</style>
<select bind:value on:change={dispatchSelect}>
	<slot />
</select>


