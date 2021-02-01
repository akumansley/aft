<script>
	export let value;
	function print(thing) {
		return JSON.stringify(thing, null, 2)
	}
	import Record from './Record.svelte';
</script>

<Record value={value.OldRecord}/> &darr; <Record value={value.NewRecord}/>
