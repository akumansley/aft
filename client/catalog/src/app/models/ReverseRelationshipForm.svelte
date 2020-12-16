<script>
	export let value;
	import { restrictToIdent } from '../../lib/util.js';

	import {HLText, HLCheckbox, HLChip} from '../../ui/form/form.js';
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
	<HLChip>&harr; {value.referencing.source[0].name}.{value.referencing.name}</HLChip>
	<div class="spacer"/>
	<HLCheckbox disabled={true} checked={value.multi}>Multiple</HLCheckbox>
</div>
