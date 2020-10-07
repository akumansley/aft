<script>
	import {HLHeader, HLContent, HLHeaderItem} from '../../ui/page/page.js';
	import HLButton from '../../ui/form/HLButton.svelte';
	import {Box} from '../../ui/spacing/spacing.js';
	import {cap} from '../../lib/util.js';

	import Name from '../Name.svelte';
	import {createEventDispatcher} from 'svelte';
	const dispatch = createEventDispatcher();

	export let value = null; 
</script>

<HLHeader>
	<HLHeaderItem>
		<Name bind:value={value.name} />
	</HLHeaderItem>
	<HLHeaderItem>
		<HLButton on:click={() => dispatch('save')}>Save</HLButton>
	</HLHeaderItem>
</HLHeader>
<Box>
	{cap(value.name)} is implemented in native code.
</Box>
