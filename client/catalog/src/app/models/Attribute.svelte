<script>
export let attribute
import HLRow from '../../ui/HLRow.svelte';
import client from '../../data/client.js';
let load = client.datatype.findOne({where:{id: attribute.datatypeId}});
</script>
<style>
dl {
	padding: 0; 
	margin:0;
}
dt {
	font-size: var(--scale--2);
	text-transform: uppercase;
	font-weight: 600;
}
dd {
	margin: 0;
}
.v-space{
	height: .5em;
}
.col {
	display: flex;
	flex-direction: column;
}
</style>
<HLRow>
	<div class="col">
	{attribute.name}
	<div class="v-space"/>
	<dl>
		<dt>Type</dt>
	{#await load}
		&nbsp;
	{:then datatype}
	<dd>
		{datatype.name}
	</dd>
	{:catch error}
		<div>Error..</div>
	{/await}
	</dl>
	</div>
</HLRow>

