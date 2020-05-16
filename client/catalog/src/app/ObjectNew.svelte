<script>
import { onMount } from 'svelte';
import client from '../data/client.js';
import HLTable from '../ui/HLTable.svelte';
import HLRowButton from '../ui/HLRowButton.svelte';
import HLRow from '../ui/HLRow.svelte';
import { AttrType } from '../data/enums.js';
import { breadcrumbStore } from './breadcrumbStore.js';
breadcrumbStore.set(
	[{
		href: "/objects",
		text: "Objects",
	}, {
		href: "/objects/new",
		text: "New",
	}]
);
const newModelOp = {
	name: "",
	attributes: {create: []},
	relationships: {create: []},
}
function addAttribute() {
	newModelOp.attributes.create = [...newModelOp.attributes.create, {
		name: "",
		attrType: 0,
	}];
}
function saveModel() {
	client.model.create({data: newModelOp})
}
function restrict(s) {
	const newVal = s.replace(/[^a-zA-Z_]/g, '');
	return newVal;
}
</script>

<style>
	.box {
		margin: 1em 1.5em; 
	}
	.v-space{
		height: .5em;
	}
	.v-space-2{
		height: 2em;
	}
	input.h1 {
		font-size: var(--scale-3);
		font-weight: 600;
		border-left: none;
		border-top: none;
		border-right: none;
	}
	input {
		border-color: var(--border-color);
		color: inherit;
		border-radius: 0;
		margin: 0;
		background: var(--background-highlight);
	}
	input.hl-row {
		border-left: none;
		border-top: none;
		border-right: none;
	}
	input::placeholder {
		color: var(--text-color-darker);
	}
	h2 {
		font-size: var(--scale--1);
		font-weight: 500;
		line-height: 1;
	}
	.hform-row {
		display: flex; 
		flex-direction: row;
	}
	.col {

	}
	.spacer {
		width: 1em;
		height: 0;
	}
	select {
		-webkit-appearance: none;
		-moz-appearance: none;
		appearance: none;
		background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='292.4' height='292.4'%3E%3Cpath fill='%23e4e1e8' d='M287 69.4a17.6 17.6 0 0 0-13-5.4H18.4c-5 0-9.3 1.8-12.9 5.4A17.6 17.6 0 0 0 0 82.2c0 5 1.8 9.3 5.4 12.9l128 127.9c3.6 3.6 7.8 5.4 12.8 5.4s9.2-1.8 12.8-5.4L287 95c3.5-3.5 5.4-7.8 5.4-12.8 0-5-1.9-9.2-5.5-12.8z'/%3E%3C/svg%3E");
		background-repeat: no-repeat, repeat;
		background-position: right .7em top 55%, 0 0;
		background-size: .5em auto, 100%;
		padding-right: 1.8em;

		background-color: var(--background-highlight);
		border-radius: 0;
		margin: 0;
		color: inherit;
		border-color: var(--border-color);
	}

</style>

<div class="box">
	<input class="h1" placeholder="Object name.." type="text" bind:value={newModelOp.name}/>
	<h2>Attributes</h2>
	<HLTable>
		{#each newModelOp.attributes.create as attr}
			<HLRow>
				<div class="hform-row">
					<div class="col">
					<input 
					     class="hl-row" 
					     placeholder="Attribute name.." 
					     type="text" 
					     on:input={(e) => { attr.name = restrict(e.target.value) }}
					     bind:value={attr.name} 
					     />
					</div>
				<div class="spacer"/>
					<div class="col">
					<select bind:value={attr.attrType}>
						{#each AttrType as id}
						<option value={id}>
							{AttrType[id]}
						</option>
						{/each}
					</select>
					</div>
				</div>
			</HLRow>
		{/each}
		<HLRowButton on:click={addAttribute}>
			+add
		</HLRowButton>
	</HLTable>
	<h2>Relationships</h2>
	<HLTable>
		<HLRowButton>+add</HLRowButton>
	</HLTable>
		<div class="v-space-2"></div>
	<HLTable>
		<HLRowButton on:click={saveModel}>
				Save
		</HLRowButton>
	</HLTable>
</div>
