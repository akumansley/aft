<script>
	import { ConnectSelect } from '../../api/api.js';

	export let value;
	import { restrictToIdent } from '../../lib/util.js';

	import HLButton from '../../ui/form/HLButton.svelte';
	import HLCheckbox from '../../ui/form/HLCheckbox.svelte';
	import HLText from '../../ui/form/HLText.svelte';
</script>

<style>
	.spacer {
		width: 1em;
	}

	.hform-row {
		display: flex; 
		flex-direction: row;
		padding: calc(var(--box-margin)/ 2) var(--box-margin);
	}
</style>

<div class="hform-row">

	<HLText bind:value={value.name} placeholder="Relationship name.."restrict={restrictToIdent} />

	<div class="spacer"/>
	<ConnectSelect bind:value={value.target} iface={"interface"} />

	<div class="spacer"/>
	<HLCheckbox bind:checked={value.multi}>Multiple</HLCheckbox>
</div>
