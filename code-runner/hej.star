def loop():
  state.set("foo", state.get("foo", 23) + 1)
  data = {
    "foo": 42,
    "msg": "enjoy ripe!",
    "state": state.get("foo"),
  }
  collect("hello ripe!")
