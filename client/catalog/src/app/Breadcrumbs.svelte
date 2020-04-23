<style>
	.breadcrumbs {
		display:flex;
		align-content: flex-start;
	}
	.breadcrumb {
	}
	.spacer {
		width: .5em;
	}
</style>
<script>
import { breadcrumbStore } from './breadcrumbStore.js';
let breadcrumbs;
breadcrumbStore.subscribe(value => {
	breadcrumbs = value;
});
</script>

<div class="breadcrumbs">
	{#each breadcrumbs as bc, ix}
		<a class="breadcrumb" href="{bc.href}">{bc.text}</a>
		{#if ix != breadcrumbs.length - 1 }
		<div class="spacer"/>
		/
		<div class="spacer"/>
		{/if}
	{/each}

</div>
