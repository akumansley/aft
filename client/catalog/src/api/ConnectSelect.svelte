<script context="module">
	import {routeStore} from '../app/stores.js';
	// cache promises
	let cache = {};

	async function getOptions(iface) {
		if (!cache[iface]) {
			cache[iface] = client.api[iface].findMany({where:{}});
		}
		return cache[iface];
	}
	routeStore.subscribe(route => {
		cache = {};
	});
</script>

<script>
	import client from '../data/client.js';
	import { onMount } from 'svelte';
	import HLSelect from '../ui/form/HLSelect.svelte';
	import {restrictToIdent, cap, isObject} from '../lib/util.js';

	export let iface = null;
	export let except = [];
	export let displayKey = "name";
	let options;

	export let value = null;

	let loaded = new Promise((resolve) => {
		onMount(async () => {
			options = await getOptions(iface);

			if (!value) {
				value = options[0];
			} else {
				for (let option of options){
					if (option.id === value.id) {
						value = option;
					}
				}
			}

			resolve();
		});
	});

</script>

{#await loaded then _}
<HLSelect bind:value={value}>
	{#each options as opt}
	{#if !except.includes(opt)}
	<option value={opt}>
		{cap(opt[displayKey])}
	</option>
	{/if}
	{/each}
</HLSelect>
{/await}
