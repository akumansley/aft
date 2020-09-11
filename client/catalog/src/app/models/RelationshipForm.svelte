<script context="module">
	import client from '../../data/client.js';
	let load = client.api.model.findMany({where:{}});
</script>

<script>
	export let value;
	let target;
	let relationship = {
		name:"",
		multi: false,
	};
	import { restrictToIdent } from '../../lib/util.js';
	import HLButton from '../../ui/form/HLButton.svelte';
	import HLCheckbox from '../../ui/form/HLCheckbox.svelte';
	import HLText from '../../ui/form/HLText.svelte';

	import ModelSelect from './ModelSelect.svelte';

	import { RelType } from '../../data/enums.js';

	import { getContext } from 'svelte';
	import { key } from '../../api/api.js';

	let operation = getContext(key);

	$: {
		if (operation === "update") {
			value = {
				where: {id: relationship.id},
				data: {
					name: relationship.name,
					target: target,
				}
			}
		} else if (operation === "create") {
			value = {
				type: "concreteRelationship",
				name: relationship.name,
				target: target,
			}

		}
	}
</script>

<style>
	.spacer {
		width: 1em;
	}

	.hform-row {
		display: flex; 
		flex-direction: row;
		padding: calc(var(--box-margin)/ 2) var(--box-margin);
	}
</style>

<div class="hform-row">

	<HLText 
	bind:value={relationship.name}
	placeholder="Relationship name.." 
	restrict={restrictToIdent}
	/>

	<div class="spacer"/>
	<ModelSelect bind:value={target} />

	<div class="spacer"/>
	<HLCheckbox bind:checked={relationship.multi}>Multiple</HLCheckbox>
</div>
