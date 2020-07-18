<script>
export let id=null;
export let value;
export let placeholder;
export let restrict;
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
		background: var(--background-highlight);
	}
</style>


<input 
	id={id} placeholder={placeholder} 
	type="text" 
	on:input={(e) => { if (restrict) {
	value = restrict(e.target.value);
	}}} 
	bind:value />
