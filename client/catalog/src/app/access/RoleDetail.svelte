<script>
export let params = null;
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
import PolicyForm from './PolicyForm.svelte';
import HLSectionTitle from '../../ui/list/HLSectionTitle.svelte';


let role = {};
let load = client.api.role.findOne({
	where: {id: params.id}, 
	include: {
		policies: {
			include: { 
				model: true 
			},
		},
	},
}).then((data) => {
	role = data;
});

const addPolicy = () => {
	role.policies = [...role.policies, {}];
};

function policyToOp(p) {
	if (p.id) {
		return {
			where: { id: p.id },
			data: {
				text: p.text,
				model: {connect: {id: p.model.id}},
			},
		}
	} else {
		return {
			text: p.text,
			model: {
				connect: {id: p.model.id},
			},
		};
	}
}

function saveRoleAndPolicies() {
	const policyOp = {
		update: [],
		create: [],
	}
	for (let p of role.policies) {
		let op = policyToOp(p)
		if (p.id) {
			policyOp.update.push(op);
		} else {
			policyOp.create.push(op);
		}
	}
	client.api.role.update({
		where: {id: role.id},
		data: {
			name: role.name,
			policies: policyOp,
		}
	});
}
</script>
<style>
	.v-space {
		height: .5em;
	}
</style>

{#await load then loaded}

<HLHeader>
	<Name id="name" placeholder="Role name.." bind:value={role.name} on:click={saveRoleAndPolicies}>
	</Name>
</HLHeader>
<HLContent>
	<HLSectionTitle>Policies</HLSectionTitle>
	{#each role.policies as policy}
		<PolicyForm bind:policy={policy} />
		<div class="v-space"/>
	{/each}
	<HLRowButton on:click={addPolicy}>+ Add</HLRowButton>
</HLContent>

{/await}
