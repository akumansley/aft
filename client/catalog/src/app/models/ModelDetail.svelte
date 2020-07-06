<script>
export let params;
import client from '../../data/client.js';
import Model from './Model.svelte';
import HLBox from '../../ui/HLBox.svelte';
import { breadcrumbStore } from '../stores.js';
import {cap} from '../util.js';

let id = params.id;
let load = client.api.model.findOne({where: {id: id}, include: {rightRelationships: true, leftRelationships: true, attributes: true}});

load.then(obj => {
breadcrumbStore.set(
	[{
		href: "/models",
		text: "Models",
	}, {
		href: "/model/" + id,
		text: cap(obj.name),
	}]
);
});
</script>

<HLBox>
	{#await load}
		&nbsp;
	{:then model}
		<Model model={model}/>
	{:catch error}
		<div>Error..</div>
	{/await}
</HLBox>

