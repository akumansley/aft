<script context="module">
let allInterfaces = [];
client.api['interface'].findMany({}).then(ifcs => {
	allInterfaces = ifcs
});
</script>

<script>
export let policy; 

$: {
	if (policy) {
		console.log(policy)
		if (!policy.text) {
			policy.text = "{}";
		}
		if (!policy.model) {
			policy.model = allInterfaces[0];
		} else {
			for (let ifc of allInterfaces) {
				if (ifc.id == policy.model.id) {
					policy.model = ifc;
					break;
				}
			}
		}
	}
}

import HLSelect from '../../ui/form/HLSelect.svelte';
import client from '../../data/client.js';
import ActionsPicker from './ActionsPicker.svelte';
import HLSection from '../../ui/form/HLSection.svelte';
import HLButton from '../../ui/form/HLButton.svelte';


let showOperations = false;
const toggleOperations = () => { showOperations = !showOperations };

let showWhere = false;
const toggleWhere = () => { showWhere = !showWhere };

</script>

<HLSection>
	<HLSelect bind:value={policy.model}>
		{#each allInterfaces as ifc}
			<option value={ifc}>{ifc.name}</option>
		{/each}
	</HLSelect>
	<HLButton on:click={toggleOperations}>Operations</HLButton>
	<HLButton on:click={toggleWhere}>Where</HLButton>
	<HLButton>Fields</HLButton>
	{#if showOperations}
		<h2>Operations</h2>
		<ActionsPicker />
	{/if}
	{#if showWhere}
		<h2>Where</h2>
		<textarea bind:value={policy.text}></textarea>
	{/if}

</HLSection>



