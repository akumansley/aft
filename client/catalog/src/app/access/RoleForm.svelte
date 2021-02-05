<script>
	import client from '../../data/client.js';

	import {HLButton} from '../../ui/form/form.js';
	import {HLHeader, HLHeaderItem, HLContent, HLSectionTitle, HLHeaderDetail} from '../../ui/page/page.js';
	import ConnectSelect from '../../api/ConnectSelect.svelte';
	import {Box, HSpace} from '../../ui/spacing/spacing.js';

	import Name from '../Name.svelte';
	import PolicyForm from './PolicyForm.svelte';
	import ModuleMultiSelect from './ModuleMultiSelect.svelte';
	import {createEventDispatcher} from 'svelte';
	
	const RPC = "4b8db42e-d084-4328-a758-a76939341ffa";

	export let value;

	function addPolicy() {
		value.policies = value.policies.add();
	};

	const dispatch = createEventDispatcher();

	let functionOptions = []
	client.api.module.findMany({
		where: {
			functions: {some: {funcType: RPC}},
		},
		order: [{"name": "asc"}],
		include: {
			functions: {where: {funcType: RPC}},
		}})
	.then(result => {
		functionOptions = result;
	})

	let showDetail = false;

</script>


<HLHeader>
	<HLHeaderItem>		
		<Name placeholder="Role name.." bind:value={value.name} />
	</HLHeaderItem>		
	<HLHeaderItem>	
		<HLButton on:click={() => dispatch('save')}>Save</HLButton>
	</HLHeaderItem>
	<HLHeaderItem>
		<HLButton on:click={() => showDetail = !showDetail}>More</HLButton>
	</HLHeaderItem>
</HLHeader>

{#if showDetail}
<HLHeaderDetail>
	<HLHeaderItem>
		Module: <HSpace/> <ConnectSelect pickDefault={(m) => m.goPackage === ""} bind:value={value.module} iface={"module"} />
	</HLHeaderItem>
</HLHeaderDetail>
{/if}


<HLContent>
	<HLSectionTitle>Models</HLSectionTitle>
	{#each value.policies as policy}
	<PolicyForm bind:value={policy} />
	{/each}
	<Box>
		<HLButton on:click={addPolicy}>+ add</HLButton>
	</Box>
	<HLSectionTitle>Functions</HLSectionTitle>
	<Box>
		<ModuleMultiSelect options={functionOptions} key={"functions"} bind:value={value.executableFunctions} />
	</Box>
</HLContent>

