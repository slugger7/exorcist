/** @import { Validator, ValidationHandler } from "../types"*/

/** @type {ValidationHandler}*/
export const handleValidation = ({ value, validators, state }) => {
  const errors = validators?.map(validator => validator(value, state)).filter(error => error !== undefined)

  return errors && errors.length > 0 ? errors : undefined
};
