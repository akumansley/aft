<script>
export let value = null;

import {createEventDispatcher} from 'svelte';
import { restrictToIdent } from '../../lib/util.js';

import ConnectSelect from '../../api/ConnectSelect.svelte';
import {HLHeader, HLHeaderItem, HLContent, HLSectionTitle, HLHeaderDetail} from '../../ui/page/page.js';
import {HalfBox, Box, HSpace} from '../../ui/spacing/spacing.js';
import {HLButton, HLText} from '../../ui/form/form.js';
import Name from '../Name.svelte';

function addEnumValue() {
	value.enumValues = value.enumValues.add();
}

const dispatch = createEventDispatcher();
let showDetail = false;

</script>

<HLHeader>
	<HLHeaderItem>
		<Name bind:value={value.name}></Name>
	</HLHeaderItem>
	<HLHeaderItem>
		<HLButton on:click={() => {dispatch('save')}}>Save</HLButton>
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

