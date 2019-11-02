# Assist pkg

This pkg contains testing helpers but for initialization.

Because of that, that type of funcs cannot contain *t testing.T arg on any function, so the errs are notified by log.Fatal (due every bad initialization causes bad processing).