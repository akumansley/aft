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
<div>{print(value.OldRecord)} &rarr {print(value.NewRecord)}</div>
<div>{print(value.ModelID)}</div>
