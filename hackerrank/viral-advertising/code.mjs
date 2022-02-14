// https://www.hackerrank.com/challenges/strange-advertising
import assert from 'assert/strict';

function viralAdvertising(n) {
  let cumulative = 0;
  let shared = 5;
  const SHARE_FRIENDS = 3;

  for (let day = 1; day <= n; day++) {
    let liked = Math.floor(shared / 2);

    shared = liked * SHARE_FRIENDS;
    cumulative += liked;

  }

  return cumulative;
}

/*
when launch new product, advertise to `5` people

day 1
half of the 5 like the ad Math.floor(5 / 2)
share with 3 friends

day 2
floor(5/2) * 3 = 2 * 3 = 6 people receive the ad

each day floor(recipients/2) like and will share with `3` friends on the next day

nobody receives the ad twice

how many people have liked the ad by end of the day

beginning at day 1

n = 5
day shared    liked   cumulative
1   5         2       2
2   6         3       5
3   9         4       9
4   12        6       15
5   18        9       24

viralAdvertising(5) = 24

cumulative = 0
SHARE_FRIENDS = 3

day 1
liked = Math.floor(shared / 2) = 2
share = liked(2) * SHARE_FRIENDS(3) = 6
cumulative += liked(2) = 2

day 2
liked = floor(shared(6) / 2) = 3
share = liked(3) * SHARE_FRIENDS(3) = 9
cumulative += liked(3) = 5

day 3
liked = floor(shared(9) / 2) = 4
share = liked(4) * SHARE_FRIENDS(3) = 12
cumulative += liked(4) = 9

day 4
liked = floor(shared(12) / 2) = 6
share = liked(6) * SHARE_FRIENDS(3) = 18
cumulative += liked(6) = 15

day 5
liked = floor(shared(18) / 2) = 9
share = liked(9) * SHARE_FRIENDS(3) = 27
cumulative += liked(9) = 24

return cumulative

*/

assert.equal(2, viralAdvertising(1))
assert.equal(5, viralAdvertising(2))
assert.equal(24, viralAdvertising(5))
assert.equal(24, viralAdvertising(5))
