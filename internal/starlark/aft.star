def makeAPIFunction(ctx, name):
	f = loadFunction(ctx, name)
	def apiFunc(modelName, body):
		return f(ctx, modelName, body)
	return apiFunc

def preamble(ctx):
    api = struct(
    	findOne=makeAPIFunction(ctx, "findOne"),
    	findMany=makeAPIFunction(ctx, "findMany"),
    	count=makeAPIFunction(ctx, "count"),
    	delete=makeAPIFunction(ctx, "delete"),
    	deleteMany=makeAPIFunction(ctx, "deleteMany"),
    	update=makeAPIFunction(ctx, "update"),
    	updateMany=makeAPIFunction(ctx, "updateMany"),
    	create=makeAPIFunction(ctx, "create"),
    	upsert=makeAPIFunction(ctx, "upsert"),
    	)
    # auth = struct(login=loadFunction(ctx, "login"))
    aft = struct(api=api)
    return aft
