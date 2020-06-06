<script>
import client from '../../data/client.js';
import {Runtime} from '../../data/enums.js';
import {cap} from '../util.js';

let load = client.datatype.findMany({	include: {
		validator: true,
	}
});

import { breadcrumbStore } from '../stores.js';
breadcrumbStore.set(
	[{
		href: "/datatypes",
		text: "Datatypes",
	}]
);

</script>

<style>
	.box {
		display: flex;
		flex-direction: row;
		flex-wrap: wrap;
	}
	.stuff {

	}
	a.object-box {
		display: flex;
		flex-direction: column;
		color: inherit;
		width: 150px;
		padding: 1em 1.5em;
	}
	a.object-box:hover {
		background: var(--background-highlight);
	}

	.spacer {
		width: 0;
	}
	.obj-title{
		font-weight: 600;
	}

</style>

<div class="box">
	{#await load}
		&nbsp;
	{:then datatypes}
		{#each datatypes as datatype}
			<a href="/datatype/{datatype.id}" class="object-box">
				<div class="obj-title">{cap(datatype.name)}</div>
				<div>{Runtime[datatype.validator.runtime]}</div>
			</a>
			<div class="spacer"/>
		{/each}
		<a href="/datatypes/new" class="object-box">
			<div>+ Add</div>
		</a>
		<div class="spacer"/>
	{:catch error}
		<div>Error..</div>
	{/await}
</div>
