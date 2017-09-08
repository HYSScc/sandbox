#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

int main(void) {
    int size = 0;
    int chunk_size = 1024 * 1024;
    void *p = NULL;

    while(1) {
        if ((p = malloc(chunk_size)) == NULL) {
            printf("Out of memory");
            break;
        }
        memset(p, 1, chunk_size);
        size += chunk_size;
        sleep(1);
    }
    return 0;
}