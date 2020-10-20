def makeFunction3(ctx, name):
	f = loadFunction(ctx, name)
	def apiFunc(modelName, body):
		return f(ctx, modelName, body)
	return apiFunc

def makeFunction2(ctx, name):
    f = loadFunction(ctx, name)
    def func(arg):
        return f(ctx, arg)
    return func

def makeFunction1(ctx, name):
    f = loadFunction(ctx, name)
    def func():
        return f(ctx)
    return func

def preamble(ctx):
    api = struct(
    	findOne=makeFunction3(ctx, "findOne"),
    	findMany=makeFunction3(ctx, "findMany"),
    	count=makeFunction3(ctx, "count"),
    	delete=makeFunction3(ctx, "delete"),
    	deleteMany=makeFunction3(ctx, "deleteMany"),
    	update=makeFunction3(ctx, "update"),
    	updateMany=makeFunction3(ctx, "updateMany"),
    	create=makeFunction3(ctx, "create"),
    	upsert=makeFunction3(ctx, "upsert"),
    	)

    auth = struct(
        authenticateAs=makeFunction2(ctx, "authenticateAs"),
        user=makeFunction1(ctx, "currentUser")
        )

    aft = struct(api=api, auth=auth)
    return aft
