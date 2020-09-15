<script context="module">
	import client from '../../data/client.js';
	let load = client.api.datatype.findMany({where:{}});
</script>

<script>
	import { onMount } from 'svelte';
	import HLSelect from '../../ui/form/HLSelect.svelte';
	import {restrictToIdent, cap, isObject} from '../../lib/util.js';

	export let init = null;
	export let value = null;
	export let op;
	let datatypes;

	let loaded = new Promise((resolve) => {
		onMount(() => {
			load.then((dts) => {
				datatypes = dts;

				if (!init) {
					value = dts[0];
				} else {
					for (let dt of dts){
						if (dt.id === init.id) {
							value = dt;
						}
					}
				}
			});
			resolve();
		});
	});


	function setOp() {
		if (init && value.id === init.id) {
			op = {};
		} else if (value) {
			op = {
				connect: {id: value.id},
			}
		}
	}
	$: setOp(value);

</script>

{#await loaded then _}
<HLSelect bind:value={value}>
	{#each Object.entries(datatypes) as dt}
	<option value={dt[1]}>
		{cap(dt[1].name)}
	</option>
	{/each}
</HLSelect>
{/await}
