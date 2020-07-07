<script>
	import { navStore } from './stores.js';
	import { canRoute } from '../app/router.js';
	import { faTerminal } from '@fortawesome/free-solid-svg-icons/faTerminal';
	import { faDatabase } from '@fortawesome/free-solid-svg-icons/faDatabase';
	import { faCalculator } from '@fortawesome/free-solid-svg-icons/faCalculator';
	import { faHdd } from '@fortawesome/free-solid-svg-icons/faHdd';
	import { faCode } from '@fortawesome/free-solid-svg-icons/faCode';
	import Icon from 'fa-svelte';


	let selected;
	navStore.subscribe(value => {
	
		selected = value;
	});
	let items = [
		{name:"Models", path:'/models', id:"model", icon:faDatabase}, 
		{name:"Datatypes", path: '/datatypes', id:"datatype", icon:faCalculator},
		{name:"Functions", path: '/rpcs',id:"rpc", icon:faCode},
		{name:"Terminal", path: '/terminal',id:"terminal", icon:faTerminal},
		{name:"Log", path:"/log",id:"log", icon:faHdd}
	];
	
</script>

<style>
	.nav {
		height: 100%;
	}
	ul {
		margin: 0;
		padding-top: 1em;
		padding-bottom: 1em;
		padding-left: 0;
		padding-right: 0;
		list-style-type: none;
	}
	li {
		padding-left: 1.5em;
		padding-bottom: .5em;
	}
	.nav-item {
		color: inherit;
		font-weight: 400;
		display:flex;
		align-items:center;
		flex-direction:row;
	}
	.active-li {
		font-weight: 600;
	}
	.non-active-icon {
		opacity: .8;
	}
	.icon {
		display:flex;
		align-items:center;
	}
	.space {
		width: .5em;
	}
	.noselect {
				   cursor:default;
	  -webkit-touch-callout: none; /* iOS Safari */
		-webkit-user-select: none; /* Safari */
		 -khtml-user-select: none; /* Konqueror HTML */
		   -moz-user-select: none; /* Old versions of Firefox */
			-ms-user-select: none; /* Internet Explorer/Edge */
				user-select: none; /* Non-prefixed version, currently
									  supported by Chrome, Edge, Opera and Firefox */
	}
</style>

<div class="nav">
	<ul>
	{#each items as item}
		<li>
			<div class="nav-item {selected ===item.id? 'active-li' : ''}">
				<div class="icon {selected ===item.id? '' : 'non-active-icon'}">
					<Icon icon={item.icon}></Icon>
				</div>
				<div class="space"/>
				<a href="{item.path}" class="noselect" on:click={canRoute}>{item.name}</a>		
			</div>
		</li>
	{/each}
	</ul>
</div>




