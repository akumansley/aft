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
const newModel = {
	name: "",
	attributes: [],
	relationships: [],
}
function addAttribute() {
	newModel.attributes = [...newModel.attributes, {
		name: "",
		attrType: 0,
	}];
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
		margin: 0;
		background: var(--background-highlight);
		color: inherit;
		border-color: var(--border-color);
		height: 2em;
	}

</style>

<div class="box">
	<input class="h1" placeholder="Object name.." type="text" bind:value={newModel.name}/>
	<h2>Attributes</h2>
	<HLTable>
		{#each newModel.attributes as attr}
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
						{#each Object.entries(AttrType) as atEntry}
						<option value={atEntry[0]}>
							{atEntry[1]}
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
</div>
