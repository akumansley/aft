<style>

	.nav {
		height: 100%;
	}
</style>
<script>
	import LinkList from '../ui/LinkList.svelte';
	let items = [
		{name:"Objects", path:'/objects'}, 
		{name:"Functions", path: '/functions'}, 
		{name:"Roles", path:"/roles"},
		{name:"Queries", path:"/queries"}
	];
</script>

<div class="nav">
	<LinkList {items} />
</div>
