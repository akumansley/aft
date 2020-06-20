<script>
import client from '../../data/client.js';
import {Runtime} from '../../data/enums.js';
import {cap} from '../util.js';

let load = client.rpc.findMany({	include: {
		code: true,
	}
});

import { breadcrumbStore } from '../stores.js';
breadcrumbStore.set(
	[{
		href: "/rpcs",
		text: "RPCs",
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
	{:then rpcs}
		{#each rpcs as rpc}
			<a href="/rpc/{rpc.id}" class="object-box">
				<div class="obj-title">{cap(rpc.name)}</div>
				<div>{Runtime[rpc.code.runtime]}</div>
			</a>
			<div class="spacer"/>
		{/each}
		<a href="/rpcs/new" class="object-box">
			<div>+ Add</div>
		</a>
		<div class="spacer"/>
	{:catch error}
		<div>Error..</div>
	{/await}
</div>
