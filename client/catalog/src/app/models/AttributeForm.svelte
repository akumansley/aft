<script>
export let attribute;
import HLRow from '../../ui/list/HLRow.svelte';
import HLSelect from '../../ui/form/HLSelect.svelte';
import HLText from '../../ui/form/HLText.svelte';
import client from '../../data/client.js';
import {afterUpdate} from 'svelte';
import {restrictToIdent, cap, isObject} from '../util.js';
let load = client.api.datatype.findMany({where:{}});

afterUpdate(() => {
	if(isObject(attribute.datatype)) {
		attribute.datatype.connect.id = attribute.datatypeId;
	}
});

</script>
<style>
.hform-row {
	display: flex; 
	flex-direction: row;
}
.spacer {
	width: 1em;
	height: 0;
}
</style>
<HLRow>
	<div class="hform-row">
		<HLText placeholder="Attribute name.." bind:value={attribute.name} restrict={restrictToIdent}/>
		<div class="spacer"/>
		{#await load then datatypes}
		<HLSelect bind:value={attribute.datatypeId}>
			{#each Object.entries(datatypes) as attr}
			<option value={attr[1].id}>
				{cap(attr[1].name)}
			</option>
			{/each}
		</HLSelect>
		{/await}
	</div>
</HLRow>
