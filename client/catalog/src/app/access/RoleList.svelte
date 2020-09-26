<script>
import { navStore } from '../stores.js';
import client from '../../data/client.js';
import { router } from '../../app/router.js';
import HLRowButton from '../../ui/list/HLRowButton.svelte';

import HLSectionTitle from '../../ui/page/HLSectionTitle.svelte';
import HLBorder from '../../ui/page/HLBorder.svelte';

import HLGrid from '../../ui/grid/HLGrid.svelte';
import HLGridItem from '../../ui/grid/HLGridItem.svelte';
import HLGridNew from '../../ui/grid/HLGridNew.svelte';

 
let load = client.api.role.findMany({ });
navStore.set("access");

let newRole = () => router.route("/roles/new");
</script>


{#await load then roles}
<HLGrid>
	<HLGridNew href={"/roles/new"}>
		Add Role
	</HLGridNew>
</HLGrid>
<HLBorder/>
<HLSectionTitle>Roles</HLSectionTitle>
<HLGrid>
	{#each roles as role}
		<HLGridItem href="/role/{role.id}" name={role.name}/>
	{/each}
	</HLGrid>
{/await}

