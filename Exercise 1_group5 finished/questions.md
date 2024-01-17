Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
>Concurrency refers to operations that are not actually done simultanously, but where the processor switches between different tasks.
    Parallel operations are actually done simultanously, but requires multiple processing units.

What is the difference between a *race condition* and a *data race*? 
> When two or more threads try to acces shared data at the same time, which leads to ambigous output since it can vary dependig on which thread accesses the data first.
>Data race is a undercategory of race conditions, that happens when two threads try to access the same variable concurrently and at least one does a write operation. This also leadsto inconsistent outputs.
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> The scheduler "decides" which thread that should be currently running. It changes between threads by: Saving the current state, then Loading the new state, switching memory context, and lastrly continuing cpu operation on the new thread.


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
> It is a way to use the same hardware resources more efficiently than if it was not used. Eventually every operation gets to a point where it has to wait for x cycles. Multithreading allows the processor to do another operation concurrently with the first one, allowing the porcessor to do productive work during wait-time.

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> They are a type of thread that is schedueled by the application as opposed to the OS. They may have less overhead for cases where many small tasks has to be done concurrently.

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> It is a way of solving the programmes problem of doing multiple things at the same time. In this regard it makes the programmers life easier, however on the other hand the concept may involve more thinking from the programmer and is thus harder. So it will probably be a combination. of harder and easier.

What do you think is best - *shared variables* or *message passing*?
> There is no "best solution". Either one can be a good way of doing it, as it depends on the context. Shared variables may be easier and rquire less overhead, but at the same time it may create bottlenecks and concurrency issues. You will not get concurrency issues with message passing, however it may add more complexity to the system than what is required.

Our questions:
Does channels make sure that a go-routine finishes before the code goes out of scope as here: https://go.dev/tour/concurrency/2