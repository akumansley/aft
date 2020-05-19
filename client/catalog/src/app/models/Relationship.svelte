<script>
export let relationship; 
export let left; 
import HLRow from '../../ui/HLRow.svelte';
</script>

<HLRow>
{#if left}
	{relationship.leftName}
{:else}
	{relationship.rightName}
{/if}
</HLRow>


