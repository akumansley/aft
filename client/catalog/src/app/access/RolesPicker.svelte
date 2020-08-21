<script>
export let roles = [];

import HLMultiSelect from '../../ui/form/HLMultiSelect.svelte';
import client from '../../data/client.js';

const load = client.api.role.findMany({});
</script>

{#await load then allRoles} 
	<HLMultiSelect options={allRoles} bind:selected={roles}/>
{/await}



