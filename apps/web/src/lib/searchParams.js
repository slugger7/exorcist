import { navigate } from "svelte-routing";

/**
  * @param {string} key
  * @param {number} def
  * @returns {number}
  */
export const getIntSearchParamOrDefault = (key, def) => {
  const params = new URLSearchParams(window.location.search);
  const param = params.get(key)
  if (!param) {
    return def
  }
  const val = parseInt(param);
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
  const url = new URL(window.location.href)
  if (val === "") {
    url.searchParams.delete(key);
  } else {
    url.searchParams.set(key, val.toString());
  }

  if (!options?.replace) {
    navigate(`${route}?${url.searchParams.toString()}`, options);
  } else {
    window.history.replaceState({}, '', url)
  }
};

/**
 * @param {string} key
 * @param {string[]} values
 * @param {string} route
 * @param {{replace?: boolean, preserveScroll?: boolean}} [options]
 */
export const setValuesAndNavigate = (key, values, route, options) => {
  const url = new URL(window.location.href)
  url.searchParams.delete(key)

  values.forEach(v => {
    url.searchParams.append(key, v)
  })


  if (!options?.replace) {
    navigate(`${route}?${url.searchParams.toString()}`, options)
  } else {
    window.history.replaceState({}, '', url)
  }
}
