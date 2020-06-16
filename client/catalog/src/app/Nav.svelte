<style>

	.nav {
		height: 100%;
	}
</style>
<script>
	import LinkList from '../ui/LinkList.svelte';
	let items = [
		{name:"Models", path:'/models'}, 
		{name:"Datatypes", path: '/datatypes'},
		{name:"Repl", path: '/repl'},
		{name:"Log", path:"/log"}
	];
</script>

<div class="nav">
	<LinkList {items} />
</div>
