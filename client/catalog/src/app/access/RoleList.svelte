<script>
	import {navStore} from '../stores.js';
	import client from '../../data/client.js';
	import {cap} from '../../lib/util.js';

	import {HLGrid, HLGridItem, HLGridNew, HLSecondary} from '../../ui/grid/grid.js';
	import {HLBorder, HLContent, HLSectionTitle, HLCallout, HLPad} from '../../ui/page/page.js';


	let appModule = null;
	let nativeModules = [];
	let load = client.api.module.findMany({
		where:{OR:[
			{roles:{some:{}}},
			{goPackage: ""},
			]},
			include: {roles: true}})
	.then(result => {
		for (let mod of result) {
			if (mod.goPackage === "") {
				appModule = mod;
			} else {
				nativeModules.push(mod);
			}
		}
	});

	function urlFor(role) {
		return "/role/" + role.id
	}

	navStore.set("access");
</script>


{#await load then roles}
<HLGrid>
	<HLGridNew href={"/roles/new"}>
		Add Role
	</HLGridNew>
</HLGrid>
<HLBorder/>

<HLContent>
	{#if appModule.roles.length}
	<HLPad>
		<HLSectionTitle>Roles</HLSectionTitle>
		<HLGrid>
			{#each appModule.roles as role}
			<HLGridItem href={urlFor(role)}>
				<div>{cap(role.name)}</div>
				<HLSecondary>{appModule.name}</HLSecondary>
			</HLGridItem>
			{/each}
		</HLGrid>
	</HLPad>
	{:else}
	<HLCallout>
		<div>No application roles yet</div>
	</HLCallout>
	{/if}

	<HLBorder spaceBottom={true} />
	<HLSectionTitle>System</HLSectionTitle>
	<HLGrid>
		{#each nativeModules as mod}
		{#each mod.roles as role}
		<HLGridItem href={urlFor(role)}>
			<div>{cap(role.name)}</div>
			<HLSecondary>{mod.name}</HLSecondary>
		</HLGridItem>
		{/each}
		{/each}
	</HLGrid>
</HLContent>
{/await}

