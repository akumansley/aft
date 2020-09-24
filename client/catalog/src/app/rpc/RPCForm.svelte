<script>
	import {HLHeader, HLContent, HLHeaderItem} from '../../ui/page/page.js';
	import HLButton from '../../ui/form/HLButton.svelte';
	import {createEventDispatcher} from 'svelte';

	import Name from '../Name.svelte';
	import CodeMirror from '../codemirror/CodeMirror.svelte';


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
	</HLHeader>
	<div class="fill-v">
		<CodeMirror bind:value={value.code}/>
	</div>
</div>
