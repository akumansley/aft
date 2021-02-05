<script>
	import {cap} from '../../lib/util.js';
	import { navStore } from '../stores.js';
	import client from '../../data/client.js';
	import {HLGrid, HLGridItem, HLGridNew, HLSecondary} from '../../ui/grid/grid.js';
	import {HLBorder, HLContent, HLSectionTitle, HLCallout, HLPad} from '../../ui/page/page.js';

	const RPC = "4b8db42e-d084-4328-a758-a76939341ffa";

	let appModule = null;
	let nativeModules = [];
	let load = client.api.module.findMany({
		where:{OR:[
			{functions:{some:{
				funcType: RPC,
			}}},
			{goPackage:""}]},
			include: {
				functions: {where:{funcType: RPC}}
			}})
	.then(result => {
		for (let mod of result) {
			if (mod.goPackage === "") {
				appModule = mod;
			} else {
				nativeModules.push(mod);
			}
		}
	});

	function urlFor(func) {
		return "/rpc/" + func.id
	}

	navStore.set("rpcs");
</script>

{#await load then load}
<HLGrid>
	<HLGridNew href={"/rpcs/new"}>Add RPC</HLGridNew>
</HLGrid>
<HLBorder/>

<HLContent>
	{#if appModule.functions.length}
	<HLPad>
		<HLSectionTitle>RPCs</HLSectionTitle>
		<HLGrid>
			{#each appModule.functions as func}
			<HLGridItem href={urlFor(func)}>
				<div>{cap(func.name)}</div>
				<HLSecondary>{appModule.name}</HLSecondary>
			</HLGridItem>
			{/each}
		</HLGrid>
	</HLPad>
	{:else}
	<HLCallout>
		<div>No application functions yet</div>
	</HLCallout>
	{/if}

	<HLBorder spaceBottom={true} />
	<HLSectionTitle>System</HLSectionTitle>
	<HLGrid>
		{#each nativeModules as mod}
		{#each mod.functions as func}
		<HLGridItem href={urlFor(func)}>
			<div>{cap(func.name)}</div>
			<HLSecondary>{mod.name}</HLSecondary>
		</HLGridItem>
		{/each}
		{/each}
	</HLGrid>
</HLContent>
{/await}