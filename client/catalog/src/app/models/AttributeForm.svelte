<script>
	import { getContext } from 'svelte';
	import { key } from '../../api/api.js';

	export let value;
	let attribute = {name:""};
	let datatype = {};

	import DatatypeSelect from './DatatypeSelect.svelte';
	import HLSelect from '../../ui/form/HLSelect.svelte';
	import HLText from '../../ui/form/HLText.svelte';
	import {restrictToIdent, cap, isObject} from '../util.js';

	let operation = getContext(key);

	$: {
		if (operation === "update") {
			value = {
				where: {id: attribute.id},
				data: {
					name: attribute.name,
					datatype: datatype,
				}
			}
		} else if (operation === "create") {
			value = {
				name: attribute.name,
				datatype: datatype,
			}

		}
	}

</script>

<style>
	.hform-row {
		display: flex; 
		flex-direction: row;
		padding: calc(var(--box-margin)/ 2) var(--box-margin);
	}
	.spacer {
		width: 1em;
		height: 0;
	}
</style>

<div class="hform-row">
	
	<HLText placeholder="Attribute name.." bind:value={attribute.name} 
	restrict={restrictToIdent}/>

	<div class="spacer"/>

	<DatatypeSelect bind:value={datatype} />
</div>