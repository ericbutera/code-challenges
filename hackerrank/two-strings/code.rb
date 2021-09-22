# https://www.hackerrank.com/challenges/two-strings

##
# Given two strings, determine if they share a common substring. 
# A substring may be as small as one character.  
# +s1+ first string
# +s2+ second string
# returns 'YES' if share a common substring, else 'NO'
def twoStrings(s1, s2)
    # since YES may be as small as one character, compare only characters, not strings
    # iterate chars in s1, test s2 to see if contains char

    # another way could be take rolling substrings of s1 and use s2.casecmp >= 0 (or not -1, nil)
    # h e l l o 
    # he el ll lo 
    # hel ell llo
    # hell ello

    s1.each_char do |char| 
        return 'YES' if s2.include?(char)
    end

    'NO'
end