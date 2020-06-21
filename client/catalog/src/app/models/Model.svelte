<script>
export let model
import HLTable from '../../ui/HLTable.svelte';
import Relationship from './Relationship.svelte';
import Attribute from './Attribute.svelte';
import {cap} from '../util.js';

</script>

<h1>{cap(model.name)}</h1>

<h2>Attributes</h2>
<HLTable>
{#each model.attributes as attr}
	<Attribute attribute={attr}/>
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
