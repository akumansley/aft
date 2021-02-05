<script>
	import {HLHeader, HLHeaderDetail, HLContent, HLHeaderItem} from '../../ui/page/page.js';
	import {HLButton, HLSelect} from '../../ui/form/form.js';
	import {HSpace } from '../../ui/spacing/spacing.js';
	import {createEventDispatcher} from 'svelte';
	import ConnectSelect from '../../api/ConnectSelect.svelte';

	import Name from '../Name.svelte';
	import CodeMirror from '../codemirror/CodeMirror.svelte';
	import client from '../../data/client.js';

	let roles = [];
	client.api.role.findMany({}).then(result => {
		roles = result;
	})

	let showDetail = false;

	const dispatch = createEventDispatcher();
	export let value = null; 
</script>
<style>
	.container {
		display: flex;
		height: 100vh;
		flex-direction: column;
	}
	.fill-v {
		flex-grow: 1;
		display: flex;
		overflow: auto;
	}
</style>

<div class="container">
	<HLHeader>
		<HLHeaderItem>
			<Name bind:value={value.name} />
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
			Role: <HSpace/> <ConnectSelect allowEmpty={true} bind:value={value.role} iface={"role"} />
		</HLHeaderItem>
		<HLHeaderItem>
			Module: <HSpace/> <ConnectSelect pickDefault={(m) => m.goPackage === ""} bind:value={value.module} iface={"module"} />
		</HLHeaderItem>
	</HLHeaderDetail>
	{/if}
	<div class="fill-v">
		<CodeMirror bind:value={value.code}/>
	</div>
</div>
