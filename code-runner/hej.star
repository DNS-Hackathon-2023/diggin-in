
def loop():

    results = measure.dig("hannig.cc")

    foo = state.get("foo", 23)

    print("Current state:", foo)

    for answer in results:
        print(answer)

    foo += 1
    state.set("foo", foo)

    collect("hannig.cc", results)

    return 42


