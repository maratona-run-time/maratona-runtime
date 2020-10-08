#include "bits/stdc++.h"
using namespace std;

int main (){
    int a[10];
    int idx = rand();
    fprintf(stderr, "acessando idx %d\n", idx);
    printf("%d\n", a[idx]);
    a[-1] = -1;
    fflush(stdout);
    return 0;
}
