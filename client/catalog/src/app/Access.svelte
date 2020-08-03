<script>
import { navStore } from './stores.js';
import client from '../data/client.js';
 
let load = client.api.policy.findMany({ });
navStore.set("access");

</script>

<style>
</style>


<h1>Policies</h1>
{#await load then policies}
	{#each policies as policy}
		{policy.name}
	{/each}
{/await}

