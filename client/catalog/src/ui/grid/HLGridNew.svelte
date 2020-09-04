<script >
	export let href = null;
	export let click = null;
	import { faPlus } from '@fortawesome/free-solid-svg-icons/faPlus';
	import Icon from 'fa-svelte';
</script>
<style>
	a.object-box {
		display:flex;
		align-items:center;
		flex-direction:row;
		color: inherit;
		width: 150px;
		padding: 1em 1.5em;
		font-weight:600;
	}
	a.object-box:hover {
		background: var(--background-highlight);
	}
	.space {
		width: .5em;
	}
	.items {
		display:flex;
		align-items:center;
	}
</style>
<a href="{href}" class="object-box" on:click={click}>
	<div class="items">
		<Icon icon={faPlus}></Icon>
		<div class="space"/>
		<slot />
	</div>
</a>
