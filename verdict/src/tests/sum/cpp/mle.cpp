#include <bits/stdc++.h>
using namespace std;

int main() {
    long long a, b; 
    scanf("%lld %lld", &a, &b);
    int sz = 1e7;
    vector<int> v(sz, 1);
    partial_sum(v.begin(), v.end(), v.begin());
    printf("%d", v.back());
}
