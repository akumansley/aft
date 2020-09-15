<script>
	import HLText from '../ui/form/HLText.svelte';
	import { onMount, getContext } from 'svelte';
	import { ctx } from './context.js';

	let init = null;
	let value = null;

	export let placeholder = null;
	export let restrict = null;
	export let key = null;

	let op;

	onMount( () => {
		op = getContext(key);
	});


	function setOp() {
		if (op && value && init && value === init) {
			op.clear(key);
		} else if (op && value) {
			op.setOp(key, value);
		}

	}
	$: setOp(init, value);

</script>

<HLText bind:value bind:placeholder bind:restrict /> 
