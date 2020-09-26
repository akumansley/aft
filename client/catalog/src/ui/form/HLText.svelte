<script>
export let id=null;
export let value;
export let placeholder;
export let restrict;
import { createEventDispatcher } from 'svelte';
const dispatch = createEventDispatcher();

</script>
<style>
	input::placeholder {
		color: var(--text-color-darker);
	}
	input {
		border-left: none;
		border-top: none;
		border-right: none;
		border-bottom: none;
		border-color: var(--border-color);
		color: inherit;
		border-radius: 0;
		margin: 0;
		padding: calc(.4em + 1px);
		background: var(--background-highlight);
	}
</style>


<input 
	id={id} placeholder={placeholder} 
	type="text" 
	bind:value 
	on:input={(e) => { if (restrict) {
		value = restrict(e.target.value);
		dispatch('input', e);
	}}} 
	/>
