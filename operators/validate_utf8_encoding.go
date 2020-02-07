package operators

import (
	"github.com/asalih/guardian/matches"
)

//var int UNICODE_ERROR_DECODING_ERROR = 0

var UNICODE_ERROR_CHARACTERS_MISSING int = -1
var UNICODE_ERROR_INVALID_ENCODING int = -2
var UNICODE_ERROR_OVERLONG_CHARACTER int = -3
var UNICODE_ERROR_RESTRICTED_CHARACTER int = -4
var UNICODE_ERROR_DECODING_ERROR int = -5

func (opMap *OperatorMap) loadValidateUtf8Encoding() {
	opMap.funcMap["validateUtf8Encoding"] = func(expression interface{}, variableData interface{}) *matches.MatchResult {
		matchResult := matches.NewMatchResult(false)
		data := variableData.(string)

		i := 0
		bytes_left := len(data)

		for i < bytes_left {
			rc := detectUtf8Character(int(data[i]), bytes_left)

			if rc <= 0 {
				return matchResult.SetMatch(true)
			}

			i += rc
			bytes_left -= rc
		}

		return matchResult
	}
}

func detectUtf8Character(p_read int, length int) int {
	unicode_len := 0
	d := 0
	c := 0

	if p_read == 0 {
		return UNICODE_ERROR_DECODING_ERROR
	}
	c = p_read

	/* If first byte begins with binary 0 it is single byte encoding */
	if (c & 0x80) == 0 {
		/* single byte unicode (7 bit ASCII equivilent) has no validation */
		return 1
	} else if (c & 0xE0) == 0xC0 {
		/* If first byte begins with binary 110 it is two byte encoding*/
		/* check we have at least two bytes */
		if length < 2 {
			unicode_len = UNICODE_ERROR_CHARACTERS_MISSING
		} else if ((p_read + 1) & 0xC0) != 0x80 {
			/* check second byte starts with binary 10 */
			unicode_len = UNICODE_ERROR_INVALID_ENCODING
		} else {
			unicode_len = 2
			/* compute character number */
			d = ((c & 0x1F) << 6) | ((p_read + 1) & 0x3F)
		}
	} else if (c & 0xF0) == 0xE0 {
		/* If first byte begins with binary 1110 it is three byte encoding */
		/* check we have at least three bytes */
		if length < 3 {
			unicode_len = UNICODE_ERROR_CHARACTERS_MISSING
		} else if ((p_read + 1) & 0xC0) != 0x80 {
			/* check second byte starts with binary 10 */
			unicode_len = UNICODE_ERROR_INVALID_ENCODING
		} else if ((p_read + 2) & 0xC0) != 0x80 {
			/* check third byte starts with binary 10 */
			unicode_len = UNICODE_ERROR_INVALID_ENCODING
		} else {
			unicode_len = 3
			/* compute character number */
			d = ((c & 0x0F) << 12) | (((p_read + 1) & 0x3F) << 6) | ((p_read + 2) & 0x3F)
		}
	} else if (c & 0xF8) == 0xF0 {
		/* If first byte begins with binary 11110 it is four byte encoding */
		/* restrict characters to UTF-8 range (U+0000 - U+10FFFF)*/
		if c >= 0xF5 {
			return UNICODE_ERROR_RESTRICTED_CHARACTER
		}
		/* check we have at least four bytes */
		if length < 4 {
			unicode_len = UNICODE_ERROR_CHARACTERS_MISSING
		} else if ((p_read + 1) & 0xC0) != 0x80 {
			unicode_len = UNICODE_ERROR_INVALID_ENCODING
		} else if ((p_read + 2) & 0xC0) != 0x80 {
			unicode_len = UNICODE_ERROR_INVALID_ENCODING
		} else if ((p_read + 3) & 0xC0) != 0x80 {
			unicode_len = UNICODE_ERROR_INVALID_ENCODING
		} else {
			unicode_len = 4
			/* compute character number */
			d = ((c & 0x07) << 18) | (((p_read + 1) & 0x3F) << 12) | (((p_read + 2) & 0x3F) << 6) | ((p_read + 3) & 0x3F)
		}
	} else {
		/* any other first byte is invalid (RFC 3629) */
		return UNICODE_ERROR_INVALID_ENCODING
	}

	/* invalid UTF-8 character number range (RFC 3629) */
	if (d >= 0xD800) && (d <= 0xDFFF) {
		return UNICODE_ERROR_RESTRICTED_CHARACTER
	}

	/* check for overlong */
	if (unicode_len == 4) && (d < 0x010000) {
		/* four byte could be represented with less bytes */
		return UNICODE_ERROR_OVERLONG_CHARACTER
	} else if (unicode_len == 3) && (d < 0x0800) {
		/* three byte could be represented with less bytes */
		return UNICODE_ERROR_OVERLONG_CHARACTER
	} else if (unicode_len == 2) && (d < 0x80) {
		/* two byte could be represented with less bytes */
		return UNICODE_ERROR_OVERLONG_CHARACTER
	}

	return unicode_len
}
