<script>
export let actions = [];

const allActions = [
	{
		id: "read",
		name: "read",
	},
	{
		id: "write",
		name: "write",
	},
	{
		id: "execute",
		name: "execute",
	},
];

import HLMultiSelect from '../../ui/form/HLMultiSelect.svelte';

</script>

<HLMultiSelect options={allActions} bind:selected={actions}/>



