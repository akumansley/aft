<script>
import client from '../../data/client.js';
import {Runtime} from '../../data/enums.js';
import {cap} from '../util.js';
import HLGrid from '../../ui/HLGrid.svelte';
import HLGridItem from '../../ui/HLGridItem.svelte';
import { breadcrumbStore } from '../stores.js';

let load = client.datatype.findMany({include: {validator: true}});

breadcrumbStore.set(
	[{
		href: "/datatypes",
		text: "Datatypes",
	}]
);
</script>
<HLGrid>
	{#await load}
		&nbsp;
	{:then datatypes}
		{#each datatypes as datatype}
			<HLGridItem href={"/datatype/" + datatype.id} name={datatype.name}>
				<div>{Runtime[datatype.validator.runtime]}</div>
			</HLGridItem>
		{/each}
		<HLGridItem href={"/datatypes/new"}>
			<div>+ Add</div>
		</HLGridItem>
	{:catch error}
		<div>Error..</div>
	{/await}
</HLGrid>
