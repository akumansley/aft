<script>
import { navStore } from '../stores.js';
import client from '../../data/client.js';
import { router } from '../../app/router.js';
import HLRow from '../../ui/list/HLRow.svelte';
import HLRowButton from '../../ui/list/HLRowButton.svelte';

 
let load = client.api.role.findMany({ });
navStore.set("access");

let newRole = () => router.route("/role/new");
</script>

<style>
</style>


<h1>Roles</h1>
{#await load then roles}
	{#each roles as role}
	<HLRow>
		{role.name}
	</HLRow>
	{/each}
	<HLRowButton on:click={newRole}>
+ add
</HLRowButton>
{/await}

