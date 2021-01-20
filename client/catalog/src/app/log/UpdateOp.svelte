<script>
	export let value;
	function print(thing) {
		return JSON.stringify(thing, null, 2)
	}
	import Record from './Record.svelte';
</script>

<Record value={value.OldRecData}/> &darr; <Record value={value.NewRecData}/>
