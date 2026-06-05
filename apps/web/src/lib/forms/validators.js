/** @import { Validator } from "../types"
 * 
 * @param {string} fieldname
 * @returns {Validator}
 */
export const numberValidator = (fieldname) => (value, _state) => {
  if (isNaN(value)) {
    return `${fieldname} should be a number`
  }
}
