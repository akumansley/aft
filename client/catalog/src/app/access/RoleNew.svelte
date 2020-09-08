<script>
import { navStore } from '../stores.js';
import {router} from '../router.js';

import client from '../../data/client.js';

import HLRowButton from '../../ui/list/HLRowButton.svelte';
import HLButton from '../../ui/form/HLButton.svelte';
import HLRow from '../../ui/list/HLRow.svelte';
import HLHeader from '../../ui/page/HLHeader.svelte';
import HLContent from '../../ui/page/HLContent.svelte';
import Name from '../Name.svelte';

const role = {
	name: "",
};

const saveAndNav = async () => {
	await client.api.role.create({data: role});
	router.route("/roles");
};

</script>

<HLHeader>
	<Name id="name" placeholder="Role name.." bind:value={role.name} on:click={saveAndNav}>
	</Name>
</HLHeader>
