<script>
	export let params = null;

	import client from '../../data/client.js';
	import { router } from '../router.js';
	import { cap } from '../../lib/util.js';

	import { navStore } from '../stores.js';
	navStore.set("datatype");
	
	import {HLHeader, HLContent, HLHeaderItem} from '../../ui/page/page.js';
	import {HLButton} from '../../ui/form/form.js';
	import {Box} from '../../ui/spacing/spacing.js';
	import Name from '../Name.svelte';


	const load = client.api.datatype.findOne({where: {id: params.id}})

	async function saveAndNav() {
		router.route("/datatypes");
	}

</script>

{#await load then datatype}
<HLHeader>
	<HLHeaderItem>
		<Name value={datatype.name}></Name>
	</HLHeaderItem>
	<HLHeaderItem>
		<HLButton on:click={saveAndNav}>Back</HLButton>
	</HLHeaderItem>
</HLHeader>
<HLContent>
	<Box>
		{cap(datatype.name)} is implemented in native code.
	</Box>
</HLContent>
{/await}