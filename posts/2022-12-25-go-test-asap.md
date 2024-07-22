# Speed up writing Go test ASAP

Imagine your project currently have 0% unit test code coverage. And your boss
keep pushing it to 80% or even 90%? What do you do? Give up?

What if I tell you there is a way? Not entirely cheating but ... you know, there
is always trade off.

If your purpose is to test carefully all path, check if all return is correctly.
Sadly this post is not for you, I guess. If you only want good number on test
coverage, with minimum effort as possible, I hope this will show you some idea
you can use :)

In my opinion, unit test is not that important (like must must have). It's just
make sure your code is running exactly as you intent it to be. If you don't
think about edge case before, unit test won't help you.

## First, rewrite the impossible (to test) out

When I learn programming, I encounter very interesting idea, which become mainly
my mindset when I dev later. I don't recall it clearly, kinda like: "Don't just
fix bugs, rewrite it so that kind of bugs will not appear again". So in our
context, there is some thing we hardly or can not write test in Go. My
suggestion is don't use that thing.

In my experience, I can list a few here:

- Read config each time call func (`viper.Get...`). You can and you should init
  all config when project starts.
- Not use Dependency Injection (DI). There are too many posts in Internet tell
  you how to do DI properly.
- Use global var (Except global var `Err...`). You should move all global var to
  fields inside some struct.

## Let the fun (writing test) begin

If you code Go long enough, you know table driven tests and how is that so
useful. You set up test data, then you test. Somewhere in the future, you change
the func, then you need to update test data, then you good!

In simple case, your func only have 2 or 3 inputs so table drive tests is still
looking good. But real world is ugly (maybe not, idk I'm just too young in this
industry). Your func can have 5 or 10 inputs, also your func call many third
party services.

Imagine having below func to upload image:

```go
type service struct {
    db DB
    redis Redis
    minio MinIO
    logService LogService
    verifyService VerifyService
}

func (s *service) Upload(ctx context.Context, req Request) error {
    // I simplify by omitting the response, only care error for now
    if err := s.verifyService.Verify(req); err != nil {
        return err
    }

    if err := s.minio.Put(req); err != nil {
        return err
    }

    if err := s.redis.Set(req); err != nil {
        return err
    }

    if err := s.db.Save(req); err != nil {
        return err
    }

    if err := s.logService.Save(req); err != nil {
        return err
    }

    return nil
}
```

With table driven test and thanks to
[stretchr/testify](https://github.com/stretchr/testify), I usually write like
this:

```go
type ServiceSuite struct {
    suite.Suite

    db DBMock
    redis RedisMock
    minio MinIOMock
    logService LogServiceMock
    verifyService VerifyServiceMock

    s service
}

func (s *ServiceSuite) SetupTest() {
    // Init mock
    // Init service
}

func (s *ServiceSuite) TestUpload() {
    tests := []struct{
        name string
        req Request
        verifyErr error
        minioErr error
        redisErr error
        dbErr error
        logErr error
        wantErr error
    }{
        {
            // Init test case
        }
    }

    for _, tc := range tests {
        s.Run(tc.name, func(){
            // Mock all error depends on test case
            if tc.verifyErr != nil {
                s.verifyService.MockVerify().Return(tc.verifyErr)
            }
            // ...

            gotErr := s.service.Upload(tc.req)
            s.Equal(wantErr, gotErr)
        })
    }
}
```

Looks good right? Be careful with this. It can go from 0 to 100 ugly real quick.

What if req is a struct with many fields? So in each test case you need to set
up req. They are almost the same, but with some error case you must alter req.
It's easy to be init with wrong value here (typing maybe ?). Also all req looks
similar, kinda duplicated.

```go
tests := []struct{
        name string
        req Request
        verifyErr error
        minioErr error
        redisErr error
        dbErr error
        logErr error
        wantErr error
    }{
        {
            req: Request {
                a: "a",
                b: {
                    c: "c",
                    d: {
                        "e": e
                    }
                }
            }
            // Other fieles
        },
         {
            req: Request {
                a: "a",
                b: {
                    c: "c",
                    d: {
                        "e": e
                    }
                }
            }
            // Other fieles
        },
         {
            req: Request {
                a: "a",
                b: {
                    c: "c",
                    d: {
                        "e": e
                    }
                }
            }
            // Other fieles
        }
    }
```

What if dependencies of service keep growing? More mock error to test data of
course.

```go
tests := []struct{
    name string
    req Request
    verifyErr error
    minioErr error
    redisErr error
    dbErr error
    logErr error
    wantErr error
    // Murr error
    aErr error
    bErr error
    cErr error
    // ...
}{
    {
        // Init test case
    }
}
```

The test file keep growing longer and longer until I feel sick about it.

See
[tektoncd/pipeline unit test](https://github.com/tektoncd/pipeline/blob/main/pkg/pod/pod_test.go)
to get a feeling about this. When I see it, `TestPodBuild` has almost 2000
lines.

The solution I propose here is simple (absolutely not perfect, but good with my
usecase) thanks to **stretchr/testify**. I init all **default** action on
**success** case. Then I **alter** request or mock error for unit test to hit on
other case. Remember if unit test is hit, code coverage is surely increased, and
that my **goal**.

```go
// Init ServiceSuite as above

func (s *ServiceSuite) TestUpload() {
    // Init success request
    req := Request{
        // ...
    }

    // Init success action
    s.verifyService.MockVerify().Return(nil)
    // ...

    gotErr := s.service.Upload(tc.req)
    s.NoError(gotErr)

    s.Run("failed", func(){
        // Alter failed request from default
        req := Request{
            // ...
        }

        gotErr := s.service.Upload(tc.req)
        s.Error(gotErr)
    })

    s.Run("another failed", func(){
        // Alter verify return
        s.verifyService.MockVerify().Return(someErr)


        gotErr := s.service.Upload(tc.req)
        s.Error(gotErr)
    })

    // ...
}
```

If you think this is not quick enough, just **ignore** the response. You only
need to check error or not if you want code coverage only.

So if request change fields or more dependencies, I need to update success case,
and maybe add corresponding error case if need.

Same idea but still with table, you can find here
[Functional table-driven tests in Go - Fatih Arslan](https://arslan.io/2022/12/04/functional-table-driven-tests-in-go/).
