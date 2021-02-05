<script >
	export let href = null;
	export let click = null;
	import {cap} from '../../lib/util.js';
</script>
<style>
	a.object-box {
		display: flex;
		flex-direction: column;
		color: inherit;
		width: 12em;
		padding: var(--box-margin);
	}
	a.object-box:hover {
		background: var(--background-highlight);
	}

</style>
<a href="{href}" class="object-box" on:click={click}>
	<slot />
</a>
