
def loop():

    results = measure.dig("hannig.cc")


    foo = state.get("foo")
    if not foo:
        foo = 1

    print("Current state:", foo)

    for answer in results:
        print(answer["message"])


    foo += 1

    state.set("foo", foo)

    return 42


