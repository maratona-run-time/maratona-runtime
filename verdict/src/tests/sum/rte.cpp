#include "bits/stdc++.h"
using namespace std;

int main (){
    long long a, b;
    scanf("%lld%lld", &a, &b);
    assert(max(a, b) < 50);
    printf("%lld\n", a + b);
    return 0;
}
