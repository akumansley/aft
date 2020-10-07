package rpc

import (
    "context"
    "fmt"
    "strconv"
    "strings"

    "awans.org/aft/internal/db"
    "awans.org/aft/internal/starlark"
    "go.starlark.net/resolve"
    gstarlark "go.starlark.net/starlark"
    "go.starlark.net/syntax"
)

var ErrInvalidInput = fmt.Errorf("Bad input:")

var terminalRPC = db.MakeNativeFunction(
    db.MakeID("180bb262-8835-4ec5-9c2b-3f455615be9a"),
    "terminal",
    2,
    terminal,
)

//     `# Oh we really really need to make this secure
// # BIG SCARY COMMENTS
// # MASSIVE NEED FOR PERMISSIONS HERE
// def main(args):
//     out, ran = exec(args["data"], "")
//     if not ran:
//         return "Starlark: " + out.strip(":")
//     return out`,
func terminal(args []interface{}) (interface{}, error) {
    ctx := args[0].(context.Context)
    input := args[1]

    rpcData := input.(map[string]interface{})
    code, ok := rpcData["data"].(string)

    if !ok {
        return "", fmt.Errorf("%w code was type %T", ErrInvalidInput, code)
    }

    sr := starlark.NewStarlarkRuntime()

    r, err := sr.Execute(code, []interface{}{ctx})
    if err != nil {
        if evale, ok := err.(*gstarlark.EvalError); ok {
            return evale.Backtrace(), nil
        }
        return fmt.Sprintf("%s", err), nil
    }
    if r == nil {
        return "", nil
    }
    return fmt.Sprintf("%v", r), nil

}

var stdStrings = map[string]struct{}{
    "aft":          struct{}{},
    "loadFunction": struct{}{},
    "re":           struct{}{},
    "urlparse":     struct{}{},
}

func parse(code interface{}) (string, bool, error) {
    if input, ok := code.(string); ok {
        f, err := syntax.Parse("", input, 0)
        if err != nil {
            return fmt.Sprintf("%s", err), false, nil
        }
        var isPredeclared = func(s string) bool {
            _, ok := stdStrings[s]
            return ok
        }
        err = resolve.File(f, isPredeclared, gstarlark.Universe.Has)
        if err != nil {
            return fmt.Sprintf("%s", err), false, nil
        }
        return "", true, nil
    }
    return "", false, fmt.Errorf("%w code was type %T", ErrInvalidInput, code)
}

func parseRPCFunc(args []interface{}) (interface{}, error) {
    input := args[1]

    rpcData := input.(map[string]interface{})
    code, ok := rpcData["data"].(string)

    if !ok {
        return "", fmt.Errorf("%w code was type %T", ErrInvalidInput, code)
    }

    msg, ok, err := parse(code)
    if err != nil {
        return nil, err
    }
    return map[string]interface{}{
        "error":  msg,
        "parsed": ok,
    }, nil
}

var parseRPC = db.MakeNativeFunction(
    db.MakeID("232d7ad5-357b-43fb-a707-a0a6ba190e7c"),
    "parse",
    2,
    parseRPCFunc,
)

func lintRPCFunc(args []interface{}) (interface{}, error) {
    input := args[1]
    rpcData := input.(map[string]interface{})
    code, ok := rpcData["data"].(string)

    if !ok {
        return "", fmt.Errorf("%w code was type %T", ErrInvalidInput, code)
    }

    msg, ok, err := parse(code)
    if ok {
        return map[string]interface{}{
            "error":  msg,
            "parsed": ok,
        }, nil
    }
    parts := strings.Split(msg, ":")
    line, err := strconv.Atoi(parts[1])
    if err != nil {
        return nil, err
    }
    start, err := strconv.Atoi(parts[2])
    if err != nil {
        return nil, err
    }

    // seems like panicy code
    return map[string]interface{}{
        "message": parts[3],
        "parsed":  ok,
        "line":    line,
        "start":   start,
        "raw":     msg,
    }, nil
}

var lintRPC = db.MakeNativeFunction(
    db.MakeID("e4be72dc-9462-49f7-bba9-3543cc6bf6c2"),
    "lint",
    2,
    lintRPCFunc,
)
