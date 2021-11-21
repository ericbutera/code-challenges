// https://www.hackerrank.com/challenges/the-birthday-bar/problem
// dotnet test --filter TheBirthdayBar

using System;
using System.Linq;
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
        /// <param name="birthday">Total sum to aim for</param>
        /// <param name="birthMonth">Number of squares to group</param>
        /// <returns>number of ways the bar can be divided</returns>
        public static int birthday(List<int> squares, int birthday, int birthMonth)
        {
            int chunkSize = birthMonth - 1; // TODO: bounds check
            int splitsFound = 0, 
                sum = 0, 
                attempt = 0;

            for (var x=0; x <= squares.Count; x++)  // iterate bar squares
            {
                sum = 0;
                for (var y=0; y <= chunkSize; y++) 
                { // iterate contiguous squares as a group of `chunkSize`
                    attempt = y + x;
                    if (attempt >= squares.Count) 
                        break;

                    sum += squares[attempt];
                }

                if (sum == birthday) 
                    splitsFound++;
            }

            return splitsFound;
        }

        public static int birthday2(List<int> squares, int birthday, int birthMonth)
        {
            int chunkSize = birthMonth
                , sum = 0
                , splitsFound = 0;

            for (var x=0; x <= squares.Count; x++)
            {
                sum = squares.Skip(x)
                    .Take(chunkSize)
                    .Sum();

                Console.WriteLine($"x:{x} sum:{sum} skip:{x} chunkSize:{chunkSize}");

                if (sum == birthday) 
                    splitsFound++; 
            }

            return splitsFound;
        }
    }

    public class UnitTests
    {
        [Fact]
        public static void Input00()
        {
            var input = new Input {
                Squares = new List<int> { 1, 2, 1, 3, 2 },
                Day = 3,
                Month = 2
            };
            Assert.Equal(2, Result.birthday(input.Squares, input.Day, input.Month));
            Assert.Equal(2, Result.birthday2(input.Squares, input.Day, input.Month));
        }

        [Fact]
        public static void Input01() 
        {
            var input = new Input {
                Squares = new List<int> { 1, 1, 1, 1, 1, 1 },
                Day = 3,
                Month = 2
            };
            Assert.Equal(0, Result.birthday(input.Squares, input.Day, input.Month));
            Assert.Equal(0, Result.birthday2(input.Squares, input.Day, input.Month));
        }

        [Fact]
        public static void Input02() 
        {
            var input = new Input {
                Squares = new List<int> { 4 },
                Day = 4,
                Month = 1
            };
            Assert.Equal(1, Result.birthday(input.Squares, input.Day, input.Month));
            Assert.Equal(1, Result.birthday2(input.Squares, input.Day, input.Month));
        }
    }
}