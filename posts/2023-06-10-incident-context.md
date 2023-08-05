# Another day another incident #02

Today's incident is all about Go context.

TLDR: context got canceled, but it shouldn't.

## The problem

Imagine a chain of APIs:

- Calling API A
- Calling API B

Normally, if API A fails, API B should not be called. But what if API A is
**optional**, whether it successes or fails, API B should be called anyway.

My buggy code is like this:

```go
if err := doA(ctx); err != nil {
    log.Error(err)
    // Skip error
}

doB(ctx)
```

The problem is `doA` taking too long, so `ctx` is canceled, and the parent of
`ctx` is canceled too. So when `doB` is called with `ctx`, it will be canceled
too (not what we want but sadly that what we got).

Example buggy code ([The Go Playground](https://go.dev/play/p/p4S27Su16VH)):

```go
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	doA(ctx)
	doB(ctx)
}

func doA(ctx context.Context) {
	ctx, ctxCancel := context.WithTimeout(ctx, 1*time.Second)
	defer ctxCancel()

	select {
	case <-time.After(2 * time.Second):
		fmt.Println("doA")
	case <-ctx.Done():
		fmt.Println("doA", ctx.Err())
	}
}

func doB(ctx context.Context) {
	ctx, ctxCancel := context.WithTimeout(ctx, 3*time.Second)
	defer ctxCancel()

	select {
	case <-time.After(2 * time.Second):
		fmt.Println("doB")
	case <-ctx.Done():
		fmt.Println("doB", ctx.Err())
	}
}
```

The output is:

```txt
doA context deadline exceeded
doB context deadline exceeded
```

As you see both `doA` and `doB` are canceled.

## The (temporary) solution

Quick Google search leads me to
[context: add WithoutCancel #40221](https://github.com/golang/go/issues/40221)
and I quote:

> This is useful in multiple frequently recurring and important scenarios:
>
> - Handling of rollback/cleanup operations in the context of an event (e.g.,
  > HTTP request) that has to continue regardless of whether the triggering
  > event is canceled (e.g., due to timeout or the client going away)
> - Handling of long-running operations triggered by an event (e.g., HTTP
  > request) that terminates before the termination of the long-running
  > operation

So beside waiting to upgrade to Go `1.21` to use `context.WithoutCancel`, you
can use this [workaround code](https://pkg.go.dev/context@master#WithoutCancel):

```go
func DisconnectContext(parent context.Context) context.Context {
	if parent == nil {
		return context.Background()
	}

	return disconnectedContext{
		parent: parent,
	}
}

type disconnectedContext struct {
	parent context.Context
}

func (ctx disconnectedContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (ctx disconnectedContext) Done() <-chan struct{} {
	return nil
}

func (ctx disconnectedContext) Err() error {
	return nil
}

func (ctx disconnectedContext) Value(key any) any {
	return ctx.parent.Value(key)
}
```

So the buggy code becomes
([The Go Playground](https://go.dev/play/p/oIU-WxEJ_F3)):

```go
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	doA(ctx)
	doB(ctx)
}

func doA(ctx context.Context) {
	ctx, ctxCancel := context.WithTimeout(ctx, 1*time.Second)
	defer ctxCancel()

	select {
	case <-time.After(2 * time.Second):
		fmt.Println("doA")
	case <-ctx.Done():
		fmt.Println("doA", ctx.Err())
	}
}

func doB(ctx context.Context) {
	ctx, ctxCancel := context.WithTimeout(DisconnectContext(ctx), 3*time.Second)
	defer ctxCancel()

	select {
	case <-time.After(2 * time.Second):
		fmt.Println("doB")
	case <-ctx.Done():
		fmt.Println("doB", ctx.Err())
	}
}
```

The output is:

```txt
doA context deadline exceeded
doB
```

As you see only `doA` is canceled, `doB` is done perfectly. And that what we
want in this case.

## Thanks

- [Using Context in Golang - Cancellation, Timeouts and Values (With Examples)](https://www.sohamkamani.com/golang/context-cancellation-and-values/)
- [Go Context timeouts can be harmful](https://uptrace.dev/blog/golang-context-timeout.html)
