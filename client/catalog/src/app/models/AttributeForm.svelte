<script>
export let attribute;
import HLRow from '../../ui/HLRow.svelte';
import HLSelect from '../../ui/HLSelect.svelte';
import HLText from '../../ui/HLText.svelte';
import client from '../../data/client.js';
import {afterUpdate} from 'svelte';
let load = client.datatype.findMany({where:{}});

function restrict(s) {
	const newVal = s.replace(/[^a-zA-Z_]/g, '');
	return newVal;
}
let cap= (s) => {
	if (!s) {
		return "";
	}
	s = s.replace(/[\w]([A-Z])/g, function(m) {
           return m[0] + " " + m[1];
       });
	return s.charAt(0).toUpperCase() + s.slice(1)
};

afterUpdate(() => {
	attribute.datatype.connect.id = attribute.datatypeId;
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
		<HLText placeholder="Attribute name.." bind:value={attribute.name} restrict={restrict}/>
		<div class="spacer"/>
		{#await load}
			&nbsp;
		{:then datatypes}
		<HLSelect bind:value={attribute.datatypeId}>
			{#each Object.entries(datatypes) as attr}
			<option value={attr[1].id}>
				{cap(attr[1].name)}
			</option>
			{/each}
		</HLSelect>
		{:catch error}
			<div>Error..</div>
		{/await}
	</div>
</HLRow>
