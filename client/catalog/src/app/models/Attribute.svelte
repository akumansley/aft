<script>
export let attribute
import {cap} from '../util.js';
import HLDetailRow from '../../ui/HLDetailRow.svelte';
import client from '../../data/client.js';
let load = client.api.datatype.findOne({where:{id: attribute.datatypeId}});
</script>

{#await load}
	&nbsp;
{:then datatype}
<HLDetailRow name={attribute.name} header={"Type"}>{cap(datatype.name)}</HLDetailRow>
{:catch error}
	<div>Error..</div>
{/await}

