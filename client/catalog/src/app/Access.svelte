<script>
import { navStore } from './stores.js';
import aft from '../data/aft.js';
 
let load = aft.api.policy.findMany({ });
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

