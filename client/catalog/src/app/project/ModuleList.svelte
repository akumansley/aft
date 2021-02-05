<script>
	import client from '../../data/client.js';
	import {cap} from '../../lib/util.js';
	import {navStore} from '../stores.js';

	import {HLGrid, HLGridItem, HLGridNew, HLSecondary} from '../../ui/grid/grid.js';
	import {HLBorder, HLContent, HLSectionTitle, HLPad} from '../../ui/page/page.js';

	let nativeModules = [];
	let localModules = [];
	let load = client.api.module.findMany({})
	.then(result => {
		for (let mod of result) {
			if (mod.goPackage === "") {
				localModules.push(mod);
			} else {
				nativeModules.push(mod);
			}
		}
	});
	navStore.set("modules");
	function elipse(str, count) {
		if (str.length > count) {
			let substr = str.slice(0, count)
			return substr + "…";
		}
		return str
	}
</script>

{#await load then load}
<HLGrid>
	<HLGridNew href={"/modules/new"}>
		Add Module
	</HLGridNew>
</HLGrid>
<HLBorder/>

<HLContent>
	<HLPad>
		<HLSectionTitle>Modules</HLSectionTitle>
		<HLGrid>
			{#each localModules as mod}
			<HLGridItem href={"/module/" + mod.id}>
				<div>{mod.name}</div>
				<HLSecondary>local</HLSecondary>
			</HLGridItem>
			{/each}
		</HLGrid>
	</HLPad>
	<HLBorder spaceBottom={true} />

	<HLSectionTitle>System</HLSectionTitle>
	<HLGrid>
		{#each nativeModules as mod}
		<HLGridItem href={"/module/" + mod.id}>
			<div>{mod.name}</div>
			<HLSecondary>{elipse(mod.goPackage, 26)}</HLSecondary>
		</HLGridItem>
		{/each}
	</HLGrid>
</HLContent>

{/await}