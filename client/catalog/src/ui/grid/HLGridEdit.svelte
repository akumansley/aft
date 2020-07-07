<script >
	export let id=null;
	export let value;
	export let placeholder="Value";
	export let restrict;
	export let remove = null;
	export let idx = null;
	import Icon from 'fa-svelte';
	import { faMinusSquare } from '@fortawesome/free-solid-svg-icons/faMinusSquare';

</script>
<style>

	.object-box {
		color: inherit;
		width:150px;
		background: var(--background-highlight);

	}
	input {
		width:calc(150px + 3em);
		background: inherit;
 		border: 1;
 		border-color: var(--background-highlight);
		color: inherit;
		padding: 1.25em 1.5em;
		margin:0;
	}
	input:hover {
		border-color: var(--border-color);
	}

	.close {
		  cursor: pointer;
		  position: absolute;
		  top: 50%;
		  width:1em;
		  left: calc(150px + 3em);
		  padding: 1em .75em;
		  transform: translate(0%, -50%); 
		  color: var(--text-color-darker);
	}
	.close:hover {
		  color: var(--text-color);
	}
	.icon {
		  opacity: 0.5;
	}
	.outer {
		position:relative;
	}
</style>
<div class="outer">
	<div class="object-box">
		<input
			id={id} placeholder={placeholder} type="text" 
			on:input={(e) => { if (restrict) {
			value = restrict(e.target.value);
			}}} bind:value />
	</div>
	<div class="close" on:click={remove(idx)}>
		<Icon class="icon" icon={faMinusSquare}></Icon>
	</div>
</div>

