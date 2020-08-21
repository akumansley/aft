<script>
import { navStore } from '../stores.js';
import {router} from '../router.js';

import client from '../../data/client.js';

import HLRowButton from '../../ui/list/HLRowButton.svelte';
import HLButton from '../../ui/form/HLButton.svelte';
import HLRow from '../../ui/list/HLRow.svelte';
import HLHeader from '../../ui/main/HLHeader.svelte';
import HLContent from '../../ui/main/HLContent.svelte';
import Name from '../Name.svelte';
import RolesPicker from './RolesPicker.svelte';
import ActionsPicker from './ActionsPicker.svelte';
import InterfacesPicker from './InterfacesPicker.svelte';

const role = {
	name: "",
};

const saveAndNav = async () => {
	await client.api.role.create({data: role});
	router.route("/access");
};

</script>

<HLHeader>
	<Name id="name" placeholder="Role name.." bind:value={role.name} click={saveAndNav}>
	</Name>
</HLHeader>
