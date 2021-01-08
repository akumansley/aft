def makeFunction2(ctx, name):
	f = loadFunction(ctx, name)
	def apiFunc(arg1, arg2):
		return f(ctx, arg1, arg2)
	return apiFunc

def makeFunction1(ctx, name):
    f = loadFunction(ctx, name)
    def func(arg):
        return f(ctx, arg)
    return func

def makeFunction0(ctx, name):
    f = loadFunction(ctx, name)
    def func():
        return f(ctx)
    return func

def preamble(ctx):
    api = struct(
    	findOne=makeFunction2(ctx, "findOne"),
    	findMany=makeFunction2(ctx, "findMany"),
    	count=makeFunction2(ctx, "count"),
    	delete=makeFunction2(ctx, "delete"),
    	deleteMany=makeFunction2(ctx, "deleteMany"),
    	update=makeFunction2(ctx, "update"),
    	updateMany=makeFunction2(ctx, "updateMany"),
    	create=makeFunction2(ctx, "create"),
    	upsert=makeFunction2(ctx, "upsert"),
    	)

    auth = struct(
        authenticateAs=makeFunction1(ctx, "authenticateAs"),
        user=makeFunction0(ctx, "currentUser"),
        checkPassword=loadFunction(ctx, "checkPassword"),  # 3 args, but don't need to curry context
        )

    aft = struct(api=api, auth=auth)
    return aft

