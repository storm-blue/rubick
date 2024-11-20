package common

import "strings"

// FindFirstParenthesesPair find next () pair
func FindFirstParenthesesPair(expression string) (left, right int) {
	// for example:
	// *)**(***(**))***
	//     |       |
	//     |       |
	//   left    right

	return FindFirstPair(expression, '(', ')')
}

// FindFirstBracesPair find next {} pair
func FindFirstBracesPair(expression string) (left, right int) {
	// for example:
	// *}**{***{**}}***
	//     |       |
	//     |       |
	//   left    right

	return FindFirstPair(expression, '{', '}')
}

// FindFirstBracketsPair find next [] pair
func FindFirstBracketsPair(expression string) (left, right int) {
	// for example:
	// *]**[***[**]]***
	//     |       |
	//     |       |
	//   left    right

	return FindFirstPair(expression, '[', ']')
}

// FindFirstAngleBracketsPair find next <> pair
func FindFirstAngleBracketsPair(expression string) (left, right int) {
	// for example:
	// *>**<***<**>>***
	//     |       |
	//     |       |
	//   left    right

	return FindFirstPair(expression, '<', '>')
}

func FindFirstPair(expression string, leftChar, rightChar int32) (left, right int) {
	left = -1
	right = -1
	braceNumber := 0
	findLeft := false

	for i, c := range expression {
		if c == leftChar {
			if !findLeft {
				left = i
				findLeft = true
			}
			braceNumber++
		} else if c == rightChar {
			if findLeft {
				braceNumber--
				if braceNumber == 0 {
					right = i
					break
				}
			}
		}
	}
	if right == -1 {
		left = -1
	}
	return
}

func UnwrapIfNeeded(expression string) string {
	return UnwrapParenthesesIfNeeded(expression)
}

func UnwrapParenthesesIfNeeded(expression string) string {
	expression = strings.TrimSpace(expression)

	left, right := FindFirstParenthesesPair(expression)
	if left == 0 && right == len(expression)-1 {
		expression = expression[1 : len(expression)-1]
		return UnwrapParenthesesIfNeeded(expression)
	} else {
		return expression
	}
}

func UnwrapBracesIfNeeded(expression string) string {
	expression = strings.TrimSpace(expression)

	left, right := FindFirstBracesPair(expression)
	if left == 0 && right == len(expression)-1 {
		expression = expression[1 : len(expression)-1]
		return UnwrapBracesIfNeeded(expression)
	} else {
		return expression
	}
}

func UnwrapBracketsIfNeeded(expression string) string {
	expression = strings.TrimSpace(expression)

	left, right := FindFirstBracketsPair(expression)
	if left == 0 && right == len(expression)-1 {
		expression = expression[1 : len(expression)-1]
		return UnwrapBracketsIfNeeded(expression)
	} else {
		return expression
	}
}

func UnwrapAngleBracketsIfNeeded(expression string) string {
	expression = strings.TrimSpace(expression)

	left, right := FindFirstAngleBracketsPair(expression)
	if left == 0 && right == len(expression)-1 {
		expression = expression[1 : len(expression)-1]
		return UnwrapAngleBracketsIfNeeded(expression)
	} else {
		return expression
	}
}

func UnwrapQuotaIfNeeded(expression string) string {
	expression = strings.TrimSpace(expression)
	if len(expression) < 2 {
		return expression
	}
	if strings.HasPrefix(expression, "\"") && strings.HasSuffix(expression, "\"") {
		expression = expression[1 : len(expression)-1]
		return expression
	} else {
		return expression
	}
}
