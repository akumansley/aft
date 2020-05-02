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
		{name:"Users", path:"/users"},
		{name:"Log", path:"/log"}
	];
</script>

<div class="nav">
	<LinkList {items} />
</div>
