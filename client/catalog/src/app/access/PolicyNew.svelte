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

const policy = {
	name: "",
	roles: [],
	actions: [],
	interfaces: [],
	conditions: [],
};

const saveAndNav = () => {
};

</script>

<HLHeader>
	<Name id="name" placeholder="Policy name.." bind:value={policy.name} click={saveAndNav}>
	</Name>
</HLHeader>
<HLContent>
	<h2>Roles</h2>
	<RolesPicker bind:roles={policy.roles}></RolesPicker>
	<h2>Actions</h2>
	<ActionsPicker bind:actions={policy.actions}></ActionsPicker>
	<h2>Models</h2>
	<h2>Conditions</h2>
	<!-- WhereClause bind:where={policy.conditions}></WhereClause> -->
</HLContent>
