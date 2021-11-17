using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;

namespace Challenges.BreakingBestWorstRecords
{
    public class Result
    {
        public static List<int> breakingRecords(List<int> scores)
        {
            var first = scores.First();

            int high = first;
            int highChanged = 0;

            int low = first;
            int lowChanged = 0;

            foreach (var score in scores) {
                if (score > high) {
                    highChanged++;
                    high = score;
                } else if (score < low) {
                    lowChanged++;
                    low = score;
                }
            }

            return new List<int>{ highChanged, lowChanged };
        }

        public class Examples : List<Examples.Example> {
            public class Example {
                public List<int> Input {get;set;}
                public List<int> Output { get;set;}
            }
        }

        public static void Test()
        {
            var example = new Examples.Example {
                Input = new List<int> { 12, 24, 10, 24 },
                Output = new List<int> { 1, 1 },
            };

            var example_00 = new Examples.Example {
                Input = new List<int>() { 10, 5, 20, 20, 4, 5, 2, 25, 1 },
                Output = new List<int> { 2, 4 }
            };
            
            var example_01 = new Examples.Example {
                Input = new List<int>() { 3, 4, 21, 36, 10, 28, 35, 5, 24, 42 },
                Output = new List<int> { 4, 0 }
            };

            var tests = new Examples
            {
                example,
                example_00,
                example_01,
            };

            List<int> output;
            foreach (var test in tests) {
                output = breakingRecords(test.Input);
                Console.WriteLine($"output {string.Join(',', output)}");
                Debug.Assert(test.Output.SequenceEqual(output));
            }

        }
    }
}

/*
Scores are in the same order as the games played. She tabulates her results as follows:
                                     Count
    Game  Score  Minimum  Maximum   Min Max
     0      12     12       12       0   0
     1      24     12       24       0   1
     2      10     10       24       1   1
     3      24     10       24       1   1
*/