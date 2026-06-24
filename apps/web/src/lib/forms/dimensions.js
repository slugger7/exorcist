/**
 * @param {number} currentHeight 
 * @param {number} currentWidth 
 * @param {number} wantedHeight 
 * @returns number
 */
export const calculateScaledWidth = (currentHeight, currentWidth, wantedHeight) => {
  if (currentHeight === 0) {
    return currentWidth
  }
  return Math.round(currentWidth / currentHeight * wantedHeight)
}

/**
 * @param {number} currentHeight 
 * @param {number} currentWidth 
 * @param {number} wantedWidth 
 * @returns number
 */
export const calculateScaledHeight = (currentHeight, currentWidth, wantedWidth) => {
  if (currentWidth === 0) {
    return currentHeight
  }
  return Math.round(currentHeight / currentWidth * wantedWidth)
}
