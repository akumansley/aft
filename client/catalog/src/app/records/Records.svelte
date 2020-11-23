<script>
	import {navStore} from '../stores.js';
	import client from '../../data/client.js';

	import {HLSelect} from '../../ui/form/form.js';
	import {HLBorder} from '../../ui/page/page.js';

	navStore.set("records");

	let selectedModel;
	let model;
	let records = [];
	let models = [];
	async function load() {
		models = await client.api.model.findMany({
			include: {
				attributes: true,
				relationships: true,
			}
		});
		selectedModel = models[0];
	}
	let init = load();
	$: loadRecords(selectedModel);

	async function loadRecords() {
		if (selectedModel) {
			const include = {}
			for (let rel of selectedModel.relationships) {
				include[rel.name] = true;
			}
			records = await client.api[selectedModel.name].findMany({
				include: include,
			});
			model = selectedModel;
		}
	}

</script>

<style>

	.header {
		margin: 1em ;
	}
	.records {
		overflow: auto;
		display: flex;
		flex-direction: column;
	}
	table {
		width: 100%;
		border-collapse: collapse;
		border-bottom: 1px solid var(--border-color);
	}
	th {
		border-right: 1px solid var(--border-color);
		border-bottom: 1px solid var(--border-color);
	}
	th:last-child {
		border-right:none;
	}
	td {
		border-right: 1px solid var(--border-color);
		border-bottom: 1px solid var(--border-color);
		padding: .25em .5em;
	}
	td:last-child {
		border-right:none;
	}
</style>

<div class="header">
	View records for:
	<HLSelect bind:value={selectedModel}>
		{#each models as m}
		<option value={m}>{m.name}</option>
		{/each}
	</HLSelect>
</div>

<HLBorder/>

{#if model}
<div class="records">
	<table>
		<tr>
			{#each model.attributes as attribute}
			<th>{attribute.name}</th>
			{/each}
			{#each model.relationships as relationship}
			<th>{relationship.name}</th>
			{/each}
		</tr>
		{#each records as record}
		<tr>
			{#each model.attributes as attribute}
			<td>
				{record[attribute.name]}
			</td>
			{/each}
			{#each model.relationships as relationship}
			{#if Array.isArray(record[relationship.name])}
			<td>
				{record[relationship.name].length} related
			</td>
			{:else}
			<td>{record[relationship.name]? record[relationship.name].id: "null"}</td>
			{/if}
			{/each}
		</tr>
		{/each}
	</table>
</div>
{/if}
