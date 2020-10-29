<script>
	export let value;
	function print(thing) {
		return JSON.stringify(thing, null, 2)
	}
</script>
<style>
	.op-closed {
		overflow: hidden;
		text-overflow: ellipsis;
	}
</style>

<div>Update</div>
<div>{print(value.OldRecFields)} &rarr {print(value.NewRecFields)}</div>
<div>{print(value.ModelID)}</div>
