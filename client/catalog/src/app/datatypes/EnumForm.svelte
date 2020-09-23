<script>
export let value = null;

import {createEventDispatcher} from 'svelte';
import { restrictToIdent } from '../../lib/util.js';

import {HLHeader, HLHeaderItem, HLContent, HLSectionTitle} from '../../ui/page/page.js';
import {HalfBox, Box} from '../../ui/spacing/spacing.js';
import {HLButton, HLText} from '../../ui/form/form.js';
import Name from '../Name.svelte';

function addEnumValue() {
	value.enumValues = value.enumValues.add();
}

const dispatch = createEventDispatcher();

</script>

<HLHeader>
	<HLHeaderItem>
		<Name bind:value={value.name}></Name>
	</HLHeaderItem>
	<HLHeaderItem>
		<HLButton on:click={() => {dispatch('save')}}>Save</HLButton>
	</HLHeaderItem>
</HLHeader>

<HLContent>
	<HLSectionTitle>Values</HLSectionTitle>
	{#each value.enumValues as enumValue, ix}
	<HalfBox>
		<HLText bind:value={enumValue.name} placeholder={"add value.."} restrict={restrictToIdent} />
	</HalfBox>
	{/each}

	<Box>
		<HLButton on:click={addEnumValue}>+ add</HLButton>
	</Box>
</HLContent>

