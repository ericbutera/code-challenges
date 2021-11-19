// https://www.hackerrank.com/challenges/the-birthday-bar/problem
// dotnet test --filter TheBirthdayBar

using System.Collections.Generic;
using Xunit;

namespace Challenges.TheBirthdayBar
{
    public class Input {
        public List<int> Squares;
        public int Day;
        public int Month;
    }

    public class Result 
    {
        /// <summary>
        /// 
        /// </summary>
        /// <param name="squares">numbers on each of the squares of chocolate</param>
        /// <param name="birthday"></param>
        /// <param name="birthMonth"></param>
        /// <returns>number of ways the bar can be divided</returns>
        public static int birthday(List<int> squares, int birthday, int birthMonth)
        {
            return 0;
        }
    }

    public class UnitTests
    {
        [Fact]
        public static void Input00()
        {
            var input00 = new Input {
                Squares = new List<int> { 1, 2, 1, 3, 2 },
                Day = 3,
                Month = 2
            };
            var output = Result.birthday(input00.Squares, input00.Day, input00.Month);
            Assert.Equal(2, output);
        }
    }
}

/*
input 00
1 2 1 3 2
3 2
output 00
2

input 01
1 1 1 1 1 1
3 2
output 01
0

input 02
4
4 1
output 02
1


*/