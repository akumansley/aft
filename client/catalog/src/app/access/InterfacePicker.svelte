<script>
export let interface = [];

import HLSelect from '../../ui/form/HLSelect.svelte';
import client from '../../data/client.js';

const load = client.api['interface'].findMany({});
</script>

{#await load then allInterfaces} 
	<HLSelect options={allInterfaces} bind:selected={interface}/>
{/await}



