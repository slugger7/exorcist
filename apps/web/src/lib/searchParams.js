import { navigate } from "svelte-routing";

/**
  * @param {string} param
  * @param {number} def
  * @returns {number}
  */
export const getIntSearchParamOrDefault = (param, def) => {
  const params = new URLSearchParams(window.location.search);
  const val = parseInt(params.get(param));
  if (isNaN(val)) {
    return def;
  }
  return val;
};

/**
 * @param {string} param
 * @param {string} def
 * @returns {string}
 */
export const getStringSearchParamOrDefault = (param, def) => {
  const params = new URLSearchParams(window.location.search)

  const val = params.get(param)

  if (!val) {
    return def
  }

  return val
}

/**
 * @param {string} param
 * @param {string[]} def
 * @returns {string[]}
 */
export const getArrayOfStringsSearchParamOrDefault = (param, def) => {
  const params = new URLSearchParams(window.location.search)

  const val = params.getAll(param)

  if (!val) {
    return def
  }

  return val
}

/**
 * @param {string} param
 * @param {boolean} def
 * @returns {boolean}
 */
export const getBoolParamOrDefault = (param, def) => {
  const params = new URLSearchParams(window.location.search)

  const val = params.get(param)

  if (val === 'false') return false
  if (val === 'true') return true

  return def
}

/** 
 * @param {string} key 
 * @param {string|boolean} val 
 * @param {string} route 
 * @param {{replace?: boolean, preserveScroll?: boolean}} [options]
 */
export const setValueAndNavigate = (key, val, route, options) => {
  const params = new URLSearchParams(window.location.search);
  if (val === "") {
    params.delete(key);
  } else {
    params.set(key, val.toString());
  }
  navigate(`${route}?${params.toString()}`, options);
};

/**
 * @param {string} key
 * @param {string[]} values
 * @param {string} route
 * @param {{replace?: boolean, preserveScroll?: boolean}} [options]
 */
export const setValuesAndNavigate = (key, values, route, options) => {
  const params = new URLSearchParams(window.location.search)
  params.delete(key)

  values.forEach(v => {
    params.append(key, v)
  })

  navigate(`${route}?${params.toString()}`, options)
}