<script>
export let relationship;
export let models;
export let modelName;
import { restrictToIdent, isObject } from '../util.js';
import HLRow from '../../ui/list/HLRow.svelte';
import HLSelect from '../../ui/form/HLSelect.svelte';
import HLButton from '../../ui/form/HLButton.svelte';
import HLText from '../../ui/form/HLText.svelte';
import { RelType } from '../../data/enums.js';


let rightModel = models[0];
for(var i = 0; i < models.length; i++) {
	if(models[i].name === relationship.rightName) {
		rightModel = models[i];
	}
}
if (isObject(relationship.rightModel)) {
	relationship.rightModel.connect.id = rightModel.id;
}

let showDetail = false;
function toggle() {
	showDetail = !showDetail;
}

</script>
<style>
.hform-row {
	display: flex; 
	flex-direction: row;
}
.spacer {
	width: 1em;
}
.v-space{
	height: 2.5em;
}

.form-grid {
	display:grid;
	column-gap: 1em;
	row-gap: .5em;
	align-items: center;
	grid-template-columns: auto auto auto;
	grid-template-rows: auto auto;
	grid-template-areas:
		"model-name rel-type rel-name"
		"model-name rel-type rel-name";
}
.model-name{
	text-align: right;
}
</style>
<HLRow>
	<HLText 
		bind:value={relationship.name}
		placeholder="Relationship name.." 
		restrict={restrictToIdent}
		/>
		<span class="spacer"/>
		<HLSelect bind:value={rightModel}>
				{#each models as m}
				<option value={m}>
					{m.name}
				</option>
				{/each}
		</HLSelect>
		<span class="spacer"/>

		{#if !showDetail}
		<HLButton on:click={e => toggle()}> More </HLButton>
		{:else}
		<HLButton on:click={e => toggle()}> Hide </HLButton>
		{/if}
</HLRow>
{#if showDetail}
<HLRow indent=1>
	<div class="form-grid">
		<div class="model-name">
			{modelName || "This model"}
		</div>
		<div class="rel-type">
			<HLSelect bind:value={relationship.leftBinding}>
					{#each Object.entries(RelType) as rt, ix}
					<option value={ix}>
						{rt[1]}
					</option>
					{/each}
			</HLSelect>
		</div>
		<div class="rel-name">
			<HLText 
				bind:value={relationship.leftName}
				placeholder="Relationship name.." 
				restrict={restrictToIdent}
				/>
		</div>
		<div class="model-name">
			{#if rightModel}
				{rightModel.name}
			{:else}
				The other model
			{/if}
		</div>
		<div class="rel-type">
			<HLSelect bind:value={relationship.rightBinding}>
				{#each Object.entries(RelType) as rt, ix}
				<option value={ix}>
					{rt[1]}
				</option>
				{/each}
			</HLSelect>
		</div>
		<div class="rel-name">
			<HLText 
				bind:value={relationship.rightName}
				placeholder="Back-reference name.." 
				restrict={restrictToIdent}
				/>
		</div>
	</div>
</HLRow>
<div class="v-space"/>
{/if}

