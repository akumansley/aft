<script>
	import {VSpace} from "../spacing/spacing.js"
	export let spaceBottom = false;
	export let spaceTop = false;
</script>
<style>
	.border {
		border-bottom: 1px solid var(--border-color);
	}
</style>
{#if spaceTop}<VSpace/>{/if}
<div class="border"></div>
{#if spaceBottom}<VSpace/>{/if}
