// Compile with `gcc foo.c -Wall -std=gnu99 -lpthread`, or use the makefile
// The executable will be named `foo` if you use the makefile, or `a.out` if you use gcc directly

#include <pthread.h>
#include <stdio.h>

pthread_mutex_t mutex;
//using mutex since it guarantees that only one thread can own the resource at a time.
int k = 0;

// Note the return type: void*
void* incrementingThreadFunction(){
    
    for (int i = 0; i < 1000000; i++)
    {
        pthread_mutex_lock(&mutex);
        k++;
        pthread_mutex_unlock(&mutex);
    }
    return NULL;
}

void* decrementingThreadFunction(){
    
    for (int i = 0; i < 1000000; i++)
    {

        pthread_mutex_lock(&mutex);
        k--;
        pthread_mutex_unlock(&mutex);
    }
    return NULL;
}

 
int main(){
    // TODO: 
    // start the two functions as their own threads using `pthread_create`
    // Hint: search the web! Maybe try "pthread_create example"?

    pthread_t thread1, thread2;
    pthread_create(&thread1, NULL, incrementingThreadFunction,"thr1");
    pthread_create(&thread2, NULL, decrementingThreadFunction, "thr2");
    
    // TODO:
    // wait for the two threads to be done before printing the final result
    // Hint: Use `pthread_join`
    pthread_join(thread1,NULL);
    pthread_join(thread2,NULL);  
    
    pthread_mutex_destroy(&mutex);
    printf("The magic number is: %d\n", k);
    return 0;
}
