/**
 * @param {string} val value
 * @param {number} def default
 * @returns {number}
 */
export const tryParseIntOrDefault = (val, def) => {
 const i = parseInt(val)
 if (isNaN(i)) {
	return def
 }
 return i
}