<script>
export let attribute;
import HLRow from '../../ui/HLRow.svelte';
import HLSelect from '../../ui/HLSelect.svelte';
import HLText from '../../ui/HLText.svelte';
import { AttrType } from '../../data/enums.js';

function restrict(s) {
	const newVal = s.replace(/[^a-zA-Z_]/g, '');
	return newVal;
}
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
		<HLSelect bind:value={attribute.attrType}>
			{#each Object.entries(AttrType) as attr, ix}
			<option value={ix}>
				{attr[1]}
			</option>
			{/each}
		</HLSelect>
	</div>
</HLRow>
