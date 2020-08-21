<script>
export let interfaces = [];

import HLMultiSelect from '../../ui/form/HLMultiSelect.svelte';
import client from '../../data/client.js';

const load = client.api['interface'].findMany({});
</script>

{#await load then allInterfaces} 
	<HLMultiSelect options={allInterfaces} bind:selected={interfaces}/>
{/await}



