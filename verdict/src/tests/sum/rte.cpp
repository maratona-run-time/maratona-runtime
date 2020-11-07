#include "bits/stdc++.h"
using namespace std;

int main (){
    long long num[20];
    for(int a=0;a<20;a++)
        num[a] = a;
    long long a, b;
    scanf("%lld%lld", &a, &b);
    printf("%lld\n", num[a] + num[b]);
    fflush(stdout);
    return 0;
}
