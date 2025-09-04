/** 
 * @param {number} size 
 * @returns {string} */
export const formatFileSize = (size) => {
  const kb = 1024
  if (size < kb) {
    return size + 'b'
  }
  if (size < kb * kb) {
    return (size / kb).toFixed(2) + 'Kb'
  }
  if (size < kb * kb * kb) {
    return (size / (kb * kb)).toFixed(2) + 'Mb'
  }
  if (size < kb * kb * kb * kb) {
    return (size / (kb * kb * kb)).toFixed(2) + 'Gb'
  }
  return size + 'b'
}