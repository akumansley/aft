<script>
	export let items = [];
</script>

<style>
	ul {
		margin: 0;
		padding-top: 1em;
		padding-bottom: 1em;
		padding-left: 0;
		padding-right: 0;
		border-right: 1px solid rgba(0,0,0,.12);
	}
	li:first-child {
	}
	li {
		list-style-type: none;
		padding-left: 1.5em;
		line-height: 2;
	}
	a {
		color: #212121;
	}
</style>

<ul>
{#each items as item}
	<li><a href="#test">{item}</a></li>
{/each}
</ul>

