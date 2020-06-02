<script>
export let model
import HLTable from '../../ui/HLTable.svelte';
import Relationship from './Relationship.svelte';
import Attribute from './Attribute.svelte';
let cap= (s) => { 
	if (!s) {
		return "";
	}
	return s.charAt(0).toUpperCase() + s.slice(1)
};
</script>
<style>
h1 {
	font-size: var(--scale-3);
	font-weight: 600;
}
h2 {
	font-size: var(--scale--1);
	font-weight: 500;
	line-height: 1;
}
.v-space{
	height: .5em;
}

</style>
<h1>{cap(model.name)}</h1>

<h2>Attributes</h2>
<HLTable>
{#each model.attributes as attr}
	<Attribute attribute={attr} />
{/each}
</HLTable>

<div class="v-space"/>
<h2>Relationships</h2>

<HLTable>
{#each model.rightRelationships as rel}
	<Relationship relationship={rel} left={false}/>
{/each}
{#each model.leftRelationships as rel}
	<Relationship relationship={rel} left={true}/>
{/each}
</HLTable>
